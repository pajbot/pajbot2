ALTER TABLE `twitch_user_permissions` CHANGE COLUMN permission permissions BIT(64) NOT NULL DEFAULT 0b0;
