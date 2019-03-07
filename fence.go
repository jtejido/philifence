package philifence

type Fence struct {
	rtree *Rtree
}

func NewFence() (*Fence, error) {
	rt, err := NewRtree()

	return &Fence{
		rtree: rt,
	}, err
}

func (r *Fence) Add(f *Feature) {
	for _, poly := range f.Geometry {
		if poly.Len() > 1 {
			r.rtree.Insert(poly, f)
		}
	}
}

func (r *Fence) Get(c Coordinate, tol float64) (matchs []*Feature) {
	nodes := r.rtree.Contains(c, tol)
	
	for _, n := range nodes {
		feature := n.Feature()
		if feature.Contains(c) {
			matchs = append(matchs, feature)
		}
	}

	return
}

func (r *Fence) Size() int {
	return r.rtree.rtree.Size()
}
