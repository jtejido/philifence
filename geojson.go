package philifence

import (
	"encoding/json"
	"github.com/kpawlik/geojson"
	"io/ioutil"
)

type Source struct {
	path string
}

func NewSource(path string) *Source {
	return &Source{path}
}

func (gj *Source) Publish() (features chan *Feature, err error) {
	collection, err := readGeoJson(gj.path)
	if err != nil {
		return
	}
	features = publishFeatureCollection(collection)
	return
}

func featureAdapter(gj *geojson.Feature) (feature *Feature, err error) {

	// TO-DO. Named or Linked Crs
	igeom, err := gj.GetGeometry()
	if igeom == nil || err != nil {
		err = errorf("Invalid geojson feature %q", gj)
		return
	}
	feature = NewFeature(igeom.GetType())

	feature.Properties = gj.Properties
	if feature.Properties != nil {
		feature.Properties["id"] = gj.Id
	}
	if feature.Properties != nil {
		feature.Crs = gj.Crs
	}

	feature.Type = igeom.GetType()
	switch geom := igeom.(type) {
	case *geojson.Point:
		poly := coordinatesAdapter(geojson.Coordinates{geom.Coordinates})
		feature.AddPoly(poly)
	case *geojson.LineString:
		poly := coordinatesAdapter(geom.Coordinates)
		feature.AddPoly(poly)
	case *geojson.MultiPoint:
		poly := coordinatesAdapter(geom.Coordinates)
		feature.AddPoly(poly)
	case *geojson.MultiLineString:
		for _, line := range geom.Coordinates {
			poly := coordinatesAdapter(line)
			feature.AddPoly(poly)
		}
	case *geojson.Polygon:
		exterior := true
		for _, line := range geom.Coordinates {
			poly := coordinatesAdapter(line)
			if exterior {
				if !poly.isClockwise() {
					poly.reverse()
				}
				exterior = false
			} else {
				if poly.isClockwise() {
					poly.reverse()
				}
			}
			feature.AddPoly(poly)
		}
	case *geojson.MultiPolygon:
		for _, multiline := range geom.Coordinates {
			exterior := true
			for _, line := range multiline {
				poly := coordinatesAdapter(line)
				if exterior {
					if !poly.isClockwise() {
						poly.reverse()
					}
					exterior = false
				} else {
					if poly.isClockwise() {
						poly.reverse()
					}
				}
				feature.AddPoly(poly)
			}
		}
	default:
		feature = nil
		err = errorf("Invalid Coordinate Type in GeoJson %q", geom)
	}
	return
}

func unmarshalFeature(raw string) (feature *geojson.Feature, err error) {
	err = json.Unmarshal([]byte(raw), &feature)
	return
}

func coordinatesAdapter(line geojson.Coordinates) (poly *Polygon) {
	poly = MakePoly(len(line))
	for i, point := range line {
		lat := float64(point[1])
		lon := float64(point[0])
		coord := Coordinate{lat: lat, lon: lon}
		poly.Coordinates[i] = coord
	}
	return
}

func readGeoJson(path string) (features *geojson.FeatureCollection, err error) {
	file, err := ioutil.ReadFile(path)
	if err != nil {
		return
	}
	err = json.Unmarshal(file, &features)
	return
}

func publishFeatureCollection(collection *geojson.FeatureCollection) (features chan *Feature) {
	features = make(chan *Feature, 10)
	go func() {
		defer close(features)
		for _, feature := range collection.Features {
			f, err := featureAdapter(feature)
			warn(err, "feature publishing")
			features <- f
		}
	}()
	return
}
