package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/zenkriztao/ayo-football-backend/internal/delivery/http/dto"
	"github.com/zenkriztao/ayo-football-backend/internal/domain/usecase"
	"github.com/zenkriztao/ayo-football-backend/pkg/response"
)

// PlayerHandler handles player related requests
type PlayerHandler struct {
	playerUseCase usecase.PlayerUseCase
}

// NewPlayerHandler creates a new instance of PlayerHandler
func NewPlayerHandler(playerUseCase usecase.PlayerUseCase) *PlayerHandler {
	return &PlayerHandler{playerUseCase: playerUseCase}
}

// Create handles player creation
// @Summary Create Player
// @Description Create a new player
// @Tags Players
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.CreatePlayerRequest true "Player details"
// @Success 201 {object} response.Response{data=dto.PlayerResponse}
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 409 {object} response.Response
// @Router /api/v1/players [post]
func (h *PlayerHandler) Create(c *gin.Context) {
	var req dto.CreatePlayerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid request body", err.Error())
		return
	}

	player, err := req.ToPlayerEntity()
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid team ID format", nil)
		return
	}

	if err := h.playerUseCase.Create(c.Request.Context(), player); err != nil {
		switch {
		case errors.Is(err, usecase.ErrTeamNotFound):
			response.Error(c, http.StatusNotFound, "Team not found", nil)
		case errors.Is(err, usecase.ErrJerseyNumberTaken):
			response.Error(c, http.StatusConflict, "Jersey number is already taken by another player in this team", nil)
		case errors.Is(err, usecase.ErrInvalidPosition):
			response.Error(c, http.StatusBadRequest, "Invalid player position", nil)
		case errors.Is(err, usecase.ErrInvalidJerseyNumber):
			response.Error(c, http.StatusBadRequest, "Jersey number must be between 1 and 99", nil)
		default:
			response.Error(c, http.StatusInternalServerError, "Failed to create player", err.Error())
		}
		return
	}

	response.Success(c, http.StatusCreated, "Player created successfully", dto.ToPlayerResponse(player))
}

// GetByID handles getting a player by ID
// @Summary Get Player
// @Description Get a player by ID
// @Tags Players
// @Accept json
// @Produce json
// @Param id path string true "Player ID"
// @Success 200 {object} response.Response{data=dto.PlayerResponse}
// @Failure 400 {object} response.Response
// @Failure 404 {object} response.Response
// @Router /api/v1/players/{id} [get]
func (h *PlayerHandler) GetByID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid player ID", nil)
		return
	}

	player, err := h.playerUseCase.GetByIDWithTeam(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, usecase.ErrPlayerNotFound) {
			response.Error(c, http.StatusNotFound, "Player not found", nil)
			return
		}
		response.Error(c, http.StatusInternalServerError, "Failed to get player", err.Error())
		return
	}

	response.Success(c, http.StatusOK, "Player retrieved successfully", dto.ToPlayerResponse(player))
}

// Update handles updating a player
// @Summary Update Player
// @Description Update an existing player
// @Tags Players
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Player ID"
// @Param request body dto.UpdatePlayerRequest true "Player details"
// @Success 200 {object} response.Response{data=dto.PlayerResponse}
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 404 {object} response.Response
// @Failure 409 {object} response.Response
// @Router /api/v1/players/{id} [put]
func (h *PlayerHandler) Update(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid player ID", nil)
		return
	}

	var req dto.UpdatePlayerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid request body", err.Error())
		return
	}

	player, err := h.playerUseCase.GetByID(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, usecase.ErrPlayerNotFound) {
			response.Error(c, http.StatusNotFound, "Player not found", nil)
			return
		}
		response.Error(c, http.StatusInternalServerError, "Failed to get player", err.Error())
		return
	}

	if err := req.UpdatePlayerEntity(player); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid request data", err.Error())
		return
	}

	if err := h.playerUseCase.Update(c.Request.Context(), player); err != nil {
		switch {
		case errors.Is(err, usecase.ErrTeamNotFound):
			response.Error(c, http.StatusNotFound, "Team not found", nil)
		case errors.Is(err, usecase.ErrJerseyNumberTaken):
			response.Error(c, http.StatusConflict, "Jersey number is already taken by another player in this team", nil)
		case errors.Is(err, usecase.ErrInvalidPosition):
			response.Error(c, http.StatusBadRequest, "Invalid player position", nil)
		case errors.Is(err, usecase.ErrInvalidJerseyNumber):
			response.Error(c, http.StatusBadRequest, "Jersey number must be between 1 and 99", nil)
		default:
			response.Error(c, http.StatusInternalServerError, "Failed to update player", err.Error())
		}
		return
	}

	response.Success(c, http.StatusOK, "Player updated successfully", dto.ToPlayerResponse(player))
}

// Delete handles deleting a player
// @Summary Delete Player
// @Description Delete a player (soft delete)
// @Tags Players
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Player ID"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 404 {object} response.Response
// @Router /api/v1/players/{id} [delete]
func (h *PlayerHandler) Delete(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid player ID", nil)
		return
	}

	if err := h.playerUseCase.Delete(c.Request.Context(), id); err != nil {
		if errors.Is(err, usecase.ErrPlayerNotFound) {
			response.Error(c, http.StatusNotFound, "Player not found", nil)
			return
		}
		response.Error(c, http.StatusInternalServerError, "Failed to delete player", err.Error())
		return
	}

	response.Success(c, http.StatusOK, "Player deleted successfully", nil)
}

// GetAll handles getting all players with pagination
// @Summary Get All Players
// @Description Get all players with pagination
// @Tags Players
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Param search query string false "Search query"
// @Param team_id query string false "Filter by team ID"
// @Success 200 {object} response.Response{data=[]dto.PlayerResponse}
// @Router /api/v1/players [get]
func (h *PlayerHandler) GetAll(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	search := c.Query("search")
	teamIDStr := c.Query("team_id")

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	var players interface{}
	var total int64
	var err error

	if teamIDStr != "" {
		teamID, parseErr := uuid.Parse(teamIDStr)
		if parseErr != nil {
			response.Error(c, http.StatusBadRequest, "Invalid team ID", nil)
			return
		}
		p, totalCount, getErr := h.playerUseCase.GetByTeamID(c.Request.Context(), teamID, page, limit)
		players = dto.ToPlayerResponseList(p)
		total = totalCount
		err = getErr
		if errors.Is(err, usecase.ErrTeamNotFound) {
			response.Error(c, http.StatusNotFound, "Team not found", nil)
			return
		}
	} else if search != "" {
		p, totalCount, searchErr := h.playerUseCase.Search(c.Request.Context(), search, page, limit)
		players = dto.ToPlayerResponseList(p)
		total = totalCount
		err = searchErr
	} else {
		p, totalCount, getAllErr := h.playerUseCase.GetAll(c.Request.Context(), page, limit)
		players = dto.ToPlayerResponseList(p)
		total = totalCount
		err = getAllErr
	}

	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to get players", err.Error())
		return
	}

	response.SuccessWithMeta(c, http.StatusOK, "Players retrieved successfully", players, response.NewMeta(page, limit, total))
}
