CREATE TABLE `Report` (
	`id` INT(11) UNSIGNED NOT NULL AUTO_INCREMENT,
    `channel_id` VARCHAR(64) NOT NULL,
    `channel_name` VARCHAR(64) NOT NULL,
    `channel_type` VARCHAR(64) NOT NULL,
    `reporter_id` VARCHAR(64) NOT NULL,
    `reporter_name` VARCHAR(64) NOT NULL,
    `target_id` VARCHAR(64) NOT NULL,
    `target_name` VARCHAR(64) NOT NULL,
    `reason` TEXT,
    `logs` TEXT,

    PRIMARY KEY (`id`)
)
COLLATE='utf8mb4_general_ci'
ENGINE=InnoDB
;
