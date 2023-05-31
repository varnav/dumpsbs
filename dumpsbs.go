package main

import (
	"bufio"
	"flag"
	"io"
	"log"
	"net"
	"os"
	"strings"
	"time"
)

var (
	writer *bufio.Writer
	f      *os.File
)

func main() {
	// Command line arguments
	host := flag.String("host", "localhost:30003", "The host to connect to")
	outputDir := flag.String("output", ".", "The output directory")
	verbose := flag.Bool("v", false, "Print not just to file but also to STDOUT")
	vnf := flag.Bool("vnf", false, "Print filtered out data to STDOUT")
	flag.Parse()

	// Connect to the TCP server
	conn, err := net.Dial("tcp", *host)
	if err != nil {
		log.Fatal(err)
	}
	println("Connected to", *host)
	defer conn.Close()

	// Read data from the TCP connection
	reader := bufio.NewReader(conn)
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			if err != io.EOF {
				log.Fatal(err)
			}
			break
		}

		// Check if we need to create a new writer
		currentFile := getCurrentFileName(*outputDir)
		if writer == nil || f.Name() != currentFile {
			if f != nil {
				f.Close()
			}
			f, writer = createWriter(currentFile)
			println("Writing to file:", currentFile)
			defer f.Close()
		}

		// Write the data to the file if necessary values are found
		linearr := strings.Split(line, ",")
		// Field numbers, 12 equals 11 in array: http://woodair.net/sbs/article/barebones42_socket_data.htm
		if len(line) > 10 && (len(linearr[11]) > 0 || len(linearr[17]) > 0 || len(linearr[10]) > 0) { // Data not empty, line is complete
			if *verbose {
				print(line)
			}
			_, err = writer.WriteString(line)
			if err != nil {
				log.Fatal(err)
			}
		} else if *vnf {
			print(line)
		}
		writer.Flush()
	}
}

func getCurrentFileName(outputDir string) string {
	t := time.Now()
	return outputDir + "/" + t.Format("20060102_15") + ".txt"
}

func createWriter(filename string) (*os.File, *bufio.Writer) {
	f, err := os.Create(filename)
	if err != nil {
		log.Fatal(err)
	}

	writer := bufio.NewWriter(f)
	return f, writer
}
