package philifence

type Properties map[string]interface{}

type PointMessage struct {
	Type       string        `json:"type"`
	Properties Properties    `json:"properties"`
	Geometry   PointGeometry `json:"geometry"`
}

type PointGeometry struct {
	Type        string    `json:"type"`
	Coordinates []float64 `json:"coordinates"`
}

type ResponseMessage struct {
	Query  PointMessage `json:"query"`
	Result []Properties `json:"result"`
}

func newPointMessage(c Coordinate, props Properties) *PointMessage {
	return &PointMessage{
		Type:       "Feature",
		Properties: props,
		Geometry: PointGeometry{
			Type:        "Point",
			Coordinates: []float64{c.lon, c.lat},
		},
	}
}

func newResponseMessage(c Coordinate, props map[string]interface{}, fences []Properties) *ResponseMessage {
	return &ResponseMessage{
		Query:  *newPointMessage(c, Properties(props)),
		Result: fences,
	}
}
