package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	dbutil "github.com/kickpick/backend/internal/db"
	"github.com/kickpick/backend/internal/db/sqlc"
)

type RedirectHandler struct {
	Queries *sqlc.Queries
}

func NewRedirectHandler(pool *pgxpool.Pool) *RedirectHandler {
	return &RedirectHandler{Queries: sqlc.New(pool)}
}

// GoToOffer logs the click-through and returns the store's affiliate URL for the
// frontend to navigate to. Section 12 PRD: POST /api/redirect/:offer_id
// (rate limited per IP, no auth required). POST is used (rather than a raw 302)
// so the click can be logged via fetch before the browser navigates away.
func (h *RedirectHandler) GoToOffer(c *fiber.Ctx) error {
	offerID, err := uuid.Parse(c.Params("offer_id"))
	if err != nil {
		return badRequest(c, "offer_id tidak valid")
	}

	offer, err := h.Queries.GetOfferByID(c.Context(), dbutil.UUID(offerID))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "penawaran tidak ditemukan"})
	}

	_ = h.Queries.InsertAuditLog(c.Context(), sqlc.InsertAuditLogParams{
		Action:    "affiliate_click",
		IpAddress: dbutil.Text(c.IP()),
		Metadata:  []byte(`{"offer_id":"` + offerID.String() + `"}`),
	})

	return c.JSON(fiber.Map{"affiliate_url": offer.AffiliateUrl})
}
