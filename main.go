package main

import (
	"flag"
	"fmt"
	"net/http"
	"path"

	"github.com/spitzfaust/file-server/logger"
)

func check(err error, logger logger.Logger) {
	if err != nil {
		logger.Error("An error occurred: %s", err)
		panic(err)
	}
}

func cleanPrefix(prefix string) string {
	cleanedPrefix := path.Clean(prefix)
	if cleanedPrefix == "." {
		cleanedPrefix = ""
	}
	cleanedPrefix += "/"

	return cleanedPrefix
}

// ProgramName is the program name.
const ProgramName = "file-server"

// Version is the program version.
const Version = "v1.0.0"

func disableCaching(next http.Handler) http.Handler {
	var etagHeaders = []string{
		"ETag",
		"If-Modified-Since",
		"If-Match",
		"If-None-Match",
		"If-Range",
		"If-Unmodified-Since",
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		for _, h := range etagHeaders {
			w.Header().Del(h)
		}
		w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
		w.Header().Set("Pragma", "no-cache")
		w.Header().Set("Expires", "0")
		next.ServeHTTP(w, r)
	})
}

func main() {
	port := flag.String("p", "8100", "port to serve on")
	prefix := flag.String("f", "", "path under which files should be exposed (e.g. /files/img)")
	directory := flag.String("d", ".", "directory of static files on host (e.g. ./documents)")
	version := flag.Bool("v", false, "display the version")
	verbose := flag.Bool("l", false, "enable detailed logs")
	flag.Parse()
	l := logger.NewLoggerWithLevel(logger.Info)
	if *verbose {
		l.SetLogLevel(logger.Debug)
	}
	l.Debug("Port set to: %s", *port)
	l.Debug("Prefix set to: %s", *prefix)
	l.Debug("Directory set to: %s", *directory)
	l.Debug("Version set to %t", *version)
	l.Debug("Verbose set to: %t", *verbose)

	if *version {
		fmt.Printf("%s: %s\n", ProgramName, Version)
		return
	}

	handler := http.FileServer(http.Dir(*directory))

	cleanedPrefix := cleanPrefix(*prefix)

	http.Handle(cleanedPrefix, l.Middleware(disableCaching(http.StripPrefix(cleanedPrefix, handler))))

	l.Info("Serving %s on HTTP port: %s\n", *directory, *port)
	l.Info("Visit: http://localhost:%s%s", *port, cleanedPrefix)
	err := http.ListenAndServe(":"+*port, nil)
	check(err, l)
}
