ALTER TABLE `pb_bot_account` MODIFY twitch_refresh_token VARCHAR(64);
ALTER TABLE `pb_user` MODIFY twitch_refresh_token VARCHAR(64);
ALTER TABLE `pb_channel` ADD twitch_access_token VARCHAR(64) NULL DEFAULT NULL;
ALTER TABLE `pb_channel` ADD twitch_refresh_token VARCHAR(64) NULL DEFAULT NULL;
