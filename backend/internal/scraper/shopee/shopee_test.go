package shopee

import (
	"context"
	"errors"
	"testing"
)

func TestConfigured(t *testing.T) {
	tests := []struct {
		name string
		cfg  Config
		want bool
	}{
		{"empty", Config{}, false},
		{"missing key", Config{PartnerID: "id"}, false},
		{"missing id", Config{PartnerKey: "key"}, false},
		{"both set", Config{PartnerID: "id", PartnerKey: "key"}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.cfg.Configured(); got != tt.want {
				t.Errorf("Configured() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestScrape_NotConfigured(t *testing.T) {
	a := New(Config{})
	_, err := a.Scrape(context.Background())
	if !errors.Is(err, ErrNotConfigured) {
		t.Errorf("expected ErrNotConfigured, got %v", err)
	}
}

func TestBrandSlug(t *testing.T) {
	a := New(Config{})
	if got := a.BrandSlug(); got != "shopee" {
		t.Errorf("BrandSlug() = %q, want %q", got, "shopee")
	}
}
