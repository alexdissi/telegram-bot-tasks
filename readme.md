# Telegram Task Management Bot

## A lightweight task management bot for Telegram that helps users manage their tasks and sends notifications for pending tasks.

### Features

Add, list, edit, complete, and delete tasks.
Notifications for tasks pending for more than 24 hours.
Enable or disable reminders for task notifications.
All of the above features can store data in a local SQLite database.

### Setup Instructions

Prerequisites
Docker and Docker Compose installed.
A Telegram Bot Token from BotFather.

### Steps to Set Up

Clone the repository.

```bash
git clone https://github.com/your-repo/telegram-task-bot.git
cd telegram-task-bot
```

GO to https://core.telegram.org/bots/tutorial to create a new bot and get the token.
Follow the instructions to create a new bot and get the token, in section "Obtain Your Bot Token"

Add the Telegram Bot Token to the .env file.

```bash
export TELEGRAM_BOT_TOKEN=your-token
```

Build and run the Docker container.

```bash
docker-compose up -d
```

#### This will:

Set up the SQLite database for local storage.
Automatically migrate the database using GORM.
Runs CRUD functions test for task

Start the bot.
Find your bot on Telegram and send /start to interact with it.

#### Features in Detail

- Task Management:
- Add tasks using /add.
- View all tasks with /list.
- Edit task descriptions using /edit.
- Mark tasks as done with /done.
- Remove tasks using /delete.
- Notifications:
- Notifications are sent hourly to users with tasks pending for over 24 hours.
- Use /enable_reminders or /disable_reminders to toggle notifications.

### Example Workflow

Add a task using /add.

```bash
User: /add Buy groceries
Bot: Task [1] added: Buy groceries
```

List all tasks using /list.

```bash
User: /list
Bot: Your tasks:
1. Buy groceries
```

Marking a Task as Completed:

```bash
User: /done 1
Bot: Task [1] marked as done: Buy groceries
```

Removing a Task:

```bash
User: /delete 1
Bot: Task [1] removed: Buy groceries
```

Enable or Disable Notifications:

```bash
User: /enable_reminders
Bot: Notifications enabled
```

```bash
User: /disable_reminders
Bot: Notifications disabled
```

For testing

```bash
docker-compose exec app go test ./...
```

### Note

For this project i used GORM to interact with the database, and the database is stored in a local SQLite file. The database is automatically migrated when the Docker container is started. The bot is built using the go-telegram-bot-api library.
