package main

import (
	"flag"
	"log"

	"github.com/RicochetStudios/registry/api/v1/util"
)

func main() {
	// Parse the command line flags.
	// https://github.com/golang/go/issues/46869#issuecomment-865648650
	flag.Parse()

	// Setup the API.
	app := util.SetupAPI()

	// Start the API.
	log.Fatal(app.Listen(util.GetEnvWithDefault("PORT", ":8443")))
}
