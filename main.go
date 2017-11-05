package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"

	shp "github.com/jonas-p/go-shp"
)

func main() {
	passedFlags := parseFlags()
	processCommand(passedFlags)
}

func parseFlags() flags {

	inputPointer := flag.String("i", "", "File for input")
	outputPointer := flag.String("o", "", "File for output")
	appendPointer := flag.Bool("a", false, "Append to output")
	testPointer := flag.Bool("t", false, "Run test mode")
	verbosePointer := flag.Bool("v", false, "Verbose output")
	centerPointer := flag.Bool("c", false, "Convert centers")
	routePointer := flag.Bool("r", false, "Convert routes")
	zipCodePointer := flag.Bool("z", false, "Convert zip codes")
	inRootPointer := flag.String("b", "", "Batch input root")
	flag.Parse()

	commandFlags := flags{
		Input:     *inputPointer,
		Output:    *outputPointer,
		Append:    *appendPointer,
		Test:      *testPointer,
		Verbose:   *verbosePointer,
		Centers:   *centerPointer,
		Routes:    *routePointer,
		ZipCodes:  *zipCodePointer,
		InputRoot: *inRootPointer}

	return commandFlags
}

func processCommand(passedFlags flags) {
	if passedFlags.InputRoot != "" {
		inputPath := passedFlags.Input
		processed := 0
		fmt.Printf("Processing root %v\n", passedFlags.InputRoot)
		if passedFlags.InputRoot != "" {
			if passedFlags.Routes {
				processed = convertAllShapes(routeFiles, passedFlags)
			} else {
				processed = convertAllShapes(zipFiles, passedFlags)
			}
		} else {
			var output io.Writer
			var err error
			fmt.Printf("Converting %v", inputPath)
			if passedFlags.Append {
				output, err = os.OpenFile(passedFlags.Output, os.O_APPEND|os.O_WRONLY, 0600)
			} else {
				output, err = os.Create(passedFlags.Output)
			}
			if err != nil {
				log.Fatal(err)
			}
			if passedFlags.Centers {
				writeCentroidOutput(inputPath, output)
			} else {
				processed = writePolygonOutput(inputPath, output)
			}
		}
		fmt.Printf("...%v processed\n", processed)
	}
}

func writeCentroidOutput(inputPath string, writer io.Writer) {
	input, err := os.Open(inputPath)
	if err != nil {
		log.Fatal(err)
	}
	reader := csv.NewReader(input)
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}

		geocode := record[0]
		latitude := record[1]
		latitude = latitude[1 : len(latitude)-1]
		longitude := record[2]

		latitude = padLeft(latitude, "0", 10)
		longitude = padLeft(longitude, "0", 10)

		//fmt.Println(zipCode, latitude, longitude)
		io.WriteString(writer, geocode+latitude+longitude+"\n")
	}
}

func writePolygonOutput(inputPath string, writer io.Writer) int {
	const GEOCODE string = "GEOCODE"
	var processed int
	processed = 0

	// open a shapefile for reading
	shape, err := shp.Open(inputPath)
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
		var line string

		for i = 0; i < pointCount; i++ {
			point := polyZ.Points[i]
			for k, f := range fields {
				processed++
				fieldName := fieldNameToString(f.Name)
				if fieldName == GEOCODE {
					geocode = shape.ReadAttribute(n, k)
					if prevGeocode != geocode {
						line = "#" + geocode
						io.WriteString(writer, line+"\n")
					}
					prevGeocode = geocode
				}
			}
			pointX := strconv.FormatFloat(point.X, 'f', 6, 32)
			pointY := strconv.FormatFloat(point.Y, 'f', 6, 32)
			line = pointX + " " + pointY
			io.WriteString(writer, line+"\n")
		}
	}
	return processed
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

func convertAllShapes(fileList [51]string, currentFlags flags) int {
	var processed int
	var totalProcessed int
	outputFile := currentFlags.Output
	output, err := os.Create(outputFile)
	if err != nil {
		log.Fatal(err)
	}
	for _, file := range fileList {
		inputFile := currentFlags.InputRoot + "/" + file
		fmt.Printf("Converting %v", file)
		processed = writePolygonOutput(inputFile, output)
		fmt.Printf("...%v processed\n", processed)
		totalProcessed += processed
	}
	return totalProcessed
}

func padLeft(str, pad string, length int) string {
	for {
		str = pad + str
		if len(str) > length {
			return str[0:length]
		}
	}
}
