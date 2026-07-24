package handler

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	dbutil "github.com/kickpick/backend/internal/db"
	"github.com/kickpick/backend/internal/db/sqlc"
	"github.com/kickpick/backend/internal/middleware"
)

type ReviewHandler struct {
	Queries *sqlc.Queries
}

func NewReviewHandler(pool *pgxpool.Pool) *ReviewHandler {
	return &ReviewHandler{Queries: sqlc.New(pool)}
}

type createReviewRequest struct {
	ProductID   string `json:"product_id"`
	Rating      int32  `json:"rating"`
	Comment     string `json:"comment"`
	FitFeedback string `json:"fit_feedback"`
}

var allowedFitFeedback = map[string]bool{
	"kekecilan": true,
	"pas":       true,
	"kebesaran": true,
}

func (h *ReviewHandler) Create(c *fiber.Ctx) error {
	uid, _ := middleware.UserIDFromContext(c)

	var req createReviewRequest
	if err := c.BodyParser(&req); err != nil {
		return badRequest(c, "invalid request body")
	}

	productID, err := uuid.Parse(req.ProductID)
	if err != nil {
		return badRequest(c, "product_id tidak valid")
	}
	if req.Rating < 1 || req.Rating > 5 {
		return badRequest(c, "rating harus antara 1 sampai 5")
	}
	req.FitFeedback = strings.ToLower(req.FitFeedback)
	if req.FitFeedback != "" && !allowedFitFeedback[req.FitFeedback] {
		return badRequest(c, "fit_feedback tidak valid")
	}

	review, err := h.Queries.CreateReview(c.Context(), sqlc.CreateReviewParams{
		ProductID:   dbutil.UUID(productID),
		UserID:      dbutil.UUID(uid),
		Rating:      req.Rating,
		Comment:     dbutil.Text(req.Comment),
		FitFeedback: dbutil.Text(req.FitFeedback),
	})
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key") {
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": "kamu sudah memberi review untuk produk ini"})
		}
		return serverError(c, err)
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"id":      dbutil.ToUUID(review.ID).String(),
		"message": "review berhasil dikirim, menunggu moderasi",
	})
}

type reportReviewRequest struct {
	Reason string `json:"reason"`
}

func (h *ReviewHandler) Report(c *fiber.Ctx) error {
	uid, _ := middleware.UserIDFromContext(c)

	reviewID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return badRequest(c, "id tidak valid")
	}

	var req reportReviewRequest
	_ = c.BodyParser(&req)

	if err := h.Queries.ReportReview(c.Context(), sqlc.ReportReviewParams{
		ReviewID:   dbutil.UUID(reviewID),
		ReportedBy: dbutil.UUID(uid),
		Reason:     dbutil.Text(req.Reason),
	}); err != nil {
		return serverError(c, err)
	}

	if err := h.Queries.FlagReviewIfReportsExceedThreshold(c.Context(), dbutil.UUID(reviewID)); err != nil {
		return serverError(c, err)
	}

	return c.JSON(fiber.Map{"message": "laporan diterima"})
}
