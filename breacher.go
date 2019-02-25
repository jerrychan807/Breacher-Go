package main

import (
	"bufio"
	"fmt"
	"github.com/akamensky/argparse"
	"github.com/levigross/grequests"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
)

var pl = fmt.Println
var pf = fmt.Printf

const ua = "Mozilla/5.0 (Windows NT 6.1; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/63.0.3239.132 Safari/537.36"

func newRequestOptions() *grequests.RequestOptions {
	return &grequests.RequestOptions{
		DialTimeout:         15 * time.Second,
		TLSHandshakeTimeout: 10 * time.Second,
		RequestTimeout:      15 * time.Second,
		InsecureSkipVerify:  true,
	}
}

func gRequestHead(url string) (*grequests.Response, error) {
	goptions := newRequestOptions()
	goptions.UserAgent = ua
	return grequests.Head(url, goptions)
}

func gRequestGet(url string) (*grequests.Response, error) {
	goptions := newRequestOptions()
	goptions.UserAgent = ua
	return grequests.Get(url, goptions)
}

func drawLabel() {
	//p(`\033[1;34m]______   ______ _______ _______ _______ _     _ _______  ______
	//|_____] |_____/ |______ |_____| |       |_____| |______ |_____/
	//|_____] |    \_ |______ |     | |_____  |     | |______ |    \_
	//
	//\033[37mGo version Made with \033[91m<3\033[37m By Jerry\033[1;m]`)
	pl("\033[1;31m--------------------------------------------------------------------------\033[1;m\n")
}

func parseArgs() (string, string, bool) {
	parser := argparse.NewParser("breacher.go", "Find the admin panel page") // Create new parser object

	url := parser.String("u", "url", &argparse.Options{Required: true, Help: "target url"})                                                                  // Create url flag
	tech_type := parser.String("t", "type", &argparse.Options{Required: false, Help: "set the website technology type i.e. html, asp, php", Default: "all"}) // Create tech_type flag
	fast_mode := parser.Flag("f", "fast", &argparse.Options{Help: "uses goroutines"})                                                                        // Create fast_mode flag

	err := parser.Parse(os.Args) // Parse input
	if err != nil {
		// In case of error print error and print usage
		// This can also be done by passing -h or --help flags
		pl(parser.Usage(err))
		os.Exit(3)
		//p("\033[1;31m[-]\033[1;m -u argument is not supplied. Enter go build breacher.go -h for help")
	}

	return *url, *tech_type, *fast_mode
}

func preHandleUrl(url string) string {
	temp_url := strings.Replace(url, "http://", "", 1) // remove http:// from the url
	temp_url1 := strings.Replace(temp_url, "/", "", 1) // removes / from url so we can have example.com and not example.com/
	start_url := "http://" + temp_url1                 // adds http:// before url, so we have a perfect URL now
	return start_url
}

func findRobotstxt(url string) {
	robotstxt_url := url + "/robots.txt"
	resp, err := http.Get(robotstxt_url)
	if err != nil {
		// handle error
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("\033[1;31m[-]\033[1;m Robots.txt not found\n ", err.Error())
	}
	if strings.Contains(string(body), "<html>") {
		pl("\033[1;31m[-]\033[1;m Robots.txt not found") // if there's an html error page then its not robots.txt
	} else {
		pl("\033[1;32m[+]\033[0m Robots.txt found. Check for any interesting entry\n")
		//pf("\033[1;32m[+]\033[0m Robots.txt url \n")
		pl("================================\n")
		pl(string(body)) // print content of robots.txt
		pl("================================\n")
	}
}

func collectPaths(tech_type string) []string {
	var paths []string
	pl("\033[1;33m[*]\033[0m Your selected technology type is " + tech_type)
	fi, err := os.Open("paths.txt")
	if err != nil {
		fmt.Printf("Error: %s\n", err)
	}
	defer fi.Close()

	br := bufio.NewReader(fi)
	for {
		tmp, _, c := br.ReadLine()
		path := strings.Replace(string(tmp), "\n", "", 1)
		if c == io.EOF {
			break
		}

		if tech_type == "asp" {
			if strings.Contains(path, "html") || strings.Contains(path, "php") {

			} else {
				paths = append(paths, path)
			}
		}
		if tech_type == "php" {
			if strings.Contains(path, "asp") || strings.Contains(path, "html") {

			} else {
				paths = append(paths, path)
			}
		}
		if tech_type == "html" {
			if strings.Contains(path, "asp") || strings.Contains(path, "php") {

			} else {
				paths = append(paths, path)
			}
		}
		if tech_type == "all" {
			paths = append(paths, path)
		}

	}
	fmt.Printf("\033[1;33m[*]\033[0m collected path length is %d\n", len(paths))
	return paths
}

func scan(url string, links []string) {
	pl("\033[1;33m[*]\033[0m Start to scan ")
	for _, link := range links {
		full_url := url + link
		sendRequest(full_url)
	}
}

func divided(links []string, goroutineNum int) [][]string {
	chunkSize := (len(links) + goroutineNum - 1) / goroutineNum

	var dividedPath [][]string
	for i := 0; i < len(links); i += chunkSize {
		end := i + chunkSize

		if end > len(links) {
			end = len(links)
		}

		dividedPath = append(dividedPath, links[i:end])
	}

	return dividedPath

}

func sendRequest(full_url string) {

	resp, err := gRequestGet(full_url)
	body := string(resp.Bytes())
	// Not the usual JSON so copy and paste from below
	if err != nil {
		pl("Unable to make request", err)
	}
	if resp.StatusCode == 200 {
		if (strings.Contains(string(body), "type=\"password\"")) {
			pf("\033[1;32m[+]\033[0m Admin panel found: %s\n", full_url)
		}else {
			pf("\033[1;31m[-]\033[1;m %s\n", full_url)
		}
	} else if resp.StatusCode == 404 {
		pf("\033[1;31m[-]\033[1;m %s\n", full_url)
	} else if resp.StatusCode == 302 {
		pf("\033[1;32m[+]\033[0m Potential EAR vulnerability found : %s\n", full_url)
	} else {
		pf("\033[1;31m[-]\033[1;m %s\n", full_url)
	}
}

func main() {

	drawLabel()
	url, tech_type, fast_mode := parseArgs()

	startTime := time.Now() // get current time
	start_url := preHandleUrl(url)
	findRobotstxt(start_url)

	collected_path := collectPaths(tech_type)
	if fast_mode {
		dividedLinks := divided(collected_path, 5)

		var wg sync.WaitGroup
		for _, link := range dividedLinks {
			wg.Add(1) // Increment the WaitGroup counter.
			go func(link []string) {
				// Launch a goroutine to fetch the link.
				scan(start_url, link)
				// Fetch the link.
				wg.Done()
			}(link)
		}
		wg.Wait() // Wait for all goroutines to finish.
	} else {
		scan(start_url, collected_path)
	}

	elapsed := time.Since(startTime)
	fmt.Println("elapsed time: ", elapsed)
}
