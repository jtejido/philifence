package philifence

import (
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/pprof"
	"strconv"
)

var fences, roads FenceIndex

func ListenAndServe(addr string, fidx, ridx FenceIndex, profile bool) error {
	info("Fencing on address %s\n", addr)
	defer info("Done Fencing\n")
	fences = fidx
	roads = ridx
	router := httprouter.New()
	router.GET("/fence", getFenceList)
	router.POST("/fence/:name/add", postFenceAdd)
	router.GET("/fence/:name/search", getFenceSearch)
	router.GET("/road", getRoadList)
	router.POST("/road/:name/add", postRoadAdd)
	router.GET("/road/:name/search", getRoadSearch)
	if profile {
		profiler(router)
		info("Profiling available at /debug/pprof/")
	}
	return http.ListenAndServe(addr, router)
}

func respond(w http.ResponseWriter, res interface{}) {
	w.Header().Set("Server", "philifence")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	writeJson(w, res)
}

func getFenceList(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	writeJson(w, fences.Keys())
}

func getRoadList(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	writeJson(w, roads.Keys())
}

func postFenceAdd(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1<<26)) // 64 MB max
	if err != nil {
		http.Error(w, "Body 64 MB max", http.StatusRequestEntityTooLarge)
		return
	}
	if err := r.Body.Close(); err != nil {
		http.Error(w, "Error closing body", http.StatusInternalServerError)
		return
	}
	name := params.ByName("name")
	g, err := unmarshalFeature(string(body))
	if err != nil {
		http.Error(w, "Unable to read geojson feature", http.StatusBadRequest)
		return
	}
	feature, err := featureAdapter(g)
	if err != nil {
		http.Error(w, "Unable to read geojson feature", http.StatusBadRequest)
		return
	}
	if err := fences.Add(name, feature); err != nil {
		http.Error(w, "Error adding feature "+err.Error(), http.StatusBadRequest)
	}
	respond(w, "success")
}

func postRoadAdd(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1<<26)) // 64 MB max
	if err != nil {
		http.Error(w, "Body 64 MB max", http.StatusRequestEntityTooLarge)
		return
	}
	if err := r.Body.Close(); err != nil {
		http.Error(w, "Error closing body", http.StatusInternalServerError)
		return
	}
	name := params.ByName("name")
	g, err := unmarshalFeature(string(body))
	if err != nil {
		http.Error(w, "Unable to read geojson feature", http.StatusBadRequest)
		return
	}
	feature, err := featureAdapter(g)
	if err != nil {
		http.Error(w, "Unable to read geojson feature", http.StatusBadRequest)
		return
	}
	if err := roads.Add(name, feature); err != nil {
		http.Error(w, "Error adding feature "+err.Error(), http.StatusBadRequest)
	}
	respond(w, "success")
}

func getFenceSearch(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	query := r.URL.Query()
	lat, err := strconv.ParseFloat(query.Get("lat"), 64)
	if err != nil {
		http.Error(w, "Query param 'lat' required as float", http.StatusBadRequest)
		return
	}
	lon, err := strconv.ParseFloat(query.Get("lon"), 64)
	if err != nil {
		http.Error(w, "Query param 'lon' required as float", http.StatusBadRequest)
		return
	}
	tol, err := strconv.ParseFloat(query.Get("tolerance"), 64)
	if err != nil {
		tol = 1 // ~1m
	}

	query.Del("lat")
	query.Del("lon")
	query.Del("tolerance")
	c := Coordinate{lat: lat, lon: lon}
	name := params.ByName("name")
	matchs, err := fences.Search(name, c, tol)
	if err != nil {
		http.Error(w, "Error search fence "+name, http.StatusBadRequest)
		return
	}
	fences := make([]Properties, len(matchs))
	for i, fence := range matchs {
		fences[i] = fence.Properties
	}
	props := make(map[string]interface{}, len(query))
	for k := range query {
		props[k] = query.Get(k)
	}

	respond(w, *newResponseMessage(c, props, fences))
}

func getRoadSearch(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	query := r.URL.Query()
	lat, err := strconv.ParseFloat(query.Get("lat"), 64)
	if err != nil {
		http.Error(w, "Query param 'lat' required as float", http.StatusBadRequest)
		return
	}
	lon, err := strconv.ParseFloat(query.Get("lon"), 64)
	if err != nil {
		http.Error(w, "Query param 'lon' required as float", http.StatusBadRequest)
		return
	}
	tol, err := strconv.ParseFloat(query.Get("tolerance"), 64)
	if err != nil {
		tol = 1 // ~1m
	}

	query.Del("lat")
	query.Del("lon")
	query.Del("tolerance")
	c := Coordinate{lat: lat, lon: lon}
	name := params.ByName("name")
	matchs, err := roads.Search(name, c, tol)
	if err != nil {
		http.Error(w, "Error search road "+name, http.StatusBadRequest)
		return
	}
	roads := make([]Properties, len(matchs))
	for i, road := range matchs {
		roads[i] = road.Properties
	}
	props := make(map[string]interface{}, len(query))
	for k := range query {
		props[k] = query.Get(k)
	}

	respond(w, *newResponseMessage(c, props, roads))
}

func writeJson(w io.Writer, msg interface{}) (err error) {
	buf, err := json.Marshal(&msg)
	_, err = w.Write(buf)
	if err != nil {
		return err
	}

	return
}

func profiler(router *httprouter.Router) {
	router.HandlerFunc("GET", "/debug/pprof/", pprof.Index)
	router.HandlerFunc("POST", "/debug/pprof/", pprof.Index)
	router.HandlerFunc("GET", "/debug/pprof/cmdline", pprof.Cmdline)
	router.HandlerFunc("POST", "/debug/pprof/cmdline", pprof.Cmdline)
	router.HandlerFunc("GET", "/debug/pprof/profile", pprof.Profile)
	router.HandlerFunc("POST", "/debug/pprof/profile", pprof.Profile)
	router.HandlerFunc("GET", "/debug/pprof/symbol", pprof.Symbol)
	router.HandlerFunc("POST", "/debug/pprof/symbol", pprof.Symbol)
	router.Handler("GET", "/debug/pprof/heap", pprof.Handler("heap"))
	router.Handler("GET", "/debug/pprof/block", pprof.Handler("block"))
	router.Handler("GET", "/debug/pprof/goroutine", pprof.Handler("goroutine"))
	router.Handler("GET", "/debug/pprof/threadcreate", pprof.Handler("threadcreate"))
}
