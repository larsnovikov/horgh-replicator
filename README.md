# Mysql binlog replicator

##### Based on [JackShadow/go-binlog-example](https://github.com/JackShadow/go-binlog-example) 

### Input
- MySQL

### Output
- MySQL

### Quick Start
- Copy `/src/.env.dist` to `/src/.env` and set credentials
- Configure your my MySQL master as `/mysql/mysql.conf`
- Start as `cd src` && `go run main.go`
- Execute `migration/master.sql` in your MySQL master and see output
