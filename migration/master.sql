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

/* This queries to test replication */
INSERT INTO test.user (`name`, `status`) VALUE ("Jack", "active");

INSERT INTO test.post (`title`, `text`) VALUE ("Title", "London is the capital of Great Britain");