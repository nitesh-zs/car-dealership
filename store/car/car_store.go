package car

import (
	"database/sql"
	"log"

	"github.com/google/uuid"

	customErrors "carAPI/custom-errors"
	"carAPI/model"
)

type store struct {
	db *sql.DB
}

//nolint:revive //store should not be exported
func New(db *sql.DB) store {
	return store{db: db}
}

func (s store) GetByBrand(brand string) ([]model.Car, error) {
	var (
		car  model.Car
		rows *sql.Rows
		err  error
		cars []model.Car
	)

	if brand == "" {
		rows, err = s.db.Query(getAllCars)
	} else {
		rows, err = s.db.Query(getCarByBrand, brand)
	}

	if err != nil {
		return nil, err
	}

	defer func() {
		rows.Close()

		err = rows.Err()
		if err != nil {
			log.Println(err)
		}
	}()

	for rows.Next() {
		err := rows.Scan(&car.ID, &car.Name, &car.YearOfManufacture, &car.Brand, &car.FuelType, &car.Engine.ID)
		if err != nil {
			return nil, err
		}

		cars = append(cars, car)
	}

	return cars, nil
}

func (s store) GetByID(id string) (*model.Car, error) {
	var car model.Car

	row := s.db.QueryRow(getCarByID, id)
	err := row.Scan(&car.ID, &car.Name, &car.YearOfManufacture, &car.Brand, &car.FuelType, &car.Engine.ID)

	if err == sql.ErrNoRows {
		return nil, customErrors.CarNotExists()
	}

	if err != nil {
		return nil, err
	}

	return &car, nil
}

func (s store) Create(car *model.Car) (*model.Car, error) {
	car.ID = uuid.NewString()

	stmt, err := s.db.Prepare(insertCar)

	if err != nil {
		return nil, err
	}

	defer stmt.Close()

	_, err = stmt.Exec(car.ID, car.Name, car.YearOfManufacture, car.Brand, car.FuelType, car.Engine.ID)
	if err != nil {
		return nil, err
	}

	return car, nil
}

func (s store) Update(car *model.Car) (*model.Car, error) {
	stmt, err := s.db.Prepare(updateCar)
	if err != nil {
		return nil, err
	}

	defer stmt.Close()

	_, err = stmt.Exec(car.Name, car.YearOfManufacture, car.Brand, car.FuelType, car.ID)
	if err != nil {
		return nil, err
	}

	return car, nil
}

func (s store) Delete(id string) error {
	stmt, err := s.db.Prepare(deleteCar)
	if err == sql.ErrNoRows {
		return customErrors.CarNotExists()
	}

	if err != nil {
		return err
	}

	defer stmt.Close()

	_, err = stmt.Exec(id)
	if err != nil {
		return err
	}

	return nil
}
