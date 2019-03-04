package philifence

import (
	"fmt"
)

type Coordinate struct {
	lat, lon float64
}

func (c Coordinate) Lon() float64 {
	return c.lon
}

func (c Coordinate) Lat() float64 {
	return c.lat
}

func (c Coordinate) String() string {
	return fmt.Sprintf("[%.5f, %.5f]", c.lat, c.lon)
}
