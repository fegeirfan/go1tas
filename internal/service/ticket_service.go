package service

import (
	"errors"

	"docger/internal/model"
	"docger/internal/repository"
)

type TicketService struct {
	ticketRepo *repository.TicketRepository
	userRepo   *repository.UserRepository
}

func NewTicketService(ticketRepo *repository.TicketRepository, userRepo *repository.UserRepository) *TicketService {
	return &TicketService{
		ticketRepo: ticketRepo,
		userRepo:   userRepo,
	}
}

func (s *TicketService) CreateTicket(req *model.CreateTicketRequest, userID int) (*model.Ticket, error) {
	ticket := &model.Ticket{
		Title:       req.Title,
		Description: req.Description,
		Status:      "open",
		Priority:    req.Priority,
		UserID:      userID,
	}

	err := s.ticketRepo.Create(ticket)
	if err != nil {
		return nil, err
	}

	return ticket, nil
}

func (s *TicketService) GetTicketByID(id int) (*model.Ticket, error) {
	return s.ticketRepo.FindByID(id)
}

func (s *TicketService) GetUserTickets(userID int) ([]model.Ticket, error) {
	return s.ticketRepo.FindByUserID(userID)
}

func (s *TicketService) GetAllTickets() ([]model.Ticket, error) {
	return s.ticketRepo.FindAll()
}

func (s *TicketService) UpdateTicket(id int, req *model.UpdateTicketRequest, currentUser *model.User) (*model.Ticket, error) {
	ticket, err := s.ticketRepo.FindByID(id)
	if err != nil {
		return nil, err
	}

	// Only ticket owner or admin can update
	if ticket.UserID != currentUser.ID && currentUser.Role != "admin" {
		return nil, errors.New("unauthorized to update this ticket")
	}

	ticket.Title = req.Title
	ticket.Description = req.Description
	ticket.Status = req.Status
	ticket.Priority = req.Priority

	// Admin can assign tickets
	if currentUser.Role == "admin" && req.AssignedTo != nil {
		ticket.AssignedTo = req.AssignedTo
	}

	err = s.ticketRepo.Update(ticket)
	if err != nil {
		return nil, err
	}

	return ticket, nil
}

func (s *TicketService) DeleteTicket(id int, currentUser *model.User) error {
	ticket, err := s.ticketRepo.FindByID(id)
	if err != nil {
		return err
	}

	// Only ticket owner or admin can delete
	if ticket.UserID != currentUser.ID && currentUser.Role != "admin" {
		return errors.New("unauthorized to delete this ticket")
	}

	return s.ticketRepo.Delete(id)
}

func (s *TicketService) AssignTicket(id int, assignedToID int, currentUser *model.User) (*model.Ticket, error) {
	// Only admin can assign tickets
	if currentUser.Role != "admin" {
		return nil, errors.New("only admin can assign tickets")
	}

	// Verify the user exists
	_, err := s.userRepo.FindByID(assignedToID)
	if err != nil {
		return nil, errors.New("assigned user not found")
	}

	err = s.ticketRepo.Assign(id, assignedToID)
	if err != nil {
		return nil, err
	}

	return s.ticketRepo.FindByID(id)
}
