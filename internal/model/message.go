package model

type From struct {
	IsFilled     bool
	LanguageCode string
}

type Message struct {
	IsFilled bool
	From     From
	Text     string
	ChatId   int64
}

type Update struct {
	Message Message
}

type InlineButton struct {
	Text string
	Url  string
}

type ParseMode string

const (
	ModeNone       ParseMode = ""
	ModeMarkdownV2 ParseMode = "MarkdownV2"
)

type Response struct {
	Text           string
	ParseMode      ParseMode
	InlineKeyboard [][]InlineButton
}

type ResponseWithChatId struct {
	Response
	ChatId int64
}
