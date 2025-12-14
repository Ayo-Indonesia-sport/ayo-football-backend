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

// TeamHandler handles team related requests
type TeamHandler struct {
	teamUseCase usecase.TeamUseCase
}

// NewTeamHandler creates a new instance of TeamHandler
func NewTeamHandler(teamUseCase usecase.TeamUseCase) *TeamHandler {
	return &TeamHandler{teamUseCase: teamUseCase}
}

// Create handles team creation
// @Summary Create Team
// @Description Create a new football team
// @Tags Teams
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.CreateTeamRequest true "Team details"
// @Success 201 {object} response.Response{data=dto.TeamResponse}
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Router /api/v1/teams [post]
func (h *TeamHandler) Create(c *gin.Context) {
	var req dto.CreateTeamRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid request body", err.Error())
		return
	}

	team := req.ToTeamEntity()
	if err := h.teamUseCase.Create(c.Request.Context(), team); err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to create team", err.Error())
		return
	}

	response.Success(c, http.StatusCreated, "Team created successfully", dto.ToTeamResponse(team))
}

// GetByID handles getting a team by ID
// @Summary Get Team
// @Description Get a team by ID
// @Tags Teams
// @Accept json
// @Produce json
// @Param id path string true "Team ID"
// @Success 200 {object} response.Response{data=dto.TeamResponse}
// @Failure 400 {object} response.Response
// @Failure 404 {object} response.Response
// @Router /api/v1/teams/{id} [get]
func (h *TeamHandler) GetByID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid team ID", nil)
		return
	}

	// Check if we should include players
	withPlayers := c.Query("with_players") == "true"

	var team interface{}
	if withPlayers {
		t, err := h.teamUseCase.GetByIDWithPlayers(c.Request.Context(), id)
		if err != nil {
			if errors.Is(err, usecase.ErrTeamNotFound) {
				response.Error(c, http.StatusNotFound, "Team not found", nil)
				return
			}
			response.Error(c, http.StatusInternalServerError, "Failed to get team", err.Error())
			return
		}
		team = dto.ToTeamResponse(t)
	} else {
		t, err := h.teamUseCase.GetByID(c.Request.Context(), id)
		if err != nil {
			if errors.Is(err, usecase.ErrTeamNotFound) {
				response.Error(c, http.StatusNotFound, "Team not found", nil)
				return
			}
			response.Error(c, http.StatusInternalServerError, "Failed to get team", err.Error())
			return
		}
		team = dto.ToTeamResponse(t)
	}

	response.Success(c, http.StatusOK, "Team retrieved successfully", team)
}

// Update handles updating a team
// @Summary Update Team
// @Description Update an existing team
// @Tags Teams
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Team ID"
// @Param request body dto.UpdateTeamRequest true "Team details"
// @Success 200 {object} response.Response{data=dto.TeamResponse}
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 404 {object} response.Response
// @Router /api/v1/teams/{id} [put]
func (h *TeamHandler) Update(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid team ID", nil)
		return
	}

	var req dto.UpdateTeamRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid request body", err.Error())
		return
	}

	team, err := h.teamUseCase.GetByID(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, usecase.ErrTeamNotFound) {
			response.Error(c, http.StatusNotFound, "Team not found", nil)
			return
		}
		response.Error(c, http.StatusInternalServerError, "Failed to get team", err.Error())
		return
	}

	req.UpdateTeamEntity(team)

	if err := h.teamUseCase.Update(c.Request.Context(), team); err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to update team", err.Error())
		return
	}

	response.Success(c, http.StatusOK, "Team updated successfully", dto.ToTeamResponse(team))
}

// Delete handles deleting a team
// @Summary Delete Team
// @Description Delete a team (soft delete)
// @Tags Teams
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Team ID"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 404 {object} response.Response
// @Router /api/v1/teams/{id} [delete]
func (h *TeamHandler) Delete(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid team ID", nil)
		return
	}

	if err := h.teamUseCase.Delete(c.Request.Context(), id); err != nil {
		if errors.Is(err, usecase.ErrTeamNotFound) {
			response.Error(c, http.StatusNotFound, "Team not found", nil)
			return
		}
		response.Error(c, http.StatusInternalServerError, "Failed to delete team", err.Error())
		return
	}

	response.Success(c, http.StatusOK, "Team deleted successfully", nil)
}

// GetAll handles getting all teams with pagination
// @Summary Get All Teams
// @Description Get all teams with pagination
// @Tags Teams
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Param search query string false "Search query"
// @Success 200 {object} response.Response{data=[]dto.TeamResponse}
// @Router /api/v1/teams [get]
func (h *TeamHandler) GetAll(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	search := c.Query("search")

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	var teams interface{}
	var total int64
	var err error

	if search != "" {
		t, totalCount, searchErr := h.teamUseCase.Search(c.Request.Context(), search, page, limit)
		teams = dto.ToTeamResponseList(t)
		total = totalCount
		err = searchErr
	} else {
		t, totalCount, getAllErr := h.teamUseCase.GetAll(c.Request.Context(), page, limit)
		teams = dto.ToTeamResponseList(t)
		total = totalCount
		err = getAllErr
	}

	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to get teams", err.Error())
		return
	}

	response.SuccessWithMeta(c, http.StatusOK, "Teams retrieved successfully", teams, response.NewMeta(page, limit, total))
}
