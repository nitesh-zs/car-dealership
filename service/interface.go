package service

import (
	"carAPI/model"
)

type CarService interface {

	// GetAll takes two params- brand and withEngine
	// if empty string is passed to brand, then cars of all brands are fetched
	GetAll(brand string, withEngine bool) ([]model.Car, error)

	// GetByID fetches a car with a given carID from DB
	GetByID(id string) (*model.Car, error)

	// Create creates a car and its underlying engine in the DB
	Create(car *model.Car) (*model.Car, error)

	// Update updates an existing car in DB
	Update(car *model.Car) (*model.Car, error)

	// Delete deletes the car with given ID from the DB
	Delete(id string) error
}
