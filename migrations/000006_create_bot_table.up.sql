CREATE TABLE `pb_bot` (
	`id` INT(11) UNSIGNED NOT NULL AUTO_INCREMENT,
	`name` VARCHAR(64) NOT NULL COMMENT 'i.e. snusbot',
	`twitch_access_token` VARCHAR(64) NULL DEFAULT NULL COMMENT 'Bot level access-token',
	`twitch_refresh_token` VARCHAR(64) NULL DEFAULT NULL COMMENT 'Bot level refresh-token',
	PRIMARY KEY (`id`)
)
COMMENT='Store available bot accouns, requires an access token with all chat/PubSub scopes'
COLLATE='utf8mb4_general_ci'
ENGINE=InnoDB
;
