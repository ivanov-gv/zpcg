package model

type InlineButton struct {
	Text string
	Url  string
}

type Response struct {
	Text, ParseMode string
	InlineKeyboard  [][]InlineButton
}
