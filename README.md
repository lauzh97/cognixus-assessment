# Todolist app
This is a very basic todo-list application. This app allow users to ADD, DELETE, LIST as well as MARK an item as done. This app is mostly executed by API.

## Basic information
This app is written in Golang for logic, with Postgresql for CRUD execution. This app also uses GRPC for the interface for API.

## Docker image
- [App image](https://hub.docker.com/r/lauzh1997/cognixus-assessment-web)
- [Database image](https://hub.docker.com/r/lauzh1997/cognixus-assessment-db)

## Before you run
- Please create a ```config.yaml``` file based on ```config.yaml.example``` in the same directory

## How to run
As requested in the assessment, simply run command:
```
docker compose run
```

## How to use
First thing a user can do is to login the app by navigating to [localhost:8081](http://localhost:8081). This will allow the user to login with their Gmail.

After logging in, a user can perform 4 different actions by using API below:

### 1. Adding a new item into todo-list
```
/v1/todo/add

method: POST
body: {
  itemName string
  itemDescription string
}
```
### 2. Listing all items in todo-list
```
/v1/todo/list

method: GET
body: no body requried
```
### 3. Deleting existing item in todo-list
```
/v1/todo/delete

method: PUT
body: {
    itemName string
}
```
### 4. Marking item as completed in todo-list
```
/v1/todo/mark

method: PUT
body: {
    itemName string
}
```

## How to build
Simply run command:
```
docker compose build
```

## How to test
There are no tests available currently.

## Version
- go 1.21.0
- protoc 24.2
- postgresql 15.4