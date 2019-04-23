/* This queries to test replication user table*/
INSERT INTO test.user (`name`, `status`, `active`) VALUE ("Jack", "active", false);
UPDATE test.user SET `name`='Tommy', `status`=true ORDER BY RAND() LIMIT 1;
DELETE FROM test.user LIMIT 1;

/* This queries to test replication post table*/
INSERT INTO test.post (`title`, `text`) VALUE ("Title", "London is the capital of Great Britain");
UPDATE test.post SET title='New title' ORDER BY RAND() LIMIT 1;
DELETE FROM test.post LIMIT 1;