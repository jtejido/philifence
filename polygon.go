package philifence

type Polygon struct {
	Coordinates []Coordinate
}

func MakePoly(length int) *Polygon {
	return &Polygon{Coordinates: make([]Coordinate, length)}
}

func NewPoly(coords ...Coordinate) *Polygon {
	return &Polygon{Coordinates: coords}
}

func (poly *Polygon) Add(c ...Coordinate) {
	poly.Coordinates = append(poly.Coordinates, c...)
}

func (poly *Polygon) computeBox() Box {

	var min, max Coordinate
	for i := 0; i < len(poly.Coordinates); i++ {
		if i == 0 {
			min, max = poly.Coordinates[0], poly.Coordinates[0]
		} else {
			c := poly.Coordinates[i]
			if c.lat < min.lat {
				min.lat = c.lat
			}
			if c.lat > max.lat {
				max.lat = c.lat
			}
			if c.lon < min.lon {
				min.lon = c.lon
			}
			if c.lon > max.lon {
				max.lon = c.lon
			}
		}
	}
	box, err := NewBox(min, max)
	check(err)
	return box

}

func (poly *Polygon) Contains(c Coordinate) bool {
	return poly.computeWindingNumber(c) != 0
}

func (poly *Polygon) reverse() {
	for i, j := 0, len(poly.Coordinates)-1; i < j; i, j = i+1, j-1 {
		poly.Coordinates[i], poly.Coordinates[j] = poly.Coordinates[j], poly.Coordinates[i]
	}
}

func (poly *Polygon) Len() int {
	return len(poly.Coordinates)
}

func (poly *Polygon) isClockwise() bool {
	coords := poly.Coordinates
	sum := 0.0
	for i, coord := range coords[:len(coords)-1] {
		next := coords[i+1]
		sum += (next.lon - coord.lon) * (next.lat + coord.lat)
	}

	return sum >= 0
}

func (poly *Polygon) computeWindingNumber(q Coordinate) (wn int) {

	for i := range poly.Coordinates[:poly.Len()-1] {
		if poly.Coordinates[i].Lat() <= q.Lat() {
			if poly.Coordinates[i+1].Lat() > q.Lat() {
				if isLeft(poly.Coordinates[i], poly.Coordinates[i+1], q) > 0 {
					wn++
				}
			}
		} else {
			if poly.Coordinates[i+1].Lat() <= q.Lat() {
				if isLeft(poly.Coordinates[i], poly.Coordinates[i+1], q) < 0 {
					wn--
				}
			}
		}
	}

	return
}

func isLeft(h, t, q Coordinate) float64 {
	return ((t.Lon() - h.Lon()) * (q.Lat() - h.Lat())) - ((q.Lon() - h.Lon()) * (t.Lat() - h.Lat()))
}

// rectangle wrapper for polygon
type Box struct {
	min, max Coordinate
}

func NewBox(min, max Coordinate) (box Box, err error) {
	if min.lat > max.lat || min.lon > max.lon {
		err = errorf("Min %v > Max %v", min, max)
		return
	}

	box = Box{min: min, max: max}

	return
}
