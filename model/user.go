package model

type User struct {
	UserID       uint   `gorm:"primaryKey" json:"user_id"`
	UserUsername string `gorm:"unique" json:"user_username"`
	UserPassword string `json:"-"`
	UserRole     string `json:"user_role"`
}

type UserResponseDTO struct {
	UserID       uint   `json:"user_id"`
	UserUsername string `json:"user_username"`
	UserRole     string `json:"user_role"`
}

func (u *User) ToDTO() UserResponseDTO {
	return UserResponseDTO{
		UserID:       u.UserID,
		UserUsername: u.UserUsername,
		UserRole:     u.UserRole,
	}
}
