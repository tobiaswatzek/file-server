# file-server

A simple go file-server that can be used to serve a local directory.

## Usage

```
file-server
  -c    enable client caching headers
  -d string
        directory of static files on host (e.g. ./documents) (default ".")
  -f string
        path under which files should be exposed (e.g. /files/img)
  -l    enable detailed logs
  -p string
        port to serve on (default "8100")
  -v    display the version
 ```
