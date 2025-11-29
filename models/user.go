package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
	Name         string `json:"name"`
	Email        string `json:"email" gorm:"uniqueIndex;not null"`
	PasswordHash string `json:"-" gorm:"not null"`
}

type RegisterRequest struct {
	// Accept either first_name+last_name (preferred) or name (legacy)
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Name      string `json:"name"`
	Email     string `json:"email" binding:"required,email"`
	Password  string `json:"password" binding:"required,min=6"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type ChangePasswordRequest struct {
	CurrentPassword string `json:"current_password" binding:"required"`
	NewPassword     string `json:"new_password" binding:"required,min=6"`
}

type UpdateResponse struct {
	Message       string `json:"message"`
	Firstname     string `json:"first_name"`
	Lastname      string `json:"last_name"`
	UpdateAt      string `json:"updated_at"`
	EmailVerified bool   `json:"email_verified"`
}

type AuthResponse struct {
	Token string `json:"token"`
}

type UserResponse struct {
	ID    uint   `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type DetailedUserResponse struct {
	ID            uint   `json:"id"`
	Email         string `json:"email"`
	Username      string `json:"username,omitempty"`
	FirstName     string `json:"first_name"`
	LastName      string `json:"last_name"`
	CreatedAt     string `json:"created_at"`
	EmailVerified bool   `json:"email_verified"`
}

type AuthData struct {
	User  DetailedUserResponse `json:"user"`
	Token string               `json:"token"`
}

type SuccessResponse struct {
	Success struct {
		Status  int         `json:"status"`
		Message string      `json:"message"`
		Data    interface{} `json:"data"`
	} `json:"success"`
}
