package usecase

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/zenkriztao/ayo-football-backend/internal/domain/entity"
	"github.com/zenkriztao/ayo-football-backend/internal/domain/repository"
	"gorm.io/gorm"
)

var (
	ErrMatchNotFound       = errors.New("match not found")
	ErrSameTeamMatch       = errors.New("home team and away team cannot be the same")
	ErrMatchAlreadyPlayed  = errors.New("match has already been played")
	ErrMatchNotCompleted   = errors.New("match has not been completed yet")
	ErrInvalidMatchStatus  = errors.New("invalid match status")
)

// MatchResultInput represents the input for recording a match result
type MatchResultInput struct {
	HomeScore int
	AwayScore int
	Goals     []GoalInput
}

// GoalInput represents a goal input
type GoalInput struct {
	PlayerID  uuid.UUID
	TeamID    uuid.UUID
	Minute    int
	IsOwnGoal bool
}

// MatchUseCase defines the interface for match operations
type MatchUseCase interface {
	Create(ctx context.Context, match *entity.Match) error
	GetByID(ctx context.Context, id uuid.UUID) (*entity.Match, error)
	GetByIDWithDetails(ctx context.Context, id uuid.UUID) (*entity.Match, error)
	Update(ctx context.Context, match *entity.Match) error
	Delete(ctx context.Context, id uuid.UUID) error
	GetAll(ctx context.Context, page, limit int) ([]entity.Match, int64, error)
	GetByDateRange(ctx context.Context, startDate, endDate time.Time, page, limit int) ([]entity.Match, int64, error)
	GetByTeamID(ctx context.Context, teamID uuid.UUID, page, limit int) ([]entity.Match, int64, error)
	GetByStatus(ctx context.Context, status entity.MatchStatus, page, limit int) ([]entity.Match, int64, error)
	RecordResult(ctx context.Context, matchID uuid.UUID, input MatchResultInput) (*entity.Match, error)
	GetCompletedMatches(ctx context.Context, page, limit int) ([]entity.Match, int64, error)
}

type matchUseCaseImpl struct {
	matchRepo  repository.MatchRepository
	teamRepo   repository.TeamRepository
	playerRepo repository.PlayerRepository
	goalRepo   repository.GoalRepository
}

// NewMatchUseCase creates a new instance of MatchUseCase
func NewMatchUseCase(
	matchRepo repository.MatchRepository,
	teamRepo repository.TeamRepository,
	playerRepo repository.PlayerRepository,
	goalRepo repository.GoalRepository,
) MatchUseCase {
	return &matchUseCaseImpl{
		matchRepo:  matchRepo,
		teamRepo:   teamRepo,
		playerRepo: playerRepo,
		goalRepo:   goalRepo,
	}
}

func (uc *matchUseCaseImpl) Create(ctx context.Context, match *entity.Match) error {
	// Validate teams exist
	homeExists, err := uc.teamRepo.Exists(ctx, match.HomeTeamID)
	if err != nil {
		return err
	}
	if !homeExists {
		return errors.New("home team not found")
	}

	awayExists, err := uc.teamRepo.Exists(ctx, match.AwayTeamID)
	if err != nil {
		return err
	}
	if !awayExists {
		return errors.New("away team not found")
	}

	// Validate teams are different
	if match.HomeTeamID == match.AwayTeamID {
		return ErrSameTeamMatch
	}

	// Set default status
	if match.Status == "" {
		match.Status = entity.MatchStatusScheduled
	}

	return uc.matchRepo.Create(ctx, match)
}

func (uc *matchUseCaseImpl) GetByID(ctx context.Context, id uuid.UUID) (*entity.Match, error) {
	match, err := uc.matchRepo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrMatchNotFound
		}
		return nil, err
	}
	return match, nil
}

func (uc *matchUseCaseImpl) GetByIDWithDetails(ctx context.Context, id uuid.UUID) (*entity.Match, error) {
	match, err := uc.matchRepo.FindByIDWithDetails(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrMatchNotFound
		}
		return nil, err
	}
	return match, nil
}

