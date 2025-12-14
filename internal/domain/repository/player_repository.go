package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/zenkriztao/ayo-football-backend/internal/domain/entity"
)

// PlayerRepository defines the interface for player data operations
type PlayerRepository interface {
	Create(ctx context.Context, player *entity.Player) error
	FindByID(ctx context.Context, id uuid.UUID) (*entity.Player, error)
	FindByIDWithTeam(ctx context.Context, id uuid.UUID) (*entity.Player, error)
	Update(ctx context.Context, player *entity.Player) error
	Delete(ctx context.Context, id uuid.UUID) error
	FindAll(ctx context.Context, page, limit int) ([]entity.Player, int64, error)
	FindByTeamID(ctx context.Context, teamID uuid.UUID, page, limit int) ([]entity.Player, int64, error)
	IsJerseyNumberTaken(ctx context.Context, teamID uuid.UUID, jerseyNumber int, excludePlayerID *uuid.UUID) (bool, error)
	Search(ctx context.Context, query string, page, limit int) ([]entity.Player, int64, error)
	Exists(ctx context.Context, id uuid.UUID) (bool, error)
	GetTopScorers(ctx context.Context, limit int) ([]PlayerGoalCount, error)
}

// PlayerGoalCount represents a player with their goal count
type PlayerGoalCount struct {
	Player    entity.Player
	GoalCount int64
}
