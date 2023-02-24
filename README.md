# up
Simple standalone file upload server with api and cli

## TODO

- decouple db and http code in Runserver()
- store ts
- implement goroutine to expire after 1d, 10m etc
- use bolt db to retrieve list of items to expire
- return a meaningful  message if a file has expired,  not just a 404,
  that is: do  remove the file when it expires  but not the associated
  db entry.
- also serve a html upload page
- add auth options (access key, users, roles, oauth2)
- add metrics
- add upctl command to remove a file
- upd: add short  uuid to files, in case multiple  files with the same
  name are being uploaded
- use global map of api endpoints like /file/get/ etc
- use separate group for /file/
