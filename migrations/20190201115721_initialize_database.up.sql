CREATE TABLE IF NOT EXISTS currencies
(
name varchar(255) primary key,
symbol varchar(255) not null,
created_at timestamp not null,
updated_at timestamp not null
);

CREATE TABLE IF NOT EXISTS rates
(
id serial primary key,
currency_name varchar(255) references currencies (name),
name varchar(255) not null,
symbol varchar(255) not null,
value numeric(12,2) not null,
created_at timestamp not null,
updated_at timestamp not null
);

CREATE TABLE IF NOT EXISTS accounts
(
id serial primary key,
currency_name varchar(255) references currencies (name),
name varchar(255) not null,
initial_balance numeric(12,2) not null,
created_at timestamp not null,
updated_at timestamp not null
);


CREATE TABLE IF NOT EXISTS transactions
(
id serial primary key, 
account_id int references accounts (id),
description varchar(255),
value numeric(12,2) not null,
type varchar(255) not null,
created_at timestamp not null,
updated_at timestamp not null
);

CREATE TABLE IF NOT EXISTS categories
(
id serial primary key,
name varchar(255),
created_at timestamp not null,
updated_at timestamp not null
);

CREATE TABLE IF NOT EXISTS transactions_categories
(
transaction_id int REFERENCES transactions ON DELETE CASCADE,
category_id int REFERENCES categories,
PRIMARY KEY (transaction_id, category_id)
);
