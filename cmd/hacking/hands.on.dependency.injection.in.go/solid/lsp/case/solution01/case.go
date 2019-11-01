package vehicle

// Liskov Substitution: Subtypes must be substitutable for their base types.

type actions interface {
	drive()
}

type poweredActions interface {
	actions
	startEngine()
	stopEngine()
}

type unpoweredActions interface {
	actions
	pushStart()
}

type Vehicle struct {}
func (Vehicle) drive() {}

type PoweredVehicle struct {
	Vehicle
}

func (v PoweredVehicle) startEngine() {}
func (v PoweredVehicle) stopEngine() {}

type Car struct {
	PoweredVehicle
}

type Buggy struct {
	Vehicle
}

func (b Buggy) pushStart() {}

func Go(vehicle actions) {
	switch conc := vehicle.(type) {
	case poweredActions:
		conc.startEngine()
	case unpoweredActions:
		conc.pushStart()
	}
	vehicle.drive()
}