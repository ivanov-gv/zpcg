package message

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

type MaybeInaccessibleMessage struct {
	Id int64
}

type Callback struct {
	Id              string
	ChatId          int64
	From            From
	Message         MaybeInaccessibleMessage
	InlineMessageId string
	Data            string
}

type UpdateType int

const (
	UnsupportedUpdateType UpdateType = iota
	MessageUpdateType
	CallbackUpdateType
)

type Update struct {
	Type     UpdateType
	Message  Message
	Callback Callback
}

type InlineButtonType int

const (
	_ InlineButtonType = iota
	UrlInlineButtonType
	CallbackInlineButtonType
)

type CallbackButton struct {
	Data string
}

type UrlButton struct {
	Url string
}

type InlineButton struct {
	Type     InlineButtonType
	Text     string
	Callback CallbackButton
	Url      UrlButton
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

type ToDelete struct {
	MessageId int64
}

type ToUpdate struct {
	MessageId int64
	Response
}

type ResponseWithChatId struct {
	Send   []Response
	Delete []ToDelete
	Update []ToUpdate
	ChatId int64
}
