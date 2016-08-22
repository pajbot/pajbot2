ALTER TABLE `pb_bot_account` MODIFY twitch_refresh_token VARCHAR(256);
ALTER TABLE `pb_user` MODIFY twitch_refresh_token VARCHAR(256);
ALTER TABLE `pb_channel` DROP COLUMN twitch_access_token;
ALTER TABLE `pb_channel` DROP COLUMN twitch_refresh_token;
