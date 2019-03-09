package philifence

import "testing"

func TestReversePolygon(t *testing.T) {
	vals := []float64{5, 4, 3, 2, 1, 0}
	poly := NewPoly()
	for i := range vals {
		f := float64(i)
		c := Coordinate{f, f}
		poly.Add(c)
	}
	poly.Exterior.reverse()
	for i, c := range poly.Exterior.Coordinates {
		if vals[i] != c.lat {
			t.Errorf("Polygon was not reversed %v", poly)
		}
	}
}

func TestClockwise(t *testing.T) {
	rcw := []Coordinate{
		{39.7435437641, -105.003612041},
		{39.7427848013, -105.003011227},
		{39.7431642838, -105.002217293},
		{39.7439067434, -105.002839565},
		{39.7435437641, -105.003612041},
	}
	poly := NewPoly(rcw...)
	if poly.Exterior.isClockwise() {
		t.Errorf("Polygon is clockwise %v", poly)
	}
	poly.Exterior.reverse()
	if !poly.Exterior.isClockwise() {
		t.Errorf("Polygon is not clockwise %v", poly)
	}
}

func TestBox(t *testing.T) {
	boxes := []Box{
		{min: cd(0, 0), max: cd(10, 10)},
	}
	points := [][]Coordinate{
		{
			cd(0, 5), cd(5, 0), cd(5, 10), cd(10, 5),
		},
	}
	for i, coords := range points {
		poly := NewPoly(coords...)
		if poly.computeBox() != boxes[i] {
			t.Errorf("Write Box test %v", poly.computeBox())

		}
	}
}

func TestContains(t *testing.T) {
	points := []Coordinate{cd(5, 5)}
	polys := [][]Coordinate{
		{
			cd(0, 5), cd(5, 0), cd(5, 10), cd(10, 5),
		},
	}
	for i, coords := range polys {
		poly := NewPoly(coords...)
		if !poly.Contains(points[i]) {
			t.Errorf("Polygon !contains %v %v", poly, points[i])
		}
	}
}

func cd(x, y float64) Coordinate {
	return Coordinate{x, y}
}
