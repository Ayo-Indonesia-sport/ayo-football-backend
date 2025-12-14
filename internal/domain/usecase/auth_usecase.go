package usecase

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/zenkriztao/ayo-football-backend/internal/domain/entity"
	"github.com/zenkriztao/ayo-football-backend/internal/domain/repository"
	"github.com/zenkriztao/ayo-football-backend/internal/infrastructure/security"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var (
	ErrInvalidCredentials = errors.New("invalid email or password")
	ErrUserAlreadyExists  = errors.New("user with this email already exists")
	ErrUserNotFound       = errors.New("user not found")
)

// AuthUseCase defines the interface for authentication operations
type AuthUseCase interface {
	Login(ctx context.Context, email, password string) (string, *entity.User, error)
	Register(ctx context.Context, name, email, password string, role entity.UserRole) (*entity.User, error)
	GetUserByID(ctx context.Context, id uuid.UUID) (*entity.User, error)
	CreateDefaultAdmin(ctx context.Context, email, password string) error
}

type authUseCaseImpl struct {
	userRepo   repository.UserRepository
	jwtService security.JWTService
}

// NewAuthUseCase creates a new instance of AuthUseCase
func NewAuthUseCase(userRepo repository.UserRepository, jwtService security.JWTService) AuthUseCase {
	return &authUseCaseImpl{
		userRepo:   userRepo,
		jwtService: jwtService,
	}
}

func (uc *authUseCaseImpl) Login(ctx context.Context, email, password string) (string, *entity.User, error) {
	user, err := uc.userRepo.FindByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", nil, ErrInvalidCredentials
		}
		return "", nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", nil, ErrInvalidCredentials
	}

	token, err := uc.jwtService.GenerateToken(user.ID, user.Email, string(user.Role))
	if err != nil {
		return "", nil, err
	}

	return token, user, nil
}

func (uc *authUseCaseImpl) Register(ctx context.Context, name, email, password string, role entity.UserRole) (*entity.User, error) {
	// Check if user already exists
	existingUser, err := uc.userRepo.FindByEmail(ctx, email)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	if existingUser != nil {
		return nil, ErrUserAlreadyExists
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &entity.User{
		Name:     name,
		Email:    email,
		Password: string(hashedPassword),
		Role:     role,
	}

	if err := uc.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

func (uc *authUseCaseImpl) GetUserByID(ctx context.Context, id uuid.UUID) (*entity.User, error) {
	user, err := uc.userRepo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	return user, nil
}

func (uc *authUseCaseImpl) CreateDefaultAdmin(ctx context.Context, email, password string) error {
	// Check if admin already exists
	_, err := uc.userRepo.FindByEmail(ctx, email)
	if err == nil {
		// Admin already exists
		return nil
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	// Create default admin
	_, err = uc.Register(ctx, "Admin", email, password, entity.RoleAdmin)
	return err
}
