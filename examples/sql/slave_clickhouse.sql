/*For ClickHouse Slave*/
CREATE DATABASE test;

CREATE TABLE test.user
(
  id Int32,
  name FixedString(40),
  status FixedString(255),
  active UInt8,
  balance Float32,
  time FixedString(255),
  date FixedString(255),
  datetime DateTime,
  created DateTime
) ENGINE = MergeTree()
PARTITION BY id
ORDER BY id;

CREATE TABLE test.post
(
  id Int32,
  title FixedString(40),
  text FixedString(255),
  created DateTime
) ENGINE = MergeTree()
PARTITION BY id
ORDER BY id;