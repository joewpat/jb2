#!/bin/bash
#copy any new/changed files to web server
rsync -av -e ssh --exclude='.git' -i ~/.ssh/jpcom ~/code/jb2/ jpatrick@34.75.110.238:/opt/joebot/

ssh -i ~/.ssh/jpcom jpatrick@34.75.110.238 'docker build -t joebot /opt/joebot/.'
ssh -i ~/.ssh/jpcom jpatrick@34.75.110.238 'docker kill $(docker ps -q -f "name=joebot")' 2>/dev/null
ssh -i ~/.ssh/jpcom jpatrick@34.75.110.238 'docker container rm $(docker container ls -qaf "name=joebot")' 2>/dev/null
ssh -i ~/.ssh/jpcom jpatrick@34.75.110.238 'docker run -d --restart=unless-stopped --name=joebot joebot'
ssh -i ~/.ssh/jpcom jpatrick@34.75.110.238 'docker system prune -f'