package generators

import (
	"reflect"
	"strings"

	"github.com/ilmaruk/esms/internal/random"
)

const (
	isPrefix int = 1
	isSuffix int = 2
)

var teamNamePrefixes = map[string]int{
	"Akademik":      isPrefix,
	"Aris":          isPrefix,
	"Atletico":      isPrefix,
	"Athletic":      isPrefix,
	"Burg":          isSuffix,
	"Dinamo":        isPrefix,
	"Dynamo":        isPrefix,
	"FC":            isPrefix | isSuffix,
	"Fortuna":       isPrefix,
	"Lokomotiv":     isPrefix,
	"Metalurg":      isPrefix,
	"Olimpia":       isPrefix,
	"Olimpija":      isPrefix,
	"Olimpik":       isPrefix,
	"Partizan":      isPrefix,
	"Pro":           isPrefix,
	"Racing":        isPrefix,
	"Real":          isPrefix,
	"Rotor":         isPrefix,
	"Slavia":        isPrefix,
	"Sparta":        isPrefix,
	"Spartak":       isPrefix,
	"Sport":         isPrefix,
	"Sporting":      isPrefix,
	"Standard":      isPrefix,
	"Torpedo":       isPrefix,
	"Union":         isPrefix,
	"United":        isSuffix,
	"Universitario": isPrefix,
	"Viktoria":      isPrefix,
	"Virtus":        isPrefix,
	"Vitoria":       isPrefix,
	"Zenit":         isPrefix,
}

