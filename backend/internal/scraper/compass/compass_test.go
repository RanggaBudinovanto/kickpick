package compass

import (
	"os"
	"testing"
)

// The fixture is a real snapshot of sepatucompass.com/collections/all-products
// captured during Plan 3 development, so this test exercises the actual JSON
// shape the site serves rather than a hand-crafted approximation of it.
func TestExtractProductNodesFromRealFixture(t *testing.T) {
	body, err := os.ReadFile("testdata/catalog_sample.html")
	if err != nil {
		t.Fatalf("failed to read fixture: %v", err)
	}

	nodes, err := extractProductNodes(body)
	if err != nil {
		t.Fatalf("extractProductNodes failed: %v", err)
	}

	if len(nodes) < 50 {
		t.Errorf("expected at least 50 products in the fixture, got %d", len(nodes))
	}

	node, ok := nodes["sepatu-compass-velocity-moc-eye_c"]
	if !ok {
		t.Fatal("expected fixture to contain handle sepatu-compass-velocity-moc-eye_c")
	}
	if node.Title != "Velocity Moc Eye_C" {
		t.Errorf("Title = %q, want %q", node.Title, "Velocity Moc Eye_C")
	}
	if node.PriceRange.MinVariantPrice.Amount == "" {
		t.Error("expected a non-empty price for the known product")
	}
	if node.PriceRange.MinVariantPrice.CurrencyCode != "IDR" {
		t.Errorf("CurrencyCode = %q, want IDR", node.PriceRange.MinVariantPrice.CurrencyCode)
	}
}

func TestExtractProductNodesRejectsUnrelatedHTML(t *testing.T) {
	_, err := extractProductNodes([]byte("<html><body>not compass</body></html>"))
	if err == nil {
		t.Error("expected an error when the remixContext marker is absent")
	}
}

func TestToScrapedProductsSkipsUnparsablePrices(t *testing.T) {
	good := productNode{Title: "Good Shoe", Handle: "good"}
	good.PriceRange.MinVariantPrice.Amount = "100000.0"
	good.PriceRange.MinVariantPrice.CurrencyCode = "IDR"

	bad := productNode{Title: "Bad Shoe", Handle: "bad"}
	bad.PriceRange.MinVariantPrice.Amount = "not-a-number"

	nodes := map[string]productNode{"good": good, "bad": bad}

	products := toScrapedProducts(nodes)

	if len(products) != 1 {
		t.Fatalf("expected 1 product (bad price skipped), got %d", len(products))
	}
	if products[0].Name != "Good Shoe" {
		t.Errorf("Name = %q, want Good Shoe", products[0].Name)
	}
	if products[0].Offers[0].Price != 100000 {
		t.Errorf("Price = %v, want 100000", products[0].Offers[0].Price)
	}
}
