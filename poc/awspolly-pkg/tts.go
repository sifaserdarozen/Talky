package main

import (
	"bufio"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/sifaserdarozen/talky/pkg/awstts"
	"github.com/sifaserdarozen/talky/pkg/tts"
)

func writeToFile(audio *tts.TtsAudio, text string, fileBase string) error {
	textFileName := fileBase + ".txt"
	err := os.WriteFile(textFileName, []byte(text), 0644)
	if err != nil {
		return err
	}

	audioFileName := fileBase + ".mp3"
	err = os.WriteFile(audioFileName, audio.Audio, 0644)
	if err != nil {
		return err
	}

	timestampsJSON, err := json.Marshal(audio.Timestamps)
	if err != nil {
		return err
	}

	timestampFileName := fileBase + ".json"
	err = os.WriteFile(timestampFileName, []byte(timestampsJSON), 0644)
	if err != nil {
		return err
	}

	return nil
}

func main() {
	dirBase := *flag.String("dir", "synthfile", "Directory to create files")

	reader := bufio.NewReader(os.Stdin)
	ctx := context.Background()

	tts, err := awstts.NewAwsTts(nil)

	if err != nil {
		log.Fatalf("failed to load config, %v", err)
	}

	log.Printf("Starting tts loop...")

	for idx := 0; true; idx++ {
		fmt.Printf("%d > ", idx)
		text, _ := reader.ReadString('\n')
		// convert CRLF to LF
		text = strings.ReplaceAll(text, "\n", "")
		if text == "" {
			continue
		}

		startTime := time.Now()
		audio, err := tts.Synthesize(ctx, text)
		ellapsedTime := time.Since(startTime)

		if err != nil {
			log.Printf("Error %v", err)
			continue
		} else {
			log.Printf("Done in %d", ellapsedTime)
		}

		if err := writeToFile(audio, text, fmt.Sprintf("%s-%d", dirBase, idx)); err != nil {
			log.Printf("Error %v", err)
		}
	}
}
