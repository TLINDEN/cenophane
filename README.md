# Cenophane
Simple standalone file upload server with expiration

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
      --frontpage string    Content or filename to be displayed on / in case someone visits (default "welcome to upload api, use /api enpoint!")
  -4, --ipv4                Only listen on ipv4
  -6, --ipv6                Only listen on ipv6
  -l, --listen string       listen to custom ip:port (use [ip]:port for ipv6) (default ":8080")
  -p, --prefork             Prefork server threads
  -s, --storagedir string   storage directory for uploaded files (default "/tmp")
      --super string        The API Context which has permissions on all contexts
  -u, --url string          HTTP endpoint w/o path
  -v, --version             Print program version
```

All flags can be set using environment variables, prefix the flag with `CENOD_` and uppercase it, eg: 
```
CENOD_LISTEN=:8080
```

In addition it is possible to set api contexts using env vars (otherwise only possible using the config file):
```
CENOD_CONTEXT_SUPPORT="support:tymag-fycyh-gymof-dysuf-doseb-puxyx"
CENOD_CONTEXT_FOOBAR="foobar:U3VuIE1hciAxOSAxMjoyNTo1NyBQTSBDRVQgMjAyMwo"
```

Configuration can also be done using a config file (searched in the following locations):
- `/etc/cenod.hcl`
- `/usr/local/etc/cenod.hcl`
- `~/.config/cenod/cenod.hcl`
- `~/.cenod`
- `$(pwd)/cenod.hcl`

Or using the flag `-c`. Sample config file:
```
listen = ":8080"
bodylimit = 10000

apicontext = [
  {
    context = "root"
    key = "0fddbff5d8010f81cd28a7d77f3e38981b13d6164c2fd6e1c3f60a4287630c37",
  },
  {
    context = "foo",
    key = "970b391f22f515d96b3e9b86a2c62c627968828e47b356994d2e583188b4190a"
  }
]

#url = "https://sokrates.daemon.de"

# this is the root context with all permissions
super = "root"
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

The client must be configured using a config file. The following locations are searched for it:
- `$(pwd)/upctl.hcl`
- `~/.config/upctl/upctl.hcl`

Sample config file for a client:
```
endpoint = "http://localhost:8080/api/v1"
apikey = "970b391f22f515d96b3e9b86a2c62c627968828e47b356994d2e583188b4190a"
```



## TODO

- also serve a html upload page
- add metrics (as in https://github.com/ansrivas/fiberprometheus)
- do not manually generate output urls, use fiber.GetRoute()
- upd: https://docs.gofiber.io/guide/error-handling/ to always use json output
- upctl: get rid of HandleResponse(), used only once anyway
- add form so that public users can upload



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
