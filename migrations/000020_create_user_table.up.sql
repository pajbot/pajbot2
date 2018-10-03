CREATE TABLE `User` (
    `id` INT(11) UNSIGNED NOT NULL AUTO_INCREMENT,
    `twitch_username` VARCHAR(64) NOT NULL,
    `twitch_userid` VARCHAR(64) NOT NULL,
    `twitch_nonce` VARCHAR(64) NOT NULL,

    PRIMARY KEY (`id`),
    UNIQUE INDEX `ui_twitch_userid` (`twitch_userid`)
)
COLLATE='utf8mb4_general_ci'
ENGINE=InnoDB
;
