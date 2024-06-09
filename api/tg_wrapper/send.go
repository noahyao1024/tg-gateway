package tg_wrapper

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
)

type SendVoiceOption struct {
	ChatID             int64
	FileName           string
	EnableNotification bool
}

func (i *Instance) SendVoice(svo *SendVoiceOption) (string, error) {
	url := fmt.Sprintf("https://api.telegram.org/bot%s/sendVoice", i.BotKey)
	method := "POST"

	payload := &bytes.Buffer{}
	writer := multipart.NewWriter(payload)
	file, errFile1 := os.Open(svo.FileName)
	defer file.Close()
	part1,
		errFile1 := writer.CreateFormFile("voice", filepath.Base(svo.FileName))
	_, errFile1 = io.Copy(part1, file)
	if errFile1 != nil {
		return "", errFile1
	}
	_ = writer.WriteField("chat_id", fmt.Sprintf("%d", svo.ChatID))
	_ = writer.WriteField("caption", "")
	_ = writer.WriteField("parse_mode", "MarkdownV2")
	if !svo.EnableNotification {
		_ = writer.WriteField("disable_notification", "true")
	}
	err := writer.Close()
	if err != nil {
		return "", err
	}

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())
	res, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return "", err
	}
	return string(body), nil
}
