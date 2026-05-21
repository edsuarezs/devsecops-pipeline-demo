package models

import (
	"errors"
	"strings"

)

// Item represents a stored item with its generated ID.
type Item struct {
	ID          int     `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
}

// ItemCreate represents the input for creating an item.
type ItemCreate struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
}

// Validate checks that the input meets business rules.
func (i *ItemCreate) Validate() error {
	i.Name = strings.TrimSpace(i.Name)
	if i.Name == "" {
		return errors.New("name is required")
	}
	if len(i.Name) > 100 {
		return errors.New("name must be 100 characters or less")
	}
	if len(i.Description) > 500 {
		return errors.New("description must be 500 characters or less")
	}
	if i.Price <= 0 {
		return errors.New("price must be greater than zero")
	}
	return nil
}
