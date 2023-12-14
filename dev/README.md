# Getting started Developing
Thank you for working on Statping! Lets get your dev environment up and running so you can start contributing right away.

## Prerequisites
Install Docker Desktop on your computer.

## Start Dev Server
We use Docker Compose for our development environment, this will allow you to not have to worry about installing anything on your computer, and everything will "just work". 

The quickest way to start is to use the "lite" Docker Compose setup.

1. Start Docker Desktop
2. Start the dev server by running:

```shell
cd dev
docker compose -f docker-compose.lite.yml up
```

3. Go to: http://localhost:8585
