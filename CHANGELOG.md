# Changelog

## Unversioned

- Remove Twitter support. (#1033)  
  Warning messages will be posted in the console if twitter tokens are configured.
- Add GitHub push event webhook support. (#1042, #1043)  
   Webhook format: `https://your-bot-domain.com/api/webhook/github/{channel_id}`  
   Example config has been updated to show Auth -> Github -> Webhook -> Secret

## v2.0.0 - 2023-08-08

- Bumped minimum Go version to 1.19. (#898)
- The nuke module now has tests for parsing parameters. (#530)
- The `Auth->Twitch->Webhook->Secret` config value is now REQUIRED. It's your own private secret you need to generate yourself, and it must be at least 10 characters and at most 100 characters long.
- The nuke module will now recognize users with global permissions. (#268)
- Message height limit no longer applies to Twitch Moderators (#89, #228)
- The version of MessageHeightTwitch was updated, which requires version 3.0 of .NET Core.
- Changed DB backend from MySQL to PostgreSQL.  
  Setting up the bot from scratch? You don't need to do anything!  
  Upgrading your already set up bot to this version? Follow the instructions in [this document](/resources/mysql-to-postgresql-transition/README.md)
