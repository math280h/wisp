<h1 align="center">
    Wisp
</h1>

<p align="center">
  Wisp is a simple, lightweight moderation bot for Discord. It is designed to be easy to use and easy to set up.
</p>


## Features

- Reports
  - Users can open reports by dm'ing the bot
  - All moderators can see/respond to reports and their identities are hidden
  - Archived reports are stored in a channel for future reference

### Coming Soon

- Moderation
  - Kick
  - Ban
  - Mute
  - Unmute
  - Warn
  - Clear

## Getting Started

Create a new bot on the [Discord Developer Portal](https://discord.com/developers/applications) and invite it to your server.

### Running the docker container

**Note:** See .env.example for environment variables

<!-- TODO:: Get actual link for docker image -->
```bash
docker run -d --env-file .env ghcr.io/username/repo:tag
```
