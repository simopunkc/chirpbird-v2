# ChirpBird

Instant messaging platform with Websocket using Golang

## how to run this project

open terminal and run command below

clone this project

```sh
git clone https://github.com/simopunkc/chirpbird-v2.git
```

if you want to run it on localhost then you have to edit .env.example file in backend folder

nano /chirpbird-v2/backend/.env.example

```sh
OAUTH_CLIENT_ID='put your Google Oauth Client ID here'
OAUTH_SECRET='put your Google Oauth Secret here'
FRONTEND_HOST='localhost:9000'
FRONTEND_DOMAIN='localhost'
FRONTEND_PROTOCOL='http://'
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
