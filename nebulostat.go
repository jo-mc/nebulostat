package main

// Copyright ©2020 J McConnell	. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE.txt file.

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"rs"
	"runtime"
	"strconv"
	"time"
)

// Readln returns a single line (without the ending \n)
// from the input buffered reader.
// An error is returned iff there is an error with the
// buffered reader.
// https://stackoverflow.com/questions/6141604/go-readline-string
func Readln(r *bufio.Reader) (string, error) {
	var (
		isPrefix bool  = true
		err      error = nil
		line, ln []byte
	)
	for isPrefix && err == nil {
		line, isPrefix, err = r.ReadLine()
		ln = append(ln, line...)
	}
	return string(ln), err
}

func main() {

	pipeorfileOK() // test we have a pipe in or a file in, will quit if neither or both.

	//  OPEN READER file or PIPE!!
	var r *bufio.Reader
	//var file *os.File
	if len(os.Args) > 1 {
		file, err := os.Open(os.Args[1])
		if err != nil {
			fmt.Println("error opening file= ", err)
			os.Exit(1)
		}
		defer file.Close()
		r = bufio.NewReader(file)
	} else {
		r = bufio.NewReader(os.Stdin)
	}

	fmt.Println("Welcome! NOTE: testing watchdog will stop after coded no of lines.")
	fmt.Println(" .... ")
	fmt.Println("Starting time is", time.Now())
	fmt.Println("ouput ...")
	//var watchdog uint64

	rS := new(rs.RStats)
	medianRq50 := new(rs.RQuant)
	lowerRq50 := new(rs.RQuant)
	upperRq50 := new(rs.RQuant)
	rs.Reinit(lowerRq50, 0.25)
	rs.Reinit(medianRq50, 0.5)
	rs.Reinit(upperRq50, 0.75)

	var linesRead uint32 = 0
	for { // j := 1; j <= ; j++                ENDLESS LOOP
		fdata, e := Readln(r)
		linesRead++
		if e == nil {
			// watchdog = watchdog + 1
			// if watchdog > 1600 {
			// 	fmt.Println("\n Break on watchdog counter, ", watchdog)
			// 	break
			// }
			// fmt.Println(fdata)
			if s, err := strconv.ParseFloat(fdata, 64); err == nil {
				rs.RollingStat(s, rS)
				rs.QuantRoller(s, medianRq50)
				rs.QuantRoller(s, lowerRq50)
				rs.QuantRoller(s, upperRq50)
			}
		} else {
			fmt.Println("end of input: Elements processed: ", linesRead)
			break
		}

	}

	fmt.Println("")
	fmt.Println("Overall results:")
	fmt.Printf("%.2f, The sample Mean\n", rS.M1)
	fmt.Printf("%.2f, Std Dev\n", math.Sqrt(rS.M2/((float64(rS.N))-1.0)))
	fmt.Printf("%.2f, The estimated variance\n", (rS.M2 / (float64(rS.N) - 1)))
	fmt.Printf("%.2f, Largest Value\n", rS.Max)
	fmt.Printf("%.2f, Smallest Value\n", rS.Min)
	fmt.Printf("%.2f, Estimated Median\n", rs.RQuantResult(medianRq50))
	fmt.Printf("%.2f, Standard Deviation of the mean.\n", (math.Sqrt(rS.M2/((float64(rS.N))-1.0)))/math.Sqrt(float64(rS.N)))
	if rS.N > 0 {
		fmt.Printf("%.2f, Skew\n", ((math.Pow(float64(rS.N)-1.0, 1.5)/float64(rS.N))*rS.M3)/(math.Pow(rS.M2, 1.5)))
	} else {
		fmt.Println("0, Skew")
	}
	if rS.N > 0 {
		fmt.Printf("%.2f, Kurtosis\n", (((float64(rS.N)-1.0)/float64(rS.N))*(float64(rS.N)-1.0))*rS.M4/(rS.M2*rS.M2)-3.0)
	} else {
		fmt.Println("0, Kurtosis")
	}
	fmt.Printf("%.2f, Estimated lower quantile\n", rs.RQuantResult(lowerRq50))
	fmt.Printf("%.2f, Estimated middle quantile (Median)\n", rs.RQuantResult(medianRq50))
	fmt.Printf("%.2f, Estimated upper quantile\n", rs.RQuantResult(upperRq50))
	fmt.Printf("%d, Number of items\n", rS.N)

	fmt.Println("Ending time is", time.Now())
	//	fmt.Println(" watchdog count: ", watchdog)

}

func pipeorfileOK() {

	var info os.FileInfo
	var err error
	var inPipe bool = false
	var inFile bool = false

	if runtime.GOOS == "windows" {
		fmt.Println("- -   -  - > Windows detected. Note: Window pipes not implemented, file argument ok.")
	} else {
		// Do we have piped input?
		info, err = os.Stdin.Stat() // standard input file descriptor
		if err != nil {
			fmt.Println("error reading stdin - exiting")
			panic(err)
		}
		if info.Mode()&os.ModeNamedPipe != 0 { // is data begin piped in?
			// we have a pipe input
			fmt.Println("we have a pipe input")
			inPipe = true
		}
	}

	// Do we have argument input?
	//var file *os.File
	if len(os.Args) > 1 { // do we have arguments : ie a file to read?
		fmt.Print(os.Args[1])
		fmt.Println(" : argument (file) input")
		file, err := os.Open(os.Args[1])
		if err != nil {
			fmt.Println("error opening file= ", err)
			os.Exit(1)
		}
		file.Close()
		inFile = true
	}

	if runtime.GOOS != "windows" {
		// Both pipe and argument? -> EXIT
		if inPipe && inFile {
			fmt.Println("- -   -  - > we have a pipe input and a file input ?? Please use one only, exiting")
			os.Exit(1)
		}
	}

	if (inPipe || inFile) == false {
		// no input
		fmt.Println("- -   -  - > No input detected ?? exiting")
		fmt.Println("- -   -  - > Usage: Pipe numbers into program (Linux only)")
		fmt.Println("- -   -  - > awk '{ print $3 }' datafile.dat | nebulostat")
		fmt.Println("- -   -  - > or use with a file argument (Linux or Windows)")
		fmt.Println("- -   -  - > nebulostat datafile.dat")
		fmt.Println("- -   -  - > or awk version")
		fmt.Println("- -   -  - > awk -f nebulostat.awk datafile.dat,   ")
		fmt.Println("- -   -  - > or pipe in:")
		fmt.Println("- -   -  - > awk '{ print $3 }' datafile.dat | awk -f nebulostat.awk")
		fmt.Println("- -   -  - > File input must consist of one number per line.")
		os.Exit(1)
	}
}
