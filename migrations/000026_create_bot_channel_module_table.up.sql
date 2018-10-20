CREATE TABLE `BotChannelModule` (
	`id` INT(11) UNSIGNED NOT NULL AUTO_INCREMENT,
	`bot_channel_id` INT(11) UNSIGNED NOT NULL,
    `module_id` VARCHAR(128) NOT NULL COMMENT 'i.e. nuke',
    `enabled` BOOLEAN NULL COMMENT 'if null, it uses the modules default enabled value',
    `settings` BLOB NULL COMMENT 'json blob with settings',

    PRIMARY KEY(`id`),

    FOREIGN KEY (bot_channel_id)
        REFERENCES BotChannel(id)
        ON DELETE CASCADE,

    UNIQUE INDEX `bot_channel_module` (bot_channel_id, module_id)
)
COLLATE='utf8mb4_general_ci'
ENGINE=InnoDB
;
