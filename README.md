# Mysql binlog replicator

##### Based on [JackShadow/go-binlog-example](https://github.com/JackShadow/go-binlog-example) 

### Input
- MySQL

### Output
- MySQL
- Yandex ClickHouse

### Quick Start
- Copy `/src/.env.dist` to `/src/.env` and set credentials.
- Configure your my MySQL master as `/mysql/mysql.conf`. 
Don't forget to set `binlog_do_db=<master_db_name>` and restart MySQL service.
- Execute `sql/structure.sql` in your MySQL master and MySQL slave.
- Execute `sql/replicator.sql` in your MySQL. It will create database for system values.
- Start Docker as `docker-compose up -d --build`
- Run as `cd src` and `go run replicator.go` in docker container.

### Testing

- Execute `sql/test.sql` in your MySQL master and see output.

  ##### OR 

- Execute `cd src` and `go run loader.go`

### Add tables to replicator

- Go to `src/models/slave`.
- Create model for your table like `model_examples/<your_output_type>/user.go` or `model_examples/<your_output_type>post.go`.
- Go to `src/models/registry.go`.
- Add your model to method `GetModel(name string) interface{ AbstractModel }`.
- Create table on MySQL slave.

### Secure columns

If you have private data in column, you can add something like `user.Status = "***"` in function `BeforeSave` in model, and all values of column will be "***".