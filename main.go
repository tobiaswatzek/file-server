package main

import (
	"flag"
	"fmt"
	"net/http"
	"path"
	"strings"
	"time"

	"github.com/justinas/alice"
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

	if !strings.HasPrefix(cleanedPrefix, "/") {
		cleanedPrefix = "/" + cleanedPrefix
	}

	return cleanedPrefix
}

// ProgramName is the program name.
const ProgramName = "file-server"

// Program version. Is set automatically in the build process.
var version = "master"
var date = time.Now().Format("2006-01-02")

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
	showVersion := flag.Bool("v", false, "display the version")
	caching := flag.Bool("c", false, "enable client caching headers")
	verbose := flag.Bool("l", false, "enable detailed logs")
	flag.Parse()
	l := logger.NewLoggerWithLevel(logger.Info)
	if *verbose {
		l.SetLogLevel(logger.Debug)
	}
	l.Debug("Port set to: %s", *port)
	l.Debug("Prefix set to: %s", *prefix)
	l.Debug("Directory set to: %s", *directory)
	l.Debug("Version set to %t", *showVersion)
	l.Debug("Caching set to: %t", *caching)
	l.Debug("Verbose set to: %t", *verbose)

	if *showVersion {
		fmt.Printf("%s\nVersion:  %s\nBuilt on: %s\n", ProgramName, version, date)
		return
	}

	finalHandler := http.FileServer(http.Dir(*directory))
	cleanedPrefix := cleanPrefix(*prefix)

	middlewareChain := alice.New(l.Middleware)
	if !*caching {
		middlewareChain = middlewareChain.Append(disableCaching)
	}

	http.Handle(cleanedPrefix, middlewareChain.Then(http.StripPrefix(cleanedPrefix, finalHandler)))

	l.Info("Serving %s on HTTP port: %s\n", *directory, *port)
	l.Info("Visit: http://localhost:%s%s", *port, cleanedPrefix)
	err := http.ListenAndServe(":"+*port, nil)
	check(err, l)
}
