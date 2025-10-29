package domain

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID           uuid.UUID `json:"id"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"` // Never expose password hash
	FirstName    string    `json:"first_name"`
	LastName     string    `json:"last_name"`
	HospitalID   string    `json:"hospital_id"`
	Specialty    string    `json:"specialty"`
	Role         string    `json:"role"` // doctor, admin
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type RegisterRequest struct {
	Email      string `json:"email" binding:"required,email"`
	Password   string `json:"password" binding:"required,min=6"`
	FirstName  string `json:"firstName" binding:"required"`
	LastName   string `json:"lastName" binding:"required"`
	HospitalID string `json:"hospitalId" binding:"required"`
	Specialty  string `json:"specialty" binding:"required"`
}

type AuthResponse struct {
	Token string `json:"token"`
	User  User   `json:"user"`
}

func NewUser(email, passwordHash, firstName, lastName, hospitalID, specialty string) *User {
	return &User{
		ID:           uuid.New(),
		Email:        email,
		PasswordHash: passwordHash,
		FirstName:    firstName,
		LastName:     lastName,
		HospitalID:   hospitalID,
		Specialty:    specialty,
		Role:         "doctor",
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}
}
