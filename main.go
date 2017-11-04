package main

import (
	"flag"
	"fmt"
	"log"
	"strconv"
	"strings"

	shp "github.com/jonas-p/go-shp"
)

type flags struct {
	Input    string
	Output   string
	Test     bool
	Verbose  bool
	Centers  bool
	Routes   bool
	ZipCodes bool
}

func main() {
	inputPointer := flag.String("i", "", "File to read")
	outputPointer := flag.String("o", "", "File to read")
	testPointer := flag.Bool("t", false, "Run test mode")
	verbosePointer := flag.Bool("v", false, "Verbose output")
	centerPointer := flag.Bool("c", false, "Convert centers")
	routePointer := flag.Bool("r", false, "Convert routes")
	zipCodePointer := flag.Bool("z", false, "Convert zip codes")

	flag.Parse()

	passedFlags := flags{
		Input:    *inputPointer,
		Output:   *outputPointer,
		Test:     *testPointer,
		Verbose:  *verbosePointer,
		Centers:  *centerPointer,
		Routes:   *routePointer,
		ZipCodes: *zipCodePointer}

	if passedFlags.Input != "" {
		writeCsv(passedFlags)
	} else if passedFlags.Centers == true {
		convertCenters()
	} else if passedFlags.Routes == true {
		convertRoutes()
	} else if passedFlags.ZipCodes == true {
		convertZips()
	}
}

func writeCsv(currentFlags flags) {
	const GEOCODE string = "GEOCODE"
	path := currentFlags.Input
	// open a shapefile for reading
	shape, err := shp.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer shape.Close()

	// fields from the attribute table (DBF)
	fields := shape.Fields()

	//w := csv.NewWriter(os.Stdout)
	var geocode string
	var prevGeocode string

	// loop through all features in the shapefile
	for shape.Next() {
		n, p := shape.Shape()

		// THIS WAS A NIGHTMARE TO FIGURE OUT!!!!!
		polyZ := p.(*shp.PolygonZ)

		// get the number of points
		pointCount := polyZ.NumPoints

		// Declare as int32 to it jives with pointCount. Better way?  Probably.
		var i int32

		for i = 0; i < pointCount; i++ {
			point := polyZ.Points[i]
			for k, f := range fields {
				fieldName := fieldNameToString(f.Name)
				if fieldName == GEOCODE {
					geocode = shape.ReadAttribute(n, k)
					if prevGeocode != geocode {
						fmt.Println("#" + geocode)
					}
					prevGeocode = geocode
				}
			}
			pointX := strconv.FormatFloat(point.X, 'f', 6, 32)
			pointY := strconv.FormatFloat(point.Y, 'f', 6, 32)
			fmt.Println(pointX + " " + pointY)
		}
	}
}

// Convert the field name to string
func fieldNameToString(bytes [11]byte) string {
	var returnString string
	// This converts bytes to a string. Intuitive, right?  :\
	returnString = string(bytes[:])
	// This removes the null values at the end of the string
	returnString = strings.Trim(returnString, "\x00")
	return returnString
}

func convertRoutes() {
	fmt.Printf("Not yet implemented")
}

func convertZips() {
	fmt.Printf("Not yet implemented")
}

func convertCenters() {
	fmt.Printf("Not yet implemented")
}
