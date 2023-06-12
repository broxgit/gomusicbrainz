package main

import (
	"github.com/broxgit/gomusicbrainz"
	"github.com/rs/zerolog/log"
)

func main() {
	// create a new WS2Client.
	client, _ := gomusicbrainz.NewWS2Client(
		"https://musicbrainz.org/ws/2",
		"A GoMusicBrainz example",
		"0.0.1-beta",
		"http://github.com/broxgit/gomusicbrainz")

	// Search for some artist(s)
	resp, _ := client.SearchArtist(`"Parov Stelar"`, nil, -1, -1)

	// Pretty print Name and score of each returned artist.
	for _, artist := range resp.Artists {
		log.Info().Msgf("Name: %-25sScore: %d\n", artist.Name, resp.Scores[artist])
	}
}
