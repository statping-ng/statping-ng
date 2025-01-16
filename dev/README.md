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
cd statping-ng # go into the root folder, it's important you run it in there.
docker compose -f dev/docker-compose.lite.yml up
```

3. [Frontend] Go to: http://localhost:8888
3. [Backend] Go to: http://localhost:8585
