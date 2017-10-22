Introduction

MapNotes is an android app allowing people to showcase events, reports and more
to their friends, family and the wider world. You can keep up to date with
what's been going on in your area, if it's last week or last year, you'll soon
know what's been happening.

Prerequisites

~~ Go version 1.9 from https://golang.org/
~~ Heroku CLI from https://www.heroku.com/
~~ Godep https://github.com/tools/godep

Installing

To set up the database once you have been granted access to the app on heroku:
1. Register the app with your Heroku CLI
`heroku git:remote -a mapnotes-backend`

2. Install Postgresql to your linux machine
Try doing
`sudo apt-get install postgresql-9.6`
If that doesn't work follow this guide for ubuntu/mint
https://www.codeproject.com/Articles/898303/Installing-and-Configuring-PostgreSQL-on-Linux-Min

3. Set up yourself as a super user for postgres on your machine
Follow this guide https://www.codeproject.com/Articles/898303/Installing-and-Configuring-PostgreSQL-on-Linux-Min

4. In the top level directory of the repository, copy DATABASE_URL from Heroku to .env
Run `heroku config:get DATABASE_URL -s  >> .env`

5. Import database from Heroku
Run `heroku pg:pull DATABASE_URL mapnotes_local_db --app mapnotes-backend`
This will pull from DATABASE_URL from the app callled `mapnotes-backend` to your machine and create a new table called `mapnotes_local_db`

6. Update .env file with your local database url
For example (replace username and password with your local account)
`DATABASE_URL=postgres://username:password@localhost:5432/mapnotes_local_db`
You may need a querystring of `?sslmode=disable`

7. Update .env file with a non-standard port for the webserver, e.g.
`PORT=9032`

8. Run Heroku locally
Run `go install` to build the go files
Run `heroku local` to start Heroku locally

Testing

Deployment

Authors & Acknowledgements
