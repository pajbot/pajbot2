CREATE TABLE IF NOT EXISTS `ModerationAction` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `ChannelID` varchar(64) NOT NULL COMMENT 'Twitch Channel owners user ID',
  `UserID` varchar(64) NOT NULL COMMENT 'Source user ID',
  `TargetID` varchar(64) NOT NULL COMMENT 'Target user ID (the user who has banned/unbanned/timed out)',
  `Action` smallint(2) NOT NULL COMMENT 'Action in int format, enums declared outside of SQL',
  `Duration` int(11) DEFAULT NULL COMMENT 'Duration of action (only used for timeouts atm)',
  `Reason` text COMMENT 'Reason for ban. Auto filled in from twich chat, but can be modified in web gui',
  `Timestamp` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'Timestamp of when the timeout occured',
  `Context` text,
  PRIMARY KEY (`id`),
  KEY `ChannelUserTarget_INDEX` (`ChannelID`,`UserID`,`TargetID`),
  KEY `ChannelTargetAction_INDEX` (`ChannelID`,`TargetID`,`Action`)
)
COLLATE='utf8mb4_general_ci'
ENGINE=InnoDB
;
