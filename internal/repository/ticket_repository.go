package repository

import (
	"database/sql"

	"docger/internal/model"
)

type TicketRepository struct {
	db *sql.DB
}

func NewTicketRepository(db *sql.DB) *TicketRepository {
	return &TicketRepository{db: db}
}

func (r *TicketRepository) Create(ticket *model.Ticket) error {
	query := `
		INSERT INTO tickets (title, description, status, priority, user_id)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at, updated_at`

	return r.db.QueryRow(query, ticket.Title, ticket.Description, ticket.Status, ticket.Priority, ticket.UserID).
		Scan(&ticket.ID, &ticket.CreatedAt, &ticket.UpdatedAt)
}

func (r *TicketRepository) FindByID(id int) (*model.Ticket, error) {
	query := `
		SELECT id, title, description, status, priority, user_id, assigned_to, created_at, updated_at
		FROM tickets WHERE id = $1`

	ticket := &model.Ticket{}
	err := r.db.QueryRow(query, id).Scan(
		&ticket.ID, &ticket.Title, &ticket.Description, &ticket.Status,
		&ticket.Priority, &ticket.UserID, &ticket.AssignedTo,
		&ticket.CreatedAt, &ticket.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return ticket, nil
}

func (r *TicketRepository) FindByUserID(userID int) ([]model.Ticket, error) {
	query := `
		SELECT id, title, description, status, priority, user_id, assigned_to, created_at, updated_at
		FROM tickets WHERE user_id = $1 ORDER BY created_at DESC`

	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tickets []model.Ticket
	for rows.Next() {
		var ticket model.Ticket
		err := rows.Scan(
			&ticket.ID, &ticket.Title, &ticket.Description, &ticket.Status,
			&ticket.Priority, &ticket.UserID, &ticket.AssignedTo,
			&ticket.CreatedAt, &ticket.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		tickets = append(tickets, ticket)
	}
	return tickets, nil
}

func (r *TicketRepository) FindAll() ([]model.Ticket, error) {
	query := `
		SELECT id, title, description, status, priority, user_id, assigned_to, created_at, updated_at
		FROM tickets ORDER BY created_at DESC`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tickets []model.Ticket
	for rows.Next() {
		var ticket model.Ticket
		err := rows.Scan(
			&ticket.ID, &ticket.Title, &ticket.Description, &ticket.Status,
			&ticket.Priority, &ticket.UserID, &ticket.AssignedTo,
			&ticket.CreatedAt, &ticket.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		tickets = append(tickets, ticket)
	}
	return tickets, nil
}

func (r *TicketRepository) Update(ticket *model.Ticket) error {
	query := `
		UPDATE tickets 
		SET title = $1, description = $2, status = $3, priority = $4, assigned_to = $5, updated_at = CURRENT_TIMESTAMP
		WHERE id = $6
		RETURNING updated_at`

	return r.db.QueryRow(query, ticket.Title, ticket.Description, ticket.Status, ticket.Priority, ticket.AssignedTo, ticket.ID).
		Scan(&ticket.UpdatedAt)
}

func (r *TicketRepository) Delete(id int) error {
	query := `DELETE FROM tickets WHERE id = $1`
	_, err := r.db.Exec(query, id)
	return err
}

func (r *TicketRepository) Assign(id int, assignedTo int) error {
	query := `
		UPDATE tickets 
		SET assigned_to = $1, updated_at = CURRENT_TIMESTAMP
		WHERE id = $2`

	_, err := r.db.Exec(query, assignedTo, id)
	return err
}
