# Cenophane
Simple standalone file upload server with expiration

## Server Usage

```
cenod -h
      --apikeys strings     Api key[s] to allow access
  -a, --apiprefix string    API endpoint path (default "/api")
  -n, --appname string      App name to say hi as (default "cenod v0.0.1")
  -b, --bodylimit int       Max allowed upload size in bytes (default 10250000000)
  -c, --config string       custom config file
  -D, --dbfile string       Bold database file to use (default "/tmp/uploads.db")
  -d, --debug               Enable debugging
  -4, --ipv4                Only listen on ipv4
  -6, --ipv6                Only listen on ipv6
  -l, --listen string       listen to custom ip:port (use [ip]:port for ipv6) (default ":8080")
  -p, --prefork             Prefork server threads
  -s, --storagedir string   storage directory for uploaded files (default "/tmp")
      --super string        The API Context which has permissions on all contexts
  -u, --url string          HTTP endpoint w/o path
  -v, --version             Print program version
```

## Client Usage

```
upctl 
Error: No command specified!
Usage:
  upctl [options] [flags]
  upctl [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  delete      Delete an upload
  describe    Describe an upload.
  download    Download a file.
  help        Help about any command
  list        List uploads
  upload      Upload files

Flags:
  -a, --apikey string     Api key to use
  -c, --config string     custom config file
  -d, --debug             Enable debugging
  -p, --endpoint string   upload api endpoint url (default "http://localhost:8080/api/v1")
  -h, --help              help for upctl
  -r, --retries int       How often shall we retry to access our endpoint (default 3)
  -v, --version           Print program version

Use "upctl [command] --help" for more information about a command.
```

## Features

- RESTful API
- Authentication and Authorization through bearer api token
- multiple tenants supported (tenant == api context)
- Each upload gets its own unique id
- download uri is public, no api required, it is intended for end users
- uploads may consist of one or multiple files
- zipped automatically
- uploads expire, either as soon as it gets downloaded or when a timer runs out
- the command line client uses the api
- configuration using HCL language
- docker container build available
- the server supports config by config file, environment variables or flags
- restrictive defaults

## TODO

- also serve a html upload page
- add metrics (as in https://github.com/ansrivas/fiberprometheus)
- add authorization checks for delete and list based on apicontext
- do not manually generate output urls, use fiber.GetRoute()
- import code from upd into upctl to avoid duplicates, like the time stuff we've now
- upd: https://docs.gofiber.io/guide/error-handling/ to always use json output
- upctl: get rid of HandleResponse(), used only once anyway
- add form so that public users can upload
- add support for custom front page

## BUGS

### upctl HTTP 413 weird behavior

- with -d reports correctly the 413, w/o it, it reports the timeout before.

## curl commands

### upload

```
curl -X POST localhost:8080/api/putfile -F "upload[]=@xxx" -F "upload[]=@yyy" -H "Content-Type: multipart/form-data"
```

### download
```
curl -O http://localhost:8080/api/v1/file/388f41f4-3f0d-41e1-a022-9132c0e9e16f/2023-02-28-18-33-xxx
```

### delete
```
curl -X DELETE http://localhost:8080/api/v1/file/388f41f4-3f0d-41e1-a022-9132c0e9e16f/
curl -X DELETE http://localhost:8080/api/v1/file/?id=388f41f4-3f0d-41e1-a022-9132c0e9e16f/
curl -X DELETE -H "Accept: application/json"  -d '{"id":"$id"}' http://localhost:8080/api/v1/file/
```
