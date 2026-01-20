package model

// UserRole represents user roles
type UserRole struct {
	RoleName string `gorm:"primaryKey;size:16" json:"role_name"`
	Users    []User `gorm:"foreignKey:Role;references:RoleName" json:"users,omitempty"`
}

var (
	CommonUserRole = UserRole{RoleName: "User"}
	AdminRole      = UserRole{RoleName: "Admin"}
	ModeratorRole  = UserRole{RoleName: "Moderator"}
)
