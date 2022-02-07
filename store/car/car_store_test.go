package car

import (
	"errors"
	"log"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	customErrors "carAPI/custom-errors"
	"carAPI/model"
)

func car() model.Car {
	return model.Car{
		ID:                uuid.NewString(),
		Name:              "Roadster",
		YearOfManufacture: 2000,
		Brand:             "Tesla",
		FuelType:          "Electric",
		Engine: model.Engine{
			ID: uuid.NewString(),
		},
	}
}

func TestStore_GetByBrand(t *testing.T) {
	car := car()

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

		assert.Equalf(t, tc.err, err, "Testcase[%v] (%v)", i, tc.desc)

		assert.Equalf(t, tc.cars, cars, "Testcase[%v] (%v)", i, tc.desc)
	}
}

func TestStore_GetByID(t *testing.T) {
	car := car()

	db, mock, err := sqlmock.New()
	if err != nil {
		log.Println(err)
	}

	defer db.Close()

	store := New(db)
	rows := sqlmock.NewRows([]string{"carID", "name", "yearOfManufacture", "brand", "fuelType", "engineId"}).
		AddRow(car.ID, car.Name, car.YearOfManufacture, car.Brand, car.FuelType, car.Engine.ID)

	mock.ExpectQuery("select \\* from cars where carId = \\?").WithArgs(car.ID).WillReturnRows(rows)
	mock.ExpectQuery("select \\* from cars where carId = \\?").WithArgs("1").WillReturnError(customErrors.CarNotExists())
	mock.ExpectQuery("select \\* from cars where carId = \\?").WithArgs("2").WillReturnError(errors.New("DB error"))

	tests := []struct {
		desc string
		id   string
		car  *model.Car
		err  error
	}{
		{"Success", car.ID, &car, nil},
		{"Not exists", "1", nil, customErrors.CarNotExists()},
		{"DB error", "2", nil, errors.New("DB error")},
	}

	for i, tc := range tests {
		car, err := store.GetByID(tc.id)

		assert.Equalf(t, tc.err, err, "Testcase[%v] (%v)", i, tc.desc)

		assert.Equalf(t, tc.car, car, "Testcase[%v] (%v)", i, tc.desc)
	}
}

func TestStore_Create(t *testing.T) {
	car := car()
	car2 := car

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
		input    *model.Car
		expected *model.Car
		err      error
	}{
		{"Success", &car, &car2, nil},
		{"DB error", &model.Car{}, nil, errors.New("DB error")},
	}

	for i, tc := range tests {
		car, err := store.Create(tc.input)

		assert.Equalf(t, tc.err, err, "Testcase[%v] (%v)", i, tc.desc)

		if car != nil {
			car.ID = tc.expected.ID
		}

		assert.Equalf(t, tc.expected, car, "Testcase[%v] (%v)", i, tc.desc)
	}
}

func TestStore_Update(t *testing.T) {
	car := car()

	db, mock, err := sqlmock.New()
	if err != nil {
		log.Println(err)
	}

	defer db.Close()

	store := New(db)

	query := "update cars set name = \\?, yearOfManufacture = \\?, brand = \\?, fuelType = \\? where carId = \\?"

	// success case

	prep := mock.ExpectPrepare(query)
	prep.ExpectExec().WithArgs("Roadster", 2000, "Tesla", "Electric", car.ID).WillReturnResult(sqlmock.NewResult(0, 1))

	// DB error

	prep = mock.ExpectPrepare(query)
	prep.ExpectExec().WillReturnError(errors.New("DB error"))

	tests := []struct {
		desc     string
		input    *model.Car
		expected *model.Car
		err      error
	}{
		{"Success", &car, &car, nil},
		{"DB error", &model.Car{ID: "1"}, nil, errors.New("DB error")},
	}

	for i, tc := range tests {
		car, err := store.Update(tc.input)

		assert.Equalf(t, tc.err, err, "Testcase[%v] (%v)", i, tc.desc)

		assert.Equalf(t, tc.expected, car, "Testcase[%v] (%v)", i, tc.desc)
	}
}

func TestStore_Delete(t *testing.T) {
	car := car()

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
	prep.ExpectExec().WithArgs(car.ID).WillReturnError(customErrors.CarNotExists())

	prep = mock.ExpectPrepare(query)
	prep.ExpectExec().WillReturnError(errors.New("DB error"))

	tests := []struct {
		desc string
		id   string
		err  error
	}{
		{"Success", car.ID, nil},
		{"Not exists", car.ID, customErrors.CarNotExists()},
		{"DB error", "", errors.New("DB error")},
	}

	for i, tc := range tests {
		err := store.Delete(tc.id)

		assert.Equalf(t, tc.err, err, "Testcase[%v] (%v)", i, tc.desc)
	}
}
