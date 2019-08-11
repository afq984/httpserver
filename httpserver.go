package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/handlers"
)

// We track the time that the server started, to make `ModTime()` of
// the served files return `serverStart` if the they have a timestamp
// that is before `serverStart`. This is due to:
// 1) Browsers use the `If-Modified-Since` header to ask the server
//    only return an actual response if the resource at the given URL
//    has been modified after the given timestamp, or otherwise the
//    server would return `304 Not Modified`.
// 2) Users are expected to run this tool to serve different
//    directories on the same address between different invocations
//    of this tool. So a browsers may compare timestamps with the
//    server while they are actually refering to different files.
//    A user switching serving directories may see its browser
//    not reloading properly.
// Never returning a time before serverStart ensures that cached resources
// belonging to an earlier server invocation will always be invalidated.
var serverStart time.Time

type afterServerStartFileInfo struct {
	os.FileInfo
}

func (fi afterServerStartFileInfo) ModTime() time.Time {
	t := fi.FileInfo.ModTime()
	if t.Before(serverStart) {
		return serverStart
	}
	return t
}

type afterServerStartFile struct {
	http.File
}

func (f afterServerStartFile) Stat() (os.FileInfo, error) {
	info, err := f.File.Stat()
	if info == nil {
		return info, err
	}
	return afterServerStartFileInfo{info}, err
}

type afterServerStartFileSystem struct {
	http.FileSystem
}

func (fs afterServerStartFileSystem) Open(name string) (http.File, error) {
	f, err := fs.FileSystem.Open(name)
	if f == nil {
		return f, err
	}
	return afterServerStartFile{f}, err
}

// NoPermanent3XX is a wrapper handler that changes all permanent (301,308) redirects
// to 307 so browsers will not remember redirects that is from a previous
// http server invocation
type NoPermanent3XX struct {
	http.Handler
}

type noPermanent3XXResponseWriter struct {
	http.ResponseWriter
}

func (h NoPermanent3XX) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.Handler.ServeHTTP(noPermanent3XXResponseWriter{w}, r)
}

func (w noPermanent3XXResponseWriter) WriteHeader(statusCode int) {
	switch statusCode {
	case http.StatusMovedPermanently, http.StatusPermanentRedirect:
		statusCode = http.StatusTemporaryRedirect
	}
	w.ResponseWriter.WriteHeader(statusCode)
}

func main() {
	port := flag.Int("port", 8000, "listen on which port")
	dir := flag.String("dir", ".", "directory to serve")
	flag.Parse()
	fmt.Printf("listening on port: %d\nserving directory: %s\n", *port, *dir)
	serverStart = time.Now()
	err := http.ListenAndServe(
		fmt.Sprintf(":%d", *port),
		handlers.LoggingHandler(
			os.Stderr,
			NoPermanent3XX{http.FileServer(afterServerStartFileSystem{http.Dir(*dir)})}))
	panic(err)
}
