# Chatsighting

Find the most hype moments in any Twitch VOD by analysing chat velocity spikes.

## What it does

- Enter a Twitch username → pick a VOD → get the top spike moments
- Shows timestamp, message count vs stream average, top words, and a link to jump straight to that moment
- Chat activity graph across the full stream

## Run locally

```bash
cp .env.example .env   # add your Twitch credentials
go run main.go
```

Open `http://localhost:3000`

## Env vars

```
TWITCH_CLIENT_ID=
TWITCH_CLIENT_SECRET=
PORT=3000
```

Get credentials at [dev.twitch.tv](https://dev.twitch.tv/console)

## Stack

- Go backend
- Vanilla HTML/CSS/JS frontend
- Twitch GQL API for chat data
- No database, no auth, stateless
