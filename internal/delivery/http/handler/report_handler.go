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

// ReportHandler handles report related requests
type ReportHandler struct {
	reportUseCase usecase.ReportUseCase
}

// NewReportHandler creates a new instance of ReportHandler
func NewReportHandler(reportUseCase usecase.ReportUseCase) *ReportHandler {
	return &ReportHandler{reportUseCase: reportUseCase}
}

// GetMatchReport handles getting a single match report
// @Summary Get Match Report
// @Description Get detailed report for a specific match
// @Tags Reports
// @Accept json
// @Produce json
// @Param id path string true "Match ID"
// @Success 200 {object} response.Response{data=dto.MatchReportResponse}
// @Failure 400 {object} response.Response
// @Failure 404 {object} response.Response
// @Router /api/v1/reports/matches/{id} [get]
func (h *ReportHandler) GetMatchReport(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid match ID", nil)
		return
	}

	report, err := h.reportUseCase.GetMatchReport(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, usecase.ErrMatchNotFound) {
			response.Error(c, http.StatusNotFound, "Match not found", nil)
			return
		}
		response.Error(c, http.StatusInternalServerError, "Failed to get match report", err.Error())
		return
	}

	response.Success(c, http.StatusOK, "Match report retrieved successfully", dto.ToMatchReportResponse(report))
}

// GetAllMatchReports handles getting all match reports
// @Summary Get All Match Reports
// @Description Get reports for all completed matches with pagination
// @Tags Reports
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Success 200 {object} response.Response{data=[]dto.MatchReportResponse}
// @Router /api/v1/reports/matches [get]
func (h *ReportHandler) GetAllMatchReports(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	reports, total, err := h.reportUseCase.GetAllMatchReports(c.Request.Context(), page, limit)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to get match reports", err.Error())
		return
	}

	response.SuccessWithMeta(c, http.StatusOK, "Match reports retrieved successfully", dto.ToMatchReportResponseList(reports), response.NewMeta(page, limit, total))
}

// GetTopScorers handles getting top scorers
// @Summary Get Top Scorers
// @Description Get top goal scorers
// @Tags Reports
// @Accept json
// @Produce json
// @Param limit query int false "Number of top scorers to return" default(10)
// @Success 200 {object} response.Response{data=[]dto.TopScorerResponse}
// @Router /api/v1/reports/top-scorers [get]
func (h *ReportHandler) GetTopScorers(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	if limit < 1 || limit > 100 {
		limit = 10
	}

	scorers, err := h.reportUseCase.GetTopScorers(c.Request.Context(), limit)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to get top scorers", err.Error())
		return
	}

	response.Success(c, http.StatusOK, "Top scorers retrieved successfully", dto.ToTopScorerResponseList(scorers))
}
