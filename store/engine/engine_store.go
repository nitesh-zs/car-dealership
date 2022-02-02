package engine

import (
	"database/sql"

	"github.com/google/uuid"

	customErrors "carAPI/custom-errors"
	"carAPI/model"
)

type engineStore struct {
	db *sql.DB
}

//nolint:revive //engineStore should not be exported
func NewEngineStore(db *sql.DB) engineStore {
	return engineStore{db: db}
}

func (s engineStore) GetByID(id string) (*model.Engine, error) {
	var engine model.Engine

	row := s.db.QueryRow("select * from engines where engineId = ?", id)
	err := row.Scan(&engine.ID, &engine.Displacement, &engine.NoOfCylinders, &engine.Range)

	if err == sql.ErrNoRows {
		return &model.Engine{}, customErrors.CarNotExists()
	}

	if err != nil {
		return &model.Engine{}, err
	}

	return &engine, nil
}

func (s engineStore) Create(engine *model.Engine) (*model.Engine, error) {
	engine.ID = uuid.NewString()
	stmt, err := s.db.Prepare("insert into engines (engineId, displacement, noOfCylinder, `range`) values (?, ?, ?, ?)")

	if err != nil {
		return &model.Engine{}, err
	}

	defer stmt.Close()

	_, err = stmt.Exec(engine.ID, engine.Displacement, engine.NoOfCylinders, engine.Range)
	if err != nil {
		return &model.Engine{}, err
	}

	return engine, nil
}

func (s engineStore) Update(engine *model.Engine) (*model.Engine, error) {
	// check if record exists in table
	_, err := s.GetByID(engine.ID)
	if err != nil {
		return &model.Engine{}, err
	}

	stmt, err := s.db.Prepare("update engines set displacement = ?, noOfCylinder = ?, `range` = ? where engineId = ?")
	if err != nil {
		return &model.Engine{}, err
	}

	defer stmt.Close()

	_, err = stmt.Exec(engine.Displacement, engine.NoOfCylinders, engine.Range, engine.ID)
	if err != nil {
		return &model.Engine{}, err
	}

	return engine, nil
}

func (s engineStore) Delete(id string) error {
	stmt, err := s.db.Prepare(`delete from engines where engineId = ?`)
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
