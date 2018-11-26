CREATE DATABASE FlaskGoat;

USE FlaskGoat;

CREATE TABLE `users` (
    `user_id` INT NOT NULL AUTO_INCREMENT,
    `username` VARCHAR(50) NOT NULL,
    `password` VARCHAR(50) NOT NULL,
    `name` VARCHAR(50),
    PRIMARY KEY (`user_id`)
);

CREATE TABLE `messages` (
    `msg_id` INT NOT NULL AUTO_INCREMENT,
    `name` VARCHAR(100) NOT NULL,
    `message` VARCHAR(1000) NOT NULL,
    `msg_time` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (`msg_id`)
);