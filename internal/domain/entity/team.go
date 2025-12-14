package entity

// Team represents a football team
type Team struct {
	BaseEntity
	Name        string   `gorm:"not null;size:255" json:"name"`
	Logo        string   `gorm:"size:500" json:"logo"`
	FoundedYear int      `gorm:"not null" json:"founded_year"`
	Address     string   `gorm:"size:500" json:"address"`
	City        string   `gorm:"not null;size:100" json:"city"`
	Players     []Player `gorm:"foreignKey:TeamID" json:"players,omitempty"`
}

// TableName returns the table name for Team entity
func (Team) TableName() string {
	return "teams"
}
