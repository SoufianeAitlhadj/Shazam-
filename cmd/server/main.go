package main

import (
	"fmt"
	"log"

	"shazam/internal/audio"
	"shazam/internal/fingerprint"
)

func main() {
	samples, sampleRate, err := audio.ReadWav("songs/wav_songs/Daft Punk - One More Time (Official Video).wav")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Sample rate:", sampleRate)
	fmt.Println("Total samples:", len(samples))

	fps := fingerprint.Extract(samples, sampleRate)

	fmt.Println("Fingerprints extracted:", len(fps))

}
