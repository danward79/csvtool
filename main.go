package main

import (
	"bufio"
	"bytes"
	"encoding/csv"
	"flag"
	"fmt"
	"log"
	"os"
)

func main() {

	/*Tools should take command line parameters of
	1. Filename (csv)
	2. Start timedate string
	3. End timedate string
	4. export path or export to existing path.
	*/

	inputFile := flag.String("f", "", "input CSV file")
	start := flag.String("s", "", "start time")
	end := flag.String("e", "", "end time")
	outputFile := flag.String("o", "", "output CSV file")
	//	header := flag.Bool("h", false, "include header row")
	//	preString := flag.String("p", "", "string to preappend to header")
	help := flag.Bool("help", false, "-help for guidance on usage")
	flag.Parse()

	if *inputFile == "" || *outputFile == "" || *start == "" || *end == "" || *help {
		fmt.Println("Tool Usage:")
		flag.PrintDefaults()
		os.Exit(1)
	}

	//var scanner *bufio.Scanner
	var r *csv.Reader
	if *inputFile != "" {
		f, err := os.Open(*inputFile)
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()

		r = csv.NewReader(f)
		//scanner = bufio.NewScanner(f)
	}

	//grep -n "11:00:00" -m 1 501_00409D8C305D_20160523.csv | cut -d : -f 1
	startPosn, err := find(r, *start)
	fmt.Println("Start:", *start, startPosn, err)
	endPosn, err := find(r, *end)
	fmt.Println("End:", endPosn, err)

	/*
		echo TimeDate | tr -d "\n" > header.csv; head -n 1 502_00409D8C3071_20160523.csv >> header.csv
		sed -n "436154,594556 p" 502_00409D8C3071_20160523.csv > data.csv
		cat header.csv data.csv > import.csv; rm data.csv
	*/
}

func find(r *csv.Reader, pat string) (int64, error) {

	count := int64(0)

	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}

		if 
	}

	for scanner.Scan() {
		count++

		if bytes.Contains(scanner.Bytes(), []byte(pat)) {
			return count, nil
		}
	}

	if err := scanner.Err(); err != nil {
		return count, err
	}

	return count, nil
}
