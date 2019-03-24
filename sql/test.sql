/* This queries to test replication user table*/
INSERT INTO test.user (`name`, `status`) VALUE ("Jack", "active");
UPDATE test.user SET `name`='Tommy' ORDER BY RAND() LIMIT 1;
DELETE FROM test.user LIMIT 1;

/* This queries to test replication post table*/
INSERT INTO test.post (`title`, `text`) VALUE ("Title", "London is the capital of Great Britain");
UPDATE test.post SET title='New title' ORDER BY RAND() LIMIT 1;
DELETE FROM test.post LIMIT 1;