ALTER TABLE `pb_user`
	ADD COLUMN `type` ENUM('bot','user') NOT NULL DEFAULT 'user' AFTER `name`,
	CHANGE COLUMN `twitch_access_token` `twitch_access_token` VARCHAR(64) NOT NULL COMMENT 'User level access-token' AFTER `type`,
	CHANGE COLUMN `twitch_refresh_token` `twitch_refresh_token` VARCHAR(256) NOT NULL AFTER `twitch_access_token`,
	ADD COLUMN `twitch_room_id` BIGINT NOT NULL AFTER `twitch_refresh_token`,
	ADD INDEX `INDEX_BY_USER_TYPE` (`type`);
