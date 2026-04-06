# CHAT-INIT

<p align="center">
  <strong>A pirate-themed terminal chat app built with Go, WebSockets, and colorful CLI screens.</strong>
</p>

<p align="center">
  Beginner-friendly to run, fun to demo, and simple enough to learn from.
</p>

<p align="center">
  <img src="https://img.shields.io/badge/Go-1.25%2B-00ADD8?style=for-the-badge&logo=go" alt="Go">
  <img src="https://img.shields.io/badge/WebSockets-Gorilla-111111?style=for-the-badge" alt="WebSockets">
  <img src="https://img.shields.io/badge/UI-Lipgloss-FF5E5B?style=for-the-badge" alt="Lipgloss">
  <img src="https://img.shields.io/badge/Style-Terminal%20First-F5B700?style=for-the-badge" alt="Terminal First">
</p>

---

## What Is CHAT-INIT?

CHAT-INIT is a real-time terminal chat project with:

- a Go server in `CHAT_INIT/`
- a Go client in `CLIENT/`
- WebSocket-powered direct messaging
- OTP-based signup and password reset support
- a playful pirate-style terminal interface

This repo is especially nice for:

- beginners who want to see a full client/server Go project
- developers learning HTTP + WebSockets together
- people who want a CLI project that feels more alive than a plain CRUD demo

---

## Screenshots

| | | |  hmm
|:---:|:---:|:---:|
| ![Screenshot 1](1.png) | ![Screenshot 2](2.png) | ![Screenshot 3](3.png) |
| ![Screenshot 4](4.png) | ![Screenshot 5](5.png) | ![Screenshot 6](6.png) |
| ![Screenshot 7](7.png) | ![Screenshot 8](8.png) | |

---

## Why It Feels Cool

- Real-time chat over WebSockets instead of fake terminal prompts
- Bright, styled CLI screens using Lipgloss
- Friend management built into the client flow
- OTP-based registration and password reset flow
- Small enough to understand, big enough to feel like a real project

---

## Project Layout

```text
CHAT_INIT/
├── CHAT_INIT/         # Main chat server
├── CLIENT/            # Main terminal chat client
├── test/              # Practice and experimental code
├── 1.png ... 8.png    # Project screenshots
└── README.md
```

Main folders you’ll actually run:

- `CHAT_INIT/` for the server
- `CLIENT/` for the client

The `test/` folder contains side experiments and practice code. It is useful for learning, but it is not the main app entrypoint.

---

## Features

- Real-time direct messaging
- Login with stored token support
- Friend requests, accept, reject, and delete flow
- Password hashing with bcrypt
- OTP email support for signup and password reset
- Terminal UI with multiple menus and styled output

---

## Beginner Quick Start

### 1. Prerequisites

Make sure you have:

- Go installed
- a terminal on macOS, Linux, or Windows
- a Gmail address with an App Password if you want OTP signup/reset to work

Check Go:

```bash
go version
```

### 2. Clone the Repo

```bash
git clone https://github.com/anzydev/CHAT_INIT.git
cd CHAT_INIT
```

### 3. Start the Server

Open terminal 1:

```bash
cd CHAT_INIT
go run .
```

Important:

- on first startup, the server will ask for `GMAIL_USER` and `GMAIL_PASS`
- it stores them in `CHAT_INIT/.env`
- this is used for OTP email features

If you do not want to use OTP features yet, you can still learn the codebase, but the current server startup flow expects those values.

### 4. Start the Client

Open terminal 2:

```bash
cd CLIENT
go run .
```

### 5. Use the App

Inside the client:

- choose `1` to log in
- choose `2` to register
- choose `3` to reset password
- type `111` to open friend management
- choose a friend number to start chatting
- type `0` to exit

---

## OTP Email Setup

CHAT-INIT uses Gmail SMTP for OTP delivery.

You need:

- `GMAIL_USER=your-email@gmail.com`
- `GMAIL_PASS=your-16-character-google-app-password`

Example:

```env
GMAIL_USER=your-email@gmail.com
GMAIL_PASS=your-app-password
```

How to get the password:

1. Enable 2-step verification on your Google account.
2. Create a Google App Password.
3. Use that app password, not your normal Gmail password.

---

## How the App Flow Works

### Server

- handles login, signup, OTP, friend actions, and WebSocket chat
- stores users in `CHAT_INIT/database.json`
- validates tokens before allowing chat actions

### Client

- talks to the server over HTTP for login and friend actions
- opens a WebSocket connection for real-time chat
- stores the session token in `CLIENT/.env`

---

## Tech Stack

- Go
- Gorilla WebSocket
- Charmbracelet Lipgloss
- bcrypt for password hashing
- JSON file storage for user data
- Gmail SMTP for OTP

---

## Commands You’ll Actually Need

Run the server:

```bash
cd CHAT_INIT
go run .
```

Run the client:

```bash
cd CLIENT
go run .
```

Format code:

```bash
gofmt -w CHAT_INIT/*.go CLIENT/*.go test/*.go test/client/*.go test/practis/*.go
```

Run tests in the practice modules:

```bash
cd test
go test ./...
```

---

## Troubleshooting

### Server asks for Gmail credentials

That is expected on first run. Enter your Gmail address and App Password so OTP can work.

### Client cannot connect

Make sure:

- the server is already running on port `4040`
- you started the client from the `CLIENT/` folder

### OTP email does not arrive

Check:

- the Gmail address is valid
- the App Password is correct
- Gmail account security settings allow App Password usage

---

## Why This Repo Is Good for Learning

- You can trace a full feature from client input to HTTP route to WebSocket message.
- The codebase uses familiar Go packages and simple JSON storage instead of heavy infrastructure.
- It shows both networking and terminal UI in one project.

If you are new to Go, this is a good repo to practice:

- structs
- JSON encoding/decoding
- HTTP handlers
- WebSocket connections
- goroutines
- terminal UX

---

## Future Improvement Ideas

- persistent database instead of JSON files
- group chat rooms
- chat history
- better error messages in the client
- Docker setup
- tests for the main server/client flow

---

## Connect

- GitHub: [mikeyoni](https://github.com/mikeyoni)
- Discord: `@ui_mikey`
- Instagram: `@ui__mikey`

---

## License

This project is licensed under the MIT License. See [LICENSE](LICENSE).
