package blacklist

import (
	"github.com/samber/lo"

	"zpcg/internal/model"
)

func NewBlackListService() *BlackListService {
	return &BlackListService{}
}

type BlackListService struct{}

func (s *BlackListService) CheckBlackList(stationIds ...model.StationId) (isBlacklisted bool, blackListed []model.BlackListedStation) {
	blackListedStationIds := lo.Filter(stationIds, IsBlackListedWithIndex)
	if len(blackListedStationIds) == 0 {
		return false, nil
	}
	return true, lo.Map(blackListedStationIds, GetBlackListedStationByIdWithIndex)
}

func IsBlackListedWithIndex(id model.StationId, _ int) bool {
	return IsBlackListed(id)
}

func IsBlackListed(id model.StationId) bool {
	return id <= 0
}

func GetBlackListedStationByIdWithIndex(id model.StationId, _ int) model.BlackListedStation {
	return GetBlackListedStationById(id)
}

func GetBlackListedStationById(id model.StationId) model.BlackListedStation {
	return BlackListedStations[-id]
}
