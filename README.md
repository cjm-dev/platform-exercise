# fender - web_api

This repo is currently in development and run in a local or development environment
This repo provides CRUD operations on a User Authentication System. It supports:

* user creation
* user deletion
* user update
* user login 
* user logout

### Requirements
mysql database 

### DB setup
```
make initdb
```

### Clean DB
```
make cleandb
```

### Run the app
```
go run main.go
```
### Usage Examples

To create users

`localhost:3000/signup`

To login a user (POST)  
`localhost:3000/login` 

To logout a user (POST)  
`localhost:3000/logout`

To update a user (PUT)  
`localhost:3000/users/{id}`

To delete a user (DELETE)  
`localhost:3000/users/{id}`

