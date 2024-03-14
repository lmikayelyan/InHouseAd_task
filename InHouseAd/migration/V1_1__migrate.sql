CREATE TABLE IF NOT EXISTS users (
    id serial not null,
    username varchar(45) not null,
    password varchar(100) not null,
    e_mail varchar(45) not null,
    phone_number varchar(45) not null,
    PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS goods (
    id serial not null,
    name varchar(45) not null,
    PRIMARY KEY (id)
);
CREATE INDEX goods_id_index ON goods(id);

CREATE TABLE IF NOT EXISTS categories (
    id serial not null,
    name varchar(45) not null,
    PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS goods_by_categories (
    id serial not null,
    goods_id int not null,
    category_id int not null,
    PRIMARY KEY (id)
);


