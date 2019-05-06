/* For Vertica slave */
/* In bigdata_db */
CREATE TABLE "user"
(
  id INT,
  name VARCHAR,
  status VARCHAR,
  active bool,
  balance FLOAT,
  time TIME,
  date DATE,
  datetime DATETIME,
  created TIMESTAMP NOT NULL
);

CREATE TABLE post
(
  id INT,
  title VARCHAR,
  text VARCHAR,
  created TIMESTAMP NOT NULL
);