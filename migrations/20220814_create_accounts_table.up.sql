CREATE TABLE accounts(
  id SERIAL not null primary key,
  username varchar(23) not null,
  password varchar(33) not null,
  nickname varchar(32),
  created_at timestamp default current_timestamp not null,
  updated_at timestamp default current_timestamp not null
);

INSERT INTO accounts(username, password) VALUES('hoge','ea703e7aa1efda0064eaa507d9e8ab7e');
INSERT INTO accounts(username, password, nickname) VALUES('fuga','c32ec965db3295bad074d2afa907b1c3', 'fugafuga');
