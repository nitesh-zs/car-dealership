package engine

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

func engine() model.Engine {
	return model.Engine{
		ID:            uuid.NewString(),
		Displacement:  0,
		NoOfCylinders: 0,
		Range:         400,
	}
}

func TestEngineStore_GetByID(t *testing.T) {
	engine := engine()

	db, mock, err := sqlmock.New()
	if err != nil {
		log.Println(err)
	}

	defer db.Close()

	store := NewEngineStore(db)
	rows := sqlmock.NewRows([]string{"engineId", "displacement", "noOfCylinder", "range"}).
		AddRow(engine.ID, engine.Displacement, engine.NoOfCylinders, engine.Range)

	mock.ExpectQuery("select \\* from engines where engineId = \\?").WithArgs(engine.ID).WillReturnRows(rows)
	mock.ExpectQuery("select \\* from engines where engineId = \\?").WithArgs("1").WillReturnError(customErrors.EngineNotExists())
	mock.ExpectQuery("select \\* from engines where engineId = \\?").WithArgs("2").WillReturnError(errors.New("DB error"))

	tests := []struct {
		desc   string
		id     string
		engine *model.Engine
		err    error
	}{
		{"Success", engine.ID, &engine, nil},
		{"Not exists", "1", &model.Engine{}, customErrors.EngineNotExists()},
		{"DB error", "2", &model.Engine{}, errors.New("DB error")},
	}

	for i, tc := range tests {
		engine, err := store.GetByID(tc.id)

		assert.Equalf(t, tc.err, err, "Testcase[%v] (%v)", i, tc.desc)

		assert.Equalf(t, tc.engine, engine, "Testcase[%v] (%v)", i, tc.desc)
	}
}

func TestEngineStore_Create(t *testing.T) {
	engine := engine()

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
		input    *model.Engine
		expected *model.Engine
		err      error
	}{
		{"Success", &engine, &engine, nil},
		{"DB error", &model.Engine{}, &model.Engine{}, errors.New("DB error")},
	}

	for i, tc := range tests {
		engine, err := store.Create(tc.input)

		assert.Equalf(t, tc.err, err, "Testcase[%v] (%v)", i, tc.desc)

		engine.ID = tc.expected.ID

		assert.Equalf(t, tc.expected, engine, "Testcase[%v] (%v)", i, tc.desc)
	}
}

func TestEngineStore_Update(t *testing.T) {
	engine := engine()

	db, mock, err := sqlmock.New()
	if err != nil {
		log.Println(err)
	}

	defer db.Close()

	store := NewEngineStore(db)

	query := "update engines set displacement = \\?, noOfCylinder = \\?, `range` = \\? where engineId = \\?"

	// Success case
	rows := sqlmock.NewRows([]string{"engineId", "displacement", "noOfCylinder", "range"}).
		AddRow(engine.ID, engine.Displacement, engine.NoOfCylinders, engine.Range)
	mock.ExpectQuery("select \\* from engines where engineId = \\?").WithArgs(engine.ID).WillReturnRows(rows)
	prep := mock.ExpectPrepare(query)
	prep.ExpectExec().WithArgs(0, 0, 400, engine.ID).WillReturnResult(sqlmock.NewResult(0, 1))

	// CarNotExists error
	mock.ExpectQuery("select \\* from engines where engineId = \\?").WithArgs("").WillReturnError(customErrors.CarNotExists())

	// DB error
	rows.AddRow(engine.ID, engine.Displacement, engine.NoOfCylinders, engine.Range)
	mock.ExpectQuery("select \\* from engines where engineId = \\?").WithArgs(engine.ID).WillReturnRows(rows)
	prep = mock.ExpectPrepare(query)
	prep.ExpectExec().WillReturnError(errors.New("DB error"))

	tests := []struct {
		desc     string
		input    *model.Engine
		expected *model.Engine
		err      error
	}{
		{"Success", &engine, &engine, nil},
		{"Not exists", &model.Engine{}, &model.Engine{}, customErrors.CarNotExists()},
		{"DB error", &engine, &model.Engine{}, errors.New("DB error")},
	}

	for i, tc := range tests {
		engine, err := store.Update(tc.input)

		assert.Equalf(t, tc.err, err, "Testcase[%v] (%v)", i, tc.desc)

		assert.Equalf(t, tc.expected, engine, "Testcase[%v] (%v)", i, tc.desc)
	}
}

func TestEngineStore_Delete(t *testing.T) {
	engine := engine()

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
	prep.ExpectExec().WithArgs(engine.ID).WillReturnError(customErrors.EngineNotExists())

	prep = mock.ExpectPrepare(query)
	prep.ExpectExec().WillReturnError(errors.New("DB error"))

	tests := []struct {
		desc string
		id   string
		err  error
	}{
		{"Success", engine.ID, nil},
		{"Not exists", engine.ID, customErrors.EngineNotExists()},
		{"DB error", "", errors.New("DB error")},
	}

	for i, tc := range tests {
		err := store.Delete(tc.id)

		assert.Equalf(t, tc.err, err, "Testcase[%v] (%v)", i, tc.desc)
	}
}
