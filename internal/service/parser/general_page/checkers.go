package general_page

import (
	"golang.org/x/net/html"

	"zpcg/internal/service/parser/utils"
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

const (
	DetailedTimetableLinkTokenData      = "a"
	DetailedTimetableLinkAttributeKey   = "title"
	DetailedTimetableLinkAttributeValue = "Detalji"
)

func IsLinkToDetailedTimetableFound(token html.Token) bool {
	return token.Type == html.StartTagToken && token.Data == DetailedTimetableLinkTokenData &&
		utils.HasAttribute(token.Attr, "", DetailedTimetableLinkAttributeKey, DetailedTimetableLinkAttributeValue)
}
