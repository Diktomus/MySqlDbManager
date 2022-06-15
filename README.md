# MySqlDbManager
MySql manager is a web service that makes CRUD (create, read, update, delete) queries via HTTP

### HTTP requests supported

GET / - get all tables in database<br>
GET /{table}?limit=20&offset=2 - get 20 entries from 2nd entry<br>
GET /{table}/{id} - get entry whith specified id<br>
PUT /{table} - create entry (columns and values are sent in request body)<br>
POST /{table}/{id} - update entry with specified id (columns and values are sent in request body)<br>
DELETE /{table}/{id} - delete entry with specified id (columns and values are sent in request body)<br>

### Config description

*login* - database login<br>
*passwd*- database password<br>
*ip* - database ip<br>
*port* - database port<br>
*db_name* - database name<br>
*max_conns* - max database connections<br>
