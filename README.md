# ChirpBird

Instant messaging platform with Websocket using Golang

## how to run this project

open terminal and run command below

clone this project

```sh
git clone https://github.com/simopunkc/chirpbird-v2.git
```

You have to edit the .env.example file in the backend folder

nano /chirpbird-v2/backend/.env.example

```sh
OAUTH_CLIENT_ID='put your Google Oauth Client ID here'
OAUTH_SECRET='put your Google Oauth Secret here'
```

if you do not have a Google Oauth Client ID then you must first create it on the Credentials menu on the Google Cloud Platform. You must also fill in the Authorized Redirect URI box with the path /oauth/google/verify. for example:

```sh
http://localhost/oauth/google/verify
```

make sure docker-compose is installed. after that run bash script with below command

```sh
cd /chirpbird-v2/
./start-server.sh
```

and then open your browser and access this URL

```sh
http://localhost/
```

You will see 1 button which you can use to login to the chat app using your google account

run the below command if you want to stop the server

```sh
cd /chirpbird-v2/
./stop-server.sh
```
