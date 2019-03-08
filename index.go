package philifence

import (
	"fmt"
	"path/filepath"
	"sync"
)

//FenceIndex is a dictionary of multiple fences. Useful if you have multiple data sets that need to be searched
type FenceIndex interface {
	Set(name string, fence *Fence)
	Get(name string) *Fence
	Add(name string, feature *Feature) error
	Search(name string, c Coordinate, tol float64) ([]*Feature, error)
	Keys() []string
}

// Returns a thread-safe FenceIndex
func NewFenceIndex() FenceIndex {
	return NewMutexFenceIndex()
}

type UnsafeFenceIndex struct {
	fences map[string]*Fence
}

func NewUnsafeFenceIndex() *UnsafeFenceIndex {
	return &UnsafeFenceIndex{fences: make(map[string]*Fence)}
}

func (idx *UnsafeFenceIndex) Set(name string, fence *Fence) {
	idx.fences[name] = fence
}

func (idx *UnsafeFenceIndex) Get(name string) (fence *Fence) {
	return idx.fences[name]
}

func (idx *UnsafeFenceIndex) Add(name string, feature *Feature) (err error) {
	fence, ok := idx.fences[name]
	if !ok {
		return fmt.Errorf("FenceIndex does not contain fence %q", name)
	}
	fence.Add(feature)
	return
}

func (idx *UnsafeFenceIndex) Search(name string, c Coordinate, tol float64) (matchs []*Feature, err error) {
	fence, ok := idx.fences[name]
	if !ok {
		err = fmt.Errorf("FenceIndex does not contain fence %q", name)
		return
	}
	info("Searching fence for latitude : %.5f, longitude : %.5f in %q", c.lat, c.lon, name)
	matchs = fence.Get(c, tol)
	return
}

func (idx *UnsafeFenceIndex) Keys() (keys []string) {
	for k := range idx.fences {
		keys = append(keys, k)
	}
	return
}

type MutexFenceIndex struct {
	fences *UnsafeFenceIndex
	sync.RWMutex
}

func NewMutexFenceIndex() *MutexFenceIndex {
	return &MutexFenceIndex{fences: NewUnsafeFenceIndex()}
}

func (idx *MutexFenceIndex) Set(name string, fence *Fence) {
	idx.Lock()
	defer idx.Unlock()
	idx.fences.Set(name, fence)
}

func (idx *MutexFenceIndex) Get(name string) *Fence {
	idx.RLock()
	defer idx.RUnlock()
	return idx.fences.Get(name)
}

func (idx *MutexFenceIndex) Add(name string, feature *Feature) error {
	idx.Lock()
	defer idx.Unlock()
	return idx.fences.Add(name, feature)
}

func (idx *MutexFenceIndex) Search(name string, c Coordinate, tol float64) ([]*Feature, error) {
	idx.RLock()
	defer idx.RUnlock()
	return idx.fences.Search(name, c, tol)
}

func (idx *MutexFenceIndex) Keys() []string {
	idx.RLock()
	defer idx.RUnlock()
	return idx.fences.Keys()
}

func LoadIndex(dir string) (fences FenceIndex, err error) {
	paths, err := filepath.Glob(filepath.Join(dir, "*json"))
	if err != nil {
		return
	}
	fences = NewFenceIndex()
	for _, path := range paths {
		key := sluggify(path)
		info("Indexing %q from %s\n", key, path)
		fence, err := NewFence()
		if err != nil {
			fatal("Error building fence for %q. ERROR: %v", key, err)
			continue
		}
		source := NewSource(path)
		features, err := source.Publish()
		if err != nil {
			return nil, err
		}
		i := 0
		for feature := range features {
			if feature.Type == "Point" {
				continue
			}
			fence.Add(feature)
			i++
		}
		info("Loaded %d features for %q\n", i, key)
		fences.Set(key, fence)
	}
	if len(fences.Keys()) < 1 {
		fences = nil
		err = fmt.Errorf("No valid geojson fences at %s", dir)
	}
	return
}
