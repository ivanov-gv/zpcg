package detailed_page

import (
	"golang.org/x/net/html"

	"zpcg/internal/service/parser/utils"
)

const (
	TimetableBeginningTokenData           = "div"
	TimetableBeginningTokenAttributeKey   = "id"
	TimetableBeginningTokenAttributeValue = "detail-stop-grid"
)

func IsTimetableReached(token html.Token) bool {
	//<div id="detail-stop-grid" class="grid-view">
	return token.Type == html.StartTagToken && token.Data == TimetableBeginningTokenData &&
		utils.HasAttribute(token.Attr, "", TimetableBeginningTokenAttributeKey, TimetableBeginningTokenAttributeValue)
}

const TableHeadEndTokenData = "thead"

func IsTableHeadEndReached(token html.Token) bool {
	//</thead>
	return token.Type == html.EndTagToken && token.Data == TableHeadEndTokenData
}
