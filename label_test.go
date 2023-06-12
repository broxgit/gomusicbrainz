package gomusicbrainz

import (
	"reflect"
	"testing"
	"time"
)

func TestSearchLabel(t *testing.T) {
	want := LabelSearchResponse{
		WS2ListResponse: WS2ListResponse{
			Count:  1,
			Offset: 0,
		},
		Labels: []*Label{
			{
				ID:             "c1c625b5-9929-4a30-8c3e-f77e109cdf07",
				Type:           "Original Production",
				Name:           "Compost Records",
				SortName:       "Compost Records",
				Disambiguation: "German record label established in 1994.",
				CountryCode:    "DE",
				LabelCode:      2518,
				Area: Area{
					ID:       "85752fda-13c4-31a3-bee5-0e5cb1f51dad",
					Name:     "Germany",
					SortName: "Germany",
				},
				Lifespan: Lifespan{
					Begin: BrainzTime{
						Time:     time.Date(1994, 1, 1, 0, 0, 0, 0, time.UTC),
						Accuracy: Year,
					},
					Ended: false,
				},
				Aliases: []*Alias{
					{
						Locale:   "ja",
						SortName: "コンポスト・レコーズ",
						Name:     "コンポスト・レコーズ",
						Type:     "Label name",
					},
				},
			},
		},
	}

	setupHTTPTesting()
	defer server.Close()
	serveTestFile("/label", "SearchLabel.xml", t)

	returned, err := client.SearchLabel(`label:"Compost%20Records"`, nil, -1, -1)
	if err != nil {
		t.Error(err)
	}

	want.Scores = ScoreMap{
		returned.Labels[0]: 100,
	}

	if !reflect.DeepEqual(*returned, want) {
		t.Error(requestDiff(&want, returned))
	}
}
