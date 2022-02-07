package service

import (
	"carAPI/model"
)

type CarService interface {

	// GetAll takes two params- brand and withEngine
	// if empty string is passed to brand, then cars of all brands should be fetched
	GetAll(brand string, withEngine bool) ([]model.Car, error)

	GetByID(id string) (*model.Car, error)
	Create(car *model.Car) (*model.Car, error)
	Update(car *model.Car) (*model.Car, error)
	Delete(id string) error
}
