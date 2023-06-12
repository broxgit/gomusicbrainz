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

	// Lookup artist by id.
	artist, err := client.LookupArtist("10adbe5e-a2c0-4bf3-8249-2b4cbf6e6ca8")
	if err != nil {
		log.Error().Err(err)
		return
	}

	log.Info().Msgf("%+v", artist)
}
