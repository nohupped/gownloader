package main

import (
	"gownloader/imports"
	"os"
	"fmt"
	"net/http"
	"strconv"
	"sync"
	"io/ioutil"
	"strings"
	"path/filepath"
)
var wg sync.WaitGroup

func main() {

	threads := 2
	switch len(os.Args) {
	case 1 :
		panic("Enter a URL to download\n")
	case 2:
		fmt.Println("Starting downloads")
	default:
		var err error
		if threads, err = strconv.Atoi(os.Args[2]); err == nil {
			 fmt.Println("Starting download")
		}
	}
	splitname := strings.Split(os.Args[1], "/")
	filename := splitname[len(splitname)-1]
	path, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		panic(err)
	}
	fmt.Println("with", threads, "threads")
	fmt.Printf("file will be saved as %s/%s\n", path, filename)
	resp, err := imports.Getlist(os.Args[1])
	if err != nil {
		panic(err)
	}
	if resp.StatusCode != http.StatusOK {
		panic(resp.Status)
	}
	fmt.Println(resp.Status)
	length, err := strconv.Atoi(resp.Header.Get("Content-Length"))
	if err != nil {
		panic(err)
	}
	fmt.Println("Content length is", length, "bytes")
	split := length/threads
	diff := length % threads
	fmt.Println(split, diff)

	for i := 0; i < threads ; i++ {
		wg.Add(1)

		first := split * i
		last := split * (i + 1)

		if (i == threads - 1) {
			last += diff
		}

		go func(min int, max int, i int) {
			slicenum := strconv.Itoa(i)
			name := "." + filename + "." + slicenum
			client := new(http.Client)
			req, _ := http.NewRequest("GET", os.Args[1], nil)
			var range_hdr string
			if  i != (threads - 1) {
				range_hdr = "bytes=" + strconv.Itoa(min) +"-" + strconv.Itoa(max-1) // Add the data for the Range header of the form "bytes=0-100"
			}else {
				range_hdr = "bytes=" + strconv.Itoa(min) +"-" + strconv.Itoa(max) // Add the data for the Range header of the form "bytes=0-100"

			}
			fmt.Println(range_hdr)
			req.Header.Add("Range", range_hdr)
			resp,_ := client.Do(req)
			defer resp.Body.Close()
			reader, _ := ioutil.ReadAll(resp.Body)
			ioutil.WriteFile("." + name, reader, 0644)
			wg.Done()
		}(first, last, i)
	}
	wg.Wait()
}



