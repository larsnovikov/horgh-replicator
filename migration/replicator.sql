CREATE DATABASE replicator;

CREATE TABLE replicator.param_values
(
  param_key varchar(255) UNIQUE PRIMARY KEY,
  param_value varchar(255) NOT NULL
) engine=InnoDB;