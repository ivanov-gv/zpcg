package render

import (
	"time"

	"golang.org/x/text/language"

	"github.com/ivanov-gv/zpcg/internal/model/render/be"
	"github.com/ivanov-gv/zpcg/internal/model/render/de"
	"github.com/ivanov-gv/zpcg/internal/model/render/en"
	"github.com/ivanov-gv/zpcg/internal/model/render/hr"
	"github.com/ivanov-gv/zpcg/internal/model/render/ru"
	"github.com/ivanov-gv/zpcg/internal/model/render/sk"
	"github.com/ivanov-gv/zpcg/internal/model/render/sr"
	"github.com/ivanov-gv/zpcg/internal/model/render/tr"
	"github.com/ivanov-gv/zpcg/internal/model/render/uk"
)

// bot general

type Bot struct {
	Name, Description, ShortDescription string
	CommandNames                        map[BotCommand]string
}

type BotCommand string

var AllCommands = []BotCommand{
	BotCommandStart,
	BotCommandHelp,
	BotCommandMap,
	BotCommandAbout,
}

const (
	BotCommandStart BotCommand = "/start"
	BotCommandHelp  BotCommand = "/help"
	BotCommandMap   BotCommand = "/map"
	BotCommandAbout BotCommand = "/about"
)

const GoogleMapWithAllStations = "https://goo.gl/maps/jYRwCAC14mmoe4L3A"

