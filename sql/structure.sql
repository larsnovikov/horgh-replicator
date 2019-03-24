/*For MySQL Master and Slave*/
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

/*For ClickHouse Slave*/
CREATE DATABASE test;

CREATE TABLE test.user (id Int32, name FixedString(40), status FixedString(255), created Int32) ENGINE = Log;

CREATE TABLE test.post (id Int32, title FixedString(40), text FixedString(255), created Int32) ENGINE = Log;