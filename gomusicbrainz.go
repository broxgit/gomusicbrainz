/*
Package gomusicbrainz implements a MusicBrainz WS2 client library.

MusicBrainz WS2 (Version 2 of the XML Web Service) supports three different requests:

# Search requests

With search requests you can search MusicBrainz´ database for all entities.
GoMusicBrainz implements one search method for every search request in the form:

	func (*WS2Client) Search<ENTITY>(searchTerm, limit, offset) (<ENTITY>SearchResponse, error)

searchTerm follows the Apache Lucene syntax and can either contain multiple
fields with logical operators or just a simple search string. Please refer to
https://lucene.apache.org/core/4_3_0/queryparser/org/apache/lucene/queryparser/classic/package-summary.html#package_description
for more details on the lucene syntax. limit defines how many entries should be
returned (1-100, default 25). offset is used for paging through more than one
page of results. To ignore limit and/or offset, set it to -1.

# Lookup requests

You can perform a lookup of an entity when you have the MBID for that entity.
GoMusicBrainz provides two ways to perform lookup requests: Either the specific
lookup method that is implemented for each entity that has a lookup endpoint
in the form

	func(*WS2Client) Lookup<ETITY>(id MBID, inc ...string) (*<ENTITY>, error)

or the common lookup method if you already have an entity (with MBID) that
implements the MBLookupEntity interface:

	func(*WS2Client) Lookup(entity MBLookupEntity, inc ...string) error

With both methods you can include inc params which affect subqueries e.g.
relationships. see
http://musicbrainz.org/doc/Development/XML_Web_Service/Version_2#inc.3D_arguments_which_affect_subqueries
Not all of them are supported yet.

# Browse requests

not supported yet.
*/
package gomusicbrainz

import (
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"
	"strconv"
	"strings"

	common "github.com/broxgit/common/http"
)

// NewWS2Client returns a new instance of WS2Client. Please provide meaningful
// information about your application as described at
// https://musicbrainz.org/doc/XML_Web_Service/Rate_Limiting#Provide_meaningful_User-Agent_strings
func NewWS2Client(wsURL, appName, version, contact string) (*WS2Client, error) {
	c := WS2Client{
		Retries: 10,
		Backoff: 2,
	}
	var err error

	c.WS2RootURL, err = url.Parse(wsURL)
	if err != nil {
		return nil, err
	}
	if !strings.HasSuffix(c.WS2RootURL.Path, "ws/2") {
		c.WS2RootURL.Path = path.Join(c.WS2RootURL.Path, "ws/2")
	}
	c.userAgentHeader = appName + "/" + version + " ( " + contact + " ) "

	return &c, nil
}

// WS2Client defines a Go client for the MusicBrainz Web Service 2.
type WS2Client struct {
	WS2RootURL      *url.URL // The API root URL
	userAgentHeader string
	Retries         int
	Backoff         int
}

func (c *WS2Client) getRequest(data interface{}, params string, endpoint string) error {
	retryClient := common.NewHTTPRetry(common.WithRetries(c.Retries), common.WithBackoff(c.Backoff))

	defaultRedirectLimit := 30

	// Preserve headers on redirect
	// See: https://github.com/golang/go/issues/4800
	retryClient.HTTPClient.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		if len(via) > defaultRedirectLimit {
			return fmt.Errorf("%d consecutive requests(redirects)", len(via))
		}
		if len(via) == 0 {
			// No redirects
			return nil
		}
		// mutate the subsequent redirect requests with the first Header
		for key, val := range via[0].Header {
			req.Header[key] = val
		}
		return nil
	}

	reqUrl := *c.WS2RootURL
	reqUrl.Path = path.Join(reqUrl.Path, endpoint)
	reqUrl.RawQuery = params

	fullURL := reqUrl.String()

	req, err := http.NewRequest(http.MethodGet, fullURL, nil)
	if err != nil {
		return err
	}

	req.Header.Set("User-Agent", c.userAgentHeader)

	resp, err := retryClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	xmlData, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	return xml.Unmarshal(xmlData, data)
}

// intParamToString returns an empty string for -1.
func intParamToString(i int) string {
	if i == -1 {
		return ""
	}
	return strconv.Itoa(i)
}

func (c *WS2Client) searchRequest(endpoint string, result interface{}, searchTerm string, limit, offset int) error {
	params := url.Values{
		"query":  {searchTerm},
		"limit":  {intParamToString(limit)},
		"offset": {intParamToString(offset)},
	}

	if err := c.getRequest(result, params.Encode(), endpoint); err != nil {
		return err
	}

	return nil
}

// searchRequestAdvanced.
func (c *WS2Client) searchRequestAdvanced(endpoint string, query string, fields map[string]string, result interface{}, limit, offset int) error {
	params := url.Values{
		"limit":  {intParamToString(limit)},
		"offset": {intParamToString(offset)},
	}

	var sb strings.Builder

	if strings.TrimSpace(query) != "" {
		sb.WriteString(fmt.Sprintf("query=%s", url.QueryEscape(query)))
	}

	for k, v := range fields {
		sb.WriteString(fmt.Sprintf("&%s=%s", k, url.QueryEscape(v)))
	}

	sb.WriteString(fmt.Sprintf("&%s", params.Encode()))

	if err := c.getRequest(result, sb.String(), endpoint); err != nil {
		return err
	}

	return nil
}

func encodeInc(inc []string) url.Values {
	if inc != nil {
		return url.Values{
			"inc": {strings.Join(inc, "+")},
		}
	}
	return nil
}

// Lookup performs a WS2 lookup request for the given entity (e.g. Artist,
// Label, ...)
func (c *WS2Client) Lookup(entity MBLookupEntity, inc ...string) error {
	if entity.Id() == "" {
		return errors.New("can't perform lookup without ID")
	}

	return c.getRequest(entity.lookupResult(), encodeInc(inc).Encode(),
		path.Join(
			entity.apiEndpoint(),
			string(entity.Id()),
		),
	)
}
