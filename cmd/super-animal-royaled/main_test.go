package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestTenantRoutesAre501WhileUnconfigured(t *testing.T) {
	s := &server{}
	for _, route := range tenantRoutes {
		w := httptest.NewRecorder()
		s.unconfigured(w, httptest.NewRequest(http.MethodPost, route, nil))
		if w.Code != http.StatusNotImplemented {
			t.Errorf("%s: got %d, want 501", route, w.Code)
		}
		var body struct {
			Code   int    `json:"code"`
			Status string `json:"status"`
		}
		if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
			t.Errorf("%s: body is not JSON: %v", route, err)
			continue
		}
		if body.Code != http.StatusNotImplemented || body.Status != "NotImplemented" {
			t.Errorf("%s: envelope %+v, want PlayFab 501 envelope", route, body)
		}
	}
}

func TestConfigFromEnv(t *testing.T) {
	t.Setenv("SAR_HTTP_ADDR", ":9999")
	t.Setenv("SAR_HEALTH_ADDR", ":9998")
	t.Setenv("PLAYFAB_TITLE_ID", "ABC12")
	cfg := configFromEnv()
	if cfg.httpAddr != ":9999" || cfg.healthAddr != ":9998" || cfg.titleID != "ABC12" {
		t.Fatalf("env config: %+v", cfg)
	}

	os.Unsetenv("SAR_HTTP_ADDR")
	os.Unsetenv("SAR_HEALTH_ADDR")
	os.Unsetenv("PLAYFAB_TITLE_ID")
	cfg = configFromEnv()
	if cfg.httpAddr != ":23620" || cfg.healthAddr != ":23621" || cfg.titleID != "" {
		t.Fatalf("default config: %+v", cfg)
	}
}
