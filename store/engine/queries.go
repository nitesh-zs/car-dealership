package engine

const (
	getAllEngines = "select * from engines"
	getEngineByID = "select * from engines where engineId = ?"
	insertEngine  = "insert into engines (engineId, displacement, noOfCylinder, `range`) values (?, ?, ?, ?)"
	updateEngine  = "update engines set displacement = ?, noOfCylinder = ?, `range` = ? where engineId = ?"
	deleteEngine  = "delete from engines where engineId = ?"
)
