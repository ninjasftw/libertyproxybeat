package main

import (
	"os"

	"github.com/elastic/beats/libbeat/beat"

	"github.com/ninjasftw/libertyproxybeat/beater"
)

func main() {
	err := beat.Run("libertyproxybeat", "", beater.New())
	if err != nil {
		os.Exit(1)
	}
}
