package service

import (
	customErrors "carAPI/custom-errors"
	"carAPI/mocks"
	"carAPI/model"
	"errors"
	"github.com/golang/mock/gomock"
	"reflect"
	"testing"
)

var car1 = model.Car{

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

var car2 = model.Car{

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

var carNotExists customErrors.CarNotExists

func TestService_GetAll(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	m := mocks.NewMockCarStore(mockCtrl)
	s := mocks.NewMockEngineStore(mockCtrl)

	gomock.InOrder(
		m.EXPECT().GetByBrand("Tesla").Return([]model.Car{
			{
				ID:                "1",
				Name:              "Roadster",
				YearOfManufacture: 2000,
				Brand:             "Tesla",
				FuelType:          "Electric",
				Engine:            model.Engine{ID: "1"},
			},
		}, nil),
		s.EXPECT().GetByID("1").Return(model.Engine{
			ID:            "1",
			Displacement:  0,
			NoOfCylinders: 0,
			Range:         400,
		}, nil),
		m.EXPECT().GetByBrand("").Return([]model.Car{
			{
				ID:                "1",
				Name:              "Roadster",
				YearOfManufacture: 2000,
				Brand:             "Tesla",
				FuelType:          "Electric",
				Engine:            model.Engine{ID: "1"},
			},
			{
				ID:                "2",
				Name:              "Abc",
				YearOfManufacture: 2020,
				Brand:             "Ferrari",
				FuelType:          "Diesel",
				Engine:            model.Engine{ID: "2"},
			},
		}, nil),
		s.EXPECT().GetByID("1").Return(model.Engine{
			ID:            "1",
			Displacement:  0,
			NoOfCylinders: 0,
			Range:         400,
		}, nil),
		s.EXPECT().GetByID("2").Return(model.Engine{
			ID:            "2",
			Displacement:  600,
			NoOfCylinders: 4,
			Range:         0,
		}, nil),
		m.EXPECT().GetByBrand("").Return(nil, errors.New("server error")),
		m.EXPECT().GetByBrand("Tesla").Return([]model.Car{
			{
				ID:                "1",
				Name:              "Roadster",
				YearOfManufacture: 2000,
				Brand:             "Tesla",
				FuelType:          "Electric",
				Engine:            model.Engine{ID: "1"},
			},
		}, nil),
		s.EXPECT().GetByID("1").Return(model.Engine{}, errors.New("server error")),
	)

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
			"",
			nil,
			errors.New("server error"),
		},
		{
			"Server error in fetching engine",
			"Tesla",
			nil,
			errors.New("server error"),
		},
	}

	svc := New(m, s)

	for i, tc := range tests {
		cars, err := svc.GetAll(tc.brand, true)
		if !reflect.DeepEqual(err, tc.err) {
			t.Errorf("Testcase[%v] failed (%v)\nExpected error: %v\nGot: %v", i, tc.desc, tc.err, err)
		}

		if !reflect.DeepEqual(cars, tc.cars) {
			t.Errorf("Testcase[%v] failed (%v)\nExpected:\n%v\nGot:\n%v", i, tc.desc, tc.cars, cars)
		}
	}
}

func TestService_GetByID(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	c := mocks.NewMockCarStore(mockCtrl)
	e := mocks.NewMockEngineStore(mockCtrl)

	gomock.InOrder(
		c.EXPECT().GetByID("1").Return(model.Car{
			ID:                "1",
			Name:              "Roadster",
			YearOfManufacture: 2000,
			Brand:             "Tesla",
			FuelType:          "Electric",
			Engine:            model.Engine{ID: "1"},
		}, nil),
		e.EXPECT().GetByID("1").Return(model.Engine{
			ID:            "1",
			Displacement:  0,
			NoOfCylinders: 0,
			Range:         400,
		}, nil),
		c.EXPECT().GetByID("2").Return(model.Car{}, errors.New("server error")),
		c.EXPECT().GetByID("3").Return(model.Car{}, carNotExists),
	)

	tests := []struct {
		desc string
		id   string
		car  model.Car
		err  error
	}{
		{"Success", "1", car1, nil},
		{"Server error", "2", model.Car{}, errors.New("server error")},
		{"Car not exists", "3", model.Car{}, carNotExists},
	}

	svc := New(c, e)

	for i, tc := range tests {
		car, err := svc.GetByID(tc.id)
		if !reflect.DeepEqual(err, tc.err) {
			t.Errorf("Testcase[%v] failed (%v)\nExpected error: %v\nGot: %v", i, tc.desc, tc.err, err)
		}

		if !reflect.DeepEqual(car, tc.car) {
			t.Errorf("Testcase[%v] failed (%v)\nExpected:\n%v\nGot:\n%v", i, tc.desc, tc.car, car)
		}
	}
}

