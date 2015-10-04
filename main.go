package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"
)

const (
	// VERSION is the current application version
	VERSION = "1.0"
)

func printFileList(resp http.ResponseWriter, req *http.Request, dirPath string) {
	log.Println("[200]: serving directory listing of ", req.URL.Path)

	// write success header
	resp.WriteHeader(200)

	// the introductory HTML code to set up the page
	resp.Write([]byte(fmt.Sprintf(`<!doctype html>
		<head><title>Directory Listing</title></head>
		<body><h2>%s</h2><div style="width: 100%%"><table style="width: 100%%"><tr><th width="10%%"></th><th width="60%%"></th><th width="30%%"></th></tr>`,
		strings.TrimSuffix(req.URL.Path, "index.html"))))

	// walk the directory, 1 layer deep only
	filepath.Walk(dirPath, func(p string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// clean up the file path
		p = strings.TrimPrefix(p, dirPath+"/")
		// check if we're trying to recurse too far
		if strings.Contains(p, "/") {
			return nil
		}

		// is it a directory?
		suffix := ""
		if info.IsDir() {
			suffix = "/"
		}

		// print the HTML snippet
		resp.Write([]byte(fmt.Sprintf(`<tr><td></td><td><a href="%s%s">%s%s</a></td><td>%d bytes</td></tr>`, info.Name(), suffix, info.Name(), suffix, info.Size())))

		return nil
	})

	// close the HTML page
	resp.Write([]byte(`</table></div></body></html>`))
}

func serveContent(resp http.ResponseWriter, req *http.Request) {
	cwd, _ := os.Getwd()

	var uri string

	// are we requesting a directory? check for index.html
	if req.URL.Path == "/" {
		uri = "/index.html"
	} else {
		// clean out any ../ to keep root directory security
		uri = path.Clean(req.URL.Path)
	}

	// stat the file to make sure it exists
	stat, err := os.Stat(path.Join(cwd, uri))
	if err != nil && uri != "/index.html" {
		// file doesn't exist and not a directory listing
		resp.WriteHeader(404)
		resp.Write([]byte(fmt.Sprintf("<html><h1><strong>404</strong></h1><h3>File Not Found</h3></html>")))
		log.Println("[404]: file not found '", uri, "'")
		return
	} else if err == nil && uri == "/index.html" {
		// index.html exists, serve it up
	} else if strings.HasSuffix(uri, "/index.html") || stat.IsDir() {
		// index.html requested, but doesn't exist
		// serve up directory listing instead
		printFileList(resp, req, path.Join(cwd, strings.TrimSuffix(uri, "index.html")))
		return
	}

	// open the file for reading
	file, err := os.Open(path.Join(cwd, uri))
	if err != nil {
		resp.WriteHeader(500)
		resp.Write([]byte(fmt.Sprintf("<html><h1><strong>500</strong></h1><h3>Internal Server Error [File Permissions]</h3></html>")))
		log.Println("[500]: unable to open requested file '", uri, "'")
		return
	}

	log.Println("[200]: serving content ", uri)

	// serve the file
	http.ServeContent(resp, req, uri, time.Unix(0, 0), file)

	return
}

func main() {
	fmt.Println("HTTP-Server v", VERSION)

	// handle all content
	http.HandleFunc("/", serveContent)

	// default serve host (*) and port (:8080)
	host := ":8080"

	// did user specify an alternative?
	if len(os.Args) > 1 {
		host = os.Args[1]
	}

	// print server info
	cwd, _ := os.Getwd()

	fmt.Println("Serving", cwd, "on", host)

	// serve content
	http.ListenAndServe(host, nil)
}
