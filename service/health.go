package service

import (
	"context"
	"log"
	"time"

	sdk "agones.dev/agones/sdks/go"
)

// doHealth sends the regular Health Pings
func DoHealth(w ServerWrapper, sdk *sdk.SDK, ctx context.Context) {
	tick := time.NewTicker(2 * time.Second)
	for {
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
