package utils

import "golang.org/x/net/html"

const tableEndTokenData = "table"

func IsTableEndReached(token html.Token) bool {
	return token.Type == html.EndTagToken && token.Data == tableEndTokenData
}

const (
	tableRowBeginningKey       = "class"
	tableRowBeginningValueOdd  = "odd"
	tableRowBeginningValueEven = "even"
)

func IsRowBeginningReached(token html.Token) bool {
	return token.Type == html.StartTagToken &&
		(HasAttribute(token.Attr, "", tableRowBeginningKey, tableRowBeginningValueOdd) ||
			HasAttribute(token.Attr, "", tableRowBeginningKey, tableRowBeginningValueEven))
}

const tableRowEndTokenData = "tr"

func IsRowEndReached(token html.Token) bool {
	return token.Type == html.EndTagToken && token.Data == tableRowEndTokenData
}

const tableCellBeginningTokenData = "td"

func IsCellBeginningReached(token html.Token) bool {
	return token.Type == html.StartTagToken && token.Data == tableCellBeginningTokenData
}

const tableCellEndTokenData = "td"

func IsCellEndReached(token html.Token) bool {
	return token.Type == html.EndTagToken && token.Data == tableCellEndTokenData
}
