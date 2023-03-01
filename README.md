# up
Simple standalone file upload server with api and cli

## TODO

- implement goroutine to expire after 1d, 10m etc
  implemented. add go routine to server, use Db.Iter()
- use bolt db to retrieve list of items to expire
- also serve a html upload page
- add auth options (access key, users, roles, oauth2)
- add metrics
- add upctl command to remove a file
- use global map of api endpoints like /file/get/ etc
- create cobra client commands (upload, list, delete, edit)


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
