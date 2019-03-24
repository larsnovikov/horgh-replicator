# Mysql binlog replicator

##### Based on [JackShadow/go-binlog-example](https://github.com/JackShadow/go-binlog-example) 

### Input
- MySQL

### Output
- MySQL

### Quick Start
- Copy `/src/.env.dist` to `/src/.env` and set credentials.
- Configure your my MySQL master as `/mysql/mysql.conf`. 
Don't forget to set `binlog_do_db=<master_db_name>` and restart MySQL service.
- Start as `cd src` and `go run main.go`.
- Execute `migration/structure.sql` in your MySQL master and MySQL slave.
- Execute `migration/test.sql` in your MySQL master and see output.