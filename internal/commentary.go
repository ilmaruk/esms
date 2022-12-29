package internal

import (
	"fmt"
	"io"
	"math/rand"
)

const (
	COMM_INJURYTIME string = "COMM_INJURYTIME"
	COMM_HALFTIME   string = "COMM_HALFTIME"
	COMM_FULLTIME   string = "COMM_FULLTIME"
	CHANCE          string = "CHANCE"
	ASSISTEDCHANCE  string = "ASSISTEDCHANCE"
)

var commentary map[string][]string

func init() {
	commentary = map[string][]string{
		COMM_INJURYTIME: {
			"The ref adds %d min. of injury time",
		},
		COMM_HALFTIME: {
			"*************  HALF TIME  ****************",
		},
		COMM_FULLTIME: {
			"*************  FULL TIME  ****************",
		},
		CHANCE: {
			"Min. %s :(%s) %s with the dribble",
			"Min. %s :(%s) %s takes possesion",
			"Min. %s :(%s) %s cuts through the defense",
			"Min. %s :(%s) %s finds a hole in the defense",
			"Min. %s :(%s) %s takes advantage of a defensive mistake",
			"Min. %s :(%s) %s finds his way through",
			"Min. %s :(%s) %s sidesteps his marker",
			"Min. %s :(%s) %s with a flashy move",
			"Min. %s :(%s) %s beats his marker",
			"Min. %s :(%s) %s with a real burst of pace",
			"Min. %s :(%s) %s bursts forward",
			"Min. %s :(%s) %s finds himself in a good position",
		},
		ASSISTEDCHANCE: {
			"Min. %2s :(%s) %s passes the ball to %s",
			"Min. %2s :(%s) %s with a smart pass to %s",
			"Min. %2s :(%s) %s finds %s in the box",
			"Min. %s :(%s) %s with a precise pass to %s",
			"Min. %s :(%s) %s heads the ball down to %s",
			"Min. %s :(%s) %s slides the ball across to %s",
			"Min. %s :(%s) %s cuts the ball back to %s",
			"Min. %s :(%s) %s with the heel pass to %s",
			"Min. %s :(%s) %s plays a long ball to %s",
			"Min. %s :(%s) %s with a glorious long pass to %s",
		},
	}
}

func PrintCommentary(w io.Writer, key string, values ...interface{}) {
	options, ok := commentary[key]
	if !ok {
		return
	}
	text := options[rand.Intn(len(options))] + "\n"
	fmt.Fprintf(w, text, values...)
}
