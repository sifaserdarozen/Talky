
# AWS Polly TTS

Experiment on [AWS Polly](https://aws.amazon.com/polly/) to get synthesis timestamps. As Polly does not provide both audio and timestamps
at the same time, there will be two api calls with corresponding **polly.SynthesizeSpeechInput**

Timestamps are generated basicly with option  
```
synthesizeInputTimestamp.SpeechMarkTypes = []types.SpeechMarkType{types.SpeechMarkTypeWord}
```

For input text  
```
What the bleep do we know?
```

Sample output will be  
```
{"time":6,"type":"word","start":28,"end":32,"value":"What"}
{"time":141,"type":"word","start":54,"end":57,"value":"the"}
{"time":213,"type":"word","start":79,"end":84,"value":"bleep"}
{"time":520,"type":"word","start":106,"end":108,"value":"do"}
{"time":656,"type":"word","start":130,"end":132,"value":"we"}
{"time":796,"type":"word","start":154,"end":158,"value":"know"}
```
