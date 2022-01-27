package customErrors

type CarNotExists string
type EngineNotExists string

func (e CarNotExists) Error() string {
	return "Car not exists"
}

func (e EngineNotExists) Error() string {
	return "Engine not exists"
}
