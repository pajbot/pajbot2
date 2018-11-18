CREATE TABLE `UserSession` (
	`id` varchar(64) NOT NULL,
	`user_id` INT(11) UNSIGNED NOT NULL,
    `expiry_date` TIMESTAMP NOT NULL DEFAULT 0,
	PRIMARY KEY (`id`),
    FOREIGN KEY (user_id)
        REFERENCES User(id)
        ON DELETE CASCADE
)
COLLATE='utf8mb4_general_ci'
ENGINE=InnoDB
;
