CREATE TABLE `BotChannel` (
	`id` INT(11) UNSIGNED NOT NULL AUTO_INCREMENT,
	`bot_id` INT(11) UNSIGNED NOT NULL,
	`twitch_channel_id` VARCHAR(64) NOT NULL COMMENT 'i.e. 11148817',
	PRIMARY KEY (`id`),
    FOREIGN KEY (bot_id)
        REFERENCES Bot(id)
        ON DELETE CASCADE,
    UNIQUE INDEX `bot_channel` (bot_id, twitch_channel_id)
)
COLLATE='utf8mb4_general_ci'
ENGINE=InnoDB
;
