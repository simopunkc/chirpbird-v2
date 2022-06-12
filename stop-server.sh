#!/bin/bash

docker-compose down
docker-compose rm -f
sleep 5
docker volume rm $(docker volume ls -q)
sleep 5
docker rmi chirpbird-v2_fe1
docker rmi chirpbird-v2_fe2
docker rmi chirpbird-v2_fe3
docker rmi chirpbird-v2_be1
docker rmi chirpbird-v2_be2
docker rmi chirpbird-v2_be3