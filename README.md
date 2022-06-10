# MySqlDbManager
MySql manager is a web service that makes CRUD (create, read, update, delete) queries via HTTP

### HTTP requests supported

GET / - get all tables in database
GET /{table}?limit=20&offset=2 - get 20 entries from 2nd entry
GET /{table}/{id} - get entry whith specified id
PUT /{table} - create entry (columns and values are sent in request body)
POST /{table}/{id} - update entry with specified id (columns and values are sent in request body)
DELETE /{table}/{id} - delete entry with specified id (columns and values are sent in request body)

### Build and launch

Build manager with `go build` command and launch:

```MySqlDbManager -login=my_login -passwd=my_passwd -ip=localhost -port=3306 -db_name=my_database```

*login* - database login
*passwd*- database password
*ip* - database ip
*port* - database port
*db_name* - database name
*max_conns* - max database connections
