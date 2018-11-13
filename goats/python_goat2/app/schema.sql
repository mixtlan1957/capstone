

CREATE TABLE IF NOT EXISTS `main`.`users`(`user_id` INTEGER PRIMARY KEY AUTOINCREMENT,  
    `name` TEXT, 
    `eid` TEXT, 
    `password` TEXT,   
    `salt` TEXT, 
    `accounting` TEXT
    );