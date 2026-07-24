package handler

import (
	"context"
	"encoding/json"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"

	dbutil "github.com/kickpick/backend/internal/db"
	"github.com/kickpick/backend/internal/db/sqlc"
)

type ProductHandler struct {
	Queries *sqlc.Queries
}

func NewProductHandler(pool *pgxpool.Pool) *ProductHandler {
	return &ProductHandler{Queries: sqlc.New(pool)}
}

type offerDTO struct {
	ID           string  `json:"id"`
	StoreName    string  `json:"store_name"`
	StoreType    string  `json:"store_type"`
	Price        float64 `json:"price"`
	Currency     string  `json:"currency"`
	InStock      bool    `json:"in_stock"`
	Size         string  `json:"size,omitempty"`
	AffiliateURL string  `json:"affiliate_url"`
}

type productDTO struct {
	ID         string          `json:"id"`
	Slug       string          `json:"slug"`
	Name       string          `json:"name"`
	Category   string          `json:"category"`
	BrandID    string          `json:"brand_id"`
	BrandName  string          `json:"brand_name"`
	BrandSlug  string          `json:"brand_slug"`
	IsLimited  bool            `json:"is_limited"`
	Attributes json.RawMessage `json:"attributes"`
	ImageURL   string          `json:"image_url"`
}

type productCardDTO struct {
	productDTO
	MinPrice float64 `json:"min_price"`
	MaxPrice float64 `json:"max_price"`
	Currency string  `json:"currency"`
	Rating   float64 `json:"rating"`
}

func toProductDTO(p sqlc.GetProductBySlugRow) productDTO {
	attrs := p.Attributes
	if len(attrs) == 0 {
		attrs = []byte("{}")
	}
	return productDTO{
		ID:         dbutil.ToUUID(p.ID).String(),
		Slug:       p.Slug,
		Name:       p.Name,
		Category:   p.Category,
		BrandID:    dbutil.ToUUID(p.BrandID).String(),
		BrandName:  p.BrandName,
		BrandSlug:  p.BrandSlug,
		IsLimited:  p.IsLimited,
		Attributes: attrs,
		ImageURL:   p.ImageUrl,
	}
}

func toProductCardDTO(p sqlc.ListProductsRow) productCardDTO {
	attrs := p.Attributes
	if len(attrs) == 0 {
		attrs = []byte("{}")
	}
	return productCardDTO{
		productDTO: productDTO{
			ID:         dbutil.ToUUID(p.ID).String(),
			Slug:       p.Slug,
			Name:       p.Name,
			Category:   p.Category,
			BrandID:    dbutil.ToUUID(p.BrandID).String(),
			BrandName:  p.BrandName,
			BrandSlug:  p.BrandSlug,
			IsLimited:  p.IsLimited,
			Attributes: attrs,
			ImageURL:   p.ImageUrl,
		},
		MinPrice: dbutil.NumericToFloat64(p.MinPrice),
		MaxPrice: dbutil.NumericToFloat64(p.MaxPrice),
		Currency: "IDR",
		Rating:   dbutil.NumericToFloat64(p.AvgRating),
	}
}

func (h *ProductHandler) ListProducts(c *fiber.Ctx) error {
	ctx := c.Context()

	page, _ := strconv.Atoi(c.Query("page", "1"))
	if page < 1 {
		page = 1
	}
	limit, _ := strconv.Atoi(c.Query("limit", "24"))
	if limit < 1 || limit > 100 {
		limit = 24
	}

	arg := sqlc.ListProductsParams{
		PageOffset: int32((page - 1) * limit),
		PageLimit:  int32(limit),
	}

	if brandSlug := c.Query("brand_id"); brandSlug != "" {
		if id, err := uuid.Parse(brandSlug); err == nil {
			arg.BrandID = dbutil.UUID(id)
		}
	}
	if brandIDs := c.Query("brand_ids"); brandIDs != "" {
		parts := strings.Split(brandIDs, ",")
		ids := make([]pgtype.UUID, 0, len(parts))
		for _, p := range parts {
			if id, err := uuid.Parse(strings.TrimSpace(p)); err == nil {
				ids = append(ids, dbutil.UUID(id))
			}
		}
		if len(ids) > 0 {
			arg.BrandIds = ids
		}
	}
	if category := c.Query("kategori"); category != "" {
		arg.Category = dbutil.Text(category)
	}
	if filter := c.Query("filter"); filter == "rare" {
		arg.IsLimited = dbutil.Bool(true)
	}
	if search := c.Query("q"); search != "" {
		arg.Search = dbutil.Text(search)
	}

	products, err := h.Queries.ListProducts(ctx, arg)
	if err != nil {
		return serverError(c, err)
	}

	dtos := make([]productCardDTO, 0, len(products))
	for _, p := range products {
		dtos = append(dtos, toProductCardDTO(p))
	}

	return c.JSON(fiber.Map{
		"data":  dtos,
		"page":  page,
		"limit": limit,
	})
}

