package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/zenkriztao/ayo-football-backend/internal/domain/entity"
)

// MatchRepository defines the interface for match data operations
type MatchRepository interface {
	Create(ctx context.Context, match *entity.Match) error
	FindByID(ctx context.Context, id uuid.UUID) (*entity.Match, error)
	FindByIDWithDetails(ctx context.Context, id uuid.UUID) (*entity.Match, error)
	Update(ctx context.Context, match *entity.Match) error
	Delete(ctx context.Context, id uuid.UUID) error
	FindAll(ctx context.Context, page, limit int) ([]entity.Match, int64, error)
	FindByDateRange(ctx context.Context, startDate, endDate time.Time, page, limit int) ([]entity.Match, int64, error)
	FindByTeamID(ctx context.Context, teamID uuid.UUID, page, limit int) ([]entity.Match, int64, error)
	FindByStatus(ctx context.Context, status entity.MatchStatus, page, limit int) ([]entity.Match, int64, error)
	Exists(ctx context.Context, id uuid.UUID) (bool, error)
	GetTeamWinCount(ctx context.Context, teamID uuid.UUID, isHome bool) (int64, error)
	GetCompletedMatches(ctx context.Context, page, limit int) ([]entity.Match, int64, error)
}
