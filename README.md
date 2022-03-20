# BudgetBot
## About
Telegram bot for archiving of purchase receipts by scanning of QR code from the receipt. Create your own archive of receipts, which you will never lose again, and will always have access to.

## Installation
### Using Docker
The simplest way to start the service is docker:

Clone the repo and run docker-compose:
```
docker-compose up -d
```

### Usage
First, you must create a Telegram bot and obtain its access token. To do that, follow the instructions <a href="https://core.telegram.org/bots">here</a>.
Next, create a `.env` file in the root of this project, and pass it this access token.