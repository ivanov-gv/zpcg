package general_page

import (
	"golang.org/x/net/html"
	"zpcg/internal/parser/utils"
)

const (
	GeneralTimetableOpeningTagData      = "div"
	GeneralTimetableOpeningTagNamespace = ""
	GeneralTimetableOpeningTagKey       = "id"
	GeneralTimetableOpeningTagValue     = "timetable-grid"
)

func IsTableReached(token html.Token) bool {
	return token.Type == html.StartTagToken &&
		token.Data == GeneralTimetableOpeningTagData &&
		utils.HasAttribute(token.Attr, GeneralTimetableOpeningTagNamespace, GeneralTimetableOpeningTagKey, GeneralTimetableOpeningTagValue)
}

func IsTableEndReached(token html.Token) bool {
	return token.Type == html.EndTagToken && token.Data == "table"
}

func IsRowBeginningReached(token html.Token) bool {
	return token.Type == html.StartTagToken &&
		(utils.HasAttribute(token.Attr, "", "class", "odd") || utils.HasAttribute(token.Attr, "", "class", "even"))
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

func IsLinkToDetailedTimetabeFound(token html.Token) bool {
	return token.Type == html.StartTagToken && token.Data == "a" &&
		utils.HasAttribute(token.Attr, "", "title", "Detalji")
}
