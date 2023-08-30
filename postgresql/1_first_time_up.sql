create schema main

create table if not exists main.todolist(
    id varchar(36) primary key,
    active boolean,
    createdOn timestamp with time zone,
    updatedOn timestamp with time zone
);

create table if not exists main.user(
    id varchar(36) primary key,
    email varchar(64),
    todoListId varchar(36),
    active boolean,
    createdOn timestamp with time zone,
    updatedOn timestamp with time zone,
    constraint fk_todoListId_user foreign key(todoListId) references main.todolist(id)
);

create table if not exists main.session(
    id varchar(36) primary key,
    userId varchar(36),
    active boolean,
    createdOn timestamp with time zone,
    expriedOn timestamp with time zone,
    updatedOn timestamp with time zone,
    constraint fk_userId foreign key(userId) references main.user(id)
);

create table if not exists main.item(
    id varchar(36) primary key,
    todoListId varchar(36),
    name varchar(64),
    description varchar(128),
    markDone boolean,
    active boolean,
    createdOn timestamp with time zone,
    updatedOn timestamp with time zone,
    constraint fk_todoListId_item foreign key(todoListId) references main.todolist(id)
);