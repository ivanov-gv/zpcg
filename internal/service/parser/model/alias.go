package model

type StationAliases struct {
	StationName string
	Aliases     []string
}

var AliasesStationsList = []StationAliases{
	{
		StationName: "Beograd Centar",
		Aliases:     []string{"belgrad"},
	},
}
