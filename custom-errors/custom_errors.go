package customerrors

import "fmt"

type EntityNotExists string

func (e EntityNotExists) Error() string {
	return fmt.Sprintf("%v not exists", string(e))
}

func CarNotExists() EntityNotExists {
	var e EntityNotExists = "Car"
	return e
}

func EngineNotExists() EntityNotExists {
	var e EntityNotExists = "Engine"
	return e
}
