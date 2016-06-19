package main

/*
CSV Tool to replace workflow needed in CLI for extracting subsets of CSV data from large 800k record files.

Old workflow, took time.
	echo TimeDate | tr -d "\n" > header.csv; head -n 1 502_00409D8C3071_20160523.csv >> header.csv
	grep -n "11:00:00" -m 1 501_00409D8C305D_20160523.csv | cut -d : -f 1
	grep -n "14:00:00" -m 1 501_00409D8C305D_20160523.csv | cut -d : -f 1
	sed -n "436154,594556 p" 502_00409D8C3071_20160523.csv > data.csv
	cat header.csv data.csv > import.csv; rm data.csv

This tool makes that easier.

Tool takes command line parameters of
	1. Filename (csv) or use Stdin
	2. Start timedate string
	3. End timedate string
	4. Export file or use Stdout
*/

import (
	"bytes"
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"sort"
	"time"
)

func main() {

	inputFile := flag.String("i", "", "Input CSV file")
	outputFile := flag.String("o", "", "Output CSV file")
	timeSpan := flag.String("t", "", "Span of time records, eg 10:00:00-16:00:00")
	recordSpan := flag.String("r", "", "Span index of records to export, eg 1-5 or 1,3-10 etc")
	allRecords := flag.Bool("all", false, "provide all records to output")
	columns := flag.String("c", "", "Which columns to export, eg 1-5 or 1,3-10 etc")
	limitColumn := flag.Int("specific", -1, "Limit search to a specific column x, default all (slow)")
	header := flag.Bool("header", false, "include header row")
	loose := flag.Bool("loose", false, "Use strict rules for length of a record")
	delimiter := flag.String("delimiter", ",", "Specifiy the delimiter to use")
	comment := flag.String("comment", "#", "Specifiy the delimiter to use")
	ignoreBlanks := flag.Int("blanks", -1, "Ignore records if this column is blank")
	help := flag.Bool("help", false, "help for guidance on usage")
	flag.Parse()

	//Open input file and if not use os.Stdin (got this idea from mandolyte)
	var r *csv.Reader

	if *inputFile == "" {
		r = csv.NewReader(os.Stdin)
	} else {
		f, err := os.Open(*inputFile)
		if err != nil {
			log.Fatal("Error creating input file:", err)
		}
		defer f.Close()
		r = csv.NewReader(f)
	}

	r.Comma = []rune(*delimiter)[0]
	r.Comment = []rune(*comment)[0]

	//CSV is not perfectly formed, ie. fields per record is inconsistenent then use loose
	if *loose {
		r.FieldsPerRecord = -1
	}

	//Create output file
	var w *csv.Writer
	if *outputFile == "" {
		w = csv.NewWriter(os.Stdout)
	} else {
		fo, err := os.Create(*outputFile)
		if err != nil {
			log.Fatal("Error creating output file:", err)
		}
		defer func() {
			w.Flush()
			fo.Close()
		}()
		w = csv.NewWriter(fo)
	}

	//generate list of columns
	var cl intList
	if *columns != "" {
		var err error
		cl, err = generateIntList(*columns)
		if err != nil {
			log.Fatal("Columnlist Error:", err)
		}
		sort.Sort(cl)
	}

	//keep the header
	if *header {
		writeHeader(r, w, cl)
	}

	switch {
	case *allRecords:
		parseAll(r, w, cl, *ignoreBlanks)

	case *recordSpan != "": //If a record range was specified generate a list of records to capture
		rl, min, max := generateRangeMap(*recordSpan)
		parseForRange(r, w, rl, min, max, cl, *ignoreBlanks)

	case *timeSpan != "": //Sort out the specified timeRange
		tr, err := formatTimeSpan(*timeSpan)
		if err != nil {
			log.Fatal("Parse Time Error:", err)
		}
		sort.Sort(tr)
		parseForTime(r, w, tr, *limitColumn, cl, *ignoreBlanks)

	case *help:
		printUsage("Tool Usage:")

	default:
		printUsage("Tool Usage:")
	}
}

