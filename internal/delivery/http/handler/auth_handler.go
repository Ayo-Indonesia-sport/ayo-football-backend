package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/zenkriztao/ayo-football-backend/internal/delivery/http/dto"
	"github.com/zenkriztao/ayo-football-backend/internal/delivery/http/middleware"
	"github.com/zenkriztao/ayo-football-backend/internal/domain/entity"
	"github.com/zenkriztao/ayo-football-backend/internal/domain/usecase"
	"github.com/zenkriztao/ayo-football-backend/pkg/response"
)

// AuthHandler handles authentication related requests
type AuthHandler struct {
	authUseCase usecase.AuthUseCase
}

// NewAuthHandler creates a new instance of AuthHandler
func NewAuthHandler(authUseCase usecase.AuthUseCase) *AuthHandler {
	return &AuthHandler{authUseCase: authUseCase}
}

// Login handles user login
// @Summary Login
// @Description Authenticate user and return JWT token
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body dto.LoginRequest true "Login credentials"
// @Success 200 {object} response.Response{data=dto.AuthResponse}
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Router /api/v1/auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid request body", err.Error())
		return
	}

	token, user, err := h.authUseCase.Login(c.Request.Context(), req.Email, req.Password)
	if err != nil {
		if errors.Is(err, usecase.ErrInvalidCredentials) {
			response.Error(c, http.StatusUnauthorized, "Invalid email or password", nil)
			return
		}
		response.Error(c, http.StatusInternalServerError, "Failed to login", err.Error())
		return
	}

	response.Success(c, http.StatusOK, "Login successful", dto.AuthResponse{
		Token: token,
		User:  dto.ToUserResponse(user),
	})
}

// Register handles user registration
// @Summary Register
// @Description Register a new user
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body dto.RegisterRequest true "Registration details"
// @Success 201 {object} response.Response{data=dto.UserResponse}
// @Failure 400 {object} response.Response
// @Failure 409 {object} response.Response
// @Router /api/v1/auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	var req dto.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid request body", err.Error())
		return
	}

	user, err := h.authUseCase.Register(c.Request.Context(), req.Name, req.Email, req.Password, entity.RoleUser)
	if err != nil {
		if errors.Is(err, usecase.ErrUserAlreadyExists) {
			response.Error(c, http.StatusConflict, "User with this email already exists", nil)
			return
		}
		response.Error(c, http.StatusInternalServerError, "Failed to register user", err.Error())
		return
	}

	response.Success(c, http.StatusCreated, "User registered successfully", dto.ToUserResponse(user))
}

// GetProfile handles getting current user profile
// @Summary Get Profile
// @Description Get current authenticated user profile
// @Tags Auth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.Response{data=dto.UserResponse}
// @Failure 401 {object} response.Response
// @Router /api/v1/auth/profile [get]
func (h *AuthHandler) GetProfile(c *gin.Context) {
	userID, exists := c.Get(middleware.UserIDKey)
	if !exists {
		response.Error(c, http.StatusUnauthorized, "User not authenticated", nil)
		return
	}

	user, err := h.authUseCase.GetUserByID(c.Request.Context(), userID.(uuid.UUID))
	if err != nil {
		if errors.Is(err, usecase.ErrUserNotFound) {
			response.Error(c, http.StatusNotFound, "User not found", nil)
			return
		}
		response.Error(c, http.StatusInternalServerError, "Failed to get user profile", err.Error())
		return
	}

	response.Success(c, http.StatusOK, "Profile retrieved successfully", dto.ToUserResponse(user))
}
