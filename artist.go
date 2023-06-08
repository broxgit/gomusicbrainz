package gomusicbrainz

import "encoding/xml"

// Artist represents generally a musician, a group of musicians, a collaboration
// of multiple musicians or other music professionals.
type Artist struct {
	ID             MBID               `xml:"id,attr"`
	Type           string             `xml:"type,attr"`
	Name           string             `xml:"name"`
	Disambiguation string             `xml:"disambiguation"`
	SortName       string             `xml:"sort-name"`
	CountryCode    string             `xml:"country"`
	Gender         string             `xml:"gender"`
	Lifespan       Lifespan           `xml:"life-span"`
	Area           Area               `xml:"area"`
	BeginArea      Area               `xml:"begin-area"`
	Aliases        []*Alias           `xml:"alias-list>alias"`
	Tags           []Tag              `xml:"tag-list>tag"`
	Relations      TargetRelationsMap `xml:"relation-list"`
}

func (mbe *Artist) lookupResult() interface{} {
	var res struct {
		XMLName xml.Name `xml:"metadata"`
		Ptr     *Artist  `xml:"artist"`
	}
	res.Ptr = mbe
	return &res
}

func (mbe *Artist) apiEndpoint() string {
	return "/artist"
}

func (mbe *Artist) Id() MBID {
	return mbe.ID
}

// LookupArtist performs an artist lookup request for the given MBID.
func (c *WS2Client) LookupArtist(id MBID, inc ...string) (*Artist, error) {
	a := &Artist{ID: id}
	err := c.Lookup(a, inc...)

	return a, err
}

// SearchArtist queries MusicBrainz´ Search Server for Artists.
//
// Possible search fields to provide in searchTerm are:
//
//	area          artist area
//	beginarea     artist begin area
//	endarea       artist end area
//	arid          MBID of the artist
//	artist        name of the artist
//	artistaccent  name of the artist with any accent characters retained
//	alias         the aliases/misspellings for the artist
//	begin         artist birth date/band founding date
//	comment       artist comment to differentiate similar artists
//	country       the two letter country code for the artist country or 'unknown'
//	end           artist death date/band dissolution date
//	ended         true if know ended even if do not know end date
//	gender        gender of the artist (“male”, “female”, “other”)
//	ipi           IPI code for the artist
//	sortname      artist sortname
//	tag           a tag applied to the artist
//	type          artist type (“person”, “group”, "other" or “unknown”)
//
// With no fields specified searchTerm searches the artist, sortname and alias
// fields. For more information visit
// http://musicbrainz.org/doc/Development/XML_Web_Service/Version_2/Search#Artist
func (c *WS2Client) SearchArtist(searchTerm string, limit, offset int) (*ArtistSearchResponse, error) {
	result := artistListResult{}
	err := c.searchRequest("/artist", &result, searchTerm, limit, offset)

	rsp := ArtistSearchResponse{}
	rsp.WS2ListResponse = result.ArtistList.WS2ListResponse
	rsp.Scores = make(ScoreMap)

	for i, v := range result.ArtistList.Artists {
		rsp.Artists = append(rsp.Artists, v.Artist)
		rsp.Scores[rsp.Artists[i]] = v.Score
	}

	return &rsp, err
}

// ArtistSearchResponse is the response type returned by the SearchArtist method.
type ArtistSearchResponse struct {
	WS2ListResponse
	Artists []*Artist
	Scores  ScoreMap
}

// ResultsWithScore returns a slice of Artists with a min score.
func (r *ArtistSearchResponse) ResultsWithScore(score int) []*Artist {
	var res []*Artist
	for _, v := range r.Artists {
		if r.Scores[v] >= score {
			res = append(res, v)
		}
	}
	return res
}

type artistListResult struct {
	ArtistList struct {
		WS2ListResponse
		Artists []struct {
			*Artist
			Score int `xml:"http://musicbrainz.org/ns/ext#-2.0 score,attr"`
		} `xml:"artist"`
	} `xml:"artist-list"`
}
