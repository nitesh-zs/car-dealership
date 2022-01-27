package store

import (
	customErrors "carAPI/custom-errors"
	"carAPI/model"
	"database/sql"
	"github.com/google/uuid"
)

type engineStore struct {
	db *sql.DB
}

func NewEngineStore(db *sql.DB) engineStore {
	return engineStore{db: db}
}

var carNotExists customErrors.CarNotExists

func (s engineStore) GetByID(ID string) (model.Engine, error) {
	var engine model.Engine

	row := s.db.QueryRow("select * from engines where engineId = ?", ID)
	err := row.Scan(&engine.ID, &engine.Displacement, &engine.NoOfCylinders, &engine.Range)
	if err == sql.ErrNoRows {
		return model.Engine{}, carNotExists
	}
	if err != nil {
		return model.Engine{}, err
	}

	return engine, nil
}

func (s engineStore) Create(engine model.Engine) (model.Engine, error) {
	engine.ID = uuid.NewString()
	stmt, err := s.db.Prepare("insert into engines (engineId, displacement, noOfCylinder, `range`) values (?, ?, ?, ?)")
	if err != nil {
		return model.Engine{}, err
	}

	defer stmt.Close()

	_, err = stmt.Exec(engine.ID, engine.Displacement, engine.NoOfCylinders, engine.Range)
	if err != nil {
		return model.Engine{}, err
	}

	return engine, nil
}

func (s engineStore) Update(engine model.Engine) (model.Engine, error) {
	stmt, err := s.db.Prepare("update engines set displacement = ?, noOfCylinder = ?, `range` = ? where engineId = ?")
	if err != nil {
		return model.Engine{}, err
	}

	defer stmt.Close()

	res, err := stmt.Exec(engine.Displacement, engine.NoOfCylinders, engine.Range, engine.ID)
	if err != nil {
		return model.Engine{}, err
	}

	rowsAff, err := res.RowsAffected()
	if rowsAff == 0 {
		return model.Engine{}, carNotExists
	}
	return engine, nil
}

func (s engineStore) Delete(ID string) error {
	stmt, err := s.db.Prepare(`delete from engines where engineId = ?`)
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
