CREATE DATABASE FlaskGoat;

CREATE TABLE `FlaskGoat`.`users` (
    `user_id` INT NOT NULL AUTO_INCREMENT,
    `username` VARCHAR(50) NOT NULL,
    `password` VARCHAR(50) NOT NULL,
    `name` VARCHAR(50),
    PRIMARY KEY (`user_id`)
);