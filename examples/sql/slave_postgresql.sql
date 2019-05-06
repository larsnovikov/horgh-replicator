/*For PostgreSQL Slave*/
CREATE DATABASE test;

/*\c test*/

CREATE TABLE public."user"
(
  id BIGSERIAL PRIMARY KEY,
  name VARCHAR(40),
  status VARCHAR(255),
  active bool,
  balance float,
  time time,
  date date,
  datetime VARCHAR(255),
  created TIMESTAMP DEFAULT NULL
);

CREATE TABLE public.post
(
  id BIGSERIAL PRIMARY KEY,
  title VARCHAR(40),
  text VARCHAR(255),
  created TIMESTAMP DEFAULT NULL
);
