# Changelog

## Unversioned

## 2.0-rc3

- Breaking: Go 1.25.8 is minimum required version to build. (#19)
- Breaking: Removed viper as a config library, moved to using only the `config.yaml` file. (#20)
- Major: Split twitch connection into separate read & write connections. (#17)
- Major: Redirect pings to bot's own channel if a channel had EventsOnlyOffline flag set to true. (#15)
- Minor: Make chat commands use tcb prefix more.
- Minor: Don't parse non-commands on PRIVMSG messages.
- Minor: Make better use of `context.Context` in the databse connector.
- Minor: Added `{login}` placeholder for streamer's name in event messages. (#14)
- Minor: Ignore notifications from channels not present in channel cache.
- Minor: Make cooldowns less strict for `!notifyme`, `!removeme` and `!subscribed` commands.
- Minor/Bugfix: Write OK status to Twitch's eventsub request before handling eventsub notification. (Probably wrong order, investigate later)
- Minor: Avoid announcing when stream is already live/offline.
- Minor: Check for command names is now case-insensitive.
- Bugfix: Make panics more clean upon failing to obtain an app access token by not deserializing nil response struct.
- Bugfix: Use new deduplicate unicode character `\u034f`.
- Bugfix: Use mutex to avoid concurrent logins map writes.
- Dev: Migrated golangci-lint to v2. (#18)
- Dev: Move `bot.Channel` and `bot.Command` defintions to their respective files, away from `bot/types.go` file.
- Dev: Remove deprecated calls to ioutil.
- Dev: Upgraded nicklaw5/helix to a major version v3.
- Dev: Upgraded gempir/go-twitch-irc to a major version v4.

## 2.0-rc2

- Minor: Reimplemented channel MOTDs. (#10)
- Minor: Added a way to subscribe to all events at once with `!notifyme all`.
- Minor: List all available events in `!notifyme` command.
- Minor: Reimplemented respecting of disabled command.
- Minor: Make periodic requests to supinic's API to signal bot being alive. (#6)
- Bugfix: Corrected response in `!bot` command.
- Bugfix: Make checks for subscription values case-insensitive again.
- Dev: Use millisecond precision in logs.
- Dev: Added licese. (#7)
- Dev: Cleaned up code for API routes.


## 2.0-rc1

- Initial release.
