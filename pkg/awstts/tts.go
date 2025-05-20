package awstts

import (
	"context"
	"encoding/json"
	"errors"
	"io"

	"github.com/sifaserdarozen/talky/pkg/tts"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/polly"
	"github.com/aws/aws-sdk-go-v2/service/polly/types"
)

type awsTimestamp struct {
	Word string `json:"value"`
	Time uint64 `json:"time"`
}

type awsTts struct {
	client *polly.Client
}

func NewAwsTts() (tts.Tts, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return nil, err
	}

	client := polly.NewFromConfig(cfg)
	if nil == client {
		return nil, errors.New("error creating Aws Polly client")
	}
	return &awsTts{client: client}, nil
}

func (awsTts awsTts) Synthesize(ctx context.Context, text string) (*tts.TtsAudio, error) {
	if text == "" {
		return nil, errors.New("nothing to synthesize")
	}

	synthesizeInputAudio := &polly.SynthesizeSpeechInput{
		Text:         &text,
		TextType:     types.TextTypeText,
		VoiceId:      types.VoiceIdJoanna,
		OutputFormat: types.OutputFormatMp3,
		LanguageCode: types.LanguageCodeEnUs,
	}
	synthesizeOutputAudio, err := awsTts.client.SynthesizeSpeech(ctx, synthesizeInputAudio)
	if err != nil {
		return nil, err
	}

	audio, err := io.ReadAll(synthesizeOutputAudio.AudioStream) // It reads everything.
	if err != nil {
		return nil, err
	}

	ttsAudio := tts.TtsAudio{Audio: audio, Timestamps: nil}

	synthesizeInputTimestamp := synthesizeInputAudio
	synthesizeInputTimestamp.OutputFormat = types.OutputFormatJson
	synthesizeInputTimestamp.SpeechMarkTypes = []types.SpeechMarkType{types.SpeechMarkTypeWord}

	synthesizeOutputTimestamp, err := awsTts.client.SynthesizeSpeech(ctx, synthesizeInputTimestamp)
	if err != nil {
		return nil, err
	}

	dec := json.NewDecoder(synthesizeOutputTimestamp.AudioStream)
	for dec.More() {
		var ts awsTimestamp
		// decode an array value (Message)
		err := dec.Decode(&ts)
		if err != nil {
			return nil, err
		}
		ttsAudio.Timestamps = append(ttsAudio.Timestamps, tts.Timestamp{Word: ts.Word, TimeInMs: ts.Time})
	}

	return &ttsAudio, nil
}
