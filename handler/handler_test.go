package handler

import (
	"bytes"
	customErrors "carAPI/custom-errors"
	"carAPI/mocks"
	"carAPI/model"
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"github.com/wI2L/jsondiff"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

var carNotExists customErrors.CarNotExists

func TestHandler_HandleGetAll(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	m := mocks.NewMockCarService(mockCtrl)
	gomock.InOrder(
		m.EXPECT().GetAll("Tesla", true).Return([]model.Car{{
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
		}}, nil),
		m.EXPECT().GetAll("", false).Return([]model.Car{{
			ID:                "1",
			Name:              "Roadster",
			YearOfManufacture: 2000,
			Brand:             "Tesla",
			FuelType:          "Electric",
			Engine:            model.Engine{},
		}, {
			ID:                "2",
			Name:              "Abc",
			YearOfManufacture: 2020,
			Brand:             "Ferrari",
			FuelType:          "Diesel",
			Engine:            model.Engine{},
		}}, nil),
		m.EXPECT().GetAll("Tesla", false).Return([]model.Car{{
			ID:                "1",
			Name:              "Roadster",
			YearOfManufacture: 2000,
			Brand:             "Tesla",
			FuelType:          "Electric",
			Engine:            model.Engine{},
		}}, nil),
		m.EXPECT().GetAll("", true).Return([]model.Car{{
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
		}, {
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
		}}, nil),
		m.EXPECT().GetAll("", false).Return(nil, errors.New("server error")),
	)

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
							"engine":{"engineId":"1","displacement":0,"noOCylinders":0,"range":500}}]`),
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
							"engine":{"engineId":"1","displacement":0,"noOCylinders":0,"range":500}},
							{"carId":"2","name":"Abc","yearOfManufacture":2020,"brand":"Ferrari","fuelType":"Diesel",
							"engine":{"engineId":"2","displacement":600,"noOCylinders":4,"range":0}}]`),
		},
		{
			"Server error",
			"",
			http.StatusInternalServerError,
			[]byte(`{"error":{"code":"DB error"}}`),
		},
	}

	for i, tc := range tests {
		h := New(m)
		r := httptest.NewRequest(http.MethodGet, "/user"+tc.params, nil)
		w := httptest.NewRecorder()
		h.HandleGetAll(w, r)
		result := w.Result()
		body, _ := io.ReadAll(result.Body)

		if result.StatusCode != tc.statusCode {
			t.Errorf("Testcase[%v] failed (%v)\nExpected status %v\tGot %v", i, tc.desc, tc.statusCode, result.StatusCode)
		}

		_, err := jsondiff.CompareJSON(tc.resp, body)
		if err != nil {
			t.Errorf("Testcase[%v] failed (%v)\nExpected:\n%v\nGot:\n%v", i, tc.desc, string(tc.resp), string(body))
		}
	}
}

