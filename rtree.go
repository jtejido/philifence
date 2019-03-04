package philifence

var (
	MinimumNodeChildren        = DefaultMinNodeEntries
	MaximumNodeChildren        = DefaultMaxNodeEntries
	Resolution                 = DefaultResolution
	dim                 uint64 = 1 << (uint(Resolution) - 1)
)

type Rtree struct {
	rtree *HRtree
}

func NewRtree() (*Rtree, error) {
	rt, err := NewTree(MinimumNodeChildren, MaximumNodeChildren, Resolution)

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

func (r *Rtree) intersections(q Rectangle) []*customRect {
	inodes := r.rtree.SearchIntersect(q)

	nodes := make([]*customRect, len(inodes))
	for i, inode := range inodes {
		nodes[i] = inode.(*customRect)
	}

	return nodes
}

func (r *Rtree) Contains(c Coordinate, tol float64) []*customRect {
	tol = tol / 100000.0 // no limit at this point
	lower := Coordinate{lon: c.lon - tol, lat: c.lat - tol}
	upper := Coordinate{lon: c.lon + tol, lat: c.lat + tol}
	poly := NewPoly(lower, upper)
	q := &customRect{polygon: poly, box: poly.computeBox()}
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
func (n *customRect) LowerLeft() Point {
	return Point{lonToUint32(n.box.min.lon), latToUint32(n.box.min.lat)}
}

func (n *customRect) UpperRight() Point {
	return Point{lonToUint32(n.box.max.lon), latToUint32(n.box.max.lat)}
}

// ensure the limit is within -180, 180
func lonToUint32(c float64) uint64 {
	return uint64(float64(dim) * ((c + 180.0) / 360.0))
}

// ensure the limit is within -90.0, 90.0
func latToUint32(c float64) uint64 {
	return uint64(float64(dim) * ((c + 90.0) / 180.0))
}
