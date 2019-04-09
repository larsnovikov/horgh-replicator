/*For MySQL Master*/
CREATE DATABASE test;

CREATE TABLE test.user
(
  id INT AUTO_INCREMENT PRIMARY KEY,
  name VARCHAR(40),
  status VARCHAR(255),
  created TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL ON UPDATE CURRENT_TIMESTAMP
) engine=InnoDB;

CREATE TABLE test.post
(
  id INT AUTO_INCREMENT PRIMARY KEY,
  title VARCHAR(40),
  text VARCHAR(255),
  created TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL ON UPDATE CURRENT_TIMESTAMP
) engine=InnoDB;

/*For MySQL Slave*/
CREATE TABLE test.user
(
  id INT AUTO_INCREMENT PRIMARY KEY,
  name VARCHAR(40),
  status VARCHAR(255),
  created TIMESTAMP NOT NULL
) engine=InnoDB;

CREATE TABLE test.post
(
  id INT AUTO_INCREMENT PRIMARY KEY,
  title VARCHAR(40),
  text VARCHAR(255),
  created TIMESTAMP NOT NULL
) engine=InnoDB;

/*For PostgreSQL Slave*/
CREATE DATABASE test;

/*\c test*/

CREATE TABLE public."user"
(
  id BIGSERIAL PRIMARY KEY,
  name VARCHAR(40),
  status VARCHAR(255),
  created TIMESTAMP DEFAULT NULL
);

CREATE TABLE public.post
(
  id BIGSERIAL PRIMARY KEY,
  title VARCHAR(40),
  text VARCHAR(255),
  created TIMESTAMP DEFAULT NULL
);

/*For ClickHouse Slave*/
CREATE DATABASE test;

CREATE TABLE test.user (id Int32, name FixedString(40), status FixedString(255), created DateTime) ENGINE = MergeTree()
PARTITION BY id
ORDER BY id;

CREATE TABLE test.post (id Int32, title FixedString(40), text FixedString(255), created DateTime) ENGINE = MergeTree()
PARTITION BY id
ORDER BY id;

/* For Vertica slave */
/* In bigdata_db */
CREATE TABLE "user"
(
  id INT,
  name VARCHAR,
  status VARCHAR,
  created TIMESTAMP NOT NULL
);

CREATE TABLE post
(
  id INT,
  title VARCHAR,
  text VARCHAR,
  created TIMESTAMP NOT NULL
);