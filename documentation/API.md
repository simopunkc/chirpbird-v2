# API

backend

## Login Oauth URL

Method

```sh
GET
```

Endpoint

```sh
/oauth/google/url
```

## Verify Login Oauth

Method

```sh
POST
```

Endpoint

```sh
/oauth/google/verify
```

Header

```sh
{
    xsrf_token: "blablabla"
}
```

Body

```sh
{
    code: "blablabla",
    state: "blablabla"
}
```

## Refresh Login Oauth

Method

```sh
POST
```

Endpoint

```sh
/oauth/google/refresh
```

Header

```sh
{
    xsrf_token: "blablabla"
}
```

## Oauth Profile

Method

```sh
GET
```

Endpoint

```sh
/oauth/google/profile
```

Header

```sh
{
    acc_token: "blablabla"
}
```

## Log Out

Method

```sh
POST
```

Endpoint

```sh
/oauth/logout
```

Header

```sh
{
    xsrf_token: "blablabla"
}
```

## Token CSRF

Method

```sh
GET
```

Endpoint

```sh
/token/csrf
```

## List Room

Method

```sh
GET
```

Endpoint

```sh
/member/room/page:PAGE_NUMBER
```

Example Path URL

```sh
/member/room/page1
```

Header

```sh
{
    acc_token: "blablabla"
}
```

## Create Room

Method

```sh
POST
```

Endpoint

```sh
/member/room/create
```

Header

```sh
{
    acc_token: "blablabla"
}
```

Body

```sh
{
    name: "blablabla"
}
```

## Join Room

Method

```sh
PUT
```

Endpoint

```sh
/member/room/join
```

Header

```sh
{
    acc_token: "blablabla"
}
```

Body

```sh
{
    token: "blablabla"
}
```

## Room Data

Method

```sh
GET
```

Endpoint

```sh
/room/:ID_ROOM
```

Example Path URL

```sh
/room/GyKeWjvydfn
```

Header

```sh
{
    acc_token: "blablabla"
}
```

## Exit Room

Method

```sh
PUT
```

Endpoint

```sh
/room/:ID_ROOM/exit
```

Header

```sh
{
    acc_token: "blablabla"
}
```

## Rename Room

Method

```sh
PUT
```

Endpoint

```sh
/room/:ID_ROOM/rename
```

Header

```sh
{
    acc_token: "blablabla"
}
```

Body

```sh
{
    name: "blablabla"
}
```

## Member Become Moderator

Method

```sh
PUT
```

Endpoint

```sh
/room/:ID_ROOM/memberToModerator
```

Header

```sh
{
    acc_token: "blablabla"
}
```

Body

```sh
{
    id_target: "blablabla"
}
```

## Cancel Moderator

Method

```sh
PUT
```

Endpoint

```sh
/room/:ID_ROOM/ModeratorToMember
```

Header

```sh
{
    acc_token: "blablabla"
}
```

Body

```sh
{
    id_target: "blablabla"
}
```

## Banned Member In Room

Method

```sh
PUT
```

Endpoint

```sh
/room/:ID_ROOM/kickMember
```

Header

```sh
{
    acc_token: "blablabla"
}
```

Body

```sh
{
    id_target: "blablabla"
}
```

## Add Member Into Room

Method

```sh
PUT
```

Endpoint

```sh
/room/:ID_ROOM/addMember
```

Header

```sh
{
    acc_token: "blablabla"
}
```

Body

```sh
{
    id_target: "blablabla"
}
```

## Publish Message Into Room

Method

```sh
POST
```

Endpoint

```sh
/room/:ID_ROOM/newChat
```

Header

```sh
{
    acc_token: "blablabla"
}
```

Body

```sh
{
    id_parent: "blablabla",
    message: "blablabla"
}
```

## Enable Room Notification

Method

```sh
PUT
```

Endpoint

```sh
/room/:ID_ROOM/enableNotif
```

Header

```sh
{
    acc_token: "blablabla"
}
```


## Disable Room Notification

Method

```sh
PUT
```

Endpoint

```sh
/room/:ID_ROOM/disableNotif
```

Header

```sh
{
    acc_token: "blablabla"
}
```

## Delete Room

Method

```sh
DELETE
```

Endpoint

```sh
/room/:ID_ROOM/deleteRoom
```

Header

```sh
{
    acc_token: "blablabla"
}
```

## List Message In Room

Method

```sh
GET
```

Endpoint

```sh
/messenger/:ID_ROOM/page{pid:[0-9]+}
```

Header

```sh
{
    acc_token: "blablabla"
}
```

## Show Single Message

Method

```sh
GET
```

Endpoint

```sh
/activity/:ID_MESSAGE
```

Header

```sh
{
    acc_token: "blablabla"
}
```

## Delete Message

Method

```sh
DELETE
```

Endpoint

```sh
/activity/:ID_MESSAGE/deleteChat
```

Header

```sh
{
    acc_token: "blablabla"
}
```
