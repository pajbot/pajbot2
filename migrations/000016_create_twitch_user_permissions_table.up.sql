CREATE TABLE `twitch_user_permissions` (
    `twitch_user_id` VARCHAR(64) NOT NULL,
    `permission` VARCHAR(64) NOT NULL,

    PRIMARY KEY (`twitch_user_id`),

    UNIQUE INDEX `user_permission` (`twitch_user_id`, `permission`)
)
COLLATE='utf8mb4_general_ci'
ENGINE=InnoDB
;
