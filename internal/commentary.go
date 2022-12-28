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
