package model

import "time"

type Ticket struct {
	ID          int       `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Status      string    `json:"status"`
	Priority    string    `json:"priority"`
	UserID      int       `json:"user_id"`
	AssignedTo  *int      `json:"assigned_to,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type CreateTicketRequest struct {
	Title       string `json:"title" binding:"required,min=3,max=255"`
	Description string `json:"description" binding:"required,min=10"`
	Priority    string `json:"priority" binding:"required,oneof=low medium high"`
}

type UpdateTicketRequest struct {
	Title       string `json:"title" binding:"required,min=3,max=255"`
	Description string `json:"description" binding:"required,min=10"`
	Status      string `json:"status" binding:"required,oneof=open in_progress resolved closed"`
	Priority    string `json:"priority" binding:"required,oneof=low medium high"`
	AssignedTo  *int   `json:"assigned_to"`
}

type AssignTicketRequest struct {
	AssignedTo int `json:"assigned_to" binding:"required"`
}
