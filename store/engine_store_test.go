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

var engine = model.Engine{
	ID:            uuid.NewString(),
	Displacement:  0,
	NoOfCylinders: 0,
	Range:         400,
}

var engineNotExists customErrors.EngineNotExists

func TestEngineStore_GetByID(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		log.Println(err)
	}

	defer db.Close()

	store := NewEngineStore(db)
	rows := sqlmock.NewRows([]string{"engineId", "displacement", "noOfCylinder", "range"}).
		AddRow(engine.ID, engine.Displacement, engine.NoOfCylinders, engine.Range)

	mock.ExpectQuery("select \\* from engines where engineId = \\?").WithArgs(engine.ID).WillReturnRows(rows)
	mock.ExpectQuery("select \\* from engines where engineId = \\?").WithArgs("1").WillReturnError(engineNotExists)
	mock.ExpectQuery("select \\* from engines where engineId = \\?").WithArgs("2").WillReturnError(errors.New("DB error"))

	tests := []struct {
		desc   string
		id     string
		engine model.Engine
		err    error
	}{
		{"Success", engine.ID, engine, nil},
		{"Not exists", "1", model.Engine{}, engineNotExists},
		{"DB error", "2", model.Engine{}, errors.New("DB error")},
	}

	for i, tc := range tests {
		engine, err := store.GetByID(tc.id)
		if !reflect.DeepEqual(err, tc.err) {
			t.Errorf("Testcase[%v] failed (%v)\nExpected error: %v\nGot: %v", i, tc.desc, tc.err, err)
		}

		if !reflect.DeepEqual(engine, tc.engine) {
			t.Errorf("Testcase[%v] failed (%v)\nExpected:\n%v\nGot:\n%v", i, tc.desc, tc.engine, engine)
		}
	}
}

func TestEngineStore_Create(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		log.Println(err)
	}

	defer db.Close()

	store := NewEngineStore(db)

	query := "insert into engines \\(engineId, displacement, noOfCylinder, `range`\\) values \\(\\?, \\?, \\?, \\?\\)"

	prep := mock.ExpectPrepare(query)
	prep.ExpectExec().WithArgs(sqlmock.AnyArg(), engine.Displacement, engine.NoOfCylinders, engine.Range).
		WillReturnResult(sqlmock.NewResult(0, 1))
	prep = mock.ExpectPrepare(query)
	prep.ExpectExec().WillReturnError(errors.New("DB error"))

	tests := []struct {
		desc     string
		input    model.Engine
		expected model.Engine
		err      error
	}{
		{"Success", engine, engine, nil},
		{"DB error", model.Engine{}, model.Engine{}, errors.New("DB error")},
	}

	for i, tc := range tests {
		engine, err := store.Create(tc.input)
		if !reflect.DeepEqual(err, tc.err) {
			t.Errorf("Testcase[%v] failed (%v)\nExpected error: %v\nGot: %v", i, tc.desc, tc.err, err)

			if engine.ID == tc.expected.ID {
				t.Errorf("Testcase[%v] failed (%v)\nNew ID was not assigned", i, tc.desc)
			}
		}

		engine.ID = tc.expected.ID

		if !reflect.DeepEqual(engine, tc.expected) {
			t.Errorf("Testcase[%v] failed (%v)\nExpected:\n%v\nGot:\n%v", i, tc.desc, tc.expected, engine)
		}
	}
}

func TestEngineStore_Update(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		log.Println(err)
	}

	defer db.Close()

	store := NewEngineStore(db)

	query := "update engines set displacement = \\?, noOfCylinder = \\?, `range` = \\? where engineId = \\?"

	prep := mock.ExpectPrepare(query)
	prep.ExpectExec().WithArgs(0, 0, 400, engine.ID).WillReturnResult(sqlmock.NewResult(0, 1))
	prep = mock.ExpectPrepare(query)
	prep.ExpectExec().WillReturnError(engineNotExists)
	prep = mock.ExpectPrepare(query)
	prep.ExpectExec().WillReturnError(errors.New("DB error"))

	tests := []struct {
		desc     string
		input    model.Engine
		expected model.Engine
		err      error
	}{
		{"Success", engine, engine, nil},
		{"Not exists", model.Engine{}, model.Engine{}, engineNotExists},
		{"DB error", model.Engine{}, model.Engine{}, errors.New("DB error")},
	}

	for i, tc := range tests {
		engine, err := store.Update(tc.input)
		if !reflect.DeepEqual(err, tc.err) {
			t.Errorf("Testcase[%v] failed (%v)\nExpected error: %v\nGot: %v", i, tc.desc, tc.err, err)
		}

		if !reflect.DeepEqual(engine, tc.expected) {
			t.Errorf("Testcase[%v] failed (%v)\nExpected:\n%v\nGot:\n%v", i, tc.desc, tc.expected, engine)
		}
	}
}

func TestEngineStore_Delete(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		log.Println(err)
	}

	defer db.Close()

	store := NewEngineStore(db)

	query := "delete from engines where engineId = \\?"

	prep := mock.ExpectPrepare(query)
	prep.ExpectExec().WithArgs(engine.ID).WillReturnResult(sqlmock.NewResult(0, 1))
	prep = mock.ExpectPrepare(query)
	prep.ExpectExec().WithArgs(engine.ID).WillReturnError(engineNotExists)
	prep = mock.ExpectPrepare(query)
	prep.ExpectExec().WillReturnError(errors.New("DB error"))

	tests := []struct {
		desc string
		id   string
		err  error
	}{
		{"Success", engine.ID, nil},
		{"Not exists", engine.ID, engineNotExists},
		{"DB error", "", errors.New("DB error")},
	}

	for i, tc := range tests {
		err := store.Delete(tc.id)
		if !reflect.DeepEqual(err, tc.err) {
			t.Errorf("Testcase[%v] failed (%v)\nExpected error: %v\nGot: %v", i, tc.desc, tc.err, err)
		}
	}
}
