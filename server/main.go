package main

import (
	"net/http"
	"io/fs"
	"os"
	"fmt"
	"strings"
)

func main() {
	fmt.Println("Starting server...")
	var sourceDirectory fs.FS = os.DirFS("./src")
	httpFileServer := http.FS(sourceDirectory)
	indexServerHandler := serveFileContents("index.html", httpFileServer)

	http.Handle("/", http.HandlerFunc(func(responseWriter http.ResponseWriter, request *http.Request) {
    indexServerHandler.ServeHTTP(responseWriter, request)
  }))
	port := 8080
	fmt.Println(fmt.Sprintf("Listening at port: %d", port))
	http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
}

func serveFileContents(file string, files http.FileSystem) http.HandlerFunc {
  return func(responseWriter http.ResponseWriter, request *http.Request) {
    // Restrict only to instances where the browser is looking for an HTML file
    if !strings.Contains(request.Header.Get("Accept"), "text/html") {
      responseWriter.WriteHeader(http.StatusNotFound)
      fmt.Fprint(responseWriter, "404 not found")
      return
    }

    file, err := files.Open(file)
    if err != nil {
      responseWriter.WriteHeader(http.StatusNotFound)
      fmt.Fprintf(responseWriter, "File: %s not found", file)
      return
    }

    fileInfo, err := file.Stat()
    if err != nil {
      responseWriter.WriteHeader(http.StatusNotFound)
      fmt.Fprintf(responseWriter, "File Info: %s not found", fileInfo)
      return
    }

    responseWriter.Header().Set("Content-Type", "text/html; charset=utf-8")
    http.ServeContent(responseWriter, request, fileInfo.Name(), fileInfo.ModTime(), file)
  }
}
