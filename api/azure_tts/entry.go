package azure_tts

import (
	"fmt"
	"io"
	"net/http"
	"strings"
)

func TTS(key, text string) ([]byte, error) {
	url := "https://eastasia.tts.speech.microsoft.com/cognitiveservices/v1"
	method := "POST"

	payload := strings.NewReader(fmt.Sprintf(`<speak version='1.0' xml:lang='en-US'><voice xml:lang='en-US' xml:gender='Male'
	  name='en-US-JennyNeural'>
		  %s
  </voice></speak>`, text))

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		return nil, err
	}
	req.Header.Add("Ocp-Apim-Subscription-Key", key)
	req.Header.Add("Content-Type", "application/ssml+xml")
	req.Header.Add("X-Microsoft-OutputFormat", "audio-24khz-160kbitrate-mono-mp3")

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}
