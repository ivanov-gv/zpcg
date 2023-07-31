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

func IsLinkToDetailedTimetabelFound(token html.Token) bool {
	return token.Type == html.StartTagToken && token.Data == "a" &&
		utils.HasAttribute(token.Attr, "", "title", "Detalji")
}
