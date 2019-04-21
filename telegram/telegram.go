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

func sendMessage(chat Replyable, text *string) (*Message, error) {
	sendMessage := method("sendMessage")
	data := struct {
		ChatID string `json:"chat_id"`
		Text   string `json:"text"`
	}{chat.GetChatID(), *text}
	message := Message{}
	if err := sendMessage.execute(&data, &message); err != nil {
		return nil, err
	}
	return &message, nil
}
