# Employee Attendance Bot

A Telegram bot that integrates with ZKTeco Biotime for real-time attendance tracking.

## Quick Start

Clone the repository:

```sh
git clone https://github.com/iteplenky/employee-attendance
cd employee-attendance
```

Start ZKT database:

```sh
make zkt-up
```

Start the bot:

```sh
make up
```

## How It Works

1. When the ZKT container starts, it runs a script that enables real-time attendance event tracking in PostgreSQL.
2. A user registers in the bot by entering their unique ID.
3. When the user checks in via ZKTeco (e.g., using Face ID), an event is triggered in the database.
4. The bot listens for new check-ins and sends instant notifications to users who have them enabled.
5. Users can view their attendance history via Telegram’s built-in keyboards.

## Commands

- `/start` – Begin using the bot.

Everything else is managed via Telegram's built-in keyboards.
