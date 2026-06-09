package tr

import "time"

// Turkish

const (
	ErrorMessage = "" +
		`Tekrar deneyin - iki istasyon, virgülle ayrılmış.Tam olarak şöyle:

Podgorica, Bar`

	StationDoesNotExistMessage   = "Bu istasyon mevcut değil"
	RailwayMapButtonTextMap      = "Karadağ Demiryolu Haritası"
	OfficialTimetableUrlText     = "Daha fazla bilgi"
	ReverseRouteInlineButtonText = "Tersine"
	AlertUpdateNotificationText  = "" +
		`Tarife zaten güncellendi
13 Haziran'dan 14 Eylül'e kadar Subotica - Belgrad - Bar yeni bir tren eklenecek

Tarifenin geri kalanı tam olarak aynı kalacak`
	SimpleUpdateNotificationText = "Bugünün tarifesi güncellendi"

	// bot description

	BotName        = "🚂 Karadağ: tren tarifesi | Montenegro train"
	BotDescription = "" +
		`> Güncel tarife
> Her istasyonu biliyor, Belgrad dahil
> Herhangi iki istasyon arasında rotaları gösterebilir, aktarma dahil

Sadece bir virgülle ayrılmış iki istasyonu yazın:

Podgorica, Bar`
	BotShortDescription = "Tüm istasyonlar ve rotalarla güncel tarife, aktarma rotaları ve Belgrad - Bar gibi uluslararası rotalar dahil"

	// bot commands

	BotCommandNameStart = "Bota başla"
	StartMessage        = "" +
		`*Karadağ Demiryolları Tarifesi*

_@Leti\_deshevle ile birlikte yapıldı_

Lütfen *bir virgülle ayrılmış iki istasyon* girin: 

>*Podgorica, Bijelo Polje*

Ve size tarifeyi göndereceğim:

>Podgorica \> Bijelo Polje
>[6100](https://zpcg.me/details?timetable=41)  ` + "`06:20 \\> 08:38`" + `
>\.\.\.

İstasyonların doğru yazımı konusunda emin değil misiniz? Sorun değil, sadece yazın, gerisini ben halledeceğim\.

Sıra sende\!
`

	// /help

	BotCommandNameHelp = "Yardım"
	HelpMessage        = "" +
		`Sıkça sorulan sorular sık sık cevaplanır:

1. Bu, Karadağ'daki tren seferlerini gösteren bir bottur. Sefer programında Sırbistan'dan Karadağ'a giden/Sırbistan'dan Karadağ'a giden trenler de yer alıyor.
2. Karadağ'ın demiryolu bağlantısı sadece Sırbistan'dan başka hiçbir ülkede yoktur.
3. İstasyon haritasını /map üzerinden kontrol edin
4. Sadece virgülle ayrılmış iki istasyonu girin: 'Podgorica, Bar' ve sefer saatlerini göreceksiniz.
5. Biletler yalnızca istasyondan veya tren içerisinde satın alınabilir. Sadece nakit, online bilet yok, bazen bazı istasyonlarda kart kabul ediliyor (evet, bazen).
6. Programın alt kısmında bulunan 'Daha fazla bilgi' bağlantısına tıklayarak fiyatı, indirimleri ve diğer ayrıntıları kontrol edin.
7. Yaz aylarında bir tren hariç, sefer saatleri yıl boyunca aynıdır. Programın geri kalanı aynı şekilde devam ediyor.
8. Soldaki "🔄 'tarih'" düğmesini kullanarak programı güncelleyin
9. Bazen trenler gecikir, özellikle yaz sezonunda.
10. Bot hakkında daha detaylı bilgi için /about adresini ziyaret edin.
`

	// /map

	BotCommandNameMap = "Tüm istasyonların haritası"
	MapMessage        = "Tüm istasyonların bulunduğu harita"

	// /about

	BotCommandNameAbout = "Bu bot hakkında"
	AboutMessage        = "" +
		`Bu bot BEERWARE lisansı altında kullanılabilir.

Bu bildirimi gördüğünüz sürece bu kod ve bu botla istediğinizi yapabilirsiniz.
Eğer bir gün karşılaşırsak ve sen şöyle düşünürsen,
Bu botun faydalı olduğunu düşünüyorsanız bana teşekkür olarak bir bira ısmarlayabilirsiniz.

Ben Niksicko tamno'yu tercih ederim.

Ben: https://github.com/ivanov-gv
Bu proje: https://github.com/ivanov-gv/zpcg

@Leti_deshevle ile birlikte yapıldı
`
)

var MonthsMap = map[time.Month]string{
	time.January:   "Ocak",
	time.February:  "Şubat",
	time.March:     "Mart",
	time.April:     "Nisan",
	time.May:       "Mayıs",
	time.June:      "Haziran",
	time.July:      "Temmuz",
	time.August:    "Ağustos",
	time.September: "Eylül",
	time.October:   "Ekim",
	time.November:  "Kasım",
	time.December:  "Aralık",
}
