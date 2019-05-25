# pajbot 2 [![CircleCI](https://circleci.com/gh/pajbot/pajbot2.svg?style=svg)](https://circleci.com/gh/pajbot/pajbot2)

A rewrite/restructuring of [pajbot 1](https://github.com/pajbot/pajbot) in Golang.

## Authors
 * [nuuls](https://github.com/nuuls)
 * [pajlada](https://github.com/pajlada)
 * [gempir](https://github.com/gempir)


## Web guide
* `cd web && npm install`
* `npm run watch` to let webpack running and compile in background

## FAQ
### After pulling the latest version, something went wrong. what should I do?
Try running `./bot fix`
### How do I add a bot?
After making yourself an admin in the config file, open up the web interface. Log in, go to `/admin`, press the "Log in as bot", then after authenticating whatever user you want as a bot, restart the bot!
