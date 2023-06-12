package gomusicbrainz

import (
	"reflect"
	"testing"
)

func TestSearchArea(t *testing.T) {
	want := AreaSearchResponse{
		WS2ListResponse: WS2ListResponse{
			Count:  1,
			Offset: 0,
		},
		Areas: []*Area{
			{
				ID:       "d79e4501-8cba-431b-96e7-bb9976f0ae76",
				Type:     "Subdivision",
				Name:     "Île-de-France",
				SortName: "Île-de-France",
				ISO31662Codes: []ISO31662Code{
					"FR-J",
				},
				Lifespan: Lifespan{
					Ended: false,
				},
				Aliases: []Alias{
					{Locale: "et", SortName: "Île-de-France", Type: "Area name", Primary: "primary", Name: "Île-de-France"},
					{Locale: "ja", SortName: "イル＝ド＝フランス地域圏", Type: "Area name", Primary: "primary", Name: "イル＝ド＝フランス地域圏"},
				},
			},
		},
	}

	setupHTTPTesting()
	defer server.Close()
	serveTestFile("/area", "SearchArea.xml", t)

	returned, err := client.SearchArea(`"Île-de-France"`, nil, -1, -1)
	if err != nil {
		t.Error(err)
	}

	want.Scores = ScoreMap{
		returned.Areas[0]: 100,
	}

	if !reflect.DeepEqual(*returned, want) {
		t.Error(requestDiff(&want, returned))
	}
}
