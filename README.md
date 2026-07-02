# Chatsighting

Find the most hype moments in any Twitch VOD by analysing chat velocity spikes.

**Live at [chatsighting.com](https://chatsighting.com)**

<img width="775" alt="Chatsighting home" src="https://github.com/user-attachments/assets/5034ee94-cd88-4123-a76f-0c36207a403a" />

## What it does

- Enter a Twitch username, pick a VOD, get the top spike moments
- Each spike shows the timestamp, message count vs stream average, top chat words, and a link to jump straight to that moment in the VOD
- Chat activity graph across the full stream
- Chat is fetched live from Twitch and never stored

**Search a streamer and pick a VOD:**

<img width="816" alt="VOD search results" src="https://github.com/user-attachments/assets/53726f56-b9ac-4271-954b-2d0ce28a21c0" />

**Chat activity graph with detected spikes:**

<img width="957" alt="Chat activity graph" src="https://github.com/user-attachments/assets/35d85df9-a0ba-4a43-be03-e619a80631f4" />

**Spike cards with top words and highlighted chat:**

<img width="1051" alt="Spike cards" src="https://github.com/user-attachments/assets/faed6bd9-7a2d-42b5-9079-7c30c46f5f0b" />

**Jump straight to the moment in the stream:**

<img width="1271" alt="Jump to moment on Twitch" src="https://github.com/user-attachments/assets/4e9975be-db5e-407b-ace4-7063467647c4" />

## How it works

- Go backend fetches VOD chat via Twitch's GQL API, using offset-based pagination to work around the API's integrity checks, with handling for empty and partial pages mid-VOD
- Spike detection compares message velocity in 30s windows against the stream average
- Per-IP rate limiting to keep the public service stable
- Runs on AWS (Lightsail, Sydney) with Let's Encrypt SSL

## Stack

- Go backend
- Vanilla HTML/CSS/JS frontend, Chart.js for the activity graph
- Twitch GQL API for chat data
- No database, no auth, stateless

*Source is in a private repo.