//parseAll output all records
func parseAll(r *csv.Reader, w *csv.Writer, cl intList, ignBlks int) {
	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal("Error Reading CSV:", err)
		}

		if ignBlks >= 0 {
			if record[ignBlks] == "" {
				continue
			}
		}

		if len(cl) > 0 {
			record = remarshallRecord(record, cl)
		}
		writeRecord(record, w)

	}
}

//parseForRange returns a range of records
func parseForRange(r *csv.Reader, w *csv.Writer, rl map[int64]bool, min int64, max int64, cl intList, ignBlks int) {

	var count int64

	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal("Error Reading CSV:", err)
		}

		//stop iterating of we've done the work we need too.
		if count < min {
			count++
			continue
		} else if count > max {
			return
		}

		//record in range list (rl)
		if _, ok := rl[count]; ok {

			if ignBlks >= 0 {
				if record[ignBlks] == "" {
					count++
					continue
				}
			}

			if len(cl) > 0 {
				record = remarshallRecord(record, cl)
			}
			writeRecord(record, w)

			count++
		}

	}
}

//parseForTime exports for a specified time range from - to
func parseForTime(r *csv.Reader, w *csv.Writer, tr timeRange, timeColumn int, cl intList, ignBlks int) {

	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal("Error Reading CSV:", err)
		}

		//Records in Range
		if recordInRange(record, timeColumn, tr) {
			if ignBlks >= 0 {
				if record[ignBlks] == "" {
					continue
				}
			}

			if len(cl) > 0 {
				record = remarshallRecord(record, cl)
			}
			writeRecord(record, w)
		}

	}
}

//generateRangeMap make a map of the range values
func generateRangeMap(rangeSpan string) (map[int64]bool, int64, int64) {
	rm := make(map[int64]bool)

	var err error
	rl, err := generateIntList(rangeSpan)
	if err != nil {
		log.Fatal("Columnlist Error:", err)
	}
	sort.Sort(rl)

	for _, v := range rl {
		rm[v] = true
	}

	return rm, rl[0], rl[len(rl)-1]
}

//writeHeader outputs header and fields as reqd
func writeHeader(r *csv.Reader, w *csv.Writer, cl intList) {
	record, err := r.Read()
	if err == io.EOF {
		return
	}
	if err != nil {
		log.Fatal("Error Reading CSV:", err)
	}

	if len(cl) > 0 {
		record = remarshallRecord(record, cl)
	}
	writeRecord(record, w)
}

//remarshallRecord takes a column list and csv record and returns a new record according to the column list
func remarshallRecord(rec []string, cl intList) []string {
	var r []string

	for _, v := range cl {
		r = append(r, rec[v])
	}

	return r
}

//recordInRange returns true if record contains time between range supplied.
func recordInRange(rec []string, column int, tr timeRange) bool {
	//No specific column specified
	if column < 0 {
		for _, f := range rec {
			v, err := stringToTime(f)
			if err == nil {
				return timeInRange(v, tr)
			}
		}
		return false
	}

	//Specific column specified
	if column >= 0 {
		v, err := stringToTime(rec[column])
		if err != nil {
			log.Println("Error parsing column time:", err)
			return false
		}
		return timeInRange(v, tr)
	}

	return false
}

//time value is between range specified
func timeInRange(v time.Time, tr timeRange) bool {
	return (v.After(tr[0]) || v.Equal(tr[0])) && (v.Before(tr[1]) || v.Equal(tr[1]))
}

//recordContains looks for the pattern provided in a record
func recordContains(rec []string, pats []string) bool {
	for _, v := range rec {
		for _, p := range pats {
			if bytes.Contains([]byte(v), []byte(p)) {
				return true
			}
		}
	}
	return false
}

//writeRecord
func writeRecord(record []string, w *csv.Writer) {
	w.Write(record)
}

//printUsage details with a custom error
func printUsage(msg string) {
	fmt.Println(msg)
	flag.PrintDefaults()
	os.Exit(1)
}
