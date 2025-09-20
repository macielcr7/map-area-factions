package models

type UserRole string

const (
	RoleAdmin        UserRole = "admin"
	RoleModerator    UserRole = "moderator"
	RoleCollaborator UserRole = "collaborator"
	RoleCitizen      UserRole = "citizen"
)

type User struct {
	BaseModel
	Name         string   `json:"name" gorm:"not null"`
	Email        string   `json:"email" gorm:"uniqueIndex;not null"`
	PasswordHash string   `json:"-" gorm:"not null"`
	Role         UserRole `json:"role" gorm:"default:'citizen'"`
	Status       string   `json:"status" gorm:"default:'active'"`
}

// TableName specifies the table name for the User model
func (User) TableName() string {
	return "users"
}