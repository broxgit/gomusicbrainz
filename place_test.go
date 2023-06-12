package gomusicbrainz

import (
	"reflect"
	"testing"
	"time"
)

func TestSearchPlace(t *testing.T) {
	want := PlaceSearchResponse{
		WS2ListResponse: WS2ListResponse{
			Count:  1,
			Offset: 0,
		},
		Places: []*Place{
			{
				ID:          "d1ab65f8-d082-492a-bd70-ce375548dabf",
				Type:        "Studio",
				Name:        "Chipping Norton Recording Studios",
				Address:     "28–30 New Street, Chipping Norton",
				Coordinates: MBCoordinates{}, // TODO cover
				Area: Area{
					ID:       "44e5e20e-8fbc-4b07-b3f2-22f2199186fd",
					Name:     "Oxfordshire",
					SortName: "Oxfordshire",
				},
				Lifespan: Lifespan{
					Begin: BrainzTime{
						Time:     time.Date(1971, 1, 1, 0, 0, 0, 0, time.UTC),
						Accuracy: Year,
					},
					End: BrainzTime{
						Time:     time.Date(1999, 10, 1, 0, 0, 0, 0, time.UTC),
						Accuracy: Month,
					},
					Ended: true,
				},
				// TODO Aliases: []*Alias
			},
		},
	}

	setupHTTPTesting()
	defer server.Close()
	serveTestFile("/place", "SearchPlace.xml", t)

	returned, err := client.SearchPlace("chipping", nil, -1, -1)
	if err != nil {
		t.Error(err)
	}

	want.Scores = ScoreMap{
		returned.Places[0]: 100,
	}

	if !reflect.DeepEqual(*returned, want) {
		t.Error(requestDiff(&want, returned))
	}
}