func (h *ProductHandler) ListTrending(c *fiber.Ctx) error {
	limit, _ := strconv.Atoi(c.Query("limit", "8"))
	if limit < 1 || limit > 50 {
		limit = 8
	}

	products, err := h.Queries.ListTrendingProducts(c.Context(), int32(limit))
	if err != nil {
		return serverError(c, err)
	}

	dtos := make([]productCardDTO, 0, len(products))
	for _, p := range products {
		dtos = append(dtos, productCardDTO{
			productDTO: productDTO{
				ID:         dbutil.ToUUID(p.ID).String(),
				Slug:       p.Slug,
				Name:       p.Name,
				Category:   p.Category,
				BrandID:    dbutil.ToUUID(p.BrandID).String(),
				BrandName:  p.BrandName,
				BrandSlug:  p.BrandSlug,
				IsLimited:  p.IsLimited,
				Attributes: emptyJSONIfBlank(p.Attributes),
				ImageURL:   p.ImageUrl,
			},
			MinPrice: dbutil.NumericToFloat64(p.MinPrice),
			MaxPrice: dbutil.NumericToFloat64(p.MaxPrice),
			Currency: "IDR",
			Rating:   dbutil.NumericToFloat64(p.AvgRating),
		})
	}

	return c.JSON(fiber.Map{"data": dtos})
}

type priceDropDTO struct {
	productDTO
	MinPrice    float64 `json:"min_price"`
	MaxPrice    float64 `json:"max_price"`
	Currency    string  `json:"currency"`
	Rating      float64 `json:"rating"`
	DropPercent int32   `json:"drop_percent"`
}

func (h *ProductHandler) ListPriceDrops(c *fiber.Ctx) error {
	limit, _ := strconv.Atoi(c.Query("limit", "8"))
	if limit < 1 || limit > 50 {
		limit = 8
	}

	products, err := h.Queries.ListPriceDropProducts(c.Context(), int32(limit))
	if err != nil {
		return serverError(c, err)
	}

	dtos := make([]priceDropDTO, 0, len(products))
	for _, p := range products {
		dtos = append(dtos, priceDropDTO{
			productDTO: productDTO{
				ID:         dbutil.ToUUID(p.ID).String(),
				Slug:       p.Slug,
				Name:       p.Name,
				Category:   p.Category,
				BrandID:    dbutil.ToUUID(p.BrandID).String(),
				BrandName:  p.BrandName,
				BrandSlug:  p.BrandSlug,
				IsLimited:  p.IsLimited,
				Attributes: emptyJSONIfBlank(p.Attributes),
				ImageURL:   p.ImageUrl,
			},
			MinPrice:    dbutil.NumericToFloat64(p.MinPrice),
			MaxPrice:    dbutil.NumericToFloat64(p.MaxPrice),
			Currency:    "IDR",
			Rating:      dbutil.NumericToFloat64(p.AvgRating),
			DropPercent: p.DropPercent,
		})
	}

	return c.JSON(fiber.Map{"data": dtos})
}

func emptyJSONIfBlank(b []byte) json.RawMessage {
	if len(b) == 0 {
		return []byte("{}")
	}
	return b
}

func (h *ProductHandler) GetProductBySlug(c *fiber.Ctx) error {
	ctx := c.Context()
	slug := c.Params("slug")

	product, err := h.Queries.GetProductBySlug(ctx, slug)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "produk tidak ditemukan"})
	}

	// Fire-and-forget: powers the Trending section (7-day view aggregation).
	go func() {
		_ = h.Queries.RecordProductView(context.Background(), product.ID)
	}()

	offers, err := h.Queries.ListOffersByProductID(ctx, product.ID)
	if err != nil {
		return serverError(c, err)
	}

	offerDTOs := make([]offerDTO, 0, len(offers))
	for _, o := range offers {
		offerDTOs = append(offerDTOs, offerDTO{
			ID:           dbutil.ToUUID(o.ID).String(),
			StoreName:    o.StoreName,
			StoreType:    o.StoreType,
			Price:        dbutil.NumericToFloat64(o.Price),
			Currency:     o.Currency,
			InStock:      o.InStock,
			Size:         o.Size.String,
			AffiliateURL: o.AffiliateUrl,
		})
	}

	reviews, err := h.Queries.ListReviewsByProduct(ctx, product.ID)
	if err != nil {
		return serverError(c, err)
	}

	reviewDTOs := make([]fiber.Map, 0, len(reviews))
	for _, r := range reviews {
		reviewDTOs = append(reviewDTOs, fiber.Map{
			"id":           dbutil.ToUUID(r.ID).String(),
			"rating":       r.Rating,
			"comment":      r.Comment.String,
			"fit_feedback": r.FitFeedback.String,
			"user_name":    r.UserName,
			"created_at":   r.CreatedAt.Time,
		})
	}

	return c.JSON(fiber.Map{
		"product": toProductDTO(product),
		"offers":  offerDTOs,
		"reviews": reviewDTOs,
	})
}

