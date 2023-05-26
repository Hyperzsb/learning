-- Products table
drop table if exists products;
create table products
(
    id          int,
    name        varchar(100),
    description varchar(1000),
    price       int, -- Without decimals, eg price 10000 = 100.00 in some currency
    inventory   int,
    create_time timestamp default current_timestamp,
    update_time timestamp default current_timestamp on update current_timestamp,
    primary key (id)
);
insert into products (id, name, description, price, inventory)
values (1, 'Product 1', 'This is the first product on this platform', 10000, 100);

-- Transaction Statuses table
drop table if exists transaction_statuses;
create table transaction_statuses
(
    id          int,
    status      varchar(20),
    create_time timestamp default current_timestamp,
    update_time timestamp default current_timestamp on update current_timestamp,
    primary key (id)
);
insert into transaction_statuses (id, status)
VALUES (1, 'Pending'),
       (2, 'Cleared'),
       (3, 'Declined'),
       (4, 'Refunded'),
       (5, 'Partially-Refunded');

-- Transactions table
drop table if exists transactions;
create table transactions
(
    id          int,
    currency    varchar(3),
    amount      int, -- Without decimals, eg price 10000 = 100.00 in some currency
    card        varchar(4),
    bank_code   varchar(100),
    status_id   int,
    create_time timestamp default current_timestamp,
    update_time timestamp default current_timestamp on update current_timestamp,
    primary key (id),
    foreign key (status_id) references transaction_statuses (id) on update cascade on delete cascade
);

-- Order Statuses table
drop table if exists order_statuses;
create table order_statuses
(
    id          int,
    status      varchar(20),
    create_time timestamp default current_timestamp,
    update_time timestamp default current_timestamp on update current_timestamp,
    primary key (id)
);
insert into order_statuses (id, status)
values (1, 'Cancelled'),
       (2, 'Cleared'),
       (3, 'Refunded');

-- Order table
drop table if exists orders;
create table orders
(
    id             int,
    product_id     int,
    transaction_id int,
    status_id      int,
    quantity       int,
    amount         int,
    create_time    timestamp default current_timestamp,
    update_time    timestamp default current_timestamp on update current_timestamp,
    primary key (id),
    foreign key (product_id) references products (id) on update cascade on delete cascade,
    foreign key (transaction_id) references transactions (id) on update cascade on delete cascade,
    foreign key (status_id) references order_statuses (id) on update cascade on delete cascade
);

-- User table
drop table if exists users;
create table users
(
    id          int,
    first_name  varchar(50),
    last_name   varchar(50),
    email       varchar(100),
    passwd      varchar(100),
    create_time timestamp default current_timestamp,
    update_time timestamp default current_timestamp on update current_timestamp,
    primary key (id)
);