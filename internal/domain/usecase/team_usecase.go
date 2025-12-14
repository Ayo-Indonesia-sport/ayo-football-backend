package usecase

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/zenkriztao/ayo-football-backend/internal/domain/entity"
	"github.com/zenkriztao/ayo-football-backend/internal/domain/repository"
	"gorm.io/gorm"
)

var (
	ErrTeamNotFound = errors.New("team not found")
)

// TeamUseCase defines the interface for team operations
type TeamUseCase interface {
	Create(ctx context.Context, team *entity.Team) error
	GetByID(ctx context.Context, id uuid.UUID) (*entity.Team, error)
	GetByIDWithPlayers(ctx context.Context, id uuid.UUID) (*entity.Team, error)
	Update(ctx context.Context, team *entity.Team) error
	Delete(ctx context.Context, id uuid.UUID) error
	GetAll(ctx context.Context, page, limit int) ([]entity.Team, int64, error)
	Search(ctx context.Context, query string, page, limit int) ([]entity.Team, int64, error)
}

type teamUseCaseImpl struct {
	teamRepo repository.TeamRepository
}

// NewTeamUseCase creates a new instance of TeamUseCase
func NewTeamUseCase(teamRepo repository.TeamRepository) TeamUseCase {
	return &teamUseCaseImpl{teamRepo: teamRepo}
}

func (uc *teamUseCaseImpl) Create(ctx context.Context, team *entity.Team) error {
	return uc.teamRepo.Create(ctx, team)
}

func (uc *teamUseCaseImpl) GetByID(ctx context.Context, id uuid.UUID) (*entity.Team, error) {
	team, err := uc.teamRepo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrTeamNotFound
		}
		return nil, err
	}
	return team, nil
}

func (uc *teamUseCaseImpl) GetByIDWithPlayers(ctx context.Context, id uuid.UUID) (*entity.Team, error) {
	team, err := uc.teamRepo.FindByIDWithPlayers(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrTeamNotFound
		}
		return nil, err
	}
	return team, nil
}

func (uc *teamUseCaseImpl) Update(ctx context.Context, team *entity.Team) error {
	exists, err := uc.teamRepo.Exists(ctx, team.ID)
	if err != nil {
		return err
	}
	if !exists {
		return ErrTeamNotFound
	}
	return uc.teamRepo.Update(ctx, team)
}

func (uc *teamUseCaseImpl) Delete(ctx context.Context, id uuid.UUID) error {
	exists, err := uc.teamRepo.Exists(ctx, id)
	if err != nil {
		return err
	}
	if !exists {
		return ErrTeamNotFound
	}
	return uc.teamRepo.Delete(ctx, id)
}

func (uc *teamUseCaseImpl) GetAll(ctx context.Context, page, limit int) ([]entity.Team, int64, error) {
	return uc.teamRepo.FindAll(ctx, page, limit)
}

func (uc *teamUseCaseImpl) Search(ctx context.Context, query string, page, limit int) ([]entity.Team, int64, error) {
	return uc.teamRepo.Search(ctx, query, page, limit)
}
