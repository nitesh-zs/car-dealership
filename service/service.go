package service

import (
	"carAPI/model"
	"carAPI/store"
)

type service struct {
	carStore    store.CarStore
	engineStore store.EngineStore
}

func New(c store.CarStore, e store.EngineStore) service {
	return service{
		carStore:    c,
		engineStore: e,
	}
}

func (s service) GetAll(brand string, withEngine bool) ([]model.Car, error) {
	cars, err := s.carStore.GetByBrand(brand)
	if err != nil {
		return nil, err
	}

	// if withEngine is true, then populate the Engine data in Cars
	if withEngine {
		for i, _ := range cars {
			car := cars[i]
			eID := car.Engine.ID
			engine, err := s.engineStore.GetByID(eID)
			if err != nil {
				return nil, err
			}
			cars[i].Engine = engine
		}
	}

	return cars, nil
}

func (s service) GetByID(ID string) (model.Car, error) {
	car, err := s.carStore.GetByID(ID)
	if err != nil {
		return model.Car{}, err
	}

	engine, err := s.engineStore.GetByID(car.Engine.ID)
	if err != nil {
		return model.Car{}, err
	}

	car.Engine = engine
	return car, nil
}

func (s service) Create(car model.Car) (model.Car, error) {
	engine, err := s.engineStore.Create(car.Engine)
	if err != nil {
		return model.Car{}, err
	}

	car.Engine.ID = engine.ID

	newCar, err := s.carStore.Create(car)
	if err != nil {
		return model.Car{}, err
	}

	newCar.Engine = engine
	return newCar, nil
}

func (s service) Update(car model.Car) (model.Car, error) {
	updatedEngine, err := s.engineStore.Update(car.Engine)
	if err != nil {
		return model.Car{}, err
	}

	updatedCar, err := s.carStore.Update(car)
	if err != nil {
		return model.Car{}, err
	}

	updatedCar.Engine = updatedEngine
	return updatedCar, nil
}

func (s service) Delete(ID string) error {
	car, err := s.carStore.GetByID(ID)
	if err != nil {
		return err
	}

	err = s.carStore.Delete(ID)
	if err != nil {
		return err
	}

	err = s.engineStore.Delete(car.Engine.ID)
	if err != nil {
		return err
	}

	return nil
}
