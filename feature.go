package philifence

import (
	"github.com/kpawlik/geojson"
	"strings"
)

type Feature struct {
	Geometry   []*Polygon
	Type       string
	Crs        *geojson.CRS
	Properties map[string]interface{}
}

func NewFeature(geometryType string, geometry ...*Polygon) *Feature {
	geometryType = strings.ToLower(geometryType)
	return &Feature{Geometry: geometry, Type: geometryType}
}

func NewPolygonFeature(geometry ...*Polygon) *Feature {
	return NewFeature("polygon", geometry...)
}

func NewLineFeature(geometry ...*Polygon) *Feature {
	return NewFeature("line", geometry...)
}

func NewPointFeature(cs ...Coordinate) *Feature {
	feature := MakeFeature(len(cs))
	feature.Type = "point"
	for i, c := range cs {
		feature.Geometry[i] = NewPoly(c)
	}
	return feature
}

func MakeFeature(length int) *Feature {
	return &Feature{Geometry: make([]*Polygon, length)}
}

func (f *Feature) AddPoly(s *Polygon) {
	f.Geometry = append(f.Geometry, s)
}

func (f *Feature) Contains(c Coordinate) bool {
	for _, poly := range f.Geometry {
		if poly.Contains(c) {
			return true
		}
	}
	return false
}
