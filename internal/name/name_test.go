package name

import (
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestUnify(t *testing.T) {
	assert.Equal(t, "niksic", Unify("Nikšić"))
	assert.Equal(t, "baresumanovica", Unify("Bare Šumanovića"))
}

var (
	AllStations = []string{
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
	UnifiedStations = lo.Map(AllStations, func(item string, index int) string { return Unify(item) })
)

func TestFindBestMatch(t *testing.T) {
	assert.Equal(t, "niksic", FindBestMatch(Unify("Nikschichsss   "), UnifiedStations))
	assert.Equal(t, "novibeograd", FindBestMatch(Unify("NoVij Belgrad"), UnifiedStations))
}
