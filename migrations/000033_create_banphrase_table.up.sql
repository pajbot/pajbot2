CREATE TABLE `Banphrase` (
	`id` INT(11) UNSIGNED NOT NULL AUTO_INCREMENT,
    `group_id` INT(11) UNSIGNED NULL DEFAULT NULL,
    `enabled` TINYINT(1) NULL DEFAULT 1 COMMENT 'NULL = Inherit from group',
	`description` TEXT NULL COMMENT 'Optional description of the banphrase, i.e. racism or banned emote',
    `phrase` TEXT NOT NULL COMMENT 'The banned phrase itself. This can be a regular expression, it all depends on the "operator" of the banphrase',
    `length` INT(11) UNSIGNED NULL DEFAULT 60 COMMENT 'NULL = Inherit from group, 0 = permaban, >0 = timeout for X seconds',
    `warning_id` INT(11) UNSIGNED NULL COMMENT 'NULL = Inherit from group, anything else is an ID to a warning "scale"',
    `case_sensitive` TINYINT(1) NULL COMMENT 'NULL = Inherit from group',
    `type` INT(11) NULL DEFAULT 0 COMMENT 'NULL = Inherit from group, 0 = contains, more IDs can be found in the go code lol xd',
    `sub_immunity` TINYINT(1) NULL DEFAULT 0 COMMENT 'NULL = Inherit from group',
    `remove_accents` TINYINT(1) NULL DEFAULT 0 COMMENT 'NULL = Inherit from group',

	PRIMARY KEY (`id`),
    FOREIGN KEY (warning_id)
        REFERENCES WarningScale(id)
        ON DELETE SET NULL,
    FOREIGN KEY (group_id)
        REFERENCES BanphraseGroup(id)
        ON DELETE SET NULL
)
COMMENT='Store banned phrases'
COLLATE='utf8mb4_general_ci'
ENGINE=InnoDB
;