var (
	BotInfoMap = map[language.Tag]Bot{
		language.English: {
			Name:             en.BotName,
			Description:      en.BotDescription,
			ShortDescription: en.BotShortDescription,
			CommandNames: map[BotCommand]string{
				BotCommandStart: en.BotCommandNameStart,
				BotCommandHelp:  en.BotCommandNameHelp,
				BotCommandMap:   en.BotCommandNameMap,
				BotCommandAbout: en.BotCommandNameAbout,
			},
		},
		language.Russian: {
			Name:             ru.BotName,
			Description:      ru.BotDescription,
			ShortDescription: ru.BotShortDescription,
			CommandNames: map[BotCommand]string{
				BotCommandStart: ru.BotCommandNameStart,
				BotCommandHelp:  ru.BotCommandNameHelp,
				BotCommandMap:   ru.BotCommandNameMap,
				BotCommandAbout: ru.BotCommandNameAbout,
			},
		},
		language.Serbian: {
			Name:             sr.BotName,
			Description:      sr.BotDescription,
			ShortDescription: sr.BotShortDescription,
			CommandNames: map[BotCommand]string{
				BotCommandStart: sr.BotCommandNameStart,
				BotCommandHelp:  sr.BotCommandNameHelp,
				BotCommandMap:   sr.BotCommandNameMap,
				BotCommandAbout: sr.BotCommandNameAbout,
			},
		},
		language.Turkish: {
			Name:             tr.BotName,
			Description:      tr.BotDescription,
			ShortDescription: tr.BotShortDescription,
			CommandNames: map[BotCommand]string{
				BotCommandStart: tr.BotCommandNameStart,
				BotCommandHelp:  tr.BotCommandNameHelp,
				BotCommandMap:   tr.BotCommandNameMap,
				BotCommandAbout: tr.BotCommandNameAbout,
			},
		},
		language.Slovak: {
			Name:             sk.BotName,
			Description:      sk.BotDescription,
			ShortDescription: sk.BotShortDescription,
			CommandNames: map[BotCommand]string{
				BotCommandStart: sk.BotCommandNameStart,
				BotCommandHelp:  sk.BotCommandNameHelp,
				BotCommandMap:   sk.BotCommandNameMap,
				BotCommandAbout: sk.BotCommandNameAbout,
			},
		},
		language.Croatian: {
			Name:             hr.BotName,
			Description:      hr.BotDescription,
			ShortDescription: hr.BotShortDescription,
			CommandNames: map[BotCommand]string{
				BotCommandStart: hr.BotCommandNameStart,
				BotCommandHelp:  hr.BotCommandNameHelp,
				BotCommandMap:   hr.BotCommandNameMap,
				BotCommandAbout: hr.BotCommandNameAbout,
			},
		},
		language.German: {
			Name:             de.BotName,
			Description:      de.BotDescription,
			ShortDescription: de.BotShortDescription,
			CommandNames: map[BotCommand]string{
				BotCommandStart: de.BotCommandNameStart,
				BotCommandHelp:  de.BotCommandNameHelp,
				BotCommandMap:   de.BotCommandNameMap,
				BotCommandAbout: de.BotCommandNameAbout,
			},
		},
		language.Ukrainian: {
			Name:             uk.BotName,
			Description:      uk.BotDescription,
			ShortDescription: uk.BotShortDescription,
			CommandNames: map[BotCommand]string{
				BotCommandStart: uk.BotCommandNameStart,
				BotCommandHelp:  uk.BotCommandNameHelp,
				BotCommandMap:   uk.BotCommandNameMap,
				BotCommandAbout: uk.BotCommandNameAbout,
			},
		},
		Belarusian: {
			Name:             be.BotName,
			Description:      be.BotDescription,
			ShortDescription: be.BotShortDescription,
			CommandNames: map[BotCommand]string{
				BotCommandStart: be.BotCommandNameStart,
				BotCommandHelp:  be.BotCommandNameHelp,
				BotCommandMap:   be.BotCommandNameMap,
				BotCommandAbout: be.BotCommandNameAbout,
			},
		},
	}
	ErrorMessageMap = map[language.Tag]string{
		language.Russian:   ru.ErrorMessage,
		language.Ukrainian: uk.ErrorMessage,
		Belarusian:         be.ErrorMessage,
		language.English:   en.ErrorMessage,
		language.German:    de.ErrorMessage,
		language.Serbian:   sr.ErrorMessage,
		language.Croatian:  hr.ErrorMessage,
		language.Slovak:    sk.ErrorMessage,
		language.Turkish:   tr.ErrorMessage,
	}
	StartMessageMap = map[language.Tag]string{
		language.Russian:   ru.StartMessage,
		language.Ukrainian: uk.StartMessage,
		Belarusian:         be.StartMessage,
		language.English:   en.StartMessage,
		language.German:    de.StartMessage,
		language.Serbian:   sr.StartMessage,
		language.Croatian:  hr.StartMessage,
		language.Slovak:    sk.StartMessage,
		language.Turkish:   tr.StartMessage,
	}
	HelpMessageMap = map[language.Tag]string{
		language.Russian:   ru.HelpMessage,
		language.Ukrainian: uk.HelpMessage,
		Belarusian:         be.HelpMessage,
		language.English:   en.HelpMessage,
		language.German:    de.HelpMessage,
		language.Serbian:   sr.HelpMessage,
		language.Croatian:  hr.HelpMessage,
		language.Slovak:    sk.HelpMessage,
		language.Turkish:   tr.HelpMessage,
	}
	MapMessageMap = map[language.Tag]string{
		language.Russian:   ru.MapMessage,
		language.Ukrainian: uk.MapMessage,
		Belarusian:         be.MapMessage,
		language.English:   en.MapMessage,
		language.German:    de.MapMessage,
		language.Serbian:   sr.MapMessage,
		language.Croatian:  hr.MapMessage,
		language.Slovak:    sk.MapMessage,
		language.Turkish:   tr.MapMessage,
	}
	AboutMessageMap = map[language.Tag]string{
		language.Russian:   ru.AboutMessage,
		language.Ukrainian: uk.AboutMessage,
		Belarusian:         be.AboutMessage,
		language.English:   en.AboutMessage,
		language.German:    de.AboutMessage,
		language.Serbian:   sr.AboutMessage,
		language.Croatian:  hr.AboutMessage,
		language.Slovak:    sk.AboutMessage,
		language.Turkish:   tr.AboutMessage,
	}
	StationDoesNotExistMessageMap = map[language.Tag]string{
		language.Russian:   ru.StationDoesNotExistMessage,
		language.Ukrainian: uk.StationDoesNotExistMessage,
		Belarusian:         be.StationDoesNotExistMessage,
		language.English:   en.StationDoesNotExistMessage,
		language.German:    de.StationDoesNotExistMessage,
		language.Serbian:   sr.StationDoesNotExistMessage,
		language.Croatian:  hr.StationDoesNotExistMessage,
		language.Slovak:    sk.StationDoesNotExistMessage,
		language.Turkish:   tr.StationDoesNotExistMessage,
	}
	RailwayMapButtonTextMap = map[language.Tag]string{
		language.Russian:   ru.RailwayMapButtonTextMap,
		language.Ukrainian: uk.RailwayMapButtonTextMap,
		Belarusian:         be.RailwayMapButtonTextMap,
		language.English:   en.RailwayMapButtonTextMap,
		language.German:    de.RailwayMapButtonTextMap,
		language.Serbian:   sr.RailwayMapButtonTextMap,
		language.Croatian:  hr.RailwayMapButtonTextMap,
		language.Slovak:    sk.RailwayMapButtonTextMap,
		language.Turkish:   tr.RailwayMapButtonTextMap,
	}
	MonthNameMap = map[language.Tag]map[time.Month]string{
		language.Russian:   ru.MonthsMap,
		language.English:   en.MonthsMap,
		language.Serbian:   sr.MonthsMap,
		language.Turkish:   tr.MonthsMap,
		language.Ukrainian: uk.MonthsMap,
		Belarusian:         be.MonthsMap,
		language.German:    de.MonthsMap,
		language.Croatian:  hr.MonthsMap,
		language.Slovak:    sk.MonthsMap,
	}
	ReverseRouteInlineButtonTextMap = map[language.Tag]string{
		language.Russian:   ru.ReverseRouteInlineButtonText,
		language.Ukrainian: uk.ReverseRouteInlineButtonText,
		Belarusian:         be.ReverseRouteInlineButtonText,
		language.English:   en.ReverseRouteInlineButtonText,
		language.German:    de.ReverseRouteInlineButtonText,
		language.Serbian:   sr.ReverseRouteInlineButtonText,
		language.Croatian:  hr.ReverseRouteInlineButtonText,
		language.Slovak:    sk.ReverseRouteInlineButtonText,
		language.Turkish:   tr.ReverseRouteInlineButtonText,
	}
	AlertUpdateNotificationTextMap = map[language.Tag]string{
		language.Russian:   ru.AlertUpdateNotificationText,
		language.Ukrainian: uk.AlertUpdateNotificationText,
		Belarusian:         be.AlertUpdateNotificationText,
		language.English:   en.AlertUpdateNotificationText,
		language.German:    de.AlertUpdateNotificationText,
		language.Serbian:   sr.AlertUpdateNotificationText,
		language.Croatian:  hr.AlertUpdateNotificationText,
		language.Slovak:    sk.AlertUpdateNotificationText,
		language.Turkish:   tr.AlertUpdateNotificationText,
	}
	SimpleUpdateNotificationTextMap = map[language.Tag]string{
		language.Russian:   ru.SimpleUpdateNotificationText,
		language.Ukrainian: uk.SimpleUpdateNotificationText,
		Belarusian:         be.SimpleUpdateNotificationText,
		language.English:   en.SimpleUpdateNotificationText,
		language.German:    de.SimpleUpdateNotificationText,
		language.Serbian:   sr.SimpleUpdateNotificationText,
		language.Croatian:  hr.SimpleUpdateNotificationText,
		language.Slovak:    sk.SimpleUpdateNotificationText,
		language.Turkish:   tr.SimpleUpdateNotificationText,
	}
	OfficialTimetableUrlTextMap = map[language.Tag]string{
		language.Russian:   ru.OfficialTimetableUrlText,
		language.Ukrainian: uk.OfficialTimetableUrlText,
		Belarusian:         be.OfficialTimetableUrlText,
		language.English:   en.OfficialTimetableUrlText,
		language.German:    de.OfficialTimetableUrlText,
		language.Serbian:   sr.OfficialTimetableUrlText,
		language.Croatian:  hr.OfficialTimetableUrlText,
		language.Slovak:    sk.OfficialTimetableUrlText,
		language.Turkish:   tr.OfficialTimetableUrlText,
	}
)
