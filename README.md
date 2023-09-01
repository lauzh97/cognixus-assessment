# Todolist app
This is a very basic todo-list application. This app allow users to ADD, DELETE, LIST as well as MARK an item as done. This app is mostly executed by API.

## Basic information
This app is written in Golang for logic, with Postgresql for CRUD execution. This app also uses GRPC for the interface for API.

## Before you run
- Please create a ```config.yaml``` file based on ```config.yaml.example``` in the same directory

## How to run (docker image)
```
$ docker network create -d bridge todo
$ docker pull lauzh1997/cognixus-assessment-web:v1.0
$ docker pull postgres:15.4
$ docker run --net todo -p 5432:5432 -e POSTGRES_USER=postgres -e POSTGRES_PASSWORD=postgres -e POSTGRES_DATABASE=postgres --hostname db --detach postgres:15.4
$ docker run --net todo -p 8081:8081 -p 8090:8090 --detach lauzh1997/cognixus-assessment-web:v1.0
```

## How to run (docker compose)
- In app directory, run command:
```
$ docker compose up
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
$ docker compose build
```

## How to test
Test file is available in internal/business folder, run command:
```
go test ./internal/business
```

## Version
- go 1.21.0
- protoc 24.2
- postgresql 15.4