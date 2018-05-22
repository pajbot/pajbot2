CREATE TABLE `twitch_users` (
    `twitch_user_id` VARCHAR(64) NOT NULL,
    `name` VARCHAR(64) NOT NULL COMMENT 'i.e. testaccount_420',
    `display_name` VARCHAR(64) NULL COMMENT 'i.e. TestAccount_420',
    PRIMARY KEY (`twitch_user_id`),
    INDEX `name` (`name`)
)
COLLATE='utf8mb4_general_ci'
ENGINE=InnoDB
;
