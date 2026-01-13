package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
)

// What this function does:
// Read data from f, extract lines and return read-only channel of strings
func getLinesChannelTest(f io.ReadCloser) <-chan string {
	str := "" //Local var accumulates partial data acts as buffering
	for {     //infinite for loop which keeps on reading from chunks stops when f.Read throws an error
		data := make([]byte, 8) // buffer allocation creating using slice of 8 byte which read 8 byte at a time simulating streaming/chunked reading
		n, err := f.Read(data)
		if err != nil {
			break
		}

		data = data[:n] // trimming of unused bytes as Read() can return fewer bytes even though data was allocated for 8 bytes

		if i := bytes.IndexByte(data, '\n'); i != -1 { //search for new line(\n), i is the index where newline is found
			str += string(data[:i]) // accumulate all the data before the new line and convert to string
			data = data[i+1:]       // move the buffer past \n
			fmt.Printf("read: %s\n", str)
			str = ""
		}
		str += string(data)
	}

	if len(str) != 0 {
		fmt.Printf("read: %s\n", str)
	}
	out := make(chan string)
	// the following code is a goroutine
	go func(s string) {
		defer close(out)
		if len(s) != 0 {
			out <- s
		}
	}(str)
	return out
}

func main() {
	f, err := os.Open("message.txt")

	if err != nil {
		log.Fatal("error", "error", err)
	}

	getLinesChannelTest(f)

	f.Close()
}
