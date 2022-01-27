package store

import (
	customErrors "carAPI/custom-errors"
	"carAPI/model"
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"log"
	"reflect"
	"testing"
)

var car = model.Car{
	ID:                uuid.NewString(),
	Name:              "Roadster",
	YearOfManufacture: 2000,
	Brand:             "Tesla",
	FuelType:          "Electric",
	Engine: model.Engine{
		ID: uuid.NewString(),
	},
}

var carNotExists customErrors.CarNotExists

func TestStore_GetByBrand(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		log.Println(err)
	}

	defer db.Close()

	store := New(db)
	rows := sqlmock.NewRows([]string{"carID", "name", "yearOfManufacture", "brand", "fuelType", "engineId"}).
		AddRow(car.ID, car.Name, car.YearOfManufacture, car.Brand, car.FuelType, car.Engine.ID)
	mock.ExpectQuery("select \\* from cars where brand = \\?").WithArgs("Tesla").WillReturnRows(rows)
	mock.ExpectQuery("select \\* from cars").WillReturnRows(rows)
	mock.ExpectQuery("select \\* from cars").WillReturnError(errors.New("DB error"))

	tests := []struct {
		desc  string
		brand string
		cars  []model.Car
		err   error
	}{
		{"Fetch all Tesla cars", "Tesla", []model.Car{car}, nil},
		{"Fetch all cars", "", nil, nil},
		{"DB error", "", nil, errors.New("DB error")},
	}

	for i, tc := range tests {
		cars, err := store.GetByBrand(tc.brand)
		if !reflect.DeepEqual(err, tc.err) {
			t.Errorf("Testcase[%v] failed (%v)\nExpected error: %v\nGot: %v", i, tc.desc, tc.err, err)
		}

		if !reflect.DeepEqual(cars, tc.cars) {
			t.Errorf("Testcase[%v] failed (%v)\nExpected:\n%v\nGot:\n%v", i, tc.desc, tc.cars, cars)
		}
	}
}

func TestStore_GetByID(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		log.Println(err)
	}

	defer db.Close()

	store := New(db)
	rows := sqlmock.NewRows([]string{"carID", "name", "yearOfManufacture", "brand", "fuelType", "engineId"}).
		AddRow(car.ID, car.Name, car.YearOfManufacture, car.Brand, car.FuelType, car.Engine.ID)

	mock.ExpectQuery("select \\* from cars where carId = \\?").WithArgs(car.ID).WillReturnRows(rows)
	mock.ExpectQuery("select \\* from cars where carId = \\?").WithArgs("1").WillReturnError(carNotExists)
	mock.ExpectQuery("select \\* from cars where carId = \\?").WithArgs("2").WillReturnError(errors.New("DB error"))

	tests := []struct {
		desc string
		id   string
		cars model.Car
		err  error
	}{
		{"Success", car.ID, car, nil},
		{"Not exists", "1", model.Car{}, carNotExists},
		{"DB error", "2", model.Car{}, errors.New("DB error")},
	}

	for i, tc := range tests {
		cars, err := store.GetByID(tc.id)
		if !reflect.DeepEqual(err, tc.err) {
			t.Errorf("Testcase[%v] failed (%v)\nExpected error: %v\nGot: %v", i, tc.desc, tc.err, err)
		}

		if !reflect.DeepEqual(cars, tc.cars) {
			t.Errorf("Testcase[%v] failed (%v)\nExpected:\n%v\nGot:\n%v", i, tc.desc, tc.cars, cars)
		}
	}
}

func TestStore_Create(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		log.Println(err)
	}

	defer db.Close()

	store := New(db)

	query := "insert into cars \\(carId, name, yearOfManufacture, brand, fuelType, engineId\\) values \\(\\?, \\?, \\?, \\?, \\?, \\?\\)"

	prep := mock.ExpectPrepare(query)
	prep.ExpectExec().WithArgs(sqlmock.AnyArg(), car.Name, car.YearOfManufacture, car.Brand, car.FuelType, car.Engine.ID).
		WillReturnResult(sqlmock.NewResult(0, 1))
	prep = mock.ExpectPrepare(query)
	prep.ExpectExec().WillReturnError(errors.New("DB error"))

	tests := []struct {
		desc     string
		input    model.Car
		expected model.Car
		err      error
	}{
		{"Success", car, car, nil},
		{"DB error", model.Car{}, model.Car{}, errors.New("DB error")},
	}

	for i, tc := range tests {
		car, err := store.Create(tc.input)
		if !reflect.DeepEqual(err, tc.err) {
			t.Errorf("Testcase[%v] failed (%v)\nExpected error: %v\nGot: %v", i, tc.desc, tc.err, err)

			if car.ID == tc.expected.ID {
				t.Errorf("Testcase[%v] failed (%v)\nNew ID was not assigned", i, tc.desc)
			}
		}

		car.ID = tc.expected.ID

		if !reflect.DeepEqual(car, tc.expected) {
			t.Errorf("Testcase[%v] failed (%v)\nExpected:\n%v\nGot:\n%v", i, tc.desc, tc.expected, car)
		}
	}
}

func TestStore_Update(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		log.Println(err)
	}

	defer db.Close()

	store := New(db)

	query := "update cars set name = \\?, yearOfManufacture = \\?, brand = \\?, fuelType = \\? where carId = \\?"

	prep := mock.ExpectPrepare(query)
	prep.ExpectExec().WithArgs("Roadster", 2000, "Tesla", "Electric", car.ID).WillReturnResult(sqlmock.NewResult(0, 1))
	prep = mock.ExpectPrepare(query)
	prep.ExpectExec().WillReturnError(carNotExists)
	prep = mock.ExpectPrepare(query)
	prep.ExpectExec().WillReturnError(errors.New("DB error"))

	tests := []struct {
		desc     string
		input    model.Car
		expected model.Car
		err      error
	}{
		{"Success", car, car, nil},
		{"Not exists", model.Car{}, model.Car{}, carNotExists},
		{"DB error", model.Car{}, model.Car{}, errors.New("DB error")},
	}

	for i, tc := range tests {
		car, err := store.Update(tc.input)
		if !reflect.DeepEqual(err, tc.err) {
			t.Errorf("Testcase[%v] failed (%v)\nExpected error: %v\nGot: %v", i, tc.desc, tc.err, err)
		}

		if !reflect.DeepEqual(car, tc.expected) {
			t.Errorf("Testcase[%v] failed (%v)\nExpected:\n%v\nGot:\n%v", i, tc.desc, tc.expected, car)
		}
	}
}

func TestStore_Delete(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		log.Println(err)
	}

	defer db.Close()

	store := New(db)

	query := "delete from cars where carId = \\?"

	prep := mock.ExpectPrepare(query)
	prep.ExpectExec().WithArgs(car.ID).WillReturnResult(sqlmock.NewResult(0, 1))
	prep = mock.ExpectPrepare(query)
	prep.ExpectExec().WithArgs(car.ID).WillReturnError(carNotExists)
	prep = mock.ExpectPrepare(query)
	prep.ExpectExec().WillReturnError(errors.New("DB error"))

	tests := []struct {
		desc string
		id   string
		err  error
	}{
		{"Success", car.ID, nil},
		{"Not exists", car.ID, carNotExists},
		{"DB error", "", errors.New("DB error")},
	}

	for i, tc := range tests {
		err := store.Delete(tc.id)
		if !reflect.DeepEqual(err, tc.err) {
			t.Errorf("Testcase[%v] failed (%v)\nExpected error: %v\nGot: %v", i, tc.desc, tc.err, err)
		}
	}
}
