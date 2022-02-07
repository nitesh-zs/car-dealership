package engine

import (
	"database/sql"
	"log"

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

func (s engineStore) GetAll() (map[string]model.Engine, error) {
	var engine model.Engine
	engines := make(map[string]model.Engine)

	rows, err := s.db.Query(getAllEngines)
	if err != nil {
		return nil, err
	}

	defer func() {
		rows.Close()

		err = rows.Err()
		if err != nil {
			log.Println(err)
		}
	}()

	for rows.Next() {
		err := rows.Scan(&engine.ID, &engine.Displacement, &engine.NoOfCylinders, &engine.Range)
		if err != nil {
			return nil, err
		}

		engines[engine.ID] = engine
	}

	return engines, nil
}

func (s engineStore) GetByID(id string) (*model.Engine, error) {
	var engine model.Engine

	row := s.db.QueryRow(getEngineByID, id)
	err := row.Scan(&engine.ID, &engine.Displacement, &engine.NoOfCylinders, &engine.Range)

	if err == sql.ErrNoRows {
		return nil, customErrors.CarNotExists()
	}

	if err != nil {
		return nil, err
	}

	return &engine, nil
}

func (s engineStore) Create(engine *model.Engine) (*model.Engine, error) {
	engine.ID = uuid.NewString()
	stmt, err := s.db.Prepare(insertEngine)

	if err != nil {
		return nil, err
	}

	defer stmt.Close()

	_, err = stmt.Exec(engine.ID, engine.Displacement, engine.NoOfCylinders, engine.Range)
	if err != nil {
		return nil, err
	}

	return engine, nil
}

func (s engineStore) Update(engine *model.Engine) (*model.Engine, error) {
	stmt, err := s.db.Prepare(updateEngine)
	if err != nil {
		return nil, err
	}

	defer stmt.Close()

	_, err = stmt.Exec(engine.Displacement, engine.NoOfCylinders, engine.Range, engine.ID)
	if err != nil {
		return nil, err
	}

	return engine, nil
}

func (s engineStore) Delete(id string) error {
	stmt, err := s.db.Prepare(deleteEngine)
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
