package tts

import (
	"context"
)

type Timestamp struct {
	Word string
	Time uint64
}

type TtsAudio struct {
	Audio      []byte
	Timestamps []Timestamp
}

type Tts interface {
	Synthesize(ctx context.Context, text string) (TtsAudio, error)
}
