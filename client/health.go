package client

import (
	"context"
	"log"
	"time"

	"github.com/RicochetStudios/registry/client/wrapper"

	sdk "agones.dev/agones/sdks/go"
)

const healthInterval = 2

// doHealth sends the regular Health Pings
func DoHealth(w wrapper.ServerWrapper, sdk *sdk.SDK, ctx context.Context) {
	tick := time.NewTicker(healthInterval * time.Second)
	for {
		// If the context is cancelled, stop the health check.
		if ctx.Err() != nil {
			return
		}

		// Check if the server is healthy.
		healthy, err := w.Healthy()
		if err != nil {
			log.Fatalf("[wrapper] Could not check health status, %v", err)
		}

		// If the server is healthy, send a health ping.
		if healthy {
			err := sdk.Health()
			if err != nil {
				log.Fatalf("[wrapper] Could not send health ping, %v", err)
			}
		}
		<-tick.C
	}
}
