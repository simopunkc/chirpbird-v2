## Traefik

we can use Traefik to help us manage SSL, manage multiple load balancer, manage multiple website with multiple port (80, 443, etc), etc. in this project I only use 1 Traefik container as reverse proxy. all microservices will depend on this container. so this container must be placed on a server with the highest availability level that is safe from various kinds of disasters such as power outages, earthquakes, and so on.

## Haproxy

Haproxy can be used as load balancer. in this project there are Front End app and Back End app. each FE and BE totalling 3 containers and handled by 2 containers load balancer to increase high availability. each load balancer can be accessed via Traefik using an SSL domain.

## Redis

In this project I am using Redis to cache the database and for the purpose of publishing/subscribing to websockets messages. i am using 3 redis server containers and for failover i am also using 3 redis sentinel containers. the main application will contact the available redis sentinel. all redis sentinels will use a voting system to determine which redis master status is available.

## MongoDB

The main database used in this project is MongoD. to increase availability and fault tolerance, this project uses 3 mongodb containers running in cluster and replication.
