# up
Simple standalone file upload server with api and cli

## TODO

- also serve a html upload page
- add metrics
- add authorization checks for delete and list based on apicontext
- change output of upload, use the same as list
- do not manually generate output urls, use fiber.GetRoute()
- import code from upd into upctl to avoid duplicates, like the time stuff we've now

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
