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
	ErrPlayerNotFound       = errors.New("player not found")
	ErrJerseyNumberTaken    = errors.New("jersey number is already taken by another player in this team")
	ErrInvalidPosition      = errors.New("invalid player position")
	ErrInvalidJerseyNumber  = errors.New("jersey number must be between 1 and 99")
)

// PlayerUseCase defines the interface for player operations
type PlayerUseCase interface {
	Create(ctx context.Context, player *entity.Player) error
	GetByID(ctx context.Context, id uuid.UUID) (*entity.Player, error)
	GetByIDWithTeam(ctx context.Context, id uuid.UUID) (*entity.Player, error)
	Update(ctx context.Context, player *entity.Player) error
	Delete(ctx context.Context, id uuid.UUID) error
	GetAll(ctx context.Context, page, limit int) ([]entity.Player, int64, error)
	GetByTeamID(ctx context.Context, teamID uuid.UUID, page, limit int) ([]entity.Player, int64, error)
	Search(ctx context.Context, query string, page, limit int) ([]entity.Player, int64, error)
}

type playerUseCaseImpl struct {
	playerRepo repository.PlayerRepository
	teamRepo   repository.TeamRepository
}

// NewPlayerUseCase creates a new instance of PlayerUseCase
func NewPlayerUseCase(playerRepo repository.PlayerRepository, teamRepo repository.TeamRepository) PlayerUseCase {
	return &playerUseCaseImpl{
		playerRepo: playerRepo,
		teamRepo:   teamRepo,
	}
}

func (uc *playerUseCaseImpl) Create(ctx context.Context, player *entity.Player) error {
	// Validate team exists
	exists, err := uc.teamRepo.Exists(ctx, player.TeamID)
	if err != nil {
		return err
	}
	if !exists {
		return ErrTeamNotFound
	}

	// Validate position
	if !entity.IsValidPosition(player.Position) {
		return ErrInvalidPosition
	}

	// Validate jersey number
	if player.JerseyNumber < 1 || player.JerseyNumber > 99 {
		return ErrInvalidJerseyNumber
	}

	// Check if jersey number is taken
	taken, err := uc.playerRepo.IsJerseyNumberTaken(ctx, player.TeamID, player.JerseyNumber, nil)
	if err != nil {
		return err
	}
	if taken {
		return ErrJerseyNumberTaken
	}

	return uc.playerRepo.Create(ctx, player)
}

func (uc *playerUseCaseImpl) GetByID(ctx context.Context, id uuid.UUID) (*entity.Player, error) {
	player, err := uc.playerRepo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrPlayerNotFound
		}
		return nil, err
	}
	return player, nil
}

func (uc *playerUseCaseImpl) GetByIDWithTeam(ctx context.Context, id uuid.UUID) (*entity.Player, error) {
	player, err := uc.playerRepo.FindByIDWithTeam(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrPlayerNotFound
		}
		return nil, err
	}
	return player, nil
}

func (uc *playerUseCaseImpl) Update(ctx context.Context, player *entity.Player) error {
	// Check player exists
	exists, err := uc.playerRepo.Exists(ctx, player.ID)
	if err != nil {
		return err
	}
	if !exists {
		return ErrPlayerNotFound
	}

	// Validate team exists
	exists, err = uc.teamRepo.Exists(ctx, player.TeamID)
	if err != nil {
		return err
	}
	if !exists {
		return ErrTeamNotFound
	}

	// Validate position
	if !entity.IsValidPosition(player.Position) {
		return ErrInvalidPosition
	}

	// Validate jersey number
	if player.JerseyNumber < 1 || player.JerseyNumber > 99 {
		return ErrInvalidJerseyNumber
	}

	// Check if jersey number is taken (exclude current player)
	taken, err := uc.playerRepo.IsJerseyNumberTaken(ctx, player.TeamID, player.JerseyNumber, &player.ID)
	if err != nil {
		return err
	}
	if taken {
		return ErrJerseyNumberTaken
	}

	return uc.playerRepo.Update(ctx, player)
}

func (uc *playerUseCaseImpl) Delete(ctx context.Context, id uuid.UUID) error {
	exists, err := uc.playerRepo.Exists(ctx, id)
	if err != nil {
		return err
	}
	if !exists {
		return ErrPlayerNotFound
	}
	return uc.playerRepo.Delete(ctx, id)
}

func (uc *playerUseCaseImpl) GetAll(ctx context.Context, page, limit int) ([]entity.Player, int64, error) {
	return uc.playerRepo.FindAll(ctx, page, limit)
}

func (uc *playerUseCaseImpl) GetByTeamID(ctx context.Context, teamID uuid.UUID, page, limit int) ([]entity.Player, int64, error) {
	// Validate team exists
	exists, err := uc.teamRepo.Exists(ctx, teamID)
	if err != nil {
		return nil, 0, err
	}
	if !exists {
		return nil, 0, ErrTeamNotFound
	}

	return uc.playerRepo.FindByTeamID(ctx, teamID, page, limit)
}

func (uc *playerUseCaseImpl) Search(ctx context.Context, query string, page, limit int) ([]entity.Player, int64, error) {
	return uc.playerRepo.Search(ctx, query, page, limit)
}
