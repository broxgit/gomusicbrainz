package gomusicbrainz

import "encoding/xml"

// Place represents a building or outdoor area used for performing or producing
// music.
type Place struct {
	ID          MBID          `xml:"id,attr"`
	Type        string        `xml:"type,attr"`
	Name        string        `xml:"name"`
	Address     string        `xml:"address"`
	Coordinates MBCoordinates `xml:"coordinates"`
	Area        Area          `xml:"area"`
	Lifespan    Lifespan      `xml:"life-span"`
	Aliases     []*Alias      `xml:"alias-list>alias"`
}

func (mbe *Place) lookupResult() interface{} {
	var res struct {
		XMLName xml.Name `xml:"metadata"`
		Ptr     *Place   `xml:"place"`
	}
	res.Ptr = mbe
	return &res
}

func (mbe *Place) apiEndpoint() string {
	return "/place"
}

func (mbe *Place) Id() MBID {
	return mbe.ID
}

// LookupPlace performs a place lookup request for the given MBID.
func (c *WS2Client) LookupPlace(id MBID, inc ...string) (*Place, error) {
	a := &Place{ID: id}
	err := c.Lookup(a, inc...)

	return a, err
}

// SearchPlace queries MusicBrainz´ Search Server for Places.
//
// Possible search fields to provide in searchTerm are:
//
//	pid       the place ID
//	address   the address of this place
//	alias     the aliases/misspellings for this area
//	area      area name
//	begin     place begin date
//	comment   disambiguation comment
//	end       place end date
//	ended     place ended
//	lat       place latitude
//	long      place longitude
//	sortname  place sort name
//	type      the aliases/misspellings for this place
//
// With no fields specified searchTerm searches the place, alias, address and
// area fields. For more information visit
// https://musicbrainz.org/doc/Development/XML_Web_Service/Version_2/Search#Place
func (c *WS2Client) SearchPlace(query string, fields map[string]string, limit, offset int) (*PlaceSearchResponse, error) {
	result := placeListResult{}
	err := c.searchRequestAdvanced("/place", query, fields, &result, limit, offset)

	rsp := PlaceSearchResponse{}
	rsp.WS2ListResponse = result.PlaceList.WS2ListResponse
	rsp.Scores = make(ScoreMap)

	for i, v := range result.PlaceList.Places {
		rsp.Places = append(rsp.Places, v.Place)
		rsp.Scores[rsp.Places[i]] = v.Score
	}

	return &rsp, err
}

// PlaceSearchResponse is the response type returned by the SearchPlace method.
type PlaceSearchResponse struct {
	WS2ListResponse
	Places []*Place
	Scores ScoreMap
}

// ResultsWithScore returns a slice of Places with a min score.
func (r *PlaceSearchResponse) ResultsWithScore(score int) []*Place {
	var res []*Place
	for _, v := range r.Places {
		if r.Scores[v] >= score {
			res = append(res, v)
		}
	}
	return res
}

type placeListResult struct {
	PlaceList struct {
		WS2ListResponse
		Places []struct {
			*Place
			Score int `xml:"http://musicbrainz.org/ns/ext#-2.0 score,attr"`
		} `xml:"place"`
	} `xml:"place-list"`
}
