create table Test.User
(
  id int auto_increment primary key,
  name varchar(40) null,
  status varchar(255) null,
  created timestamp default CURRENT_TIMESTAMP not null on update CURRENT_TIMESTAMP
)
  engine=InnoDB;


INSERT Into Test.User (`name`) VALUE ("Jack");
UPDATE Test.User SET name="Jonh" WHERE id=1;