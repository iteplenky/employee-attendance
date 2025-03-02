## **Employee Attendance Bot**

### Description

Employee Attendance Bot is a Telegram bot designed to track employee attendance. It integrates with the time tracking system (ZKTeco BioTime) and notifies users in real time when their arrival is registered.

The bot uses a PostgreSQL database to store user information, Redis for caching subscribers, and asynchronous event processing. A database trigger mechanism ensures an instant response to changes.

**Features**

- Registration of users via IIN (Individual Identification Number)
- Subscription management for attendance notifications
- Instant Telegram notifications upon arrival registration
- High-speed processing using caching and asynchronous execution
- Logging of user actions for audit purposes

**Architecture**

The system consists of the following components:

- **PostgreSQL** – stores user data, subscriptions, and attendance logs
- **Redis** – caches subscriptions to speed up event processing
- **Telegram Bot** – processes user commands and sends notifications
- **Event Processing Service** – listens for database events and sends them to Redis

**Getting Started**

**Requirements**

- Docker and Docker Compose
- Telegram account

**Installation and Setup**

Clone the repository:

```shell
git clone https://github.com/iteplenky/employee-attendance.git
```

Create a .env file and specify the bot API token, database, and Redis connection details:

```env
TOKEN=your_telegram_bot_token
DATABASE_URL=postgres://user:password@postgres:5432/employee-attendance?sslmode=disable
DATABASE_ATTENDANCE_URL=postgres://user:password@employee-attendance-postgres-1:5432/biotime?sslmode=disable
REDIS_URL=redis://redis:6379
```

Start the containers using Docker Compose and Makefile:

```shell
make up
```

The bot will automatically start and be ready for use.

### How to Use
- Start the bot in Telegram.
- Register by sending your IIN.
- Receive instant messages upon each arrival at work.

### Details
- Connection to BioTime is established via an INSERT trigger in the attendance_log table.
- PostgreSQL automatically sends the event to the attendance_events channel using LISTEN/NOTIFY.
- Redis receives this event and forwards it to the bot processor.
- The bot sends notifications to subscribed users.
- All attendance records and user actions are logged in the database.