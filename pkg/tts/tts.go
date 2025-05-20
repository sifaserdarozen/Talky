package tts

import (
	"context"
)

type Timestamp struct {
	Word     string `json:"value"`
	TimeInMs uint64 `json:"time"`
}

type TtsAudio struct {
	Audio      []byte
	Timestamps []Timestamp
}

type Tts interface {
	Synthesize(ctx context.Context, text string) (*TtsAudio, error)
}
