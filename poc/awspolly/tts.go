package main

import (
	"bufio"
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/polly"
	"github.com/aws/aws-sdk-go-v2/service/polly/types"
)

func SynthesizeSSML(ctx context.Context, client *polly.Client, ssml string, outputFile string) error {
	if ssml == "" {
		return errors.New("nothing to synthesize")
	}

	textFileName := outputFile + ".txt"
	err := os.WriteFile(textFileName, []byte(ssml), 0644)
	if err != nil {
		return err
	}

	synthesizeInputAudio := &polly.SynthesizeSpeechInput{
		Text:         &ssml,
		TextType:     types.TextTypeSsml,
		VoiceId:      types.VoiceIdJoanna,
		OutputFormat: types.OutputFormatMp3,
		LanguageCode: types.LanguageCodeEnUs,
	}
	synthesizeOutputAudio, err := client.SynthesizeSpeech(ctx, synthesizeInputAudio)
	if err != nil {
		return err
	}

	audio, err := io.ReadAll(synthesizeOutputAudio.AudioStream) // It reads everything.
	if err != nil {
		return err
	}

	audioFileName := outputFile + ".mp3"
	err = os.WriteFile(audioFileName, audio, 0644)
	if err != nil {
		return err
	}

	synthesizeInputTimestamp := synthesizeInputAudio
	synthesizeInputTimestamp.OutputFormat = types.OutputFormatJson
	synthesizeInputTimestamp.SpeechMarkTypes = []types.SpeechMarkType{types.SpeechMarkTypeWord}

	synthesizeOutputTimestamp, err := client.SynthesizeSpeech(ctx, synthesizeInputTimestamp)
	if err != nil {
		return err
	}

	timestampFileName := outputFile + ".json"

	timestamps, err := io.ReadAll(synthesizeOutputTimestamp.AudioStream) // It reads everything.
	if err != nil {
		return err
	}
	err = os.WriteFile(timestampFileName, timestamps, 0644)
	if err != nil {
		return err
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

	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatalf("failed to load config, %v", err)
	}

	client := polly.NewFromConfig(cfg)
	if nil == client {
		log.Fatalf("Error creating client")
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

		ssml := convertToSSML(text)
		log.Printf(" with ssml %s", ssml)

		startTime := time.Now()
		err := SynthesizeSSML(ctx, client, ssml, fmt.Sprintf("%s-%d", dir, idx))
		ellapsedTime := time.Since(startTime)

		if err != nil {
			log.Printf("Error %s", err.Error())
		} else {
			log.Printf("Done in %d", ellapsedTime)
		}
	}
}

/*
func main() {

	if nil == dbCfg.Url {
		cfg, err = config.LoadDefaultConfig(context.Background(), config.WithRetryer(func() aws.Retryer {
			return retry.NewStandard(func(o *retry.StandardOptions) {
				o.RateLimiter = ratelimit.None
				o.MaxAttempts = 1
			})
		}))
	} else {
		provider := credentials.NewStaticCredentialsProvider("someaccess", "somesecret", "sometoken")
		cfg, err = config.LoadDefaultConfig(context.Background(), config.WithCredentialsProvider(provider), config.WithRegion("local"), config.WithBaseEndpoint(*dbCfg.Url))
	}

}
*/
