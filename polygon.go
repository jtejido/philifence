package philifence

type PolyRing struct {
	Coordinates []Coordinate
	Box 		Box
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

func (pr *PolyRing) Len() int {
	return len(pr.Coordinates)
}

func (pr *PolyRing) Add(c ...Coordinate) {
	pr.Coordinates = append(pr.Coordinates, c...)
}

func (pr *PolyRing) reverse() {
	for i, j := 0, len(pr.Coordinates)-1; i < j; i, j = i+1, j-1 {
		pr.Coordinates[i], pr.Coordinates[j] = pr.Coordinates[j], pr.Coordinates[i]
	}
}

// we'll stick to winding number algorithm
func (pr *PolyRing) Contains(c Coordinate) bool {
	return pr.computeWindingNumber(c) != 0
}

// Raycast shows an implementation of the ray casting point-in-polygon
// (PNPoly) algorithm for testing if a point is inside a closed polygon.
// Also known as the crossing number or the even-odd rule algorithm.
//
// https://www.ecse.rpi.edu/Homepages/wrf/Research/Short_Notes/pnpoly.html
func (pr *PolyRing) inside(pt Coordinate) bool {
    if len(pr.Coordinates) < 3 {
        return false
    }
    in := rayIntersectsSegment(pt, pr.Coordinates[len(pr.Coordinates)-1], pr.Coordinates[0])
    for i := 1; i < len(pr.Coordinates); i++ {
        if rayIntersectsSegment(pt, pr.Coordinates[i-1], pr.Coordinates[i]) {
            in = !in
        }
    }
    return in
}
 
func rayIntersectsSegment(p, a, b Coordinate) bool {
    return (a.lat > p.lat) != (b.lat > p.lat) &&
        p.lon < (b.lon-a.lon)*(p.lat-a.lat)/(b.lat-a.lat)+a.lon
}

// The Winding Number method which counts the number of times the polygon winds around the point P. 
// The point is outside only when this "winding number" wn = 0; otherwise, the point is inside.
//
// http://geomalgorithms.com/a03-_inclusion.html
func (pr *PolyRing) computeWindingNumber(q Coordinate) (wn int) {

	for i := range pr.Coordinates[:pr.Len()-1] {
		if pr.Coordinates[i].lat <= q.lat {
			if pr.Coordinates[i+1].lat > q.lat {
				if isLeft(pr.Coordinates[i], pr.Coordinates[i+1], q) > 0 {
					wn++
				}
			}
		} else {
			if pr.Coordinates[i+1].lat <= q.lat {
				if isLeft(pr.Coordinates[i], pr.Coordinates[i+1], q) < 0 {
					wn--
				}
			}
		}
	}

	return
}

func isLeft(h, t, q Coordinate) float64 {
	return ((t.lon - h.lon) * (q.lat - h.lat)) - ((q.lon - h.lon) * (t.lat - h.lat))
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

type Polygon struct {
	Coordinates *PolyRing
	Holes []*PolyRing
}

func MakePoly(length int) *Polygon {
	return &Polygon{Coordinates: MakePolyRing(length)}
}

func NewPoly(coords ...Coordinate) *Polygon {
	return &Polygon{Coordinates: NewPolyRing(coords...)}
}

func (poly *Polygon) Add(c ...Coordinate) {
	poly.Coordinates.Add(c...)
}

func (poly *Polygon) Contains(c Coordinate) (ok bool) {
	ok = poly.Coordinates.computeWindingNumber(c) != 0

	if ok {
		for _, hole := range poly.Holes {
			if hole.computeWindingNumber(c) != 0 {
				ok = false
				break
			}
		}
	}

	return
}

func (poly *Polygon) Len() int {
	return poly.Coordinates.Len()
}

func (poly *Polygon) computeBox() Box {
	return poly.Coordinates.computeBox()
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