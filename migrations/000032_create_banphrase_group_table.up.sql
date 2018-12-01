CREATE TABLE `BanphraseGroup` (
	`id` INT(11) UNSIGNED NOT NULL AUTO_INCREMENT,
    `enabled` TINYINT(1) NOT NULL DEFAULT 1,
    `name` VARCHAR(64) NOT NULL,
	`description` TEXT NULL COMMENT 'Optional description of the banphrase group, i.e. racism or banned emote',
    `length` INT(11) UNSIGNED NOT NULL DEFAULT 60 COMMENT '0 = permaban, >0 = timeout for X seconds',
    `warning_id` INT(11) UNSIGNED NULL DEFAULT NULL COMMENT 'ID to a warning "scale"',
    `case_sensitive` TINYINT(1) NOT NULL DEFAULT 0,
    `type` INT(11) NOT NULL DEFAULT 0 COMMENT '0 = contains, more IDs can be found in the go code lol xd',
    `sub_immunity` TINYINT(1) NOT NULL DEFAULT 0,
    `remove_accents` TINYINT(1) NOT NULL DEFAULT 0,

	PRIMARY KEY (`id`),
    FOREIGN KEY (warning_id)
        REFERENCES WarningScale(id)
        ON DELETE SET NULL,
    UNIQUE INDEX `group_name` (name)
)
COMMENT='Store banphrase groups. this will make it easier to manage multiple banphrases at the same time'
COLLATE='utf8mb4_general_ci'
ENGINE=InnoDB
;
