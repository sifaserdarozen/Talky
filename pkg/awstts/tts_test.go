package awstts

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/sifaserdarozen/talky/pkg/tts"

	"github.com/aws/aws-sdk-go-v2/service/polly"

	"github.com/stretchr/testify/require"
)

type requestKey struct {
	outputFormat string
	text         string
}

func TestAwsSynthesize(t *testing.T) {
	text := "What the bleep do we know"

	expectedAudio := []byte{10, 20, 30, 40}

	timestamps := []awsTimestamp{
		{Word: "What", Time: 6, Type: "word", Start: 28, End: 32},
		{Word: "the", Time: 141, Type: "word", Start: 54, End: 57},
		{Word: "bleep", Time: 213, Type: "word", Start: 79, End: 84},
		{Word: "do", Time: 520, Type: "word", Start: 106, End: 108},
		{Word: "we", Time: 656, Type: "word", Start: 130, End: 132},
		{Word: "know", Time: 796, Type: "word", Start: 154, End: 158}}

	expectedTimestamps := []tts.Timestamp{
		{Word: "What", TimeInMs: 6},
		{Word: "the", TimeInMs: 141},
		{Word: "bleep", TimeInMs: 213},
		{Word: "do", TimeInMs: 520},
		{Word: "we", TimeInMs: 656},
		{Word: "know", TimeInMs: 796}}

	buf := new(bytes.Buffer)
	enc := json.NewEncoder(buf)
	for _, v := range timestamps {
		err := enc.Encode(v)
		require.NoError(t, err, "Should encode timestamps")
	}

	bytes, err := io.ReadAll(buf)
	require.NoError(t, err, "Readall should be successful")

	responseMap := map[requestKey][]byte{
		{outputFormat: "json", text: text}: bytes,
		{outputFormat: "mp3", text: text}:  expectedAudio,
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		synthesizeInputAudio := polly.SynthesizeSpeechInput{}
		dec := json.NewDecoder(r.Body)
		dec.DisallowUnknownFields()
		err := dec.Decode(&synthesizeInputAudio)
		if err != nil {
			t.Errorf("Unable to decode speech synthesis request %v", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		response, isOk := responseMap[requestKey{outputFormat: string(synthesizeInputAudio.OutputFormat), text: *synthesizeInputAudio.Text}]

		if !isOk {
			fmt.Println(*synthesizeInputAudio.Text)
			fmt.Println(synthesizeInputAudio.TextType)
			fmt.Println(synthesizeInputAudio.OutputFormat)
			t.Errorf("Unable to find request in the map %v", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		w.WriteHeader(http.StatusOK)
		_, err = w.Write(response)
		if err != nil {
			t.Errorf("Unable to write data %v", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

	}))

	defer server.Close()

	tts, err := NewAwsTts(&server.URL)
	require.NoError(t, err, "Client should be created")

	audio, err := tts.Synthesize(context.TODO(), text)
	require.NoError(t, err, "Synthesis response should be received")

	require.Equal(t, expectedTimestamps, audio.Timestamps)
	require.Equal(t, expectedAudio, audio.Audio)

}
