package gomusicbrainz

import (
	"reflect"
	"testing"
)

func TestSearchCDStub(t *testing.T) {
	want := CDStubSearchResponse{
		WS2ListResponse: WS2ListResponse{
			Count:  1,
			Offset: 0,
		},
		CDStubs: []*CDStub{
			{
				ID:      "vi44WFVS5zRT2svM.PORcEm9LJk-",
				Title:   "Silent Conflict (Live @ The Hard Rock Cafe)",
				Artist:  "Bonobo",
				Barcode: "634479355059",
				Comment: "CD Baby id:bonobo",
				TrackList: struct {
					Count int `xml:"count,attr"`
				}{
					Count: 3,
				},
			},
		},
	}

	setupHTTPTesting()
	defer server.Close()
	serveTestFile("/cdstub", "SearchCDStub.xml", t)

	returned, err := client.SearchCDStub(`bonobo`, -1, -1)
	if err != nil {
		t.Error(err)
	}

	want.Scores = ScoreMap{
		returned.CDStubs[0]: 100,
	}

	if !reflect.DeepEqual(*returned, want) {
		t.Error(requestDiff(&want, returned))
	}
}
