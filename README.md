# goBot for Telegram

Simplest GoBot. Can store statistics into SQLite and generate TOP20 flooder list. 
Don't be serious, I have fun and practice!

### Prepairing 
 - Have a talk with @FatherBot in Telegram
 - Get APIKey. Looks like <NUM>:<HASH>
 - Set Bot's "Privacy" to Disable and "Join to Channel" Enable
 - Set Bot's user pic
 - Set Commands as "/status" and "/help"
 - Invite your new Bot to channel

## Start from command line
```
./goBot -h

Avaliable options:

Usage of /home/goBot:
  -apikey string
        Bot ApiKey. See @BotFather messages for details
  -c string
        default config file. See settings.ini.example (default "settings.ini")
  -dbdriver string
        Driver DB. (default "sqlite3")
  -dbname string
        Database of users (default "empty.db")
  -dir string
        Working directory (default "./")
```

### Simple settings.ini

```
[main]
ApiKey = 
ChatRoomID = 
SQLITE_DB = ./empty.db
```

## Start with Docker:
```
sudo docker run -v /etc/ssl/certs/:/etc/ssl/certs \
-v /tmp/empty.db:/home/empty.db \
-t doctornkz/gobot \
-apikey <YOUR_KEY> \
-dbname /home/empty.db
```

# Using 

If you did everything right, you can use menu options in chat - /status and /help
Bot parses all activity in channel and shows TOP20 list of flooders:
```
-= TOP LIST =- 
1. Doctornkz = 22
2. Bobby = 6
3. Billy = 3
```
