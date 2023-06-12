package gomusicbrainz

import (
	"reflect"
	"testing"
)

func TestSearchRecording(t *testing.T) {
	want := RecordingSearchResponse{
		WS2ListResponse: WS2ListResponse{
			Count:  1,
			Offset: 0,
		},
		Recordings: []*Recording{
			{
				ID:     "07339604-c19c-4efe-9195-f9c3b127a458",
				Title:  "Fred",
				Length: 473000,
				ArtistCredit: ArtistCredit{
					NameCredits: []NameCredit{
						{
							Artist{
								ID:       "695e75b5-c6db-43ee-abeb-2f3e50d96c3e",
								Name:     "Imperiet",
								SortName: "Imperiet",
							},
						},
					},
				},
			},
		},
	}

	setupHTTPTesting()
	defer server.Close()
	serveTestFile("/recording", "SearchRecording.xml", t)

	returned, err := client.SearchRecording("Fred", nil, -1, -1)
	if err != nil {
		t.Error(err)
	}

	want.Scores = ScoreMap{
		returned.Recordings[0]: 100,
	}

	if !reflect.DeepEqual(*returned, want) {
		t.Error(requestDiff(&want, returned))
	}
}
