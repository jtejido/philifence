package main

import (
	"fmt"
	"github.com/codegangsta/cli"
	"github.com/jtejido/philifence"
	"log"
	"os"
)

var version = "0.0.1"

func client(args []string) {
	app := cli.NewApp()
	app.Name = "PhiliFence"
	app.Usage = "Putting up the White Picket-Fences and laying Yellow-Bricked roads around you."
	app.Version = version
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "port, p",
			Value: "8080",
			Usage: "Port to bind to",
		},
		cli.StringFlag{
			Name:  "road-path, road",
			Value: "../osm_philippine_roads_wgs84_2012/",
			Usage: "Path for roads",
		},
		cli.StringFlag{
			Name:  "fence-path, fence",
			Value: "../gadm_philippine_cities_wgs84_v2/",
			Usage: "Path for city boundaries",
		},
		cli.BoolFlag{
			Name:  "with-profiler",
			Usage: "Profiling endpoints",
		},
	}
	app.Action = func(c *cli.Context) {
		log.Println("Starting PhiliFence")
		fencePath := fmt.Sprintf("%s", c.String("fence-path"))
		fences, err := philifence.LoadIndex(fencePath)
		if err != nil {
			die(c, err.Error())
		}
		roadPath := fmt.Sprintf("%s", c.String("road-path"))
		roads, err := philifence.LoadIndex(roadPath)
		if err != nil {
			die(c, err.Error())
		}
		prof := c.Bool("with-profiler")
		port := fmt.Sprintf(":%s", c.String("port"))
		err = philifence.ListenAndServe(port, fences, roads, prof)
		die(c, err.Error())
	}
	app.Run(args)
}

func main() {
	client(os.Args)
}

func die(c *cli.Context, msg string) {
	cli.ShowAppHelp(c)
	fmt.Println("Done!")
	fmt.Println(msg)
	os.Exit(1)
}
