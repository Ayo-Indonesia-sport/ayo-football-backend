package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/zenkriztao/ayo-football-backend/internal/domain/entity"
)

// GoalRepository defines the interface for goal data operations
type GoalRepository interface {
	Create(ctx context.Context, goal *entity.Goal) error
	CreateBatch(ctx context.Context, goals []entity.Goal) error
	FindByID(ctx context.Context, id uuid.UUID) (*entity.Goal, error)
	Update(ctx context.Context, goal *entity.Goal) error
	Delete(ctx context.Context, id uuid.UUID) error
	FindByMatchID(ctx context.Context, matchID uuid.UUID) ([]entity.Goal, error)
	FindByPlayerID(ctx context.Context, playerID uuid.UUID) ([]entity.Goal, error)
	DeleteByMatchID(ctx context.Context, matchID uuid.UUID) error
	GetTopScorers(ctx context.Context, limit int) ([]TopScorerResult, error)
}

// TopScorerResult represents a player with their goal statistics
type TopScorerResult struct {
	PlayerID   uuid.UUID
	PlayerName string
	TeamID     uuid.UUID
	TeamName   string
	GoalCount  int64
}
