/*For MySQL Slave*/
CREATE TABLE test.user
(
  id INT AUTO_INCREMENT PRIMARY KEY,
  name VARCHAR(40),
  status VARCHAR(255),
  active bool,
  balance float,
  time time,
  date date,
  datetime datetime,
  created TIMESTAMP NOT NULL
) engine=InnoDB;

CREATE TABLE test.post
(
  id INT AUTO_INCREMENT PRIMARY KEY,
  title VARCHAR(40),
  text VARCHAR(255),
  created TIMESTAMP NOT NULL
) engine=InnoDB;