func (h *ProductHandler) GetPriceHistory(c *fiber.Ctx) error {
	ctx := c.Context()
	slug := c.Params("slug")

	product, err := h.Queries.GetProductBySlug(ctx, slug)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "produk tidak ditemukan"})
	}

	days, _ := strconv.Atoi(c.Query("days", "30"))
	if days != 90 {
		days = 30
	}

	history, err := h.Queries.ListPriceHistoryByProductID(ctx, sqlc.ListPriceHistoryByProductIDParams{
		ProductID:    product.ID,
		RecordedDate: dbutil.Date(time.Now().AddDate(0, 0, -days)),
	})
	if err != nil {
		return serverError(c, err)
	}

	points := make([]fiber.Map, 0, len(history))
	for _, p := range history {
		points = append(points, fiber.Map{
			"date":  p.RecordedDate.Time.Format("2006-01-02"),
			"price": dbutil.NumericToFloat64(p.Price),
		})
	}

	return c.JSON(fiber.Map{"data": points})
}

func (h *ProductHandler) GetSizeConversion(c *fiber.Ctx) error {
	ctx := c.Context()
	slug := c.Params("slug")

	refBrandSlug := c.Query("reference_brand")
	size := c.Query("size")
	if refBrandSlug == "" || size == "" {
		return badRequest(c, "reference_brand dan size wajib diisi")
	}

	product, err := h.Queries.GetProductBySlug(ctx, slug)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "produk tidak ditemukan"})
	}

	refBrand, err := h.Queries.GetBrandBySlug(ctx, refBrandSlug)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "brand referensi tidak ditemukan"})
	}

	rows, err := h.Queries.GetSizeConversion(ctx, sqlc.GetSizeConversionParams{
		ReferenceBrandID: refBrand.ID,
		TargetBrandID:    product.BrandID,
		ReferenceSize:    size,
	})
	if err != nil {
		return serverError(c, err)
	}

	if len(rows) == 0 {
		return c.JSON(fiber.Map{"data": nil, "message": "data konversi belum tersedia untuk kombinasi ini"})
	}

	return c.JSON(fiber.Map{"data": rows[0].TargetSize})
}

func (h *ProductHandler) ListBrands(c *fiber.Ctx) error {
	brands, err := h.Queries.ListBrands(c.Context())
	if err != nil {
		return serverError(c, err)
	}

	dtos := make([]fiber.Map, 0, len(brands))
	for _, b := range brands {
		dtos = append(dtos, fiber.Map{
			"id":       dbutil.ToUUID(b.ID).String(),
			"name":     b.Name,
			"slug":     b.Slug,
			"logo_url": b.LogoUrl.String,
			"is_local": b.IsLocal,
		})
	}

	return c.JSON(fiber.Map{"data": dtos})
}

func (h *ProductHandler) GetBrandBySlug(c *fiber.Ctx) error {
	ctx := c.Context()
	slug := c.Params("slug")

	brand, err := h.Queries.GetBrandBySlug(ctx, slug)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "brand tidak ditemukan"})
	}

	products, err := h.Queries.ListProducts(ctx, sqlc.ListProductsParams{
		BrandID:    brand.ID,
		PageLimit:  100,
		PageOffset: 0,
	})
	if err != nil {
		return serverError(c, err)
	}

	dtos := make([]productCardDTO, 0, len(products))
	for _, p := range products {
		dtos = append(dtos, toProductCardDTO(p))
	}

	return c.JSON(fiber.Map{
		"brand": fiber.Map{
			"id":       dbutil.ToUUID(brand.ID).String(),
			"name":     brand.Name,
			"slug":     brand.Slug,
			"logo_url": brand.LogoUrl.String,
			"is_local": brand.IsLocal,
		},
		"products": dtos,
	})
}

func (h *ProductHandler) Autocomplete(c *fiber.Ctx) error {
	ctx := c.Context()
	query := c.Query("q")
	if query == "" {
		return c.JSON(fiber.Map{"products": []any{}, "brands": []any{}})
	}

	products, err := h.Queries.SearchProductsAutocomplete(ctx, query)
	if err != nil {
		return serverError(c, err)
	}
	brands, err := h.Queries.SearchBrandsAutocomplete(ctx, query)
	if err != nil {
		return serverError(c, err)
	}

	productResults := make([]fiber.Map, 0, len(products))
	for _, p := range products {
		productResults = append(productResults, fiber.Map{
			"id":   dbutil.ToUUID(p.ID).String(),
			"name": p.Name,
			"slug": p.Slug,
		})
	}
	brandResults := make([]fiber.Map, 0, len(brands))
	for _, b := range brands {
		brandResults = append(brandResults, fiber.Map{
			"id":   dbutil.ToUUID(b.ID).String(),
			"name": b.Name,
			"slug": b.Slug,
		})
	}

	return c.JSON(fiber.Map{"products": productResults, "brands": brandResults})
}
