package name

import (
	"testing"

	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
)

func TestUnify(t *testing.T) {
	assert.Equal(t, "niksic", Unify("Nikšić"))
	assert.Equal(t, "baresumanovica", Unify("Bare Šumanovića"))
}

var (
	allStationsList = []string{
		"Aerodrom",
		"Bačka Topola",
		"Bar",
		"Bare Šumanovića",
		"Beograd Centar",
		"Beška",
		"Bijelo Polje",
		"Bioče",
		"Branešci",
		"Bratonožići",
		"Brodarevo",
		"Čačak",
		"Crmnica",
		"Dabovići",
		"Danilovgrad",
		"Golubovci",
		"Inđija",
		"Kolašin",
		"Kos",
		"Kosjerić",
		"Kragujevac",
		"Kraljevo",
		"Kruševački Potok",
		"Kruševo",
		"Lajkovac",
		"Lapovo",
		"Lazarevac",
		"Lješnica",
		"Ljutotuk",
		"Lovćenac",
		"Lutovo",
		"Mateševo",
		"Mijatovo Kolo",
		"Mojkovac",
		"Morača",
		"Nikšić",
		"Nova Pazova",
		"Novi Beograd",
		"Novi Sad",
		"Oblutak",
		"Ostrog",
		"Padež",
		"Podgorica",
		"Požega",
		"Priboj",
		"Pričelje",
		"Prijepolje",
		"Prijepolje teretna",
		"Rakovica",
		"Ravna Rijeka",
		"Selište",
		"Slap",
		"Slijepač Most",
		"Šobajići",
		"Spuž",
		"Stara Pazova",
		"Štitarička Rijeka",
		"Stubica",
		"Subotica",
		"Šušanj",
		"Sutomore",
		"Trebaljevo",
		"Trebešica",
		"Užice",
		"Valjevo",
		"Velika Plana",
		"Virpazar",
		"Vranjina",
		"Vrbas",
		"Vrbnica",
		"Žari",
		"Zemun",
		"Zeta",
		"Zlatica",
		"Zmajevo",
	}
	unifiedStationsNameList = lo.Map(allStationsList, func(item string, index int) string { return Unify(item) })
)

func TestFindBestMatch(t *testing.T) {
	NewStationNameResolver(nil, unifiedStationsNameList)
	// niksic
	match, err := findBestMatch(Unify("Nikschichsss   "), unifiedStationsNameList)
	assert.NoError(t, err)
	assert.Equal(t, "niksic", match)
	// novi beograd
	match, err = findBestMatch(Unify("NoVij    Belgrad"), unifiedStationsNameList)
	assert.NoError(t, err)
	assert.Equal(t, "novibeograd", match)
	// indija
	match, err = findBestMatch(Unify("indija"), unifiedStationsNameList)
	assert.NoError(t, err)
	assert.Equal(t, "inđija", match)
}
