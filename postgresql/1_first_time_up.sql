create schema main

create table if not exists main.todolist(
    id varchar(36) primary key,
    active boolean default true,
    createdOn timestamp with time zone default current_timestamp,
    updatedOn timestamp with time zone default current_timestamp
);

create table if not exists main.user(
    id varchar(36) primary key,
    email varchar(64),
    todoListId varchar(36),
    active boolean default true,
    createdOn timestamp with time zone default current_timestamp,
    updatedOn timestamp with time zone default current_timestamp,
    constraint fk_todoListId_user foreign key(todoListId) references main.todolist(id)
);

create table if not exists main.item(
    id varchar(36) primary key,
    todoListId varchar(36),
    name varchar(64),
    description varchar(128),
    markDone boolean default false,
    active boolean default true,
    createdOn timestamp with time zone default current_timestamp,
    updatedOn timestamp with time zone default current_timestamp,
    constraint fk_todoListId_item foreign key(todoListId) references main.todolist(id)
);