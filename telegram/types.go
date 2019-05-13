package telegram

import "strconv"

type Replyable interface {
	GetChatID() string
}

type User struct {
	ID        int32  `json:"id"`
	IsBot     bool   `json:"is_bot"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Username  string `json:"username"`
}

func (user *User) GetChatID() string {
	return strconv.FormatInt(int64(user.ID), 10)
}

type Chat struct {
	ID       int32  `json:"id"`
	Type     string `json:"type"`
	Username string `json:"username"`
}

func (chat *Chat) GetChatID() string {
	if chat.Username != "" {
		return "@" + chat.Username
	} else {
		return strconv.FormatInt(int64(chat.ID), 10)
	}
}

type Update struct {
	UpdateID int32   `json:"update_id"`
	Message  Message `json:"message"`
}

type Message struct {
	MessageID             int32  `json:"message_id"`
	Text                  string `json:"text"`
	From                  User   `json:"from"`
	Chat                  Chat   `json:"chat"`
	DisableWebPagePreview bool   `json:"disable_web_page_preview"`
	DisableNotifications  bool   `json:"disable_notification"`
}
