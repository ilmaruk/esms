package generators

import (
	"github.com/ilmaruk/esms/internal/random"
	"github.com/pallinder/go-randomdata"
)

var teamNamePrefixes = []string{
	"Atletico",
	"Athletic",
	"Dinamo",
	"Dynamo",
	"FC",
	"Fortuna",
	"Lokomotiv",
	"Olimpia",
	"Olimpija",
	"Partizan",
	"Pro",
	"Racing",
	"Real",
	"Sparta",
	"Spartak",
	"Sport",
	"Sporting",
	"Standard",
	"Torpedo",
	"Union",
	"Universitario",
	"Viktoria",
	"Virtus",
	"Vitoria",
}

func GenerateTeamName(rnd random.Randomiser) string {
	prefix := teamNamePrefixes[rnd.UniformRandom(len(teamNamePrefixes))]
	return prefix + " " + randomdata.City()
}
