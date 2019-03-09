package philifence

import (
	"github.com/jtejido/hrtree"
	"math"
)

var (
	MinimumNodeChildren        = 50
	MaximumNodeChildren        = 200
	Resolution                 = 32 // 64-bit resolution for hilbert curve
	dim                 uint64 = 1 << (uint(Resolution) - 1)
)

const (
	earthRadius = 6371e3 // assume WGS84?
	radians     = math.Pi / 180
	degrees     = 180 / math.Pi
)

type Rtree struct {
	rtree *hrtree.HRtree
}

func NewRtree() (*Rtree, error) {
	rt, err := hrtree.NewTree(MinimumNodeChildren, MaximumNodeChildren, Resolution)

	return &Rtree{
		rtree: rt,
	}, err
}

func (r *Rtree) Size() int {
	return r.rtree.Size()
}

func (r *Rtree) Insert(s *Polygon, data interface{}) {
	node := &customRect{s, s.computeBox(), data}
	r.rtree.Insert(node)
}

func (r *Rtree) intersections(q hrtree.Rectangle) []*customRect {
	inodes := r.rtree.SearchIntersect(q)

	nodes := make([]*customRect, len(inodes))
	for i, inode := range inodes {
		nodes[i] = inode.(*customRect)
	}

	return nodes
}

func (r *Rtree) Contains(c Coordinate, tol float64) []*customRect {
	q := rectFromCenter(c, tol)
	return r.intersections(q)
}

// implements Spatial
type customRect struct {
	polygon *Polygon
	box     Box // precomputed box
	data    interface{}
}

func (n *customRect) Feature() *Feature {
	return n.Value().(*Feature)
}

func (n *customRect) Value() interface{} {
	return n.data
}

// implements Rectangle
func (n *customRect) LowerLeft() hrtree.Point {
	return hrtree.Point{lonToUint32(n.box.min.lon), latToUint32(n.box.min.lat)}
}

func (n *customRect) UpperRight() hrtree.Point {
	return hrtree.Point{lonToUint32(n.box.max.lon), latToUint32(n.box.max.lat)}
}

// ensure the limit is within -180, 180
func lonToUint32(c float64) uint64 {
	return uint64(float64(dim) * ((c + 180.0) / 360.0))
}

// ensure the limit is within -90.0, 90.0
func latToUint32(c float64) uint64 {
	return uint64(float64(dim) * ((c + 90.0) / 180.0))
}

// from http://janmatuschek.de/LatitudeLongitudeBoundingCoordinates#Latitude
func rectFromCenter(c Coordinate, meters float64) *customRect {

	c.lat *= radians
	c.lon *= radians

	r := meters / earthRadius

	minLat := c.lat - r
	maxLat := c.lat + r

	latT := math.Asin(math.Sin(c.lat) / math.Cos(r))
	lonΔ := math.Acos((math.Cos(r) - math.Sin(latT)*math.Sin(c.lat)) / (math.Cos(latT) * math.Cos(c.lat)))

	minLon := c.lon - lonΔ
	maxLon := c.lon + lonΔ

	if maxLat > math.Pi/2 {
		minLon = -math.Pi
		maxLat = math.Pi / 2
		maxLon = math.Pi
	}

	if minLat < -math.Pi/2 {
		minLat = -math.Pi / 2
		minLon = -math.Pi
		maxLon = math.Pi
	}

	if minLon < -math.Pi || maxLon > math.Pi {
		minLon = -math.Pi
		maxLon = math.Pi
	}

	minLon = math.Mod(minLon+3*math.Pi, 2*math.Pi) - math.Pi // normalise to -180..+180°
	maxLon = math.Mod(maxLon+3*math.Pi, 2*math.Pi) - math.Pi

	minLat *= degrees
	minLon *= degrees
	maxLat *= degrees
	maxLon *= degrees
	lower := Coordinate{lon: minLon, lat: minLat}
	upper := Coordinate{lon: maxLon, lat: maxLat}
	poly := NewPoly(lower, upper)
	return &customRect{polygon: poly, box: poly.computeBox()}

}
