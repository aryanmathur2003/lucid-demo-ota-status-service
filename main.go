// Command ota-status-service reports the progress of over-the-air update
// campaigns to internal dashboards.
package main

import (
	"log"
	"net/http"
	"os"

	"github.com/aryanmathur2003/lucid-demo-ota-status-service/internal/ota"
)

func main() {
	addr := os.Getenv("LISTEN_ADDR")
	if addr == "" {
		addr = ":8081"
	}
	log.Printf("ota-status-service listening on %s", addr)
	log.Fatal(http.ListenAndServe(addr, ota.NewServer(ota.SeedCampaigns()).Routes()))
}
