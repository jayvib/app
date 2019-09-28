package shapes

import "testing"

func TestPerimeter(t *testing.T) {
	tests := []struct{
		name string
		shape Perimeterer
		want float64
	}{
		{"rectangle", Rectangle{Width: 10.0, Height: 10.0}, 40.0},
		{"circle", Circle{Radius: 10.0}, 62.83185307179586}	,
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T){
			assertPerimeter(t, tc.shape, tc.want)
		})
	}

}

func TestArea(t *testing.T) {
	tests := []struct{
		name string
		shape Areaer
		want float64
	}	{
		{
			"rectangle",
			Rectangle{Width: 10.0, Height: 6.0},
			60.0,
		},
		{
			"circle",
			Circle{Radius: 10},
			314.1592653589793,
		},
		{
			"triangle",
			Triangle{Height: 10.0, Base: 5.0},
			25.0,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T){
			assertArea(t, tc.shape, tc.want)
		})
	}

}

func assertFloat(t *testing.T, want, got float64) {
	t.Helper()
	if got != want {
		t.Errorf("got '%g' want '%g'", got, want)
	}
}


// ISP client should not depend on the method that they didn't use.
func assertArea(t *testing.T, shape Areaer, want float64) { // what about applying the interface segregation principle
	t.Helper()
	got := shape.Area()
	assertFloat(t, want, got)
}

func assertPerimeter(t *testing.T, shape Perimeterer, want float64) {
	t.Helper()
	got := shape.Perimeter()
	assertFloat(t, got, want)
}
