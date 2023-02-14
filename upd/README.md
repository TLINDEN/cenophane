## Test the server

### single file upload

curl -X POST localhost:8080/api/putfile -F "upload[]=@/home/scip/2023-02-06_10-51.png" -H "Content-Type: multipart/form-data"

get with:

curl -o x http://localhost:8080/api/getfile/05988587-cc1c-4590-9dc0-564b8a912686/2023-02-06_10-51.png

### multiple files upload

curl -X POST localhost:8080/api/putfile -F "upload[]=@/home/scip/2023-02-06_10-51.png" -F "upload[]=@/home/scip/pgstat.png" -H "Content-Type: multipart/form-data" 

## TODO

- return HTTP errors as JSON as well
- add auto expire customization, e.g. expire after 1 day
- add cleanup storage per api call (single id and everything)
- add oauth2
