package main

import (
	"bufio"
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	texttospeech "google.golang.org/api/texttospeech/v1beta1"
)

func SynthesizeSSML(ctx context.Context, service *texttospeech.TextService, ssml string, outputFile string) error {
	if ssml == "" {
		return errors.New("nothing to synthesize")
	}
	req := texttospeech.SynthesizeSpeechRequest{
		Input: &texttospeech.SynthesisInput{
			Ssml: string(ssml),
		},
		Voice: &texttospeech.VoiceSelectionParams{
			LanguageCode: "en-US",
			SsmlGender:   "FEMALE",
		},
		AudioConfig: &texttospeech.AudioConfig{
			AudioEncoding: "MP3",
		},
		EnableTimePointing: []string{"SSML_MARK"},
	}

	textFileName := outputFile + ".txt"
	err := os.WriteFile(textFileName, []byte(ssml), 0644)
	if err != nil {
		return err
	}

	resp, err := service.Synthesize(&req).Do()
	if err != nil {
		return err
	}

	audio, err := base64.StdEncoding.DecodeString(resp.AudioContent)
	if err != nil {
		return err
	}

	audioFileName := outputFile + ".mp3"
	err = os.WriteFile(audioFileName, audio, 0644)
	if err != nil {
		return err
	}

	for _, tp := range resp.Timepoints {
		jsonBytes, err := json.Marshal(tp)
		if err != nil {
			return err
		}
		fmt.Println(string(jsonBytes))
	}

	return nil
}

// Helper to convert raw text into ssml format requesting timestamps
func convertToSSML(text string) string {
	if text == "" {
		return ""
	}

	words := strings.Fields(text)

	ssml := "<speak>"
	idx := 0
	for ; idx < len(words)-1; idx++ {
		ssml += fmt.Sprintf("<mark name=\"word_%d\"/>", idx)
		ssml += words[idx]
		ssml += " "
	}
	ssml += fmt.Sprintf("<mark name=\"word_%d\"/>", idx)
	ssml += words[idx]

	ssml += "</speak>"
	return ssml
}

func main() {
	dir := *flag.String("dir", "synthfile", "Directory to create files")

	reader := bufio.NewReader(os.Stdin)
	ctx := context.Background()

	service, err := texttospeech.NewService(ctx)
	if err != nil {
		log.Fatalf("Error %s", err.Error())
	}

	texttospeechService := texttospeech.NewTextService(service)

	log.Printf("Starting tts loop...")

	for idx := 0; true; idx++ {
		fmt.Printf("%d > ", idx)
		text, _ := reader.ReadString('\n')
		// convert CRLF to LF
		text = strings.Replace(text, "\n", "", -1)
		if text == "" {
			continue
		}

		ssml := convertToSSML(text)
		log.Printf(" with ssml %s", ssml)

		startTime := time.Now()
		err := SynthesizeSSML(ctx, texttospeechService, ssml, fmt.Sprintf("%s-%d", dir, idx))
		ellapsedTime := time.Since(startTime)

		if err != nil {
			log.Printf("Error %s", err.Error())
		} else {
			log.Printf("Done in %d", ellapsedTime)
		}
	}
}
