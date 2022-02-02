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
		rows, err = s.db.Query("select * from cars")
	} else {
		rows, err = s.db.Query("select * from cars where brand = ?", brand)
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

	row := s.db.QueryRow("select * from cars where carId = ?", id)
	err := row.Scan(&car.ID, &car.Name, &car.YearOfManufacture, &car.Brand, &car.FuelType, &car.Engine.ID)

	if err == sql.ErrNoRows {
		return &model.Car{}, customErrors.CarNotExists()
	}

	if err != nil {
		return &model.Car{}, err
	}

	return &car, nil
}

func (s store) Create(car *model.Car) (*model.Car, error) {
	car.ID = uuid.NewString()

	stmt, err := s.db.Prepare(`insert into cars (carId, name, yearOfManufacture, brand, fuelType, engineId)
										values (?, ?, ?, ?, ?, ?)`)

	if err != nil {
		return &model.Car{}, err
	}

	defer stmt.Close()

	_, err = stmt.Exec(car.ID, car.Name, car.YearOfManufacture, car.Brand, car.FuelType, car.Engine.ID)
	if err != nil {
		return &model.Car{}, err
	}

	return car, nil
}

func (s store) Update(car *model.Car) (*model.Car, error) {
	// check if record exists in table
	carFromDB, err := s.GetByID(car.ID)
	if err != nil {
		return &model.Car{}, err
	}

	stmt, err := s.db.Prepare(`update cars set name = ?, yearOfManufacture = ?, brand = ?, fuelType = ? where carId = ?`)

	if err != nil {
		return &model.Car{}, err
	}

	defer stmt.Close()

	_, err = stmt.Exec(car.Name, car.YearOfManufacture, car.Brand, car.FuelType, car.ID)
	if err != nil {
		return &model.Car{}, err
	}

	car.Engine.ID = carFromDB.Engine.ID

	return car, nil
}

func (s store) Delete(id string) error {
	stmt, err := s.db.Prepare(`delete from cars where carId = ?`)
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
