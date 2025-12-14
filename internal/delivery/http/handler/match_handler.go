package handler

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/zenkriztao/ayo-football-backend/internal/delivery/http/dto"
	"github.com/zenkriztao/ayo-football-backend/internal/domain/entity"
	"github.com/zenkriztao/ayo-football-backend/internal/domain/usecase"
	"github.com/zenkriztao/ayo-football-backend/pkg/response"
)

// MatchHandler handles match related requests
type MatchHandler struct {
	matchUseCase usecase.MatchUseCase
}

// NewMatchHandler creates a new instance of MatchHandler
func NewMatchHandler(matchUseCase usecase.MatchUseCase) *MatchHandler {
	return &MatchHandler{matchUseCase: matchUseCase}
}

// Create handles match creation
// @Summary Create Match
// @Description Create a new match schedule
// @Tags Matches
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.CreateMatchRequest true "Match details"
// @Success 201 {object} response.Response{data=dto.MatchResponse}
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Router /api/v1/matches [post]
func (h *MatchHandler) Create(c *gin.Context) {
	var req dto.CreateMatchRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid request body", err.Error())
		return
	}

	match, err := req.ToMatchEntity()
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid request data", err.Error())
		return
	}

	if err := h.matchUseCase.Create(c.Request.Context(), match); err != nil {
		if errors.Is(err, usecase.ErrSameTeamMatch) {
			response.Error(c, http.StatusBadRequest, "Home team and away team cannot be the same", nil)
			return
		}
		response.Error(c, http.StatusInternalServerError, "Failed to create match", err.Error())
		return
	}

	response.Success(c, http.StatusCreated, "Match created successfully", dto.ToMatchResponse(match))
}

// GetByID handles getting a match by ID
// @Summary Get Match
// @Description Get a match by ID
// @Tags Matches
// @Accept json
// @Produce json
// @Param id path string true "Match ID"
// @Success 200 {object} response.Response{data=dto.MatchResponse}
// @Failure 400 {object} response.Response
// @Failure 404 {object} response.Response
// @Router /api/v1/matches/{id} [get]
func (h *MatchHandler) GetByID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid match ID", nil)
		return
	}

	match, err := h.matchUseCase.GetByIDWithDetails(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, usecase.ErrMatchNotFound) {
			response.Error(c, http.StatusNotFound, "Match not found", nil)
			return
		}
		response.Error(c, http.StatusInternalServerError, "Failed to get match", err.Error())
		return
	}

	response.Success(c, http.StatusOK, "Match retrieved successfully", dto.ToMatchResponse(match))
}

// Update handles updating a match
// @Summary Update Match
// @Description Update an existing match
// @Tags Matches
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Match ID"
// @Param request body dto.UpdateMatchRequest true "Match details"
// @Success 200 {object} response.Response{data=dto.MatchResponse}
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 404 {object} response.Response
// @Router /api/v1/matches/{id} [put]
func (h *MatchHandler) Update(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid match ID", nil)
		return
	}

	var req dto.UpdateMatchRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid request body", err.Error())
		return
	}

	match, err := h.matchUseCase.GetByID(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, usecase.ErrMatchNotFound) {
			response.Error(c, http.StatusNotFound, "Match not found", nil)
			return
		}
		response.Error(c, http.StatusInternalServerError, "Failed to get match", err.Error())
		return
	}

	if err := req.UpdateMatchEntity(match); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid request data", err.Error())
		return
	}

	if err := h.matchUseCase.Update(c.Request.Context(), match); err != nil {
		if errors.Is(err, usecase.ErrSameTeamMatch) {
			response.Error(c, http.StatusBadRequest, "Home team and away team cannot be the same", nil)
			return
		}
		response.Error(c, http.StatusInternalServerError, "Failed to update match", err.Error())
		return
	}

	response.Success(c, http.StatusOK, "Match updated successfully", dto.ToMatchResponse(match))
}

// Delete handles deleting a match
// @Summary Delete Match
// @Description Delete a match (soft delete)
// @Tags Matches
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Match ID"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 404 {object} response.Response
// @Router /api/v1/matches/{id} [delete]
func (h *MatchHandler) Delete(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid match ID", nil)
		return
	}

	if err := h.matchUseCase.Delete(c.Request.Context(), id); err != nil {
		if errors.Is(err, usecase.ErrMatchNotFound) {
			response.Error(c, http.StatusNotFound, "Match not found", nil)
			return
		}
		response.Error(c, http.StatusInternalServerError, "Failed to delete match", err.Error())
		return
	}

	response.Success(c, http.StatusOK, "Match deleted successfully", nil)
}

