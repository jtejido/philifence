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
		poly := exteriorAdapter(geojson.Coordinates{geom.Coordinates})
		feature.AddPoly(poly)
	case *geojson.LineString:
		poly := exteriorAdapter(geom.Coordinates)
		feature.AddPoly(poly)
	case *geojson.MultiPoint:
		poly := exteriorAdapter(geom.Coordinates)
		feature.AddPoly(poly)
	case *geojson.MultiLineString:
		for _, line := range geom.Coordinates {
			poly := exteriorAdapter(line)
			feature.AddPoly(poly)
		}
	case *geojson.Polygon:
		poly := multilineAdapter(geom.Coordinates)
		feature.AddPoly(poly)
	case *geojson.MultiPolygon:
		for _, multiline := range geom.Coordinates {
			poly := multilineAdapter(multiline)
			feature.AddPoly(poly)
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

func coordinateAdapter(line geojson.Coordinates, ring *PolyRing) {
	for i, point := range line {
		lat := float64(point[1])
		lon := float64(point[0])
		coord := Coordinate{lat: lat, lon: lon}
		ring.Coordinates[i] = coord
	}
}

func exteriorAdapter(line geojson.Coordinates) (poly *Polygon) {
	poly = MakePoly(len(line))
	coordinateAdapter(line, poly.Coordinates)
	return
}

func multilineAdapter(coordinates geojson.MultiLine) (poly *Polygon) {
	
	exterior := true
	ctr := 0
	for _, line := range coordinates {
		if exterior {
			// first is the exterior ring
			poly = exteriorAdapter(line)
			// https://tools.ietf.org/html/rfc7946#section-3.1.6
			if poly.Coordinates.isClockwise() {
				poly.Coordinates.reverse()
			}

			exterior = false
		} else {
			if poly.Holes == nil {
				poly.Holes = make([]*PolyRing, len(coordinates))
			}

			if poly.Holes[ctr] == nil {
				poly.Holes[ctr] = MakePolyRing(len(line))
			}
			
			coordinateAdapter(line, poly.Holes[ctr])

			if !poly.Holes[ctr].isClockwise() {
				poly.Holes[ctr].reverse()
			}

			ctr++
		}
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