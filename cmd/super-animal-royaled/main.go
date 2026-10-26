// Command super-animal-royaled answers the PlayFab endpoints Super Animal Royale's
// Switch build (0100DB3011E8C000) calls when a console is redirected to OpenPak.
//
// The title is PlayFab-only: UnitySDK Client login (LoginWithNintendoSwitchDeviceId,
// LoginWithNintendoServiceAccount), Multiplayer Servers + matchmaking/QoS, and
// PlayFab Party native for voice/chat/networks. There is no Photon realtime and no
// NEX/PRUDP in this title. See docs/protocol.md for the wire facts.
//
// This service is a scaffold: every PlayFab tenant route answers 501 until the
// title's PlayFab TitleId is recovered (it lives in a compressed
// PlayFabSharedSettings asset inside data.unity3d) and a tenant is configured.
// Party relay itself belongs to platforms/switch/games/playfab
// (internal/party); this daemon only fronts the title-specific HTTP surface.
package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
)

// tenantRoutes are the PlayFab API calls the client is known to make, from the
// recon of the Switch executable (see docs/protocol.md). Each answers 501 while
// no tenant is configured.
var tenantRoutes = []string{
	"/Client/LoginWithNintendoSwitchDeviceId",
	"/Client/LoginWithNintendoServiceAccount",
	"/Match/CreateMatchmakingTicket",
	"/Match/GetMatchmakingTicket",
	"/Match/CancelMatchmakingTicket",
	"/Match/CreateServerBackfillTicket",
	"/MultiplayerServer/GetMultiplayerServerDetails",
	"/MultiplayerServer/GetGameServerRegions",
	"/MultiplayerServer/ListPartyQosServers",
	"/MultiplayerServer/GetRemoteLoginEndpoint",
}

type config struct {
	httpAddr   string
	healthAddr string
	titleID    string
}

func configFromEnv() config {
	return config{
		httpAddr:   envOr("SAR_HTTP_ADDR", ":23620"),
		healthAddr: envOr("SAR_HEALTH_ADDR", ":23621"),
		titleID:    os.Getenv("PLAYFAB_TITLE_ID"),
	}
}

func envOr(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func main() {
	cfg := configFromEnv()
	s := &server{titleID: cfg.titleID}

	go func() {
		h := http.NewServeMux()
		h.HandleFunc("/healthz", func(w http.ResponseWriter, _ *http.Request) {
			w.Header().Set("Content-Type", "text/plain")
			w.Write([]byte("ok\n"))
		})
		log.Printf("[super-animal-royale] health on %s", cfg.healthAddr)
		log.Fatal(http.ListenAndServe(cfg.healthAddr, h))
	}()

	mux := http.NewServeMux()
	for _, route := range tenantRoutes {
		mux.HandleFunc(route, s.unconfigured)
	}
	log.Printf("[super-animal-royale] serving HTTP on %s (title=%q)", cfg.httpAddr, cfg.titleID)
	log.Fatal(http.ListenAndServe(cfg.httpAddr, mux))
}

type server struct {
	titleID string
}

// unconfigured answers every tenant route until a PlayFab TitleId is set. It
// uses the PlayFab SDK envelope (code/status/data) so the client sees a
// well-formed error rather than a connection failure.
func (s *server) unconfigured(w http.ResponseWriter, r *http.Request) {
	log.Printf("[super-animal-royale] %s %s: no tenant configured", r.Method, r.URL.Path)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNotImplemented)
	json.NewEncoder(w).Encode(map[string]any{
		"code":   http.StatusNotImplemented,
		"status": "NotImplemented",
		"error":  "PlayFab tenant not configured (PLAYFAB_TITLE_ID unset)",
		"data":   nil,
	})
}