func (uc *matchUseCaseImpl) Update(ctx context.Context, match *entity.Match) error {
	// Check match exists
	exists, err := uc.matchRepo.Exists(ctx, match.ID)
	if err != nil {
		return err
	}
	if !exists {
		return ErrMatchNotFound
	}

	// Validate teams
	if match.HomeTeamID == match.AwayTeamID {
		return ErrSameTeamMatch
	}

	homeExists, err := uc.teamRepo.Exists(ctx, match.HomeTeamID)
	if err != nil {
		return err
	}
	if !homeExists {
		return errors.New("home team not found")
	}

	awayExists, err := uc.teamRepo.Exists(ctx, match.AwayTeamID)
	if err != nil {
		return err
	}
	if !awayExists {
		return errors.New("away team not found")
	}

	return uc.matchRepo.Update(ctx, match)
}

func (uc *matchUseCaseImpl) Delete(ctx context.Context, id uuid.UUID) error {
	exists, err := uc.matchRepo.Exists(ctx, id)
	if err != nil {
		return err
	}
	if !exists {
		return ErrMatchNotFound
	}
	return uc.matchRepo.Delete(ctx, id)
}

func (uc *matchUseCaseImpl) GetAll(ctx context.Context, page, limit int) ([]entity.Match, int64, error) {
	return uc.matchRepo.FindAll(ctx, page, limit)
}

func (uc *matchUseCaseImpl) GetByDateRange(ctx context.Context, startDate, endDate time.Time, page, limit int) ([]entity.Match, int64, error) {
	return uc.matchRepo.FindByDateRange(ctx, startDate, endDate, page, limit)
}

func (uc *matchUseCaseImpl) GetByTeamID(ctx context.Context, teamID uuid.UUID, page, limit int) ([]entity.Match, int64, error) {
	exists, err := uc.teamRepo.Exists(ctx, teamID)
	if err != nil {
		return nil, 0, err
	}
	if !exists {
		return nil, 0, ErrTeamNotFound
	}
	return uc.matchRepo.FindByTeamID(ctx, teamID, page, limit)
}

func (uc *matchUseCaseImpl) GetByStatus(ctx context.Context, status entity.MatchStatus, page, limit int) ([]entity.Match, int64, error) {
	return uc.matchRepo.FindByStatus(ctx, status, page, limit)
}

func (uc *matchUseCaseImpl) RecordResult(ctx context.Context, matchID uuid.UUID, input MatchResultInput) (*entity.Match, error) {
	// Get existing match
	match, err := uc.matchRepo.FindByIDWithDetails(ctx, matchID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrMatchNotFound
		}
		return nil, err
	}

	// Check if match is already completed
	if match.Status == entity.MatchStatusCompleted {
		// Delete existing goals and record new ones
		if err := uc.goalRepo.DeleteByMatchID(ctx, matchID); err != nil {
			return nil, err
		}
	}

	// Update match scores
	match.HomeScore = &input.HomeScore
	match.AwayScore = &input.AwayScore
	match.Status = entity.MatchStatusCompleted

	if err := uc.matchRepo.Update(ctx, match); err != nil {
		return nil, err
	}

	// Record goals
	goals := make([]entity.Goal, len(input.Goals))
	for i, g := range input.Goals {
		// Validate player exists
		exists, err := uc.playerRepo.Exists(ctx, g.PlayerID)
		if err != nil {
			return nil, err
		}
		if !exists {
			return nil, ErrPlayerNotFound
		}

		goals[i] = entity.Goal{
			MatchID:   matchID,
			PlayerID:  g.PlayerID,
			TeamID:    g.TeamID,
			Minute:    g.Minute,
			IsOwnGoal: g.IsOwnGoal,
		}
	}

	if len(goals) > 0 {
		if err := uc.goalRepo.CreateBatch(ctx, goals); err != nil {
			return nil, err
		}
	}

	// Fetch updated match with all details
	return uc.matchRepo.FindByIDWithDetails(ctx, matchID)
}

func (uc *matchUseCaseImpl) GetCompletedMatches(ctx context.Context, page, limit int) ([]entity.Match, int64, error) {
	return uc.matchRepo.GetCompletedMatches(ctx, page, limit)
}
