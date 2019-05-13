package telegram

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/the-fusy/rentit/config"
)

type method string

func (method *method) execute(data interface{}, answer interface{}) error {
	telegramURL := fmt.Sprintf("https://api.telegram.org/bot%v/%v", config.BotToken, string(*method))

	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	res, err := http.Post(telegramURL, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	defer res.Body.Close()

	var result struct {
		Ok          bool        `json:"ok"`
		Description string      `json:"description"`
		Result      interface{} `json:"result"`
	}
	result.Result = answer

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return err
	}

	if !result.Ok {
		return fmt.Errorf("Telegram api error: %v", result.Description)
	}

	return nil
}

func GetUpdates() (*[]Update, error) {
	getUpdates := method("getUpdates")
	updates := []Update{}
	if err := getUpdates.execute(nil, &updates); err != nil {
		return nil, err
	}
	return &updates, nil
}

func SendMessage(chat Replyable, text *string, pagePreview, notifications bool) (*Message, error) {
	sendMessage := method("sendMessage")
	data := struct {
		ChatID                string `json:"chat_id"`
		Text                  string `json:"text"`
		DisableWebPagePreview bool   `json:"disable_web_page_preview"`
		DisableNotifications  bool   `json:"disable_notification"`
	}{chat.GetChatID(), *text, !pagePreview, !notifications}
	message := Message{}
	if err := sendMessage.execute(&data, &message); err != nil {
		return nil, err
	}
	return &message, nil
}

func SendPhotos(chat Replyable, photos *[]string) (*Message, error) {
	sendPhotos := method("sendMediaGroup")
	type mediaStruct struct {
		Type    string `json:"type"`
		Media   string `json:"media"`
		Caption string `json:"caption"`
	}
	medias := make([]mediaStruct, 0)
	for i := range *photos {
		if i == 10 {
			break
		}
		medias = append(medias, mediaStruct{"photo", (*photos)[i], "KEKOS"})
	}
	data := struct {
		ChatID string        `json:"chat_id"`
		Media  []mediaStruct `json:"media"`
	}{chat.GetChatID(), medias}
	message := Message{}
	if err := sendPhotos.execute(&data, &message); err != nil {
		return nil, err
	}
	return &message, nil
}
