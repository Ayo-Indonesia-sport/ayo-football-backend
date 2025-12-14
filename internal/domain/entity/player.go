package entity

import "github.com/google/uuid"

// PlayerPosition represents the position of a player
type PlayerPosition string

const (
	PositionForward    PlayerPosition = "forward"     // Penyerang
	PositionMidfielder PlayerPosition = "midfielder"  // Gelandang
	PositionDefender   PlayerPosition = "defender"    // Bertahan
	PositionGoalkeeper PlayerPosition = "goalkeeper"  // Penjaga Gawang
)

// Player represents a football player
type Player struct {
	BaseEntity
	TeamID       uuid.UUID      `gorm:"type:uuid;not null;index" json:"team_id"`
	Name         string         `gorm:"not null;size:255" json:"name"`
	Height       float64        `gorm:"not null" json:"height"` // in cm
	Weight       float64        `gorm:"not null" json:"weight"` // in kg
	Position     PlayerPosition `gorm:"type:varchar(20);not null" json:"position"`
	JerseyNumber int            `gorm:"not null" json:"jersey_number"`
	Team         *Team          `gorm:"foreignKey:TeamID" json:"team,omitempty"`
}

// TableName returns the table name for Player entity
func (Player) TableName() string {
	return "players"
}

// ValidPositions returns all valid player positions
func ValidPositions() []PlayerPosition {
	return []PlayerPosition{
		PositionForward,
		PositionMidfielder,
		PositionDefender,
		PositionGoalkeeper,
	}
}

// IsValidPosition checks if a position is valid
func IsValidPosition(position PlayerPosition) bool {
	for _, p := range ValidPositions() {
		if p == position {
			return true
		}
	}
	return false
}
