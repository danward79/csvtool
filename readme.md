# CSVTool

This is a tool I put together to allow manipulation of very large CSV data files, I have to deal with at Work. The files have 111 fields and are often almost 500k records....

The data has subsets of data so grabbing what I need by field and time makes life easier.

It is also possible to retrieve a range/list of records, using -r (line number) instead of -t (time range). Using the -all will output all records according to the other parameters set.

You can use the tool with Stdin and Stdout, to pipe from one tool to another or you can specify an input and or output file.

### Usage
Command switches are as follows:
```go
Tool Usage:
  -all
      provide all records to output
  -blanks int
    	Ignore records if this column is blank (default -1)
  -c string
    	Which columns to export, eg 1-5 or 1,3-10 etc
  -comment string
    	Specifiy the delimiter to use (default "#")
  -delimiter string
    	Specifiy the delimiter to use (default ",")
  -header
    	include header row
  -help
    	help for guidance on usage
  -i string
    	Input CSV file
  -loose
    	Use strict rules for length of a record
  -o string
    	Output CSV file
  -r string
    	Span index of records to export, eg 1-5 or 1,3-10 etc
  -specific int
    	Limit search to a specific column x, default all (slow) (default -1)
  -t string
    	Span of time records, eg 10:00:00-16:00:00
```

for example, the following will read the file and output to the specified file, with a header ignoring the record lengths, columns 0,32 to 85, 96 to 110. In addition it will only match the time on column 0 and use column 32 for ignoring blank lines. Provided the data is between the time span.

```bash
./csvtool -i 502_00409D8C3071_20160524.csv -t "24/05/2016 06:00:00.000 +1000-24/05/2016 18:59:59.999 +1000" -loose  -o subsecondDataTraction.csv -header -specific 0 -c 0,32-85,96-110 -blanks 32
```
