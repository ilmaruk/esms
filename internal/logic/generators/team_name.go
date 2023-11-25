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
	"Atletico":      isPrefix,
	"Athletic":      isPrefix,
	"Burg":          isSuffix,
	"Dinamo":        isPrefix,
	"Dynamo":        isPrefix,
	"FC":            isPrefix | isSuffix,
	"Fortuna":       isPrefix,
	"Lokomotiv":     isPrefix,
	"Olimpia":       isPrefix,
	"Olimpija":      isPrefix,
	"Partizan":      isPrefix,
	"Pro":           isPrefix,
	"Racing":        isPrefix,
	"Real":          isPrefix,
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
}

var azgaarElvenBurgs = "Adrindest,Aethel,Afranthemar,Aiqua,Alari,Allanar,Almalian,Alora,Alyanasari,Alyelona,Alyran,Ammar,Anyndell,Arasari,Aren,Ashmebel,Aymlume,Bel-Didhel,Brinorion,Caelora,Chaulssad,Chaundra,Cyhmel,Cyrang,Dolarith,Dolonde,Draethe,Dranzan,Draugaust,E'ana,Eahil,Edhil,Eebel,Efranluma,Eld-Sinnocrin,Elelthyr,Ellanalin,Ellena,Ellorthond,Eltaesi,Elunore,Emyranserine,Entheas,Eriargond,Esari,Esath,Eserius,Eshsalin,Eshthalas,Evraland,Faellenor,Famelenora,Filranlean,Filsaqua,Gafetheas,Gaf Serine,Geliene,Gondorwin,Guallu,Haeth,Hanluna,Haulssad,Heloriath,Himlarien,Himliene,Hinnead,Hlinas,Hloireenil,Hluihei,Hlurthei,Hlynead,Iaenarion,Iaron,Illanathaes,Illfanora,Imlarlon,Imyse,Imyvelian,Inferius,Inlurth,innsshe,Iralserin,Irethtalos,Irholona,Ishal,Ishlashara,Ithelion,Ithlin,Iulil,Jaal,Jamkadi,Kaalume,Kaansera,Karanthanil,Karnosea,Kasethyr,Keatheas,Kelsya,Keth Aiqua,Kmlon,Kyathlenor,Kyhasera,Lahetheas,Lefdorei,Lelhamelle,Lilean,Lindeenil,Lindoress,Litys,Llaughei,Lya,Lyfa,Lylharion,Lynathalas,Machei,Masenoris,Mathethil,Mathentheas,Meethalas,Menyamar,Mithlonde,Mytha,Mythsemelle,Mythsthas,Naahona,Nalore,Nandeedil,Nasad Ilaurth,Nasin,Nathemar,Neadar,Neilon,Nelalon,Nellean,Nelnetaesi,Nilenathyr,Nionande,Nylm,Nytenanas,Nythanlenor,O'anlenora,Obeth,Ofaenathyr,Ollmnaes,Ollsmel,Olwen,Olyaneas,Omanalon,Onelion,Onelond,Orlormel,Ormrion,Oshana,Oshvamel,Raethei,Rauguall,Reisera,Reslenora,Ryanasera,Rymaserin,Sahnor,Saselune,Sel-Zedraazin,Selananor,Sellerion,Selmaluma,Shaeras,Shemnas,Shemserin,Sheosari,Sileltalos,Siriande,Siriathil,Srannor,Sshanntyr,Sshaulu,Syholume,Sylharius,Sylranbel,Taesi,Thalor,Tharenlon,Thelethlune,Thelhohil,Themar,Thene,Thilfalean,Thilnaenor,Thvethalas,Thylathlond,Tiregul,Tlauven,Tlindhe,Ulal,Ullve,Ulmetheas,Ulssin,Umnalin,Umye,Umyheserine,Unanneas,Unarith,Undraeth,Unysarion,Vel-Shonidor,Venas,Vin Argor,Wasrion,Wlalean,Yaeluma,Yeelume,Yethrion,Ymserine,Yueghed,Yuerran,Yuethin"

func GenerateTeamName(rnd random.Randomiser) string {
	// Burg
	burgs := strings.Split(azgaarElvenBurgs, ",")
	burg := burgs[rnd.UniformRandom(len(burgs))]

	// Prefix/Suffix
	keys := reflect.ValueOf(teamNamePrefixes).MapKeys()
	value := keys[rnd.UniformRandom(len(keys))].Interface().(string)

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
