# Receipt bot

## About

Telegram bot for archiving of purchase receipts by scanning the QR code from the receipt. Create your own archive of
receipts, which you will never lose again, and will always have access to.
Currently, only <b>slovak</b> receipts are supported!

## Usage
First, you will need to create a Telegram bot, and obtain its access token.
To do that, follow the
instructions <a href="https://core.telegram.org/bots">here</a>.

## Environment variables

To use the bot correctly, you need to set two environment variables into `.env` file in the root of the project. There
is an example `.env.example` file to get you started. Just rename this file to `.env`, and put in the correct values.

### Required values

| Name           | Description                                |
|----------------|--------------------------------------------|
| TELEGRAM_TOKEN | The access token of your telegram bot      |
| POSTGRES_URL   | The connection string to your SQL database |

## Installation

### Using Docker

The simplest way to start the service is docker:

Clone the repo and run docker-compose:

```shell
docker-compose up -d
```

### Locally
You can also use the bot locally by running:

```shell
go build
go run main.go
```
However, it is presumed that you already have a running database, and you put its connection string into the `.env` file.

## Features
### Receipt QR code scanner
You can simply take a picture of a receipt, and the bot will take care of the rest. Your receipts will be archived in the database, as well as locally on the machine running the bot.

![ Alt text](assets/qr_scan.gif) / ! [](assets/qr_scan.gif)

If you prefer to input the QR code manually, you can do this by using the command `/qr yourQrCode`.

![ Alt text](assets/qr_command.gif) / ! [](assets/qr_command.gif)

## TODOs, future features (hopefully)
- add attachment to a receipt (document, photo,...)
- add a description to an existing receipt
- send receipt via e-mail or external service
- delete individual receipts