func TestHandler_HandleGetById(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	m := mocks.NewMockCarService(mockCtrl)
	gomock.InOrder(
		m.EXPECT().GetByID("1").Return(model.Car{
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
		}, nil),
		m.EXPECT().GetByID("2").Return(model.Car{}, carNotExists),
		m.EXPECT().GetByID("3").Return(model.Car{}, errors.New("server error")),
	)

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
							"engine":{"engineId":"1","displacement":0,"noOCylinders":0,"range":500}}`),
		},

		{
			"User not exists",
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
		r := httptest.NewRequest(http.MethodGet, "/user", nil)
		m := make(map[string]string)

		m["id"] = tc.id
		r = mux.SetURLVars(r, m)

		w := httptest.NewRecorder()

		h.HandleGetByID(w, r)
		result := w.Result()

		body, _ := io.ReadAll(result.Body)

		if result.StatusCode != tc.statusCode {
			t.Errorf("Testcase[%v] failed (%v)\nExpected status %v\tGot %v", i, tc.desc, tc.statusCode, result.StatusCode)
		}

		_, err := jsondiff.CompareJSON(tc.resp, body)
		if err != nil {
			t.Errorf("Testcase[%v] failed (%v)\nExpected:\n%v\nGot:\n%v", i, tc.desc, string(tc.resp), string(body))
		}
	}
}

func TestHandler_HandleCreate(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	m := mocks.NewMockCarService(mockCtrl)

	gomock.InOrder(
		m.EXPECT().Create(model.Car{
			Name:              "Roadster",
			YearOfManufacture: 2000,
			Brand:             "Tesla",
			FuelType:          "Electric",
			Engine: model.Engine{
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
			Engine: model.Engine{
				ID:            "1",
				Displacement:  0,
				NoOfCylinders: 0,
				Range:         400,
			},
		}, nil),
		m.EXPECT().Create(model.Car{
			Name:              "Roadster",
			YearOfManufacture: 2000,
			Brand:             "Tesla",
			FuelType:          "Electric",
			Engine: model.Engine{
				Displacement:  0,
				NoOfCylinders: 0,
				Range:         400,
			},
		}).Return(model.Car{}, errors.New("server error")),
	)

	tests := []struct {
		desc       string
		body       io.Reader
		statusCode int
		resp       []byte
	}{
		{
			"Success",
			bytes.NewReader([]byte(`{"name":"Roadster","yearOfManufacture":2000,"brand":"Tesla","fuelType":"Electric",
							"engine":{"range":400}}`)),
			http.StatusCreated,
			[]byte(`{"carId":"1","name":"Roadster","yearOfManufacture":2000,"brand":"Tesla","fuelType":"Electric",
							"engine":{"engineId":"1","displacement":0,"noOCylinders":0,"range":400}}`),
		},
		{
			"Server Error",
			bytes.NewReader([]byte(`{"name":"Roadster","yearOfManufacture":2000,"brand":"Tesla","fuelType":"Electric",
							"engine":{"range":400}}`)),
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
		h.HandleCreate(w, r)
		result := w.Result()
		body, _ := io.ReadAll(result.Body)

		if result.StatusCode != tc.statusCode {
			t.Errorf("Testcase[%v] failed (%v)\nExpected status %v\tGot %v", i, tc.desc, tc.statusCode, result.StatusCode)
		}

		_, err := jsondiff.CompareJSON(tc.resp, body)
		if err != nil {
			t.Errorf("Testcase[%v] failed (%v)\nExpected:\n%v\nGot:\n%v", i, tc.desc, string(tc.resp), string(body))
		}
	}
}

func TestHandler_HandleUpdate(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	m := mocks.NewMockCarService(mockCtrl)

	gomock.InOrder(
		m.EXPECT().Update(model.Car{
			ID:                "1",
			Name:              "Roadster",
			YearOfManufacture: 2000,
			Brand:             "Tesla",
			FuelType:          "Electric",
			Engine: model.Engine{
				Displacement:  0,
				NoOfCylinders: 0,
				Range:         450,
			},
		}).Return(model.Car{
			ID:                "1",
			Name:              "Roadster",
			YearOfManufacture: 2000,
			Brand:             "Tesla",
			FuelType:          "Electric",
			Engine: model.Engine{
				ID:            "1",
				Displacement:  0,
				NoOfCylinders: 0,
				Range:         450,
			},
		}, nil),
		m.EXPECT().Update(model.Car{
			ID:                "2",
			Name:              "Roadster",
			YearOfManufacture: 2000,
			Brand:             "Tesla",
			FuelType:          "Electric",
			Engine: model.Engine{
				Displacement:  0,
				NoOfCylinders: 0,
				Range:         400,
			},
		}).Return(model.Car{}, errors.New("server error")),
	)

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
							"engine":{"range":450}}`)),
			http.StatusCreated,
			[]byte(`{"carId":"1","name":"Roadster","yearOfManufacture":2000,"brand":"Tesla","fuelType":"Electric",
							"engine":{"engineId":"1","displacement":0,"noOCylinders":0,"range":450}}`),
		},
		{
			"Server Error",
			"2",
			bytes.NewReader([]byte(`{"name":"Roadster","yearOfManufacture":2000,"brand":"Tesla","fuelType":"Electric",
							"engine":{"range":400}}`)),
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
			bytes.NewReader([]byte(`{"name":"Roadster","yearOfManufacture":2100,"brand":"Tesla","fuelType":"Electric",
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

		h.HandleUpdate(w, r)

		result := w.Result()
		body, _ := io.ReadAll(result.Body)

		if result.StatusCode != tc.statusCode {
			t.Errorf("Testcase[%v] failed (%v)\nExpected status %v\tGot %v", i, tc.desc, tc.statusCode, result.StatusCode)
		}

		_, err := jsondiff.CompareJSON(tc.resp, body)
		if err != nil {
			t.Errorf("Testcase[%v] failed (%v)\nExpected:\n%v\nGot:\n%v", i, tc.desc, string(tc.resp), string(body))
		}
	}
}

func TestHandler_HandleDelete(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	m := mocks.NewMockCarService(mockCtrl)

	gomock.InOrder(
		m.EXPECT().Delete("1").Return(nil),
		m.EXPECT().Delete("2").Return(carNotExists),
		m.EXPECT().Delete("3").Return(errors.New("server error")),
	)

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
			"User not exists",
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
		r := httptest.NewRequest(http.MethodGet, "/user", nil)
		w := httptest.NewRecorder()
		m := make(map[string]string)

		m["id"] = tc.id
		r = mux.SetURLVars(r, m)

		h.HandleDelete(w, r)

		result := w.Result()
		body, _ := io.ReadAll(result.Body)

		if result.StatusCode != tc.statusCode {
			t.Errorf("Testcase[%v] failed (%v)\nExpected status %v\tGot %v", i, tc.desc, tc.statusCode, result.StatusCode)
		}

		if !reflect.DeepEqual(tc.resp, body) {
			t.Errorf("Testcase[%v] failed (%v)\nExpected:\n%v\nGot:\n%v", i, tc.desc, string(tc.resp), string(body))
		}
	}
}
