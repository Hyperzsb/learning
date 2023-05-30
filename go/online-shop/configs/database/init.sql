-- Drop all tables in dependency order
drop table if exists users;
drop table if exists orders;
drop table if exists order_statuses;
drop table if exists transactions;
drop table if exists transaction_statuses;
drop table if exists customers;
drop table if exists products;

-- Products table
-- drop table if exists products;
create table products
(
    id          int auto_increment,
    name        varchar(127),
    description varchar(1023),
    price       int, -- Without decimals, eg price 10000 = 100.00 in some currency
    inventory   int,
    create_time timestamp default current_timestamp,
    update_time timestamp default current_timestamp on update current_timestamp,
    primary key (id)
);
insert into products (id, name, description, price, inventory)
values (1, 'Product 1', 'This is the first product on this platform', 10000, 100);

-- Customers table
-- drop table if exists customers;
create table customers
(
    id          int auto_increment,
    name        varchar(127),
    email       varchar(127),
    create_time timestamp default current_timestamp,
    update_time timestamp default current_timestamp on update current_timestamp,
    primary key (id)
);

-- Transaction Statuses table
-- drop table if exists transaction_statuses;
create table transaction_statuses
(
    id          int auto_increment,
    status      varchar(31),
    create_time timestamp default current_timestamp,
    update_time timestamp default current_timestamp on update current_timestamp,
    primary key (id)
);
insert into transaction_statuses (id, status)
VALUES (1, 'Pending'),
       (2, 'Declined'),
       (3, 'Cleared'),
       (4, 'Refunded'),
       (5, 'Partially-Refunded');

-- Transactions table
-- drop table if exists transactions;
create table transactions
(
    id             int auto_increment,
    payment_intent varchar(255),
    payment_method varchar(255),
    currency       varchar(3),
    amount         int, -- Without decimals, eg price 10000 = 100.00 in some currency
    card           varchar(4),
    bank_code      varchar(127),
    status_id      int,
    create_time    timestamp default current_timestamp,
    update_time    timestamp default current_timestamp on update current_timestamp,
    primary key (id),
    foreign key (status_id) references transaction_statuses (id) on update cascade on delete cascade
);

-- Order Statuses table
-- drop table if exists order_statuses;
create table order_statuses
(
    id          int auto_increment,
    status      varchar(31),
    create_time timestamp default current_timestamp,
    update_time timestamp default current_timestamp on update current_timestamp,
    primary key (id)
);
insert into order_statuses (id, status)
values (1, 'Pending'),
       (2, 'Cancelled'),
       (3, 'Cleared'),
       (4, 'Refunded'),
       (5, 'Partially-Refunded');

-- Order table
-- drop table if exists orders;
create table orders
(
    id             int auto_increment,
    product_id     int,
    transaction_id int,
    customer_id    int,
    status_id      int,
    quantity       int,
    amount         int,
    create_time    timestamp default current_timestamp,
    update_time    timestamp default current_timestamp on update current_timestamp,
    primary key (id),
    foreign key (product_id) references products (id) on update cascade on delete cascade,
    foreign key (transaction_id) references transactions (id) on update cascade on delete cascade,
    foreign key (customer_id) references customers (id) on update cascade on delete cascade,
    foreign key (status_id) references order_statuses (id) on update cascade on delete cascade
);

-- User table
-- drop table if exists users;
create table users
(
    id          int auto_increment,
    first_name  varchar(63),
    last_name   varchar(63),
    email       varchar(127),
    password    varchar(127),
    create_time timestamp default current_timestamp,
    update_time timestamp default current_timestamp on update current_timestamp,
    primary key (id)
);
insert into users (first_name, last_name, email, password)
values ('admin', 'root', 'admin@root.com', '$2a$10$FTrKW04AaYKsylmQkQCFoOQKeDaG723R6/5BLjW9aBTlog6RnzYMm');