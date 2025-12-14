package entity

import (
	"time"

	"github.com/google/uuid"
)

// MatchStatus represents the status of a match
type MatchStatus string

const (
	MatchStatusScheduled MatchStatus = "scheduled"
	MatchStatusOngoing   MatchStatus = "ongoing"
	MatchStatusCompleted MatchStatus = "completed"
	MatchStatusCancelled MatchStatus = "cancelled"
)

// Match represents a football match between two teams
type Match struct {
	BaseEntity
	MatchDate    time.Time   `gorm:"not null;index" json:"match_date"`
	MatchTime    string      `gorm:"not null;size:10" json:"match_time"` // Format: HH:MM
	HomeTeamID   uuid.UUID   `gorm:"type:uuid;not null;index" json:"home_team_id"`
	AwayTeamID   uuid.UUID   `gorm:"type:uuid;not null;index" json:"away_team_id"`
	HomeScore    *int        `gorm:"default:null" json:"home_score"`
	AwayScore    *int        `gorm:"default:null" json:"away_score"`
	Status       MatchStatus `gorm:"type:varchar(20);default:'scheduled'" json:"status"`
	HomeTeam     *Team       `gorm:"foreignKey:HomeTeamID" json:"home_team,omitempty"`
	AwayTeam     *Team       `gorm:"foreignKey:AwayTeamID" json:"away_team,omitempty"`
	Goals        []Goal      `gorm:"foreignKey:MatchID" json:"goals,omitempty"`
}

// TableName returns the table name for Match entity
func (Match) TableName() string {
	return "matches"
}

// MatchResult represents the result of a match
type MatchResult string

const (
	ResultHomeWin MatchResult = "home_win"
	ResultAwayWin MatchResult = "away_win"
	ResultDraw    MatchResult = "draw"
)

// GetResult returns the result of the match
func (m *Match) GetResult() MatchResult {
	if m.HomeScore == nil || m.AwayScore == nil {
		return ""
	}
	if *m.HomeScore > *m.AwayScore {
		return ResultHomeWin
	}
	if *m.AwayScore > *m.HomeScore {
		return ResultAwayWin
	}
	return ResultDraw
}

// GetResultDisplay returns a human-readable result string
func (m *Match) GetResultDisplay() string {
	result := m.GetResult()
	switch result {
	case ResultHomeWin:
		return "Home Team Win"
	case ResultAwayWin:
		return "Away Team Win"
	case ResultDraw:
		return "Draw"
	default:
		return "Not Played"
	}
}