func TestService_Create(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	c := mocks.NewMockCarStore(mockCtrl)
	e := mocks.NewMockEngineStore(mockCtrl)

	gomock.InOrder(
		e.EXPECT().Create(model.Engine{
			Displacement:  0,
			NoOfCylinders: 0,
			Range:         400,
		}).Return(model.Engine{
			ID:            "1",
			Displacement:  0,
			NoOfCylinders: 0,
			Range:         400,
		}, nil),
		c.EXPECT().Create(model.Car{
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
		}).Return(model.Car{
			ID:                "1",
			Name:              "Roadster",
			YearOfManufacture: 2000,
			Brand:             "Tesla",
			FuelType:          "Electric",
			Engine:            model.Engine{ID: "1"},
		}, nil),
		e.EXPECT().Create(model.Engine{}).Return(model.Engine{}, errors.New("server error")),
	)

	tests := []struct {
		desc  string
		input model.Car
		car   model.Car
		err   error
	}{
		{
			"Success",
			model.Car{
				Name:              "Roadster",
				YearOfManufacture: 2000,
				Brand:             "Tesla",
				FuelType:          "Electric",
				Engine:            model.Engine{Range: 400},
			},
			car1,
			nil},
		{
			"Server error",
			model.Car{},
			model.Car{},
			errors.New("server error")},
	}

	svc := New(c, e)

	for i, tc := range tests {
		car, err := svc.Create(tc.input)
		if !reflect.DeepEqual(err, tc.err) {
			t.Errorf("Testcase[%v] failed (%v)\nExpected error: %v\nGot: %v", i, tc.desc, tc.err, err)
		}

		if !reflect.DeepEqual(car, tc.car) {
			t.Errorf("Testcase[%v] failed (%v)\nExpected:\n%v\nGot:\n%v", i, tc.desc, tc.car, car)
		}
	}
}

func TestService_Update(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	c := mocks.NewMockCarStore(mockCtrl)
	e := mocks.NewMockEngineStore(mockCtrl)

	gomock.InOrder(
		e.EXPECT().Update(model.Engine{
			Displacement:  0,
			NoOfCylinders: 0,
			Range:         400,
		}).Return(model.Engine{
			ID:            "1",
			Displacement:  0,
			NoOfCylinders: 0,
			Range:         400,
		}, nil),
		c.EXPECT().Update(model.Car{
			Name:              "Roadster",
			YearOfManufacture: 2000,
			Brand:             "Tesla",
			FuelType:          "Electric",
			Engine:            model.Engine{Range: 400},
		}).Return(model.Car{
			ID:                "1",
			Name:              "Roadster",
			YearOfManufacture: 2000,
			Brand:             "Tesla",
			FuelType:          "Electric",
			Engine:            model.Engine{ID: "1"},
		}, nil),
		e.EXPECT().Update(model.Engine{}).Return(model.Engine{}, errors.New("server error")),
	)

	tests := []struct {
		desc  string
		input model.Car
		car   model.Car
		err   error
	}{
		{
			"Success",
			model.Car{
				Name:              "Roadster",
				YearOfManufacture: 2000,
				Brand:             "Tesla",
				FuelType:          "Electric",
				Engine:            model.Engine{Range: 400},
			},
			car1,
			nil},
		{
			"Server error",
			model.Car{},
			model.Car{},
			errors.New("server error"),
		},
	}

	svc := New(c, e)

	for i, tc := range tests {
		car, err := svc.Update(tc.input)
		if !reflect.DeepEqual(err, tc.err) {
			t.Errorf("Testcase[%v] failed (%v)\nExpected error: %v\nGot: %v", i, tc.desc, tc.err, err)
		}

		if !reflect.DeepEqual(car, tc.car) {
			t.Errorf("Testcase[%v] failed (%v)\nExpected:\n%v\nGot:\n%v", i, tc.desc, tc.car, car)
		}
	}
}

func TestService_Delete(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	c := mocks.NewMockCarStore(mockCtrl)
	e := mocks.NewMockEngineStore(mockCtrl)

	gomock.InOrder(
		c.EXPECT().GetByID("1").Return(car1, nil),
		e.EXPECT().Delete("1").Return(nil),
		c.EXPECT().Delete("1").Return(nil),
		c.EXPECT().GetByID("").Return(model.Car{}, nil),
		e.EXPECT().Delete("").Return(carNotExists),
		c.EXPECT().GetByID("").Return(model.Car{}, nil),
		e.EXPECT().Delete("").Return(errors.New("server error")),
	)

	tests := []struct {
		desc string
		id   string
		err  error
	}{
		{"Success", "1", nil},
		{"Car not exists", "", carNotExists},
		{"Server error", "", errors.New("server error")},
	}

	svc := New(c, e)

	for i, tc := range tests {
		err := svc.Delete(tc.id)
		if !reflect.DeepEqual(err, tc.err) {
			t.Errorf("Testcase[%v] failed (%v)\nExpected error: %v\nGot: %v", i, tc.desc, tc.err, err)
		}
	}
}
