package main

import (
	"flag"
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

func main() {
	port := flag.String("p", "8100", "port to serve on")
	prefix := flag.String("f", "", "the path under which files should be exposed (e.g. /files/img)")
	directory := flag.String("d", ".", "the directory of static files on host (e.g. ./documents)")
	verbose := flag.Bool("v", false, "if set the output is more detailed")
	flag.Parse()
	l := logger.NewLoggerWithLevel(logger.Info)
	if *verbose {
		l.SetLogLevel(logger.Debug)
	}
	l.Debug("Port set to: %s", *port)
	l.Debug("Prefix set to: %s", *prefix)
	l.Debug("Directory set to: %s", *directory)
	l.Debug("Verbose set to: %t", *verbose)

	handler := http.FileServer(http.Dir(*directory))

	cleanedPrefix := cleanPrefix(*prefix)

	http.Handle(cleanedPrefix, http.StripPrefix(cleanedPrefix, handler))

	l.Info("Serving %s on HTTP port: %s\n", *directory, *port)
	l.Info("Visit: http://localhost:%s%s", *port, cleanedPrefix)
	e := http.ListenAndServe(":"+*port, nil)
	l.Error("%s", e)
}
