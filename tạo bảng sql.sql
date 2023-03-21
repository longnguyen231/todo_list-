CREATE TABLE tasks (
                      id int primary key auto_increment,
                      user_id int,
                      description varchar(55),
                      complete boolean default false,
                      created_at  timestamp default current_timestamp
);
create table users(
                    id int primary key auto_increment,
                    name varchar(30),
                    email varchar(40),
                    password varchar(40),
                    age int,
                    avatar varchar(40)
);