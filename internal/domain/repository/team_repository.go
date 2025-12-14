package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/zenkriztao/ayo-football-backend/internal/domain/entity"
)

// TeamRepository defines the interface for team data operations
type TeamRepository interface {
	Create(ctx context.Context, team *entity.Team) error
	FindByID(ctx context.Context, id uuid.UUID) (*entity.Team, error)
	FindByIDWithPlayers(ctx context.Context, id uuid.UUID) (*entity.Team, error)
	Update(ctx context.Context, team *entity.Team) error
	Delete(ctx context.Context, id uuid.UUID) error
	FindAll(ctx context.Context, page, limit int) ([]entity.Team, int64, error)
	Search(ctx context.Context, query string, page, limit int) ([]entity.Team, int64, error)
	Exists(ctx context.Context, id uuid.UUID) (bool, error)
}
