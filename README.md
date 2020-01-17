# piholebot
Pi Hole Ad-Blocker Tweet Bot

### dev
Build for Pi
============
```
GOOS=linux GOARCH=arm GOARM=5 go build -v
```
Copy over
=========
```
scp files/etc/piholebot/piholebot.production.ini pi@danfe:~/
scp piholebot pi@danfe:~/
```

### pihole: danfe - one-time setup
Copy bot config
===============
```
sudo mkdir /etc/piholebot
sudo cp piholebot.production.ini /etc/piholebot/
```
Symlink bot binary
==================
```
sudo ln -sf /home/pi/piholebot /usr/local/bin/
```

Crontab Entry
=============
```
55 23 * * * PIENV=production /usr/local/bin/piholebot >> /var/log/piholebot/piholebot.log 2>&1
```
