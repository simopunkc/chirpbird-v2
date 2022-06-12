#!/bin/bash

mongosh <<EOF
rs.status();
var config = {
    "_id": "dbrs",
    "version": 1,
    "members": [
        {
            "_id": 0,
            "host": "172.16.0.2:27017",
            "priority": 3
        },
        {
            "_id": 1,
            "host": "172.16.0.3:27017",
            "priority": 2
        },
        {
            "_id": 2,
            "host": "172.16.0.4:27017",
            "priority": 1
        }
    ]
};
rs.initiate(config, { force: true });
rs.status();
rs.secondaryOk();
use chirpbird-v2;
db.createCollection("members");
db.createCollection("rooms");
db.createCollection("room_activities");
EOF