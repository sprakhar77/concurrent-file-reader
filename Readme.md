## Problem statement

Given a cookie log file in the following format:

```
cookie,timestamp
AtY0laUfhglK3lC7,2018-12-09T14:19:00+00:00
SAZuXPGUrfbcn5UA,2018-12-09T10:13:00+00:00
5UAVanZf6UtGyKVS,2018-12-09T07:25:00+00:00
AtY0laUfhglK3lC7,2018-12-09T06:19:00+00:00
SAZuXPGUrfbcn5UA,2018-12-08T22:03:00+00:00
4sMM2LxV07bPJzwf,2018-12-08T21:30:00+00:00
fbcn5UAVanZf6UtG,2018-12-08T09:30:00+00:00
4sMM2LxV07bPJzwf,2018-12-07T23:30:00+00:00
```
### Solution:

This distributed log file reader concurrently reads chunks of the file (chunk size can be specified) via a user specified thread pool. All the results read are then summed up to find the most active cookie. The program can also be run across diferent nodes to provide multi layer distribution of work. For eg. The program can run in map-reduce style on 20 nodes, each of which run 1000 go routines to read chunks of a file that is in terabytes or even petabytes.

### Assumptions:
- If multiple cookies meet that criteria, all of them are returned on separate lines.
- You can assume -d parameter takes date in UTC time zone.
- Cookies in the log file are sorted by timestamp (most recent occurrence is the first line of the file).


### How to build?
Note - Make sure you have go installed

`go build -o a cmd/main.go`

### How to run tests?
`go test ./...`

### How to run?

`./a -f <path-to-file> -d <date>`

For eg - to test run this project on the sample input provided, run the command below

`./a -f test.txt -d 2018-12-09`