var azgaarBurgs = []string{
	// Human
	"Amberglen,Angelhand,Arrowden,Autumnband,Autumnkeep,Basinfrost,Basinmore,Bayfrost,Beargarde,Bearmire,Bellcairn,Bellport,Bellreach,Blackwatch,Bleakward,Bonemouth,Boulder,Bridgefalls,Bridgeforest,Brinepeak,Brittlehelm,Bronzegrasp,Castlecross,Castlefair,Cavemire,Claymond,Claymouth,Clearguard,Cliffgate,Cliffshear,Cliffshield,Cloudbay,Cloudcrest,Cloudwood,Coldholde,Cragbury,Crowgrove,Crowvault,Crystalrock,Crystalspire,Cursefield,Curseguard,Cursespell,Dawnforest,Dawnwater,Deadford,Deadkeep,Deepcairn,Deerchill,Demonfall,Dewglen,Dewmere,Diredale,Direden,Dirtshield,Dogcoast,Dogmeadow,Dragonbreak,Dragonhold,Dragonward,Dryhost,Dustcross,Dustwatch,Eaglevein,Earthfield,Earthgate,Earthpass,Ebonfront,Edgehaven,Eldergate,Eldermere,Embervault,Everchill,Evercoast,Falsevale,Faypond,Fayvale,Fayyard,Fearpeak,Flameguard,Flamewell,Freyshell,Ghostdale,Ghostpeak,Gloomburn,Goldbreach,Goldyard,Grassplains,Graypost,Greeneld,Grimegrove,Grimeshire,Heartfall,Heartford,Heartvault,Highbourne,Hillpass,Hollowstorm,Honeywater,Houndcall,Houndholde,Iceholde,Icelight,Irongrave,Ironhollow,Knightlight,Knighttide,Lagoonpass,Lakecross,Lastmere,Laststar,Lightvale,Limeband,Littlehall,Littlehold,Littlemire,Lostcairn,Lostshield,Loststar,Madfair,Madham,Midholde,Mightglen,Millstrand,Mistvault,Mondpass,Moonacre,Moongulf,Moonwell,Mosshand,Mosstide,Mosswind,Mudford,Mudwich,Mythgulch,Mythshear,Nevercrest,Neverfront,Newfalls,Nighthall,Oakenbell,Oakenrun,Oceanstar,Oldreach,Oldwall,Oldwatch,Oxbrook,Oxlight,Pearlhaven,Pinepond,Pondfalls,Pondtown,Pureshell,Quickbell,Quickpass,Ravenside,Roguehaven,Roseborn,Rosedale,Rosereach,Rustmore,Saltmouth,Sandhill,Scorchpost,Scorchstall,Shadeforest,Shademeadow,Shadeville,Shimmerrun,Shimmerwood,Shroudrock,Silentkeep,Silvercairn,Silvergulch,Smallmire,Smoothcliff,Smoothgrove,Smoothtown,Snakemere,Snowbay,Snowshield,Snowtown,Southbreak,Springmire,Springview,Stagport,Steammouth,Steamwall,Steepmoor,Stillhall,Stoneguard,Stonespell,Stormhand,Stormhorn,Sungulf,Sunhall,Swampmaw,Swangarde,Swanwall,Swiftwell,Thorncairn,Thornhelm,Thornyard,Timberside,Tradewick,Westmeadow,Westpoint,Whiteshore,Whitvalley,Wildeden,Wildwell,Wildyard,Winterhaven,Wolfpass",
	// Elven
	"Adrindest,Aethel,Afranthemar,Aiqua,Alari,Allanar,Almalian,Alora,Alyanasari,Alyelona,Alyran,Ammar,Anyndell,Arasari,Aren,Ashmebel,Aymlume,Bel-Didhel,Brinorion,Caelora,Chaulssad,Chaundra,Cyhmel,Cyrang,Dolarith,Dolonde,Draethe,Dranzan,Draugaust,E'ana,Eahil,Edhil,Eebel,Efranluma,Eld-Sinnocrin,Elelthyr,Ellanalin,Ellena,Ellorthond,Eltaesi,Elunore,Emyranserine,Entheas,Eriargond,Esari,Esath,Eserius,Eshsalin,Eshthalas,Evraland,Faellenor,Famelenora,Filranlean,Filsaqua,Gafetheas,Gaf Serine,Geliene,Gondorwin,Guallu,Haeth,Hanluna,Haulssad,Heloriath,Himlarien,Himliene,Hinnead,Hlinas,Hloireenil,Hluihei,Hlurthei,Hlynead,Iaenarion,Iaron,Illanathaes,Illfanora,Imlarlon,Imyse,Imyvelian,Inferius,Inlurth,innsshe,Iralserin,Irethtalos,Irholona,Ishal,Ishlashara,Ithelion,Ithlin,Iulil,Jaal,Jamkadi,Kaalume,Kaansera,Karanthanil,Karnosea,Kasethyr,Keatheas,Kelsya,Keth Aiqua,Kmlon,Kyathlenor,Kyhasera,Lahetheas,Lefdorei,Lelhamelle,Lilean,Lindeenil,Lindoress,Litys,Llaughei,Lya,Lyfa,Lylharion,Lynathalas,Machei,Masenoris,Mathethil,Mathentheas,Meethalas,Menyamar,Mithlonde,Mytha,Mythsemelle,Mythsthas,Naahona,Nalore,Nandeedil,Nasad Ilaurth,Nasin,Nathemar,Neadar,Neilon,Nelalon,Nellean,Nelnetaesi,Nilenathyr,Nionande,Nylm,Nytenanas,Nythanlenor,O'anlenora,Obeth,Ofaenathyr,Ollmnaes,Ollsmel,Olwen,Olyaneas,Omanalon,Onelion,Onelond,Orlormel,Ormrion,Oshana,Oshvamel,Raethei,Rauguall,Reisera,Reslenora,Ryanasera,Rymaserin,Sahnor,Saselune,Sel-Zedraazin,Selananor,Sellerion,Selmaluma,Shaeras,Shemnas,Shemserin,Sheosari,Sileltalos,Siriande,Siriathil,Srannor,Sshanntyr,Sshaulu,Syholume,Sylharius,Sylranbel,Taesi,Thalor,Tharenlon,Thelethlune,Thelhohil,Themar,Thene,Thilfalean,Thilnaenor,Thvethalas,Thylathlond,Tiregul,Tlauven,Tlindhe,Ulal,Ullve,Ulmetheas,Ulssin,Umnalin,Umye,Umyheserine,Unanneas,Unarith,Undraeth,Unysarion,Vel-Shonidor,Venas,Vin Argor,Wasrion,Wlalean,Yaeluma,Yeelume,Yethrion,Ymserine,Yueghed,Yuerran,Yuethin",
	// Dwarven
	"Addundad,Ahagzad,Ahazil,Akil,Akzizad,Anumush,Araddush,Arar,Arbhur,Badushund,Baragzig,Baragzund,Barakinb,Barakzig,Barakzinb,Barakzir,Baramunz,Barazinb,Barazir,Bilgabar,Bilgatharb,Bilgathaz,Bilgila,Bilnaragz,Bilnulbar,Bilnulbun,Bizaddum,Bizaddush,Bizanarg,Bizaram,Bizinbiz,Biziram,Bunaram,Bundinar,Bundushol,Bundushund,Bundushur,Buzaram,Buzundab,Buzundush,Gabaragz,Gabaram,Gabilgab,Gabilgath,Gabizir,Gabunal,Gabunul,Gabuzan,Gatharam,Gatharbhur,Gathizdum,Gathuragz,Gathuraz,Gila,Giledzir,Gilukkhath,Gilukkhel,Gunala,Gunargath,Gunargil,Gundumunz,Gundusharb,Gundushizd,Kharbharbiln,Kharbhatharb,Kharbhela,Kharbilgab,Kharbuzadd,Khatharbar,Khathizdin,Khathundush,Khazanar,Khazinbund,Khaziragz,Khaziraz,Khizdabun,Khizdusharbh,Khizdushath,Khizdushel,Khizdushur,Kholedzar,Khundabiln,Khundabuz,Khundinarg,Khundushel,Khuragzig,Khuramunz,Kibarak,Kibilnal,Kibizar,Kibunarg,Kibundin,Kibuzan,Kinbadab,Kinbaragz,Kinbarakz,Kinbaram,Kinbizah,Kinbuzar,Nala,Naledzar,Naledzig,Naledzinb,Naragzah,Naragzar,Naragzig,Narakzah,Narakzar,Naramunz,Narazar,Nargabad,Nargabar,Nargatharb,Nargila,Nargundum,Nargundush,Nargunul,Narukthar,Narukthel,Nula,Nulbadush,Nulbaram,Nulbilnarg,Nulbunal,Nulbundab,Nulbundin,Nulbundum,Nulbuzah,Nuledzah,Nuledzig,Nulukkhaz,Nulukkhund,Nulukkhur,Sharakinb,Sharakzar,Sharamunz,Sharbarukth,Shatharbhizd,Shatharbiz,Shathazah,Shathizdush,Shathola,Shaziragz,Shizdinar,Shizdushund,Sholukkharb,Shundinulb,Shundushund,Shurakzund,Shuramunz,Tumunzadd,Tumunzan,Tumunzar,Tumunzinb,Tumunzir,Ukthad,Ulbirad,Ulbirar,Ulunzar,Ulur,Umunzad,Undalar,Undukkhil,Undun,Undur,Unduzur,Unzar,Unzathun,Usharar,Zaddinarg,Zaddushur,Zaharbad,Zaharbhizd,Zarakib,Zarakzar,Zaramunz,Zarukthel,Zinbarukth,Zirakinb,Zirakzir,Ziramunz,Ziruktharbh,Zirukthur,Zundumunz",
	// Draconic
	"Aaronarra,Adalon,Adamarondor,Aeglyl,Aerosclughpalar,Aghazstamn,Aglaraerose,Agoshyrvor,Alduin,Alhazmabad,Altagos,Ammaratha,Amrennathed,Anaglathos,Andrathanach,Araemra,Araugauthos,Arauthator,Arharzel,Arngalor,Arveiaturace,Athauglas,Augaurath,Auntyrlothtor,Azarvilandral,Azhaq,Balagos,Baratathlaer,Bleucorundum,BrazzPolis,Canthraxis,Capnolithyl,Charvekkanathor,Chellewis,Chelnadatilar,Cirrothamalan,Claugiyliamatar,Cragnortherma,Dargentum,Dendeirmerdammarar,Dheubpurcwenpyl,Domborcojh,Draconobalen,Dragansalor,Dupretiskava,Durnehviir,Eacoathildarandus,Eldrisithain,Enixtryx,Eormennoth,Esmerandanna,Evenaelorathos,Faenphaele,Felgolos,Felrivenser,Firkraag,Fll'Yissetat,Furlinastis,Galadaeros,Galglentor,Garnetallisar,Garthammus,Gaulauntyr,Ghaulantatra,Glouroth,Greshrukk,Guyanothaz,Haerinvureem,Haklashara,Halagaster,Halaglathgar,Havarlan,Heltipyre,Hethcypressarvil,Hoondarrh,Icehauptannarthanyx,Iiurrendeem,Ileuthra,Iltharagh,Ingeloakastimizilian,Irdrithkryn,Ishenalyr,Iymrith,Jaerlethket,Jalanvaloss,Jharakkan,Kasidikal,Kastrandrethilian,Khavalanoth,Khuralosothantar,Kisonraathiisar,Kissethkashaan,Kistarianth,Klauth,Klithalrundrar,Krashos,Kreston,Kriionfanthicus,Krosulhah,Krustalanos,Kruziikrel,Kuldrak,Lareth,Latovenomer,Lhammaruntosz,Llimark,Ma'fel'no'sei'kedeh'naar,MaelestorRex,Magarovallanthanz,Mahatnartorian,Mahrlee,Malaeragoth,Malagarthaul,Malazan,Maldraedior,Maldrithor,MalekSalerno,Maughrysear,Mejas,Meliordianix,Merah,Mikkaalgensis,Mirmulnir,Mistinarperadnacles,Miteach,Mithbarazak,Morueme,Moruharzel,Naaslaarum,Nahagliiv,Nalavarauthatoryl,Naxorlytaalsxar,Nevalarich,Nolalothcaragascint,Nurvureem,Nymmurh,Odahviing,Olothontor,Ormalagos,Otaaryliakkarnos,Paarthurnax,Pelath,Pelendralaar,Praelorisstan,Praxasalandos,Protanther,Qiminstiir,Quelindritar,Ralionate,Rathalylaug,Rathguul,Rauglothgor,Raumorthadar,Relonikiv,Ringreemeralxoth,Roraurim,Rynnarvyx,Sablaxaahl,Sahloknir,Sahrotaar,Samdralyrion,Saryndalaghlothtor,Sawaka,Shalamalauth,Shammagar,Sharndrel,Shianax,Skarlthoon,Skurge,Smergadas,Ssalangan,Sssurist,Sussethilasis,Sylvallitham,Tamarand,Tantlevgithus,Tarlacoal,Tenaarlaktor,Thalagyrt,Tharas'kalagram,Thauglorimorgorus,Thoklastees,Thyka,Tsenshivah,Ueurwen,Uinnessivar,Urnalithorgathla,Velcuthimmorhar,Velora,Vendrathdammarar,Venomindhar,Viinturuth,Voaraghamanthar,Voslaarum,Vr'tark,Vrondahorevos,Vuljotnaak,Vulthuryol,Wastirek,Worlathaugh,Xargithorvar,Xavarathimius,Yemere,Ylithargathril,Ylveraasahlisar,Za-Jikku,Zarlandris,Zellenesterex,Zilanthar,Zormapalearath,Zundaerazylym,Zz'Pzora",
}

func GenerateTeamName(rnd random.Randomiser) string {
	// Burg
	burgsSet := azgaarBurgs[rnd.UniformRandom(len(azgaarBurgs)-1)]
	burgs := strings.Split(burgsSet, ",")
	burg := burgs[rnd.UniformRandom(len(burgs)-1)]

	// Prefix/Suffix
	keys := reflect.ValueOf(teamNamePrefixes).MapKeys()
	value := keys[rnd.UniformRandom(len(keys)-1)].Interface().(string)

	prefixSuffix := teamNamePrefixes[value]
	isP := (prefixSuffix & isPrefix) != 0
	isS := (prefixSuffix & isSuffix) != 0

	if isP && isS {
		isP = rnd.ThrowWithProb(50)
	}

	if isP {
		return value + " " + burg
	}

	return burg + " " + value
}
