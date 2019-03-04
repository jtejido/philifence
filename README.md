# philifence
![openstreetmap](http://2.bp.blogspot.com/-xanwG0Mtg18/UjhzWnBWQDI/AAAAAAAAR-c/GTA4EgA1GKg/s1600/ianlopez_temp.jpg=250x250)

PhiliFence is a GeoFencing and Route-listing REST-API service, aimed at finding the boundary (district boundaries as defined by [GADM](https://gadm.org/) and national roads (as defined by [OSM](https://www.openstreetmap.org/about)) from a given location and boundary tolerance.

This aims (though not in due time...I think...) to be a complete [route-finding](https://wiki.openstreetmap.org/wiki/OpenRouteService) and [geo-fencing](https://mediavision2020.com/25-top-geofencing-companies/) in-memory service aimed at systems providing geo-push notifications (geo-targeted marketing (Facebook's promoted ads, Spatially, ThinkNear, etc.), notifications at location-enter/exit (Apple's Notifications), proximity alert (ProximiT), Uber/Lyft-like driver push notifications on users searching for cars, Child-location services (when kids left a location), location-based SMS (NDRRMC, Seismic/Atmospheric location-related events), Content Localization (Netflix's geographic restriction), etc.), and of course, route-finding.

This makes use of an R-Tree variant called a [Hilbert R-Tree](https://en.wikipedia.org/wiki/Hilbert_R-tree) for geometrical indexing, which provides better coordinate ordering, and thus, better compression.

Currently, I have a [not-so-updated](http://philgis.org/general-country-datasets/country-basemaps) administrative boundaries and national roads data, but is API-ready for adding new fences and roads as needed.


## Usage:


### Starting the Service:

Unzip GeoJson sources in "../osm_philippine_roads_wgs84_2012/" and "../gadm_philippine_cities_wgs84_v2/".

Go to the cli/ directory and run ./cli -h

```bash
NAME:
   PhiliFence - Putting up the White Picket-Fences and laying Yellow-Bricked roads around you.

USAGE:
   cli [global options] command [command options] [arguments...]

VERSION:
   0.0.1

COMMANDS:
     help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --port value, -p value             Port to bind to (default: "8080")
   --road-path value, --road value    Path for roads (default: "../osm_philippine_roads_wgs84_2012/")
   --fence-path value, --fence value  Path for city boundaries (default: "../gadm_philippine_cities_wgs84_v2/")
   --with-profiler                    Profiling endpoints
   --help, -h                         show help
   --version, -v                      print the version
```

Simply starting the service should index all geojson from a given path.


```bash
$ ./cli -port=8383
2019/03/04 21:12:10 Starting PhiliFence
2019/03/04 21:12:10 INFO: Putting up fence "philippine-cities" from ../gadm_philippine_cities_wgs84_v2/philippine_cities.json
2019/03/04 21:12:11 INFO: Loaded 1647 features for "philippine-cities"
2019/03/04 21:12:11 INFO: Putting up fence "philippine-roads" from ../osm_philippine_roads_wgs84_2012/philippine_roads.json
2019/03/04 21:12:31 INFO: Loaded 276780 features for "philippine-roads"
2019/03/04 21:12:31 INFO: Fencing on address :8383
```

### Using the Service:


***Get fence that intersects a given location (e.g., as a user inside Fort San Pedro Church, Cebu)***

```
http://localhost:8181/fence/philippine-cities/search?lat=10.2925&lon=123.9056
http://localhost:8181/fence/philippine-cities/search?lat=10.2925&lon=123.9056&tolerance=1
```

***Get roads that intersects a given location (e.g., as a user in Elliptical Road, Q.C.)***

```
http://localhost:8181/road/philippine-roads/search?lat=14.6503&lon=121.0520
http://localhost:8181/road/philippine-roads/search?lat=14.6503&lon=121.0520&tolerance=10
```

**note:** tolerance is the bounding box around the given point, this value is in meters (it creates a bounded box around the point)

***Load All fence indices***

```
http://localhost:8181/fence
```

***Load All road indices***

```
http://localhost:8181/road
```

***Add new fence index given a name***

```
POST json at http://localhost:8181/fence/{name}/add
```

***Add new road index given a name***

```
POST json at http://localhost:8181/road/{name}/add
```


## To-Do:

1. Object insertion at a given fence (e.g., nearest restaurants within query's fence boundary).
2. Objects within roads A.K.A. actual best-route-finding stuff. (Isoline routing, etc.).
3. Merge.
4. K-NearestNeighbours inside a fence (see [Geodesy-PHP](https://github.com/jtejido/geodesy-php)).
5. G-NearestNeighbours.
6. Data Persistence and Data reload on-start
7. Scalable, Distributed type ([SD-Rtree](http://cedric.cnam.fr/~dumouza/EnsPubli/icde07.pdf) implem. on Hilbert RTree?).
8. Tests (Hilbert R-Tree is [tested] though).
