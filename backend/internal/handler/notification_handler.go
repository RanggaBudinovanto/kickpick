package handler

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	dbutil "github.com/kickpick/backend/internal/db"
	"github.com/kickpick/backend/internal/db/sqlc"
	"github.com/kickpick/backend/internal/middleware"
)

type NotificationHandler struct {
	Queries *sqlc.Queries
}

func NewNotificationHandler(pool *pgxpool.Pool) *NotificationHandler {
	return &NotificationHandler{Queries: sqlc.New(pool)}
}

func (h *NotificationHandler) List(c *fiber.Ctx) error {
	uid, _ := middleware.UserIDFromContext(c)

	page, _ := strconv.Atoi(c.Query("page", "1"))
	if page < 1 {
		page = 1
	}
	limit := 20

	items, err := h.Queries.ListNotificationsByUser(c.Context(), sqlc.ListNotificationsByUserParams{
		UserID: dbutil.UUID(uid),
		Limit:  int32(limit),
		Offset: int32((page - 1) * limit),
	})
	if err != nil {
		return serverError(c, err)
	}

	dtos := make([]fiber.Map, 0, len(items))
	for _, n := range items {
		dtos = append(dtos, fiber.Map{
			"id":         dbutil.ToUUID(n.ID).String(),
			"type":       n.Type,
			"title":      n.Title,
			"body":       n.Body,
			"action_url": n.ActionUrl.String,
			"is_read":    n.IsRead,
			"created_at": n.CreatedAt.Time,
		})
	}

	return c.JSON(fiber.Map{"data": dtos})
}

func (h *NotificationHandler) UnreadCount(c *fiber.Ctx) error {
	uid, _ := middleware.UserIDFromContext(c)

	count, err := h.Queries.CountUnreadNotifications(c.Context(), dbutil.UUID(uid))
	if err != nil {
		return serverError(c, err)
	}

	return c.JSON(fiber.Map{"count": count})
}

func (h *NotificationHandler) MarkRead(c *fiber.Ctx) error {
	uid, _ := middleware.UserIDFromContext(c)

	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return badRequest(c, "id tidak valid")
	}

	if err := h.Queries.MarkNotificationRead(c.Context(), sqlc.MarkNotificationReadParams{
		ID:     dbutil.UUID(id),
		UserID: dbutil.UUID(uid),
	}); err != nil {
		return serverError(c, err)
	}

	return c.JSON(fiber.Map{"message": "notifikasi ditandai terbaca"})
}
