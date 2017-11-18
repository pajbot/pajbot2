CREATE TABLE `pb_twitch_user` (
	`id` INT(11) UNSIGNED NOT NULL AUTO_INCREMENT,
	`twitch_user_id` VARCHAR(20) NULL DEFAULT NULL COMMENT 'i.e. 12345678',
	`twitch_access_token` VARCHAR(64) NULL DEFAULT NULL,
	`twitch_refresh_token` VARCHAR(64) NULL DEFAULT NULL,
	PRIMARY KEY (`id`),
	UNIQUE INDEX `IN_twitch_user_id` (`twitch_user_id`)
)
COLLATE='utf8mb4_general_ci'
ENGINE=InnoDB
;
