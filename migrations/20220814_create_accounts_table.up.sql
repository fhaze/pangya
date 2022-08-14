CREATE TABLE accounts(
  id SERIAL not null primary key,
  username varchar(23) not null,
  password varchar(33) not null,
  created_at timestamp default current_timestamp not null,
  updated_at timestamp default current_timestamp not null
);

INSERT INTO accounts(username, password) VALUES('hoge','ea703e7aa1efda0064eaa507d9e8ab7e');
