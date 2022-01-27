package store

import (
	"carAPI/model"
	"database/sql"
	"github.com/google/uuid"
)

type store struct {
	db *sql.DB
}

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
		rows.Err()
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

func (s store) GetByID(ID string) (model.Car, error) {
	var car model.Car

	row := s.db.QueryRow("select * from cars where carId = ?", ID)
	err := row.Scan(&car.ID, &car.Name, &car.YearOfManufacture, &car.Brand, &car.FuelType, &car.Engine.ID)
	if err == sql.ErrNoRows {
		return model.Car{}, carNotExists
	}
	if err != nil {
		return model.Car{}, err
	}

	return car, nil
}

func (s store) Create(car model.Car) (model.Car, error) {
	car.ID = uuid.NewString()
	stmt, err := s.db.Prepare(`insert into cars (carId, name, yearOfManufacture, brand, fuelType, engineId)
										values (?, ?, ?, ?, ?, ?)`)
	if err != nil {
		return model.Car{}, err
	}

	defer stmt.Close()

	_, err = stmt.Exec(car.ID, car.Name, car.YearOfManufacture, car.Brand, car.FuelType, car.Engine.ID)
	if err != nil {
		return model.Car{}, err
	}

	return car, nil
}

func (s store) Update(car model.Car) (model.Car, error) {
	stmt, err := s.db.Prepare(`update cars set name = ?, yearOfManufacture = ?, brand = ?, fuelType = ? where carId = ?`)

	if err != nil {
		return model.Car{}, err
	}

	defer stmt.Close()

	res, err := stmt.Exec(car.Name, car.YearOfManufacture, car.Brand, car.FuelType, car.ID)
	if err != nil {
		return model.Car{}, err
	}

	rowsAff, err := res.RowsAffected()
	if rowsAff == 0 {
		return model.Car{}, carNotExists
	}

	return car, nil
}

func (s store) Delete(ID string) error {
	stmt, err := s.db.Prepare(`delete from cars where carId = ?`)
	if err == sql.ErrNoRows {
		return carNotExists
	}
	if err != nil {
		return err
	}

	defer stmt.Close()

	_, err = stmt.Exec(ID)
	if err != nil {
		return err
	}

	return nil
}
