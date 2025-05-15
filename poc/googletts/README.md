
# Google TTS

Experiment on [Google TTS](https://cloud.google.com/text-to-speech?_gl=1*3zuf70*_up*MQ..&gclid=Cj0KCQjwoZbBBhDCARIsAOqMEZX2yCo9e4MjZL3VUU9YVXjRlgKdsSBAn_wvTNXIu-2oiPUW0jUjLQQaAlv6EALw_wcB&gclsrc=aw.ds&hl=en) to get synthesis timestamps. For this functionality input text should be provided as [ssml](https://cloud.google.com/text-to-speech/docs/ssml) type and requested timepoints should be [**mark**ed](https://cloud.google.com/text-to-speech/docs/ssml#ssml_timepoints)

Timestamps are generated basicly with option  
```
EnableTimePointing: []string{"SSML_MARK"}
```

to texttospeech.SynthesizeSpeechRequest  

For input text  
```
What the bleep do we know?
```
and corresponding ssml translation
```
'<speak><s>
         <mark name="timepoint_1"/>What 
         <mark name="timepoint_2"/>the 
         <mark name="timepoint_3"/>bleep 
         <mark name="timepoint_4"/>do 
         <mark name="timepoint_5"/>i 
         <mark name="timepoint_6"/>know?
         </s></speak>'
```

Sample output will be  
```
  "timepoints": [
    {
      "markName": "timepoint_1",
      "timeSeconds": 0.014999999664723873
    },
    {
      "markName": "timepoint_2",
      "timeSeconds": 0.25499999523162842
    },
    {
      "markName": "timepoint_3",
      "timeSeconds": 0.37000000476837158
    },
    {
      "markName": "timepoint_4",
      "timeSeconds": 0.7371666431427002
    },
    {
      "markName": "timepoint_5",
      "timeSeconds": 0.89358323812484741
    },
    {
      "markName": "timepoint_6",
      "timeSeconds": 1.0185831785202026
    }
  ],
```





