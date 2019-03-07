package philifence

type PolyRing struct {
	Coordinates []Coordinate
}

func MakePolyRing(length int) *PolyRing {
	return &PolyRing{
		Coordinates: make([]Coordinate, length),
	}
}

func NewPolyRing(coords ...Coordinate) *PolyRing {
	return &PolyRing{
		Coordinates: coords,
	}
}

func (pr *PolyRing) Add(c ...Coordinate) {
	pr.Coordinates = append(pr.Coordinates, c...)
}

func (pr *PolyRing) Len() int {
	return len(pr.Coordinates)
}

func (pr *PolyRing) computeBox() Box {

	var min, max Coordinate
	for i := 0; i < len(pr.Coordinates); i++ {
		if i == 0 {
			min, max = pr.Coordinates[0], pr.Coordinates[0]
		} else {
			c := pr.Coordinates[i]
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

func (pr *PolyRing) isClockwise() bool {
	coords := pr.Coordinates
	sum := 0.0
	for i, coord := range coords[:len(coords)-1] {
		next := coords[i+1]
		sum += (next.lon - coord.lon) * (next.lat + coord.lat)
	}

	return sum >= 0
}

func (pr *PolyRing) computeWindingNumber(q Coordinate) (wn int) {

	for i := range pr.Coordinates[:pr.Len()-1] {
		if pr.Coordinates[i].Lat() <= q.Lat() {
			if pr.Coordinates[i+1].Lat() > q.Lat() {
				if isLeft(pr.Coordinates[i], pr.Coordinates[i+1], q) > 0 {
					wn++
				}
			}
		} else {
			if pr.Coordinates[i+1].Lat() <= q.Lat() {
				if isLeft(pr.Coordinates[i], pr.Coordinates[i+1], q) < 0 {
					wn--
				}
			}
		}
	}

	return
}

type Polygon struct {
	Coordinates *PolyRing
	Holes []*PolyRing
}

func MakePoly(length int) *Polygon {
	return &Polygon{
		Coordinates: MakePolyRing(length),
	}
}

func NewPoly(coords ...Coordinate) *Polygon {
	return &Polygon{
		Coordinates: NewPolyRing(coords...),
	}
}

func (poly *Polygon) Add(c ...Coordinate) {
	poly.Coordinates.Add(c...)
}

func (poly *Polygon) computeBox() Box {
	return poly.Coordinates.computeBox()
}

func (poly *Polygon) Contains(c Coordinate) (ok bool) {
	ok = poly.Coordinates.computeWindingNumber(c) != 0

	if ok && poly.Holes != nil {
		// if point is in poly but inside a hole, then return false
		for _, h := range poly.Holes {
			if h.computeWindingNumber(c) != 0 {
				ok = false
				break
			}
		}
	}

	return
}

// func (poly *Polygon) reverse() {
// 	for i, j := 0, len(poly.Coordinates)-1; i < j; i, j = i+1, j-1 {
// 		poly.Coordinates[i], poly.Coordinates[j] = poly.Coordinates[j], poly.Coordinates[i]
// 	}
// }

func (poly *Polygon) Len() int {
	return poly.Coordinates.Len()
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
