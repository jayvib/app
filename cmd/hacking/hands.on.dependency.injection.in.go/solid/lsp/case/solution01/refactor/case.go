package vehicle

// Liskov Substitution: Subtypes must be substitutable for their base types.

type driver interface {
	drive()
}

type actions interface {
	driver
	start()
}

type Vehicle struct {}
func (Vehicle) drive() {}

type PoweredVehicle struct {
	driver
}

func (v PoweredVehicle) startEngine() {}
func (v PoweredVehicle) stopEngine() {}

type Car struct {
	PoweredVehicle
}

func(c Car) start() {
	c.PoweredVehicle.startEngine()
}

type Buggy struct {
	driver
}

func(b Buggy) start() {
	b.pushStart()
}

func (b Buggy) pushStart() {}

func Go(vehicle actions) {
	vehicle.start()
	vehicle.drive()
}