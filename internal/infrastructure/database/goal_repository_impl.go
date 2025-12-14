package database

import (
	"context"

	"github.com/google/uuid"
	"github.com/zenkriztao/ayo-football-backend/internal/domain/entity"
	"github.com/zenkriztao/ayo-football-backend/internal/domain/repository"
	"gorm.io/gorm"
)

type goalRepositoryImpl struct {
	db *gorm.DB
}

// NewGoalRepository creates a new instance of GoalRepository
func NewGoalRepository(db *gorm.DB) repository.GoalRepository {
	return &goalRepositoryImpl{db: db}
}

func (r *goalRepositoryImpl) Create(ctx context.Context, goal *entity.Goal) error {
	return r.db.WithContext(ctx).Create(goal).Error
}

func (r *goalRepositoryImpl) CreateBatch(ctx context.Context, goals []entity.Goal) error {
	if len(goals) == 0 {
		return nil
	}
	return r.db.WithContext(ctx).Create(&goals).Error
}

func (r *goalRepositoryImpl) FindByID(ctx context.Context, id uuid.UUID) (*entity.Goal, error) {
	var goal entity.Goal
	err := r.db.WithContext(ctx).
		Preload("Player").
		Preload("Team").
		First(&goal, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &goal, nil
}

func (r *goalRepositoryImpl) Update(ctx context.Context, goal *entity.Goal) error {
	return r.db.WithContext(ctx).Save(goal).Error
}

func (r *goalRepositoryImpl) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&entity.Goal{}, "id = ?", id).Error
}

func (r *goalRepositoryImpl) FindByMatchID(ctx context.Context, matchID uuid.UUID) ([]entity.Goal, error) {
	var goals []entity.Goal
	err := r.db.WithContext(ctx).
		Preload("Player").
		Preload("Team").
		Where("match_id = ?", matchID).
		Order("minute ASC").
		Find(&goals).Error
	return goals, err
}

func (r *goalRepositoryImpl) FindByPlayerID(ctx context.Context, playerID uuid.UUID) ([]entity.Goal, error) {
	var goals []entity.Goal
	err := r.db.WithContext(ctx).
		Preload("Match").
		Preload("Team").
		Where("player_id = ?", playerID).
		Order("created_at DESC").
		Find(&goals).Error
	return goals, err
}

func (r *goalRepositoryImpl) DeleteByMatchID(ctx context.Context, matchID uuid.UUID) error {
	return r.db.WithContext(ctx).
		Where("match_id = ?", matchID).
		Delete(&entity.Goal{}).Error
}

func (r *goalRepositoryImpl) GetTopScorers(ctx context.Context, limit int) ([]repository.TopScorerResult, error) {
	var results []repository.TopScorerResult

	err := r.db.WithContext(ctx).
		Table("goals").
		Select("goals.player_id, players.name as player_name, players.team_id, teams.name as team_name, COUNT(goals.id) as goal_count").
		Joins("JOIN players ON players.id = goals.player_id AND players.deleted_at IS NULL").
		Joins("JOIN teams ON teams.id = players.team_id AND teams.deleted_at IS NULL").
		Where("goals.deleted_at IS NULL AND goals.is_own_goal = false").
		Group("goals.player_id, players.name, players.team_id, teams.name").
		Order("goal_count DESC").
		Limit(limit).
		Scan(&results).Error

	return results, err
}
