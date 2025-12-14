package database

import (
	"context"

	"github.com/google/uuid"
	"github.com/zenkriztao/ayo-football-backend/internal/domain/entity"
	"github.com/zenkriztao/ayo-football-backend/internal/domain/repository"
	"gorm.io/gorm"
)

type playerRepositoryImpl struct {
	db *gorm.DB
}

// NewPlayerRepository creates a new instance of PlayerRepository
func NewPlayerRepository(db *gorm.DB) repository.PlayerRepository {
	return &playerRepositoryImpl{db: db}
}

func (r *playerRepositoryImpl) Create(ctx context.Context, player *entity.Player) error {
	return r.db.WithContext(ctx).Create(player).Error
}

func (r *playerRepositoryImpl) FindByID(ctx context.Context, id uuid.UUID) (*entity.Player, error) {
	var player entity.Player
	err := r.db.WithContext(ctx).First(&player, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &player, nil
}

func (r *playerRepositoryImpl) FindByIDWithTeam(ctx context.Context, id uuid.UUID) (*entity.Player, error) {
	var player entity.Player
	err := r.db.WithContext(ctx).
		Preload("Team").
		First(&player, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &player, nil
}

func (r *playerRepositoryImpl) Update(ctx context.Context, player *entity.Player) error {
	return r.db.WithContext(ctx).Save(player).Error
}

func (r *playerRepositoryImpl) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&entity.Player{}, "id = ?", id).Error
}

func (r *playerRepositoryImpl) FindAll(ctx context.Context, page, limit int) ([]entity.Player, int64, error) {
	var players []entity.Player
	var total int64

	offset := (page - 1) * limit

	err := r.db.WithContext(ctx).Model(&entity.Player{}).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = r.db.WithContext(ctx).
		Preload("Team").
		Offset(offset).
		Limit(limit).
		Order("created_at DESC").
		Find(&players).Error
	if err != nil {
		return nil, 0, err
	}

	return players, total, nil
}

func (r *playerRepositoryImpl) FindByTeamID(ctx context.Context, teamID uuid.UUID, page, limit int) ([]entity.Player, int64, error) {
	var players []entity.Player
	var total int64

	offset := (page - 1) * limit

	err := r.db.WithContext(ctx).
		Model(&entity.Player{}).
		Where("team_id = ?", teamID).
		Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = r.db.WithContext(ctx).
		Where("team_id = ?", teamID).
		Offset(offset).
		Limit(limit).
		Order("jersey_number ASC").
		Find(&players).Error
	if err != nil {
		return nil, 0, err
	}

	return players, total, nil
}

func (r *playerRepositoryImpl) IsJerseyNumberTaken(ctx context.Context, teamID uuid.UUID, jerseyNumber int, excludePlayerID *uuid.UUID) (bool, error) {
	var count int64
	query := r.db.WithContext(ctx).
		Model(&entity.Player{}).
		Where("team_id = ? AND jersey_number = ?", teamID, jerseyNumber)

	if excludePlayerID != nil {
		query = query.Where("id != ?", *excludePlayerID)
	}

	err := query.Count(&count).Error
	return count > 0, err
}

func (r *playerRepositoryImpl) Search(ctx context.Context, query string, page, limit int) ([]entity.Player, int64, error) {
	var players []entity.Player
	var total int64

	offset := (page - 1) * limit
	searchQuery := "%" + query + "%"

	err := r.db.WithContext(ctx).
		Model(&entity.Player{}).
		Where("name ILIKE ?", searchQuery).
		Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = r.db.WithContext(ctx).
		Preload("Team").
		Where("name ILIKE ?", searchQuery).
		Offset(offset).
		Limit(limit).
		Order("created_at DESC").
		Find(&players).Error
	if err != nil {
		return nil, 0, err
	}

	return players, total, nil
}

func (r *playerRepositoryImpl) Exists(ctx context.Context, id uuid.UUID) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&entity.Player{}).
		Where("id = ?", id).
		Count(&count).Error
	return count > 0, err
}

func (r *playerRepositoryImpl) GetTopScorers(ctx context.Context, limit int) ([]repository.PlayerGoalCount, error) {
	var results []repository.PlayerGoalCount

	err := r.db.WithContext(ctx).
		Table("goals").
		Select("players.*, COUNT(goals.id) as goal_count").
		Joins("JOIN players ON players.id = goals.player_id").
		Where("goals.deleted_at IS NULL AND players.deleted_at IS NULL").
		Group("players.id").
		Order("goal_count DESC").
		Limit(limit).
		Scan(&results).Error

	return results, err
}
