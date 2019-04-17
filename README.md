# Mysql binlog replicator

##### Based on [JackShadow/go-binlog-example](https://github.com/JackShadow/go-binlog-example) 

### Master types
- MySQL

### Slave types
- MySQL
- PostgreSQL
- Yandex ClickHouse
- HP Vertica

### Quick Start
- Copy `/src/.env.dist` to `/src/.env` and set credentials.
- Configure your my MySQL master as `/mysql/mysql.conf`. 
Don't forget to set `binlog_do_db=<master_db_name>` and restart MySQL service.
- Execute `sql/structure.sql` in your MySQL master and slave.
- Execute `sql/replicator.sql` in your MySQL. It will create database for system values.
- Start Docker as `docker-compose up -d --build`
- Run as `cd src` and `go run main.go listen` in docker container.

### Testing

- Copy `examples/user.json` and `examples/post.json` to `src/configs`
- Execute `sql/test.sql` in your MySQL master and see output.

  ##### OR 

- Execute `cd src` and `go run main.go load`

### Add tables to replicator

- Create json config for your table like `examples/user.json` or `examples.post.json`.
- Create table on slave.

### Custom handlers for field value

- Create `plugins/user/<plugin_name>/handler.go` like `create plugins/system/set_value/handler.go`
- Execute `go build -buildmode=plugin -o plugins/user/<plugin_name>/handler.so plugins/user/<plugin_name>/handler.go`
- Add to field description in your `<model>.json`

```
"beforeSave": {
  "method": "user/<plugin_name>",
  "params": [
    "***"
  ]
}
```