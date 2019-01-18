ALTER TABLE `Bot`
ADD COLUMN `twitch_userid` VARCHAR(64) NOT NULL AFTER `id`,
ADD COLUMN `twitch_access_token_expiry` DATETIME NOT NULL;
