<p align="center">
  <img src="./assets/logo.png" width="240px" height="240px" />
</p>
<h1 align="center">
    Wisp
</h1>
<p align="center">
  Wisp is a Discord bot that helps moderate servers. Wisp is designed to be easy to use and easy to configure.
</p>

<p align="center">
    Wisp is also fully open source and can be self-hosted by anyone. (No more sketchy bots!)
</p>

<p align="center">
  <img src="https://github.com/math280h/wisp/actions/workflows/pr.yaml/badge.svg" />
</p>

## Features

- Reports
  - Users can open reports by dm'ing the bot
  - All moderators can see/respond to reports and their identities are hidden
  - Archived reports are stored in a channel for future reference
- Detain
  - Detained users will be given a configurable role
  - When a user is detained, they will get a notification and a new report will open if one is not already open.
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
- Report Messages
  - Users can report messages by reacting to them with the :octagonal_sign: (🛑) emoji
  - Moderators will get an alert in the alerting channel with a link to the message.
- Message Logs
  - When a message is edited an embed containing the original message and the new message is sent to the message logs channel
  - When a message is deleted an embed containing the deleted message is sent to the message logs channel

## Getting Started

Create a new bot on the [Discord Developer Portal](https://discord.com/developers/applications) and invite it to your server.

### Running the docker container

**NOTE:** In the .env file, it's important none of the values is sourrounded by quotes.
  If you are using qoutes docker will escape the values as \"value\" and the bot will not work.
  (*We will make a fix for this in the future*)

```bash
docker compose up
```

## Configuration

Wisp can be configured either through parameters passed to the bot or through environment variables.

You can find all env variables in the [.env.example](.env.example) file.

You can find all command line arguments [flags.go](internal/shared/flags.go)

## Development

### Prisma

```bash
go run github.com/steebchen/prisma-client-go generate
go run github.com/steebchen/prisma-client-go db push
```

### Running the bot

Use [air](https://github.com/air-verse/air) to run the bot with hot reloading.

### Windows

```bash
air --build.cmd "go build -o /tmp/main.exe cmd\main.go"
```

### Linux

```bash
air --build.cmd "go build -o /tmp/main cmd/main.go"
```
