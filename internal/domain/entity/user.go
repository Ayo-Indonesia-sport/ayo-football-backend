package entity

// UserRole represents the role of a user
type UserRole string

const (
	RoleAdmin UserRole = "admin"
	RoleUser  UserRole = "user"
)

// User represents a system user
type User struct {
	BaseEntity
	Email    string   `gorm:"uniqueIndex;not null;size:255" json:"email"`
	Password string   `gorm:"not null;size:255" json:"-"`
	Name     string   `gorm:"not null;size:255" json:"name"`
	Role     UserRole `gorm:"type:varchar(20);default:'user'" json:"role"`
}

// TableName returns the table name for User entity
func (User) TableName() string {
	return "users"
}

// IsAdmin checks if user has admin role
func (u *User) IsAdmin() bool {
	return u.Role == RoleAdmin
}
