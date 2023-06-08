package gomusicbrainz

import (
	"reflect"
	"testing"
)

func TestSearchReleaseGroup(t *testing.T) {
	want := ReleaseGroupSearchResponse{
		WS2ListResponse: WS2ListResponse{
			Count:  1,
			Offset: 0,
		},
		ReleaseGroups: []*ReleaseGroup{
			{
				ID:          "70664047-2545-4e46-b75f-4556f2a7b83e",
				Type:        "Single",
				Title:       "Main Tenance",
				PrimaryType: "Single",
				ArtistCredit: ArtistCredit{
					NameCredits: []NameCredit{
						{
							Artist{
								ID:             "a8fa58d8-f60b-4b83-be7c-aea1af11596b",
								Name:           "Fred Giannelli",
								SortName:       "Giannelli, Fred",
								Disambiguation: "US electronic artist",
							},
						},
					},
				},
				Releases: []*Release{
					{
						ID:    "9168f4cc-a852-4ba5-bf85-602996625651",
						Title: "Main Tenance",
					},
				},
				Tags: []*Tag{
					{
						Count: 1,
						Name:  "electronic",
					},
					{
						Count: 1,
						Name:  "electronica",
					},
				},
			},
		},
	}

	setupHTTPTesting()
	defer server.Close()
	serveTestFile("/release-group", "SearchReleaseGroup.xml", t)

	returned, err := client.SearchReleaseGroup("Tenance", -1, -1)
	if err != nil {
		t.Error(err)
	}

	want.Scores = ScoreMap{
		returned.ReleaseGroups[0]: 100,
	}

	if !reflect.DeepEqual(*returned, want) {
		t.Error(requestDiff(&want, returned))
	}
}
