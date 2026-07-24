package handler

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"

	dbutil "github.com/kickpick/backend/internal/db"
	"github.com/kickpick/backend/internal/db/sqlc"
)

// logAudit records a sensitive action per Section 14 PRD (login, password change,
// account deletion, etc). userID is optional (nil for actions before a user is known).
func logAudit(ctx context.Context, q *sqlc.Queries, userID *uuid.UUID, action, ip string) {
	var pgUserID pgtype.UUID
	if userID != nil {
		pgUserID = dbutil.UUID(*userID)
	}

	_ = q.InsertAuditLog(ctx, sqlc.InsertAuditLogParams{
		UserID:    pgUserID,
		Action:    action,
		IpAddress: dbutil.Text(ip),
		Metadata:  []byte("{}"),
	})
}
