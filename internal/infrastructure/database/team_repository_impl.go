package database

import (
	"context"

	"github.com/google/uuid"
	"github.com/zenkriztao/ayo-football-backend/internal/domain/entity"
	"github.com/zenkriztao/ayo-football-backend/internal/domain/repository"
	"gorm.io/gorm"
)

type teamRepositoryImpl struct {
	db *gorm.DB
}

// NewTeamRepository creates a new instance of TeamRepository
func NewTeamRepository(db *gorm.DB) repository.TeamRepository {
	return &teamRepositoryImpl{db: db}
}

func (r *teamRepositoryImpl) Create(ctx context.Context, team *entity.Team) error {
	return r.db.WithContext(ctx).Create(team).Error
}

func (r *teamRepositoryImpl) FindByID(ctx context.Context, id uuid.UUID) (*entity.Team, error) {
	var team entity.Team
	err := r.db.WithContext(ctx).First(&team, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &team, nil
}

func (r *teamRepositoryImpl) FindByIDWithPlayers(ctx context.Context, id uuid.UUID) (*entity.Team, error) {
	var team entity.Team
	err := r.db.WithContext(ctx).
		Preload("Players").
		First(&team, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &team, nil
}

func (r *teamRepositoryImpl) Update(ctx context.Context, team *entity.Team) error {
	return r.db.WithContext(ctx).Save(team).Error
}

func (r *teamRepositoryImpl) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&entity.Team{}, "id = ?", id).Error
}

func (r *teamRepositoryImpl) FindAll(ctx context.Context, page, limit int) ([]entity.Team, int64, error) {
	var teams []entity.Team
	var total int64

	offset := (page - 1) * limit

	err := r.db.WithContext(ctx).Model(&entity.Team{}).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = r.db.WithContext(ctx).
		Offset(offset).
		Limit(limit).
		Order("created_at DESC").
		Find(&teams).Error
	if err != nil {
		return nil, 0, err
	}

	return teams, total, nil
}

func (r *teamRepositoryImpl) Search(ctx context.Context, query string, page, limit int) ([]entity.Team, int64, error) {
	var teams []entity.Team
	var total int64

	offset := (page - 1) * limit
	searchQuery := "%" + query + "%"

	err := r.db.WithContext(ctx).
		Model(&entity.Team{}).
		Where("name ILIKE ? OR city ILIKE ?", searchQuery, searchQuery).
		Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = r.db.WithContext(ctx).
		Where("name ILIKE ? OR city ILIKE ?", searchQuery, searchQuery).
		Offset(offset).
		Limit(limit).
		Order("created_at DESC").
		Find(&teams).Error
	if err != nil {
		return nil, 0, err
	}

	return teams, total, nil
}

func (r *teamRepositoryImpl) Exists(ctx context.Context, id uuid.UUID) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&entity.Team{}).
		Where("id = ?", id).
		Count(&count).Error
	return count > 0, err
}
