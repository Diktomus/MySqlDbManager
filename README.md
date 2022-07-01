#dbmanager
dbmanager is a web service that makes CRUD (create, read, update, delete) queries via HTTP

### HTTP requests supported

GET / - get all tables in database<br>
GET /{table}?limit=20&offset=2 - get 20 entries from 2nd entry<br>
GET /{table}/{id} - get entry whith specified id<br>
PUT /{table} - create entry (columns and values are sent in request body)<br>
POST /{table}/{id} - update entry with specified id (columns and values are sent in request body)<br>
DELETE /{table}/{id} - delete entry with specified id (columns and values are sent in request body)<br>

### Config description

*LOGIN* - database login<br>
*PASSWD*- database password<br>
*IP* - database ip<br>
*PORT* - database port<br>
*DB_NAME* - database name<br>
*MAX_CONNS* - max database connections<br>
