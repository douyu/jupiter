CREATE TABLE `user` (
`id`  int NOT NULL AUTO_INCREMENT ,
`username`  varchar(255) NOT NULL DEFAULT '' ,
`password`  varchar(255) NOT NULL DEFAULT '' ,
`nickname`  varchar(255) NOT NULL DEFAULT '' ,
`address`  varchar(255) NOT NULL DEFAULT '' ,
PRIMARY KEY (`id`)
)ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;
insert into user (id,username,password,nickname,address)VALUES(null,"admin","123456","rose","WUHAN");