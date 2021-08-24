<center>
![Statping-ng](https://raw.githubusercontent.com/statping-ng/statping-ng/dev/frontend/public/img/banner.png =60)

**Statping-ng** - *Web and App Status Monitoring for Any Type of Project*

[Website](https://statping-ng.github.io) | [Wiki](https://github.com/statping-ng/statping-ng/wiki)

[Linux](https://github.com/statping-ng/statping-ng/wiki/Linux) | [Windows](https://github.com/statping-ng/statping-ng/wiki/Windows) | [Mac](https://github.com/statping-ng/statping-ng/wiki/Mac) | [Docker](https://github.com/statping-ng/statping-ng/wiki/Docker)
</center>

# Statping-ng - Status Page & Monitoring Server

An easy to use Status Page for your websites and applications. Statping will automatically fetch the application and render a beautiful status page with tons of features for you to build an even better status page. This Status Page generator allows you to use MySQL, Postgres, or SQLite on multiple operating systems.

Statping-ng aims to be an updated drop-in replacement of statping after development stopped on the original fork.


[![License](https://img.shields.io/github/license/statping-ng/statping-ng?color=success&style=for-the-badge&logo)](https://github.com/statping-ng/statping-ng/blob/stable/LICENSE)

![GitHub last commit](https://img.shields.io/github/last-commit/statping-ng/statping-ng?style=for-the-badge&logo=github) ![Unstable Build](https://img.shields.io/github/workflow/status/statping-ng/statping-ng/1.%20Development%20Build%20and%20Testing?label=Dev&style=for-the-badge&logo=github) ![Unstable Build](https://img.shields.io/github/workflow/status/statping-ng/statping-ng/2.%20Unstable%20Build,%20Test%20and%20Deploy?label=Unstable&style=for-the-badge&logo=github) ![Stable Build](https://img.shields.io/github/workflow/status/statping-ng/statping-ng/3.%20Stable%20Build,%20Test%20and%20Deploy?label=Stable&style=for-the-badge&logo=github)

[![Docker Pulls](https://img.shields.io/docker/pulls/adamboutcher/statping-ng?style=for-the-badge&logo=docker)](https://hub.docker.com/r/adamboutcher/statping-ng) [![Docker Image Size](https://img.shields.io/docker/image-size/adamboutcher/statping-ng/latest?style=for-the-badge&logo=docker)](https://hub.docker.com/r/adamboutcher/statping-ng)

 ![Go Version](https://img.shields.io/github/go-mod/go-version/statping-ng/statping-ng?style=for-the-badge) [![Go Report Card](https://goreportcard.com/badge/github.com/statping-ng/statping-ng?style=for-the-badge)](https://goreportcard.com/badge/github.com/statping-ng/statping-ng)

---
# About Statping-ng

![Statping-ng example](https://statping-ng.github.io/assets/external/statupsiterun.gif =320x235)

## A Future-Proof Status Page
Statping-ng strives to remain future-proof and remain intact if a failure is created. Your Statping-ng service should not be running on the same instance you're trying to monitor. If your server crashes your Status Page should still remaining online to notify your users of downtime.

[![Play with Docker](https://statping-ng.github.io/assets/external/docker-pwd.png)](https://labs.play-with-docker.com/?stack=https://raw.githubusercontent.com/statping-ng/statping-ng/stable/dev/pwd-stack.yml) - Login is `admin`, password `admin`.

## No Requirements
Statping-ng is built in Go Language so all you need is the pre-compiled binary based on your operating system. You won't need to install anything extra once you have the Statping binary installed. You can even run Statping-ng on a Raspberry Pi.

<center>

[![Linux](https://statping-ng.github.io/assets/external/linux.png)](https://github.com/statping-ng/statping-ng/wiki/Linux) [![Windows](https://statping-ng.github.io/assets/external/windows.png)](https://github.com/statping-ng/statping-ng/wiki/Windows) [![Mac](https://statping-ng.github.io/assets/external/apple.png)](https://github.com/statping-ng/statping-ng/wiki/Mac) [![Docker](https://statping-ng.github.io/assets/external/Docker.png)](https://hub.docker.com/r/adamboutcher/statping-ng) [![Android](https://statping-ng.github.io/assets/external/android.png)](https://play.google.com/store/apps/details?id=com.statping) [![iOS](https://statping-ng.github.io/assets/external/appstore.png)](https://itunes.apple.com/us/app/apple-store/id1445513219)

</center>


<img align="right" width="320" height="235" src="https://gitimgs.s3-us-west-2.amazonaws.com/slack-notifer.png">
<h2>Lightweight and Fast</h2>
Statping-ng is a very lightweight application and is available for Linux, Mac, and Windows. The Docker image is only ~16Mb so you know that this application won't be filling up your hard drive space.
The Status binary for all other OS's is ~17Mb at most.
<br><br><br><br><br><br>

<img align="left" width="320" height="235" src="https://img.cjx.io/statping_iphone_bk.png">
<h2>Mobile App is Gorgeous</h2>
The Statping app is available on the App Store and Google Play for free. The app will allow you to view services, receive notifications when a service is offline, update groups, users, services, messages, and more! Start your own Statping-ng server and then connect it to the app by scanning the QR code in settings.

<p align="center">
<a href="https://play.google.com/store/apps/details?id=com.statping"><img src="https://img.cjx.io/google-play.svg"></a>
<a href="https://itunes.apple.com/us/app/apple-store/id1445513219"><img src="https://img.cjx.io/app-store-badge.svg"></a>
</p>

<br><br>

## Run on Any Server
Want to run it on your own Docker server? Awesome! Statping-ng has multiple docker-compose.yml files to work with. Statping-ng can automatically create a SSL Certification for your status page.
<br><br><br><br>

<img align="left" width="320" height="205" src="https://img.cjx.io/statping_theme.gif">
<h2>Custom SASS Styling</h2>
Statping-ng will allow you to completely customize your Status Page using SASS styling with easy to use variables. The Docker image actually contains a prebuilt SASS binary so you won't even need to setup anything!
<br><br><br><br>

## Slack, Email, Twilio and more
Statping-ng includes email notification via SMTP and Slack integration using [Incoming Webhook](https://api.slack.com/incoming-webhooks). Insert the webhook URL into the Settings page in Statping-ng and enable the Slack integration. Anytime a service fails, the channel that you specified on Slack will receive a message.
<br><br><br><br>

<h2>User Created Notifiers</h2>
View the [Plugin Wiki](https://github.com/statping-ng/statping-ng/wiki/Statping-Plugins) to see detailed information about Golang Plugins. Statping-ng isn't just another Status Page for your applications, it's a framework that allows you to create your own plugins to interact with every element of your status page. [Notifier's](https://github.com/statping-ng/statping-ng/wiki/Notifiers) can also be create with only 1 golang file.
<br><br><br><br>

<img align="center" width="100%" height="250" src="https://img.cjx.io/statupsc2.png">

<br><br><br><br>

<img align="right" width="320" height="235" src="https://img.cjx.io/statping_settings.gif">
<h2>Easy to use Dashboard</h2>
Having a straight forward dashboard makes Statping-ng that much better. Monitor your websites and applications with a basic HTTP GET request, or add a POST request with your own JSON to post to the endpoint.
<br><br><br><br>

## Run on Docker
Use the [Statping Docker Image](https://hub.docker.com/r/adamboutcher/statping-ng) to create a status page in seconds. Checkout the [Docker Wiki](https://github.com/statping-ng/statping-ng/wiki/Docker) to view more details on how to get started using Docker.
```bash
docker run -it -p 8080:8080 adamboutcher/statping-ng
```
There are multiple ways to startup a Statping-ng server. You want to make sure Statping-ng is on it's own instance that is not on the same server as the applications you wish to monitor. It doesn't look good when your Status Page goes down.
<br><br><br><br>

## Docker Compose
In this folder there is a standard docker-compose file that include nginx, postgres, and Statping-ng.
```bash
docker-compose up -d
```
<br><br><br><br>

## Docker Compose with Automatic SSL
You can automatically start a Statping-ng server with automatic SSL encryption using this docker-compose file. First point your domain's DNS to the Statping-ng server, and then run this docker-compose command with DOMAIN and EMAIL. Email is for letsencrypt services.
```bash
LETSENCRYPT_HOST=mydomain.com \
    LETSENCRYPT_EMAIL=info@mydomain.com \
    docker-compose -f docker-compose-ssl.yml up -d
```
Once your instance has started, it will take a moment to get your SSL certificate. Make sure you have a A or CNAME record on your domain that points to the IP/DNS of your server running Statping-ng.
<br><br><br><br>

## Prometheus Exporter
Statping-ng includes a [Prometheus Exporter](https://github.com/statping-ng/statping-ng/wiki/Prometheus-Exporter) so you can have even more monitoring power with your services. The Prometheus exporter can be seen on `/metrics`, simply create another exporter in your prometheus config. Use your Statping-ng API Secret for the Authorization Bearer header, the `/metrics` URL is dedicated for Prometheus and requires the correct API Secret has `Authorization` header.
```yaml
scrape_configs:
  - job_name: 'statping'
    bearer_token: 'MY API SECRET HERE'
    static_configs:
      - targets: ['statping:8080']
```
<br><br><br><br>

## Contributing
Statping-ng accepts Push Requests to the `dev` branch! Feel free to add your own features and notifiers. You probably want to checkout the [Notifier Wiki](https://github.com/statping-ng/statping-ng/wiki/Notifiers) to get a better understanding on how to create your own notification methods for failing/successful services. Testing on Statping-ng will test each function on MySQL, Postgres, and SQLite. I recommend running MySQL and Postgres Docker containers for testing. You can find multiple docker-compose files in the dev directory.
