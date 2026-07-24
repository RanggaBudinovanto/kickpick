package db

import (
	"testing"

	"github.com/google/uuid"
)

// Float64ToNumeric silently produced an invalid (NULL) Numeric when Scan was
// given a raw float64 instead of a string — caught during Plan 3 scraper
// testing when every price insert failed with a NOT NULL violation. This test
// guards against that regression.
func TestFloat64ToNumericRoundTrip(t *testing.T) {
	cases := []float64{0, 1, 1698000, 999999.99, 0.5}

	for _, want := range cases {
		n := Float64ToNumeric(want)
		if !n.Valid {
			t.Fatalf("Float64ToNumeric(%v) produced an invalid Numeric", want)
		}
		got := NumericToFloat64(n)
		if got != want {
			t.Errorf("round trip mismatch: got %v, want %v", got, want)
		}
	}
}

func TestUUIDRoundTrip(t *testing.T) {
	want := uuid.New()
	pg := UUID(want)
	if !pg.Valid {
		t.Fatal("UUID() produced an invalid pgtype.UUID")
	}
	if got := ToUUID(pg); got != want {
		t.Errorf("round trip mismatch: got %v, want %v", got, want)
	}
}

func TestTextEmptyStringIsInvalid(t *testing.T) {
	if Text("").Valid {
		t.Error("Text(\"\") should be invalid (NULL), not an empty string value")
	}
	if got := Text("hello"); !got.Valid || got.String != "hello" {
		t.Errorf("Text(\"hello\") = %+v, want valid with String=hello", got)
	}
}
