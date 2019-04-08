# Mysql binlog replicator

##### Based on [JackShadow/go-binlog-example](https://github.com/JackShadow/go-binlog-example) 

### Master types
- MySQL

### Slave types
- MySQL
- PostgreSQL
- Yandex ClickHouse

### Quick Start
- Copy `/src/.env.dist` to `/src/.env` and set credentials.
- Configure your my MySQL master as `/mysql/mysql.conf`. 
Don't forget to set `binlog_do_db=<master_db_name>` and restart MySQL service.
- Execute `sql/structure.sql` in your MySQL master and slave.
- Execute `sql/replicator.sql` in your MySQL. It will create database for system values.
- Start Docker as `docker-compose up -d --build`
- Run as `cd src` and `go run replicator.go` in docker container.

### Testing

- Copy `examples/user.json` and `examples/post.json` to `src/configs`
- Execute `sql/test.sql` in your MySQL master and see output.

  ##### OR 

- Execute `cd src` and `go run loader.go`

### Add tables to replicator

- Create json config for your table like `examples/user.json` or `examples.post.json`.
- Create table on slave.