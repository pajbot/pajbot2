CREATE TABLE `twitch_user_channel_permissions` (
    `twitch_user_id` VARCHAR(64) NOT NULL,
    `channel_id` VARCHAR(64) NOT NULL,
    `permissions` BIT(64) NOT NULL DEFAULT 0b0,

    PRIMARY KEY (`twitch_user_id`),

    UNIQUE INDEX `user_channel_permission` (`twitch_user_id`, `channel_id`)
)
COLLATE='utf8mb4_general_ci'
ENGINE=InnoDB
;
