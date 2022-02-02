package handler

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/wI2L/jsondiff"

	customErrors "carAPI/custom-errors"
	"carAPI/mocks"
	"carAPI/model"
)

func car1() *model.Car {
	return &model.Car{
		ID:                "1",
		Name:              "Roadster",
		YearOfManufacture: 2000,
		Brand:             "Tesla",
		FuelType:          "Electric",
		Engine: model.Engine{
			ID:            "1",
			Displacement:  0,
			NoOfCylinders: 0,
			Range:         500,
		},
	}
}

func car2() *model.Car {
	return &model.Car{
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

func car3() *model.Car {
	return &model.Car{
		ID:                "1",
		Name:              "Roadster",
		YearOfManufacture: 2000,
		Brand:             "Tesla",
		FuelType:          "Electric",
		Engine:            model.Engine{},
	}
}

func car4() *model.Car {
	return &model.Car{
		ID:                "2",
		Name:              "Abc",
		YearOfManufacture: 2020,
		Brand:             "Ferrari",
		FuelType:          "Diesel",
		Engine:            model.Engine{},
	}
}

func TestHandler_Get(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	m := mocks.NewMockCarService(mockCtrl)

	m.EXPECT().GetAll("Tesla", true).Return([]model.Car{*car1()}, nil)
	m.EXPECT().GetAll("", false).Return([]model.Car{*car3(), *car4()}, nil)
	m.EXPECT().GetAll("Tesla", false).Return([]model.Car{*car3()}, nil)
	m.EXPECT().GetAll("", true).Return([]model.Car{*car1(), *car2()}, nil)
	m.EXPECT().GetAll("BMW", false).Return(nil, errors.New("server error"))

	tests := []struct {
		desc       string
		params     string
		statusCode int
		resp       []byte
	}{
		{
			"Fetch Tesla cars with engine",
			"?brand=Tesla&withEngine=true",
			http.StatusOK,
			[]byte(`[{"carId":"1","name":"Roadster","yearOfManufacture":2000,"brand":"Tesla","fuelType":"Electric",
							"engine":{"engineId":"1","displacement":0,"noOfCylinders":0,"range":500}}]`),
		},

		{
			"Fetch all cars without engine",
			"?brand=&withEngine=false",
			http.StatusOK,
			[]byte(`[{"carId":"1","name":"Roadster","yearOfManufacture":2000,"brand":"Tesla","fuelType":"Electric",
							"engine":{}},{"carId":"2","name":"Abc","yearOfManufacture":2020,"brand":"Ferrari",
							"fuelType":"Diesel","engine":{}}]`),
		},
		{
			"Fetch Tesla cars without engine",
			"?brand=Tesla",
			http.StatusOK,
			[]byte(`[{"carId":"1","name":"Roadster","yearOfManufacture":2000,"brand":"Tesla","fuelType":"Electric",
							"engine":{}}]`),
		},
		{
			"Fetch all cars with engine",
			"?withEngine=true",
			http.StatusOK,
			[]byte(`[{"carId":"1","name":"Roadster","yearOfManufacture":2000,"brand":"Tesla","fuelType":"Electric",
							"engine":{"engineId":"1","displacement":0,"noOfCylinders":0,"range":500}},
							{"carId":"2","name":"Abc","yearOfManufacture":2020,"brand":"Ferrari","fuelType":"Diesel",
							"engine":{"engineId":"2","displacement":600,"noOfCylinders":4,"range":0}}]`),
		},
		{
			"Server error",
			"?brand=BMW",
			http.StatusInternalServerError,
			[]byte(`{"error":{"code":"DB error"}}`),
		},
		{
			"Invalid value of withEngine",
			"?withEngine=abc",
			http.StatusBadRequest,
			[]byte(`{"error":{"code":"invalid value of withEngine","message":"withEngine must be true or false"}}`),
		},
	}

	for i, tc := range tests {
		h := New(m)
		r := httptest.NewRequest(http.MethodGet, "/car"+tc.params, nil)
		w := httptest.NewRecorder()
		h.Get(w, r)
		result := w.Result()
		body, _ := io.ReadAll(result.Body)

		result.Body.Close()

		if result.StatusCode != tc.statusCode {
			t.Errorf("Testcase[%v] failed (%v)\nExpected status %v\tGot %v", i, tc.desc, tc.statusCode, result.StatusCode)
		}

		_, err := jsondiff.CompareJSON(tc.resp, body)
		if err != nil {
			t.Errorf("Testcase[%v] failed (%v)\nExpected:\n%v\nGot:\n%v", i, tc.desc, string(tc.resp), string(body))
		}
	}
}

func TestHandler_GetById(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	m := mocks.NewMockCarService(mockCtrl)

	m.EXPECT().GetByID("1").Return(car1(), nil)
	m.EXPECT().GetByID("2").Return(&model.Car{}, customErrors.CarNotExists())
	m.EXPECT().GetByID("3").Return(&model.Car{}, errors.New("server error"))

	tests := []struct {
		desc       string
		id         string
		statusCode int
		resp       []byte
	}{
		{
			"Success",
			"1",
			http.StatusOK,
			[]byte(`{"carId":"1","name":"Roadster","yearOfManufacture":2000,"brand":"Tesla","fuelType":"Electric",
							"engine":{"engineId":"1","displacement":0,"noOfCylinders":0,"range":500}}`),
		},

		{
			"Car not exists",
			"2",
			http.StatusNotFound,
			[]byte(`{"error":{"code":"entity not found", "id":"2"}}`),
		},

		{
			"Server Error",
			"3",
			http.StatusInternalServerError,
			[]byte(`{"error":{"code":"DB error"}}`),
		},
	}

	h := New(m)

	for i, tc := range tests {
		r := httptest.NewRequest(http.MethodGet, "/car", nil)
		m := make(map[string]string)

		m["id"] = tc.id
		r = mux.SetURLVars(r, m)

		w := httptest.NewRecorder()

		h.GetByID(w, r)
		result := w.Result()

		body, _ := io.ReadAll(result.Body)

		result.Body.Close()

		if result.StatusCode != tc.statusCode {
			t.Errorf("Testcase[%v] failed (%v)\nExpected status %v\tGot %v", i, tc.desc, tc.statusCode, result.StatusCode)
		}

		_, err := jsondiff.CompareJSON(tc.resp, body)
		if err != nil {
			t.Errorf("Testcase[%v] failed (%v)\nExpected:\n%v\nGot:\n%v", i, tc.desc, string(tc.resp), string(body))
		}
	}
}

func TestHandler_Create(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	m := mocks.NewMockCarService(mockCtrl)

	m.EXPECT().Create(car1()).Return(car1(), nil)
	m.EXPECT().Create(car2()).Return(&model.Car{}, errors.New("server error"))

	tests := []struct {
		desc       string
		body       io.Reader
		statusCode int
		resp       []byte
	}{
		{
			"Success",
			bytes.NewReader([]byte(`{"carId":"1","name":"Roadster","yearOfManufacture":2000,"brand":"Tesla","fuelType":"Electric",
							"engine":{"engineId":"1","displacement":0,"noOfCylinders":0,"range":500}}`)),
			http.StatusCreated,
			[]byte(`{"carId":"1","name":"Roadster","yearOfManufacture":2000,"brand":"Tesla","fuelType":"Electric",
							"engine":{"engineId":"1","displacement":0,"noOfCylinders":0,"range":500}}`),
		},
		{
			"Server Error",
			bytes.NewReader([]byte(`{"carId":"2","name":"Abc","yearOfManufacture":2020,"brand":"Ferrari","fuelType":"Diesel",
							"engine":{"engineId":"2","displacement":600,"noOfCylinders":4,"range":0}}`)),
			http.StatusInternalServerError,
			[]byte(`{"error":{"code":"DB error"}}`),
		},
		{
			"Unmarshal Error",
			bytes.NewReader([]byte("Invalid")),
			http.StatusBadRequest,
			[]byte(`{"error":{"code":"invalid body", "message":"cannot parse given body"}}`),
		},
		{
			"Validation Error",
			bytes.NewReader([]byte("{}")),
			http.StatusBadRequest,
			[]byte(`{"error":{"code":"missing param(s)", "requiredParams":["name", "yearOfManufacture","brand",
							"fuelType", "engine"],"engineParams":"either range or displacement and noOfCylinders must be passed"}}`),
		},
		{
			"Invalid Year",
			bytes.NewReader([]byte(`{"name":"Roadster","yearOfManufacture":2100,"brand":"Tesla","fuelType":"Electric",
							"engine":{"range":400}}`)),
			http.StatusBadRequest,
			[]byte(`{"error":{"code":"invalid body", "message":"invalid year of manufacture"}}`),
		},
		{
			"Invalid Brand",
			bytes.NewReader([]byte(`{"name":"Roadster","yearOfManufacture":2000,"brand":"Pesla","fuelType":"Electric",
							"engine":{"range":400}}`)),
			http.StatusBadRequest,
			[]byte(`{"error":{"code":"invalid body", "message":"supported brands are Tesla, Ferrari, Porsche and BMW"}}`),
		},
		{
			"Invalid FuelType",
			bytes.NewReader([]byte(`{"name":"Roadster","yearOfManufacture":200,"brand":"Tesla","fuelType":"CNG",
							"engine":{"range":400}}`)),
			http.StatusBadRequest,
			[]byte(`{"error":{"code":"invalid body", "message":"fuelType must be Electric, Petrol or Diesel"}}`),
		},
	}

	h := New(m)

	for i, tc := range tests {
		r := httptest.NewRequest(http.MethodGet, "/car", tc.body)
		w := httptest.NewRecorder()
		h.Create(w, r)
		result := w.Result()

		body, _ := io.ReadAll(result.Body)

		result.Body.Close()

		if result.StatusCode != tc.statusCode {
			t.Errorf("Testcase[%v] failed (%v)\nExpected status %v\tGot %v", i, tc.desc, tc.statusCode, result.StatusCode)
		}

		_, err := jsondiff.CompareJSON(tc.resp, body)
		if err != nil {
			t.Errorf("Testcase[%v] failed (%v)\nExpected:\n%v\nGot:\n%v", i, tc.desc, string(tc.resp), string(body))
		}
	}
}

func TestHandler_Update(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	m := mocks.NewMockCarService(mockCtrl)

	m.EXPECT().Update(car1()).Return(car1(), nil)
	m.EXPECT().Update(car2()).Return(&model.Car{}, errors.New("server error"))

	tests := []struct {
		desc       string
		id         string
		body       io.Reader
		statusCode int
		resp       []byte
	}{
		{
			"Success",
			"1",
			bytes.NewReader([]byte(`{"name":"Roadster","yearOfManufacture":2000,"brand":"Tesla","fuelType":"Electric",
							"engine":{"engineId":"1","range":500}}`)),
			http.StatusOK,
			[]byte(`{"carId":"1","name":"Roadster","yearOfManufacture":2000,"brand":"Tesla","fuelType":"Electric",
							"engine":{"engineId":"1","displacement":0,"noOfCylinders":0,"range":500}}`),
		},
		{
			"Server Error",
			"2",
			bytes.NewReader([]byte(`{"carId":"2","name":"Abc","yearOfManufacture":2020,"brand":"Ferrari","fuelType":"Diesel",
							"engine":{"engineId":"2","displacement":600,"noOfCylinders":4,"range":0}}`)),
			http.StatusInternalServerError,
			[]byte(`{"error":{"code":"DB error"}}`),
		},
		{
			"Unmarshal Error",
			"1",
			bytes.NewReader([]byte("Invalid")),
			http.StatusBadRequest,
			[]byte(`{"error":{"code":"invalid body", "message":"cannot parse given body"}}`),
		},
		{
			"Validation Error",
			"1",
			bytes.NewReader([]byte("{}")),
			http.StatusBadRequest,
			[]byte(`{"error":{"code":"missing param(s)", "requiredParams":["name", "yearOfManufacture","brand",
							"fuelType", "engine"],"engineParams":"either range or displacement and noOfCylinders must be passed"}}`),
		},
		{
			"Invalid Year",
			"1",
			bytes.NewReader([]byte(`{"name":"Roadster","yearOfManufacture":2100,"fuelType":"Electric",
							"engine":{"range":400}}`)),
			http.StatusBadRequest,
			[]byte(`{"error":{"code":"invalid body", "message":"invalid year of manufacture"}}`),
		},
		{
			"Invalid Brand",
			"1",
			bytes.NewReader([]byte(`{"name":"Roadster","yearOfManufacture":2000,"brand":"Pesla","fuelType":"Electric",
							"engine":{"range":400}}`)),
			http.StatusBadRequest,
			[]byte(`{"error":{"code":"invalid body", "message":"supported brands are Tesla, Ferrari, Porsche and BMW"}}`),
		},
		{
			"Invalid FuelType",
			"1",
			bytes.NewReader([]byte(`{"name":"Roadster","yearOfManufacture":200,"brand":"Tesla","fuelType":"CNG",
							"engine":{"range":400}}`)),
			http.StatusBadRequest,
			[]byte(`{"error":{"code":"invalid body", "message":"fuelType must be Electric, Petrol or Diesel"}}`),
		},
	}

	h := New(m)

	for i, tc := range tests {
		r := httptest.NewRequest(http.MethodGet, "/car", tc.body)
		w := httptest.NewRecorder()
		m := make(map[string]string)

		m["id"] = tc.id

		r = mux.SetURLVars(r, m)

		h.Update(w, r)

		result := w.Result()
		body, _ := io.ReadAll(result.Body)

		result.Body.Close()

		if result.StatusCode != tc.statusCode {
			t.Errorf("Testcase[%v] failed (%v)\nExpected status %v\tGot %v", i, tc.desc, tc.statusCode, result.StatusCode)
		}

		_, err := jsondiff.CompareJSON(tc.resp, body)
		if err != nil {
			t.Errorf("Testcase[%v] failed (%v)\nExpected:\n%v\nGot:\n%v", i, tc.desc, string(tc.resp), string(body))
		}
	}
}

func TestHandler_Delete(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	m := mocks.NewMockCarService(mockCtrl)

	m.EXPECT().Delete("1").Return(nil)
	m.EXPECT().Delete("2").Return(customErrors.CarNotExists())
	m.EXPECT().Delete("3").Return(errors.New("server error"))

	tests := []struct {
		desc       string
		id         string
		statusCode int
		resp       []byte
	}{
		{
			"Success",
			"1",
			http.StatusNoContent,
			[]byte(""),
		},
		{
			"Car not exists",
			"2",
			http.StatusNotFound,
			[]byte(`{"error":{"code":"entity not found","id":"2"}}`),
		},
		{
			"Server Error",
			"3",
			http.StatusInternalServerError,
			[]byte(`{"error":{"code":"DB error"}}`),
		},
	}

	h := New(m)

	for i, tc := range tests {
		r := httptest.NewRequest(http.MethodGet, "/car", nil)
		w := httptest.NewRecorder()
		m := make(map[string]string)

		m["id"] = tc.id
		r = mux.SetURLVars(r, m)

		h.Delete(w, r)

		result := w.Result()
		body, _ := io.ReadAll(result.Body)

		result.Body.Close()

		assert.Equalf(t, tc.statusCode, result.StatusCode, "Testcase[%v] (%v)", i, tc.desc)

		assert.Equalf(t, tc.resp, body, "Testcase[%v] (%v)", i, tc.desc)
	}
}