// GetAll handles getting all matches with pagination
// @Summary Get All Matches
// @Description Get all matches with pagination
// @Tags Matches
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Param team_id query string false "Filter by team ID"
// @Param status query string false "Filter by status (scheduled, ongoing, completed, cancelled)"
// @Param start_date query string false "Start date filter (YYYY-MM-DD)"
// @Param end_date query string false "End date filter (YYYY-MM-DD)"
// @Success 200 {object} response.Response{data=[]dto.MatchResponse}
// @Router /api/v1/matches [get]
func (h *MatchHandler) GetAll(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	teamIDStr := c.Query("team_id")
	status := c.Query("status")
	startDateStr := c.Query("start_date")
	endDateStr := c.Query("end_date")

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	var matches interface{}
	var total int64
	var err error

	if teamIDStr != "" {
		teamID, parseErr := uuid.Parse(teamIDStr)
		if parseErr != nil {
			response.Error(c, http.StatusBadRequest, "Invalid team ID", nil)
			return
		}
		m, totalCount, getErr := h.matchUseCase.GetByTeamID(c.Request.Context(), teamID, page, limit)
		matches = dto.ToMatchResponseList(m)
		total = totalCount
		err = getErr
		if errors.Is(err, usecase.ErrTeamNotFound) {
			response.Error(c, http.StatusNotFound, "Team not found", nil)
			return
		}
	} else if status != "" {
		m, totalCount, getErr := h.matchUseCase.GetByStatus(c.Request.Context(), entity.MatchStatus(status), page, limit)
		matches = dto.ToMatchResponseList(m)
		total = totalCount
		err = getErr
	} else if startDateStr != "" && endDateStr != "" {
		startDate, parseErr := time.Parse("2006-01-02", startDateStr)
		if parseErr != nil {
			response.Error(c, http.StatusBadRequest, "Invalid start date format", nil)
			return
		}
		endDate, parseErr := time.Parse("2006-01-02", endDateStr)
		if parseErr != nil {
			response.Error(c, http.StatusBadRequest, "Invalid end date format", nil)
			return
		}
		m, totalCount, getErr := h.matchUseCase.GetByDateRange(c.Request.Context(), startDate, endDate, page, limit)
		matches = dto.ToMatchResponseList(m)
		total = totalCount
		err = getErr
	} else {
		m, totalCount, getAllErr := h.matchUseCase.GetAll(c.Request.Context(), page, limit)
		matches = dto.ToMatchResponseList(m)
		total = totalCount
		err = getAllErr
	}

	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to get matches", err.Error())
		return
	}

	response.SuccessWithMeta(c, http.StatusOK, "Matches retrieved successfully", matches, response.NewMeta(page, limit, total))
}

// RecordResult handles recording a match result
// @Summary Record Match Result
// @Description Record the result of a completed match
// @Tags Matches
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Match ID"
// @Param request body dto.RecordMatchResultRequest true "Match result"
// @Success 200 {object} response.Response{data=dto.MatchResponse}
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 404 {object} response.Response
// @Router /api/v1/matches/{id}/result [post]
func (h *MatchHandler) RecordResult(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid match ID", nil)
		return
	}

	var req dto.RecordMatchResultRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid request body", err.Error())
		return
	}

	// Convert goals
	goals := make([]usecase.GoalInput, len(req.Goals))
	for i, g := range req.Goals {
		playerID, parseErr := uuid.Parse(g.PlayerID)
		if parseErr != nil {
			response.Error(c, http.StatusBadRequest, "Invalid player ID in goals", nil)
			return
		}
		teamID, parseErr := uuid.Parse(g.TeamID)
		if parseErr != nil {
			response.Error(c, http.StatusBadRequest, "Invalid team ID in goals", nil)
			return
		}
		goals[i] = usecase.GoalInput{
			PlayerID:  playerID,
			TeamID:    teamID,
			Minute:    g.Minute,
			IsOwnGoal: g.IsOwnGoal,
		}
	}

	input := usecase.MatchResultInput{
		HomeScore: req.HomeScore,
		AwayScore: req.AwayScore,
		Goals:     goals,
	}

	match, err := h.matchUseCase.RecordResult(c.Request.Context(), id, input)
	if err != nil {
		if errors.Is(err, usecase.ErrMatchNotFound) {
			response.Error(c, http.StatusNotFound, "Match not found", nil)
			return
		}
		if errors.Is(err, usecase.ErrPlayerNotFound) {
			response.Error(c, http.StatusNotFound, "One or more players not found", nil)
			return
		}
		response.Error(c, http.StatusInternalServerError, "Failed to record match result", err.Error())
		return
	}

	response.Success(c, http.StatusOK, "Match result recorded successfully", dto.ToMatchResponse(match))
}
