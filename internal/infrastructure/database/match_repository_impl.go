package database

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/zenkriztao/ayo-football-backend/internal/domain/entity"
	"github.com/zenkriztao/ayo-football-backend/internal/domain/repository"
	"gorm.io/gorm"
)

type matchRepositoryImpl struct {
	db *gorm.DB
}

// NewMatchRepository creates a new instance of MatchRepository
func NewMatchRepository(db *gorm.DB) repository.MatchRepository {
	return &matchRepositoryImpl{db: db}
}

func (r *matchRepositoryImpl) Create(ctx context.Context, match *entity.Match) error {
	return r.db.WithContext(ctx).Create(match).Error
}

func (r *matchRepositoryImpl) FindByID(ctx context.Context, id uuid.UUID) (*entity.Match, error) {
	var match entity.Match
	err := r.db.WithContext(ctx).First(&match, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &match, nil
}

func (r *matchRepositoryImpl) FindByIDWithDetails(ctx context.Context, id uuid.UUID) (*entity.Match, error) {
	var match entity.Match
	err := r.db.WithContext(ctx).
		Preload("HomeTeam").
		Preload("AwayTeam").
		Preload("Goals").
		Preload("Goals.Player").
		Preload("Goals.Team").
		First(&match, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &match, nil
}

func (r *matchRepositoryImpl) Update(ctx context.Context, match *entity.Match) error {
	return r.db.WithContext(ctx).Save(match).Error
}

func (r *matchRepositoryImpl) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&entity.Match{}, "id = ?", id).Error
}

func (r *matchRepositoryImpl) FindAll(ctx context.Context, page, limit int) ([]entity.Match, int64, error) {
	var matches []entity.Match
	var total int64

	offset := (page - 1) * limit

	err := r.db.WithContext(ctx).Model(&entity.Match{}).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = r.db.WithContext(ctx).
		Preload("HomeTeam").
		Preload("AwayTeam").
		Offset(offset).
		Limit(limit).
		Order("match_date DESC, match_time DESC").
		Find(&matches).Error
	if err != nil {
		return nil, 0, err
	}

	return matches, total, nil
}

func (r *matchRepositoryImpl) FindByDateRange(ctx context.Context, startDate, endDate time.Time, page, limit int) ([]entity.Match, int64, error) {
	var matches []entity.Match
	var total int64

	offset := (page - 1) * limit

	err := r.db.WithContext(ctx).
		Model(&entity.Match{}).
		Where("match_date BETWEEN ? AND ?", startDate, endDate).
		Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = r.db.WithContext(ctx).
		Preload("HomeTeam").
		Preload("AwayTeam").
		Where("match_date BETWEEN ? AND ?", startDate, endDate).
		Offset(offset).
		Limit(limit).
		Order("match_date ASC, match_time ASC").
		Find(&matches).Error
	if err != nil {
		return nil, 0, err
	}

	return matches, total, nil
}

func (r *matchRepositoryImpl) FindByTeamID(ctx context.Context, teamID uuid.UUID, page, limit int) ([]entity.Match, int64, error) {
	var matches []entity.Match
	var total int64

	offset := (page - 1) * limit

	err := r.db.WithContext(ctx).
		Model(&entity.Match{}).
		Where("home_team_id = ? OR away_team_id = ?", teamID, teamID).
		Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = r.db.WithContext(ctx).
		Preload("HomeTeam").
		Preload("AwayTeam").
		Where("home_team_id = ? OR away_team_id = ?", teamID, teamID).
		Offset(offset).
		Limit(limit).
		Order("match_date DESC, match_time DESC").
		Find(&matches).Error
	if err != nil {
		return nil, 0, err
	}

	return matches, total, nil
}

func (r *matchRepositoryImpl) FindByStatus(ctx context.Context, status entity.MatchStatus, page, limit int) ([]entity.Match, int64, error) {
	var matches []entity.Match
	var total int64

	offset := (page - 1) * limit

	err := r.db.WithContext(ctx).
		Model(&entity.Match{}).
		Where("status = ?", status).
		Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = r.db.WithContext(ctx).
		Preload("HomeTeam").
		Preload("AwayTeam").
		Where("status = ?", status).
		Offset(offset).
		Limit(limit).
		Order("match_date ASC, match_time ASC").
		Find(&matches).Error
	if err != nil {
		return nil, 0, err
	}

	return matches, total, nil
}

func (r *matchRepositoryImpl) Exists(ctx context.Context, id uuid.UUID) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&entity.Match{}).
		Where("id = ?", id).
		Count(&count).Error
	return count > 0, err
}

func (r *matchRepositoryImpl) GetTeamWinCount(ctx context.Context, teamID uuid.UUID, isHome bool) (int64, error) {
	var count int64
	var query *gorm.DB

	if isHome {
		query = r.db.WithContext(ctx).
			Model(&entity.Match{}).
			Where("home_team_id = ? AND status = ? AND home_score > away_score", teamID, entity.MatchStatusCompleted)
	} else {
		query = r.db.WithContext(ctx).
			Model(&entity.Match{}).
			Where("away_team_id = ? AND status = ? AND away_score > home_score", teamID, entity.MatchStatusCompleted)
	}

	err := query.Count(&count).Error
	return count, err
}

func (r *matchRepositoryImpl) GetCompletedMatches(ctx context.Context, page, limit int) ([]entity.Match, int64, error) {
	var matches []entity.Match
	var total int64

	offset := (page - 1) * limit

	err := r.db.WithContext(ctx).
		Model(&entity.Match{}).
		Where("status = ?", entity.MatchStatusCompleted).
		Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = r.db.WithContext(ctx).
		Preload("HomeTeam").
		Preload("AwayTeam").
		Preload("Goals").
		Preload("Goals.Player").
		Where("status = ?", entity.MatchStatusCompleted).
		Offset(offset).
		Limit(limit).
		Order("match_date DESC, match_time DESC").
		Find(&matches).Error
	if err != nil {
		return nil, 0, err
	}

	return matches, total, nil
}
