package blacklist

import (
	"github.com/samber/lo"

	"github.com/ivanov-gv/zpcg/internal/model/timetable"
)

func NewBlackListService() *BlackListService {
	return &BlackListService{}
}

type BlackListService struct{}

func (s *BlackListService) CheckBlackList(stationIds ...timetable.StationId) (isBlacklisted bool, blackListed []timetable.BlackListedStation) {
	blackListedStationIds := lo.Filter(stationIds, IsBlackListedWithIndex)
	if len(blackListedStationIds) == 0 {
		return false, nil
	}
	return true, lo.Map(blackListedStationIds, GetBlackListedStationByIdWithIndex)
}

func IsBlackListedWithIndex(id timetable.StationId, _ int) bool {
	return IsBlackListed(id)
}

func IsBlackListed(id timetable.StationId) bool {
	return id <= 0
}

func GetBlackListedStationByIdWithIndex(id timetable.StationId, _ int) timetable.BlackListedStation {
	return GetBlackListedStationById(id)
}

func GetBlackListedStationById(id timetable.StationId) timetable.BlackListedStation {
	return BlackListedStations[-id]
}
