## Telegram Bot for Notifier-Service

This Telegram bot is an extension of the [Notifier-Service](https://github.com/NIC-Golang/notifier-service). It allows users to receive notifications and interact with the system via Telegram.

### Features
- User authentication via Telegram
- Receiving real-time notifications
- Managing notification preferences
- Interactive commands and buttons for ease of use

### Installation

#### Prerequisites
- Go 1.23
- Docker & Docker Compose (optional)
- MySQL (if running locally)

#### Clone the Repository
```sh
git clone https://github.com/maximusprimeavenger/telegram-bot.git
cd telegram-bot
```

#### Clone the Notifier Service
```sh
git clone https://github.com/NIC-Golang/notifier-service.git
cd notifier-service
cp .env.example .env
```

#### Configure Environment Variables
Create a `.env` file with the following variables:
```
SECRET=your_telegram_bot_token
DATABASE_URL=postgres://user:password@localhost:5432/notifier
API_HOST=http://notifier-service:8080
PORT_BOT=8000
PORT = 3306
PASSWORD = your_password
USER = your_user
ROOT_PASSWORD = your_root_password
```

### Running the Bot

#### With Docker
```sh
docker compose up
```

#### Without Docker
```sh
go mod tidy
go run cmd/main.go
```

### Usage
The bot supports the following commands:
- `/start` – Start interaction with the bot
- `/my_orders` – View user-specific notifications
- `/register` – Register in the system

### API Integration
This bot communicates with the Notifier-Service through REST API. It sends and retrieves user data, ensuring real-time synchronization with the main service.

### Contributing
Feel free to contribute by submitting issues or pull requests to improve the bot.

