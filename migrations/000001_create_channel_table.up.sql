CREATE TABLE `pb_channel` (
	`id` INT(11) UNSIGNED NOT NULL AUTO_INCREMENT,
	`channel` VARCHAR(64) NOT NULL COMMENT 'i.e. forsenlol',
	`nickname` VARCHAR(64) NULL DEFAULT NULL COMMENT 'i.e. Forsen',
	`twitch_channel_id` BIGINT(20) NULL DEFAULT NULL COMMENT 'i.e. 12345678',
	`twitch_access_token` VARCHAR(64) NULL DEFAULT NULL,
	`twitch_refresh_token` VARCHAR(64) NULL DEFAULT NULL,
	PRIMARY KEY (`id`),
	UNIQUE INDEX `channel` (`channel`)
)
COLLATE='utf8mb4_general_ci'
ENGINE=InnoDB
;
