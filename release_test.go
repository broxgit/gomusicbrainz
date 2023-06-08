package gomusicbrainz

import (
	"reflect"
	"testing"
	"time"
)

func TestSearchRelease(t *testing.T) {
	want := ReleaseSearchResponse{
		WS2ListResponse: WS2ListResponse{
			Count:  1,
			Offset: 0,
		},
		Releases: []*Release{
			{
				ID:     "9ab1b03e-6722-4ab8-bc7f-a8722f0d34c1",
				Title:  "Fred Schneider & The Shake Society",
				Status: "official",
				TextRepresentation: TextRepresentation{
					Language: "eng",
					Script:   "latn",
				},
				ArtistCredit: ArtistCredit{
					NameCredits: []NameCredit{
						{
							Artist{
								ID:       "43bcca8b-9edc-4997-8343-122350e790bf",
								Name:     "Fred Schneider",
								SortName: "Schneider, Fred",
							},
						},
					},
				},
				ReleaseGroup: ReleaseGroup{
					Type: "Album",
				},
				Date: BrainzTime{
					Time:     time.Date(1991, 4, 30, 0, 0, 0, 0, time.UTC),
					Accuracy: Day,
				},
				CountryCode: "us",
				Barcode:     "075992659222",
				Asin:        "075992659222",
				LabelInfos: []LabelInfo{
					{
						CatalogNumber: "9 26592-2",
						Label: &Label{
							Name: "Reprise Records",
						},
					},
				},
				Mediums: []*Medium{
					{
						Format: "cd",
					},
				},
			},
		},
	}

	setupHTTPTesting()
	defer server.Close()
	serveTestFile("/release", "SearchRelease.xml", t)

	returned, err := client.SearchRelease("Fred", -1, -1)
	if err != nil {
		t.Error(err)
	}

	want.Scores = ScoreMap{
		returned.Releases[0]: 100,
	}

	if !reflect.DeepEqual(*returned, want) {
		t.Error(requestDiff(&want, returned))
	}
}
