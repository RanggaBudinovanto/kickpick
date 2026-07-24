package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	dbutil "github.com/kickpick/backend/internal/db"
	"github.com/kickpick/backend/internal/db/sqlc"
	"github.com/kickpick/backend/internal/middleware"
)

type WishlistHandler struct {
	Queries *sqlc.Queries
}

func NewWishlistHandler(pool *pgxpool.Pool) *WishlistHandler {
	return &WishlistHandler{Queries: sqlc.New(pool)}
}

func (h *WishlistHandler) List(c *fiber.Ctx) error {
	uid, _ := middleware.UserIDFromContext(c)

	items, err := h.Queries.ListWishlistByUser(c.Context(), dbutil.UUID(uid))
	if err != nil {
		return serverError(c, err)
	}

	dtos := make([]fiber.Map, 0, len(items))
	for _, w := range items {
		dtos = append(dtos, fiber.Map{
			"id":           dbutil.ToUUID(w.ID).String(),
			"product_id":   dbutil.ToUUID(w.ProductID).String(),
			"product_name": w.ProductName,
			"product_slug": w.ProductSlug,
			"alert_active": w.AlertActive,
			"alert_type":   w.AlertType.String,
		})
	}

	return c.JSON(fiber.Map{"data": dtos})
}

type addWishlistRequest struct {
	ProductID string `json:"product_id"`
}

func (h *WishlistHandler) Add(c *fiber.Ctx) error {
	uid, _ := middleware.UserIDFromContext(c)

	var req addWishlistRequest
	if err := c.BodyParser(&req); err != nil {
		return badRequest(c, "invalid request body")
	}

	productID, err := uuid.Parse(req.ProductID)
	if err != nil {
		return badRequest(c, "product_id tidak valid")
	}

	item, err := h.Queries.AddWishlistItem(c.Context(), sqlc.AddWishlistItemParams{
		UserID:      dbutil.UUID(uid),
		ProductID:   dbutil.UUID(productID),
		AlertActive: false,
	})
	if err != nil {
		return serverError(c, err)
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"id":         dbutil.ToUUID(item.ID).String(),
		"product_id": dbutil.ToUUID(item.ProductID).String(),
	})
}

func (h *WishlistHandler) Remove(c *fiber.Ctx) error {
	uid, _ := middleware.UserIDFromContext(c)

	itemID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return badRequest(c, "id tidak valid")
	}

	if err := h.Queries.RemoveWishlistItem(c.Context(), sqlc.RemoveWishlistItemParams{
		ID:     dbutil.UUID(itemID),
		UserID: dbutil.UUID(uid),
	}); err != nil {
		return serverError(c, err)
	}

	return c.SendStatus(fiber.StatusNoContent)
}

type setAlertRequest struct {
	AlertActive bool   `json:"alert_active"`
	AlertType   string `json:"alert_type"`
}

func (h *WishlistHandler) SetAlert(c *fiber.Ctx) error {
	uid, _ := middleware.UserIDFromContext(c)

	itemID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return badRequest(c, "id tidak valid")
	}

	var req setAlertRequest
	if err := c.BodyParser(&req); err != nil {
		return badRequest(c, "invalid request body")
	}

	item, err := h.Queries.SetWishlistAlert(c.Context(), sqlc.SetWishlistAlertParams{
		ID:          dbutil.UUID(itemID),
		UserID:      dbutil.UUID(uid),
		AlertActive: req.AlertActive,
		AlertType:   dbutil.Text(req.AlertType),
	})
	if err != nil {
		return serverError(c, err)
	}

	return c.JSON(fiber.Map{
		"id":           dbutil.ToUUID(item.ID).String(),
		"alert_active": item.AlertActive,
		"alert_type":   item.AlertType.String,
	})
}
