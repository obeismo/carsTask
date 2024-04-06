CREATE TABLE People
(
    id serial not null unique,
    name varchar(255) not null,
    surname varchar(255) not null,
    patronymic varchar(255)
);

CREATE TABLE Car
(
    id serial not null unique,
    reg_num varchar(255) not null unique,
    mark varchar(255) not null,
    model varchar(255) not null,
    year integer,
    owner integer not null,
    foreign key (owner) references People(id)
);