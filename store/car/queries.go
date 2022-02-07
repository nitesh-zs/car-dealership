package car

const (
	getAllCars    = "select * from cars"
	getCarByBrand = "select * from cars where brand = ?"
	getCarByID    = "select * from cars where carId = ?"
	insertCar     = `insert into cars (carId, name, yearOfManufacture, brand, fuelType, engineId)
					values (?, ?, ?, ?, ?, ?)`
	updateCar = `update cars set name = ?, yearOfManufacture = ?, brand = ?, fuelType = ? where carId = ?`
	deleteCar = `delete from cars where carId = ?`
)
