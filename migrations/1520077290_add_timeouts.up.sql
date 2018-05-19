CREATE TABLE `ModerationAction` (
  `id` INT(11) UNSIGNED NOT NULL AUTO_INCREMENT,
  `ChannelID` VARCHAR(64) NOT NULL COMMENT 'Twitch Channel owners user ID',
  `UserID` VARCHAR(64) NOT NULL COMMENT 'Source user ID',
  `TargetID` VARCHAR(64) NOT NULL COMMENT 'Target user ID (the user who has banned/unbanned/timed out)',
  `Action` SMALLINT(2) NOT NULL COMMENT 'Action in int format, enums declared outside of SQL',
  `Duration` INT(11) NULL COMMENT 'Duration of action (only used for timeouts atm)',
  `Reason` TEXT NULL COMMENT 'Reason for ban. Auto filled in from twich chat, but can be modified in web gui',
  `Timestamp` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'Timestamp of when the timeout occured',
  PRIMARY KEY (`id`),
  INDEX `ChannelUserTarget_INDEX` (`ChannelID`, `UserID`, `TargetID`),
  INDEX `ChannelTargetAction_INDEX` (`ChannelID`, `TargetID`, `Action`)
)
COLLATE='utf8mb4_general_ci'
;
