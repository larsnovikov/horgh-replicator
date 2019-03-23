CREATE DATABASE Test;

CREATE TABLE Test.User
(
  id INT AUTO_INCREMENT PRIMARY KEY,
  name VARCHAR(40),
  status VARCHAR(255),
  created TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL ON UPDATE CURRENT_TIMESTAMP
) engine=InnoDB;

CREATE TABLE Test.Post
(
  id INT AUTO_INCREMENT PRIMARY KEY,
  title VARCHAR(40),
  text VARCHAR(255),
  created TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL ON UPDATE CURRENT_TIMESTAMP
) engine=InnoDB;

/* This queries to Test replication */
INSERT INTO Test.User (`name`, `status`) VALUE ("Jack", "active");

INSERT INTO Test.Post (`title`, `text`) VALUE ("Title", "London is the capital of Great Britain");