CREATE DATABASE replicator;

CREATE TABLE replicator.param_values
(
  id int auto_increment PRIMARY KEY,
  param_key varchar(255) NOT NULL DEFAULT '',
  param_value varchar(255) NOT NULL
) engine=InnoDB;

INSERT INTO replicator.param_values(param_key, param_value) VALUES('last_position_name', '');
INSERT INTO replicator.param_values(param_key, param_value) VALUES('last_position_pos', '');