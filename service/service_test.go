package service

import (
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	customErrors "carAPI/custom-errors"
	"carAPI/mocks"
	"carAPI/model"
)

func car1() model.Car {
	return model.Car{

		ID:                "1",
		Name:              "Roadster",
		YearOfManufacture: 2000,
		Brand:             "Tesla",
		FuelType:          "Electric",
		Engine: model.Engine{
			ID:            "1",
			Displacement:  0,
			NoOfCylinders: 0,
			Range:         400,
		},
	}
}

func car2() model.Car {
	return model.Car{

		ID:                "2",
		Name:              "Abc",
		YearOfManufacture: 2020,
		Brand:             "Ferrari",
		FuelType:          "Diesel",
		Engine: model.Engine{
			ID:            "2",
			Displacement:  600,
			NoOfCylinders: 4,
			Range:         0,
		},
	}
}

func car3() model.Car {
	return model.Car{
		ID:                "1",
		Name:              "Roadster",
		YearOfManufacture: 2000,
		Brand:             "Tesla",
		FuelType:          "Electric",
		Engine:            model.Engine{ID: "1"},
	}
}

func car4() model.Car {
	return model.Car{
		ID:                "2",
		Name:              "Abc",
		YearOfManufacture: 2020,
		Brand:             "Ferrari",
		FuelType:          "Diesel",
		Engine:            model.Engine{ID: "2"},
	}
}

func TestService_GetAll(t *testing.T) {
	car1 := car1()
	car2 := car2()

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	m := mocks.NewMockCarStore(mockCtrl)
	s := mocks.NewMockEngineStore(mockCtrl)

	m.EXPECT().GetByBrand("Tesla").Return([]model.Car{car3()}, nil).AnyTimes()

	m.EXPECT().GetByBrand("").Return([]model.Car{car3(), car4()}, nil).AnyTimes()

	m.EXPECT().GetByBrand("Jaguar").Return(nil, errors.New("server error"))

	s.EXPECT().GetAll().Return(map[string]model.Engine{
		"1": {
			ID:            "1",
			Displacement:  0,
			NoOfCylinders: 0,
			Range:         400,
		},
		"2": {
			ID:            "2",
			Displacement:  600,
			NoOfCylinders: 4,
			Range:         0,
		},
	}, nil).AnyTimes()

	tests := []struct {
		desc  string
		brand string
		cars  []model.Car
		err   error
	}{
		{
			"Fetch Tesla cars",
			"Tesla",
			[]model.Car{car1},
			nil,
		},
		{
			"Fetch all cars",
			"",
			[]model.Car{car1, car2},
			nil,
		}, {
			"Server error in fetching car",
			"Jaguar",
			nil,
			errors.New("server error"),
		},
	}

	svc := New(m, s)

	for i, tc := range tests {
		cars, err := svc.GetAll(tc.brand, true)

		assert.Equalf(t, tc.err, err, "Testcase[%v] (%v)", i, tc.desc)

		assert.Equalf(t, tc.cars, cars, "Testcase[%v] (%v)", i, tc.desc)
	}
}

func TestService_GetByID(t *testing.T) {
	car1 := car1()

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	c := mocks.NewMockCarStore(mockCtrl)
	e := mocks.NewMockEngineStore(mockCtrl)

	c.EXPECT().GetByID("1").Return(&model.Car{
		ID:                "1",
		Name:              "Roadster",
		YearOfManufacture: 2000,
		Brand:             "Tesla",
		FuelType:          "Electric",
		Engine:            model.Engine{ID: "1"},
	}, nil)

	e.EXPECT().GetByID("1").Return(&model.Engine{
		ID:            "1",
		Displacement:  0,
		NoOfCylinders: 0,
		Range:         400,
	}, nil)

	c.EXPECT().GetByID("2").Return(nil, errors.New("server error"))

	c.EXPECT().GetByID("3").Return(nil, customErrors.CarNotExists())

	tests := []struct {
		desc string
		id   string
		car  *model.Car
		err  error
	}{
		{"Success", "1", &car1, nil},
		{"Server error", "2", nil, errors.New("server error")},
		{"Car not exists", "3", nil, customErrors.CarNotExists()},
	}

	svc := New(c, e)

	for i, tc := range tests {
		car, err := svc.GetByID(tc.id)

		assert.Equalf(t, tc.err, err, "Testcase[%v] (%v)", i, tc.desc)

		assert.Equalf(t, tc.car, car, "Testcase[%v] (%v)", i, tc.desc)
	}
}

