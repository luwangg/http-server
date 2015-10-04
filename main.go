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
	VERSION = "0.2"
)

func printFileList(resp http.ResponseWriter, req *http.Request, dirPath string) {
	log.Println("[200]: serving directory listing of ", req.URL.Path)

	resp.WriteHeader(200)
	resp.Write([]byte(fmt.Sprintf(`<!doctype html>
		<head><title>Directory Listing</title></head>
		<body><h2>%s</h2><div style="width: 100%%"><table style="width: 100%%"><tr><th width="10%"></th><th width="60%%"></th><th width="30%%"></th></tr>`,
		strings.TrimSuffix(req.URL.Path, "index.html"))))

	filepath.Walk(dirPath, func(p string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		p = strings.TrimPrefix(p, dirPath+"/")
		if strings.Contains(p, "/") {
			return nil
		}

		suffix := ""
		if info.IsDir() {
			suffix = "/"
		}

		resp.Write([]byte(fmt.Sprintf(`<tr><td></td><td><a href="%s%s">%s%s</a></td><td>%d bytes</td></tr>`, info.Name(), suffix, info.Name(), suffix, info.Size())))

		return nil
	})

	resp.Write([]byte(`</table></div></body></html>`))
}

func serveContent(resp http.ResponseWriter, req *http.Request) {
	cwd, _ := os.Getwd()

	var uri string

	if req.URL.Path == "/" {
		uri = "/index.html"
	} else {
		uri = path.Clean(req.URL.Path)
	}

	stat, err := os.Stat(path.Join(cwd, uri))
	if err != nil && uri != "/index.html" {
		resp.WriteHeader(404)
		resp.Write([]byte(fmt.Sprintf("<html><h1><strong>404</strong></h1><h3>File Not Found</h3></html>")))

		log.Println("[404]: file not found '", uri, "'")
		return
	} else if err == nil && uri == "/index.html" {

	} else if strings.HasSuffix(uri, "/index.html") || stat.IsDir() {
		printFileList(resp, req, path.Join(cwd, strings.TrimSuffix(uri, "index.html")))
		return
	}

	file, err := os.Open(path.Join(cwd, uri))
	if err != nil {
		resp.WriteHeader(500)
		resp.Write([]byte(fmt.Sprintf("<html><h1><strong>500</strong></h1><h3>Internal Server Error [File Permissions]</h3></html>")))

		log.Println("[500]: unable to open requested file '", uri, "'")

		return
	}

	log.Println("[200]: serving content ", uri)

	http.ServeContent(resp, req, uri, time.Unix(0, 0), file)

	return
}

func main() {
	fmt.Println("HTTP-Server v", VERSION)

	http.HandleFunc("/", serveContent)

	host := ":8080"

	if len(os.Args) > 1 {
		host = os.Args[1]
	}

	cwd, _ := os.Getwd()

	fmt.Println("Serving", cwd, "on", host)

	http.ListenAndServe(host, nil)
}
