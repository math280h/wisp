<p align="center">
  <img src="./assets/logo.png" width="240px" height="240px" />
</p>
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
- Warns and Strikes
  - Users can be warned and given strikes
  - Warn and strike gives the user "points" which when they reach a certain threshold, they are automatically banned 
    - Both the points given and the threshold are configurable
  - Notifies the user when they are warned or given a strike
  - Ability to historically view all warns and strikes for a user
- Ban / Kick
  - Users can be banned or kicked from the server
  - Points
    - When a user is banned, they are given 100% of the threshold points
    - When a user is kicked, they are given 50% of the threshold points
  - Ban and kick reasons are logged
- Suggestions
  - Users can make suggestions to the server
  - All users that can view the suggestions channel can vote on suggestions
    - Users can only vote once per suggestion
  - Suggestions can be marked as Accepted or Denied
    - Users with Administator permissions can mark suggestions as completed (TBD if the permission will be configurable)

## Getting Started

Create a new bot on the [Discord Developer Portal](https://discord.com/developers/applications) and invite it to your server.

### Running the docker container

**Note:** See .env.example for environment variables

<!-- TODO:: Get actual link for docker image -->
```bash
docker run -d --env-file .env ghcr.io/math280h/wisp/wisp:latest
```