func TestService_Create(t *testing.T) {
	car1 := car1()

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	c := mocks.NewMockCarStore(mockCtrl)
	e := mocks.NewMockEngineStore(mockCtrl)

	e.EXPECT().Create(&model.Engine{
		Displacement:  0,
		NoOfCylinders: 0,
		Range:         400,
	}).Return(&model.Engine{
		ID:            "1",
		Displacement:  0,
		NoOfCylinders: 0,
		Range:         400,
	}, nil)

	c.EXPECT().Create(&model.Car{
		Name:              "Roadster",
		YearOfManufacture: 2000,
		Brand:             "Tesla",
		FuelType:          "Electric",
		Engine: model.Engine{
			ID:            "1",
			Displacement:  0,
			NoOfCylinders: 0,
			Range:         400,
		},
	}).Return(&model.Car{
		ID:                "1",
		Name:              "Roadster",
		YearOfManufacture: 2000,
		Brand:             "Tesla",
		FuelType:          "Electric",
		Engine:            model.Engine{ID: "1"},
	}, nil)

	e.EXPECT().Create(&model.Engine{}).Return(nil, errors.New("server error"))

	tests := []struct {
		desc  string
		input *model.Car
		car   *model.Car
		err   error
	}{
		{
			"Success",
			&model.Car{
				Name:              "Roadster",
				YearOfManufacture: 2000,
				Brand:             "Tesla",
				FuelType:          "Electric",
				Engine:            model.Engine{Range: 400},
			},
			&car1,
			nil},
		{
			"Server error",
			&model.Car{},
			nil,
			errors.New("server error")},
	}

	svc := New(c, e)

	for i, tc := range tests {
		car, err := svc.Create(tc.input)

		assert.Equalf(t, tc.err, err, "Testcase[%v] (%v)", i, tc.desc)

		assert.Equalf(t, tc.car, car, "Testcase[%v] (%v)", i, tc.desc)
	}
}

func TestService_Update(t *testing.T) {
	car1 := car1()
	car2 := car2()

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	c := mocks.NewMockCarStore(mockCtrl)
	e := mocks.NewMockEngineStore(mockCtrl)

	c.EXPECT().GetByID(car1.ID).Return(&car1, nil)
	c.EXPECT().Update(&car1).Return(&model.Car{
		ID:                "1",
		Name:              "Roadster",
		YearOfManufacture: 2000,
		Brand:             "Tesla",
		FuelType:          "Electric",
		Engine:            model.Engine{ID: "1"},
	}, nil)

	e.EXPECT().Update(&model.Engine{
		ID:            "1",
		Displacement:  0,
		NoOfCylinders: 0,
		Range:         400,
	}).Return(&model.Engine{
		ID:            "1",
		Displacement:  0,
		NoOfCylinders: 0,
		Range:         400,
	}, nil)

	c.EXPECT().GetByID(car2.ID).Return(&car2, nil)
	c.EXPECT().Update(&car2).Return(nil, errors.New("server error"))

	tests := []struct {
		desc  string
		input *model.Car
		car   *model.Car
		err   error
	}{
		{
			"Success",
			&car1,
			&car1,
			nil},
		{
			"Server error",
			&car2,
			nil,
			errors.New("server error"),
		},
	}

	svc := New(c, e)

	for i, tc := range tests {
		car, err := svc.Update(tc.input)

		assert.Equalf(t, tc.err, err, "Testcase[%v] (%v)", i, tc.desc)

		assert.Equalf(t, tc.car, car, "Testcase[%v] (%v)", i, tc.desc)
	}
}

func TestService_Delete(t *testing.T) {
	car1 := car1()

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	c := mocks.NewMockCarStore(mockCtrl)
	e := mocks.NewMockEngineStore(mockCtrl)

	c.EXPECT().GetByID("1").Return(&car1, nil)
	c.EXPECT().Delete("1").Return(nil)
	e.EXPECT().Delete("1").Return(nil)

	c.EXPECT().GetByID("2").Return(&model.Car{ID: "2"}, nil)
	c.EXPECT().Delete("2").Return(customErrors.CarNotExists())

	c.EXPECT().GetByID("3").Return(&model.Car{ID: "3"}, nil)
	c.EXPECT().Delete("3").Return(errors.New("server error"))

	tests := []struct {
		desc string
		id   string
		err  error
	}{
		{"Success", "1", nil},
		{"Car not exists", "2", customErrors.CarNotExists()},
		{"Server error", "3", errors.New("server error")},
	}

	svc := New(c, e)

	for i, tc := range tests {
		err := svc.Delete(tc.id)

		assert.Equalf(t, tc.err, err, "Testcase[%v] (%v)", i, tc.desc)
	}
}
