package gomusicbrainz

// CDStub represents an anonymously submitted track list.
type CDStub struct {
	ID        string `xml:"id,attr"` // seems not to be a valid MBID (UUID)
	Title     string `xml:"title"`
	Artist    string `xml:"artist"`
	Barcode   string `xml:"barcode"`
	Comment   string `xml:"comment"`
	TrackList struct {
		Count int `xml:"count,attr"`
	} `xml:"track-list"`
}

// SearchCDStub queries MusicBrainz´ Search Server for CDStubs.
//
// Possible search fields to provide in searchTerm are:
//
//	artist   artist name
//	title    release name
//	barcode  release barcode
//	comment  general comments about the release
//	tracks   number of tracks on the CD stub
//	discid   disc ID of the CD
//
// With no fields specified searchTerm searches only the artist Field. For more
// information visit
// https://musicbrainz.org/doc/Development/XML_Web_Service/Version_2/Search#CDStubs
func (c *WS2Client) SearchCDStub(query string, fields map[string]string, limit, offset int) (*CDStubSearchResponse, error) {
	result := cdStubListResult{}
	err := c.searchRequestAdvanced("/cdstub", query, fields, &result, limit, offset)

	rsp := CDStubSearchResponse{}
	rsp.WS2ListResponse = result.CDStubList.WS2ListResponse
	rsp.Scores = make(ScoreMap)

	for i, v := range result.CDStubList.CDStubs {
		rsp.CDStubs = append(rsp.CDStubs, v.CDStub)
		rsp.Scores[rsp.CDStubs[i]] = v.Score
	}

	return &rsp, err
}

// CDStubSearchResponse is the response type returned by the SearchCDStub method.
type CDStubSearchResponse struct {
	WS2ListResponse
	CDStubs []*CDStub
	Scores  ScoreMap
}

// ResultsWithScore returns a slice of CDStubs with a min score.
func (r *CDStubSearchResponse) ResultsWithScore(score int) []*CDStub {
	var res []*CDStub
	for _, v := range r.CDStubs {
		if r.Scores[v] >= score {
			res = append(res, v)
		}
	}
	return res
}

type cdStubListResult struct {
	CDStubList struct {
		WS2ListResponse
		CDStubs []struct {
			*CDStub
			Score int `xml:"http://musicbrainz.org/ns/ext#-2.0 score,attr"`
		} `xml:"cdstub"`
	} `xml:"cdstub-list"`
}
