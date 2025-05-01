package entity

import "time"

type CallDTO struct {
	ClientName  string `json:"client_name" binding:"required"`
	PhoneNumber string `json:"phone_number" binding:"required"`
	Description string `json:"description" binding:"required"`
}

type UpdateCallStatusDTO struct {
	Status string `json:"status" binding:"required"`
}

type CallResponse struct {
	ID          int64     `json:"id"`
	ClientName  string    `json:"client_name"`
	PhoneNumber string    `json:"phone_number"`
	Description string    `json:"description"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
}

type Call struct {
	ID          int64     `json:"id"`
	ClientName  string    `json:"client_name"`
	PhoneNumber string    `json:"phone_number"`
	Description string    `json:"description"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	UserID      int64     `json:"user_id"`
}

type AuthRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}
