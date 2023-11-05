package utils

import "golang.org/x/net/html"

func IsTableEndReached(token html.Token) bool {
	return token.Type == html.EndTagToken && token.Data == "table"
}

func IsRowBeginningReached(token html.Token) bool {
	return token.Type == html.StartTagToken &&
		(HasAttribute(token.Attr, "", "class", "odd") || HasAttribute(token.Attr, "", "class", "even"))
}

func IsRowEndReached(token html.Token) bool {
	return token.Type == html.EndTagToken && token.Data == "tr"
}

func IsCellBeginningReached(token html.Token) bool {
	return token.Type == html.StartTagToken && token.Data == "td"
}

func IsCellEndReached(token html.Token) bool {
	return token.Type == html.EndTagToken && token.Data == "td"
}
