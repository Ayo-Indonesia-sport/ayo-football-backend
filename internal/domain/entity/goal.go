package entity

import "github.com/google/uuid"

// Goal represents a goal scored in a match
type Goal struct {
	BaseEntity
	MatchID   uuid.UUID `gorm:"type:uuid;not null;index" json:"match_id"`
	PlayerID  uuid.UUID `gorm:"type:uuid;not null;index" json:"player_id"`
	TeamID    uuid.UUID `gorm:"type:uuid;not null;index" json:"team_id"`
	Minute    int       `gorm:"not null" json:"minute"` // Minute when goal was scored
	IsOwnGoal bool      `gorm:"default:false" json:"is_own_goal"`
	Match     *Match    `gorm:"foreignKey:MatchID" json:"match,omitempty"`
	Player    *Player   `gorm:"foreignKey:PlayerID" json:"player,omitempty"`
	Team      *Team     `gorm:"foreignKey:TeamID" json:"team,omitempty"`
}

// TableName returns the table name for Goal entity
func (Goal) TableName() string {
	return "goals"
}
