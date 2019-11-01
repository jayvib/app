package vehicle

// Liskov Substitution: Subtypes must be substitutable for their base types.

type actions interface {
	drive()
	startEngine()
}

type Vehicle struct {}
func (Vehicle) drive() {}
func (Vehicle) startEngine() {}
func (Vehicle) stopEngine() {}

type Car struct {
	Vehicle
}

type Sled struct {
	Vehicle
}

func (s Sled) startEngine() {}
func (s Sled) stopEngine() {}
func (s Sled) pushStart() {}

func Go(vehicle actions) {
	if sled, ok := vehicle.(Sled); ok {
		sled.pushStart()
	} else {
		vehicle.startEngine()
	}
	vehicle.drive()
}


