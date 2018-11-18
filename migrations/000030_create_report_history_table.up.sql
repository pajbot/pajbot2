CREATE TABLE `ReportHistory` (
	`id` INT(11) UNSIGNED NOT NULL COMMENT 'report id, same as in Report',
    `channel_id` VARCHAR(64) NOT NULL COMMENT 'twitch ID of channel user was reported in',
    `channel_name` VARCHAR(64) NOT NULL COMMENT 'twitch username of channel the user was reported in',
    `channel_type` VARCHAR(64) NOT NULL,
    `reporter_id` VARCHAR(64) NOT NULL COMMENT 'twitch user ID of reporter',
    `reporter_name` VARCHAR(64) NOT NULL COMMENT 'twitch user name of reporter',
    `target_id` VARCHAR(64) NOT NULL COMMENT 'twitch user ID of person being reported',
    `target_name` VARCHAR(64) NOT NULL COMMENT 'twitch user name of person being reported',
    `reason` TEXT,
    `logs` TEXT,
    `time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'time report was added',

    `handler_id` VARCHAR(64) NOT NULL COMMENT 'twitch user ID of person who handled the report',
    `handler_name` VARCHAR(64) NOT NULL COMMENT 'twitch user name of person who handled the report',

    `action` TINYINT UNSIGNED NOT NULL COMMENT 'number constant for what action was taken for the report. 1 = ban, 2 = timeout, 3 = dismiss',
    `action_duration` INT(11) UNSIGNED NOT NULL DEFAULT 0 COMMENT 'number of seconds for the action. only relevant for timeouts',

    `time_handled` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,

    PRIMARY KEY (`id`)
)
COLLATE='utf8mb4_general_ci'
ENGINE=InnoDB
;
