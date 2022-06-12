# Schema Databases

I am using MongoDB database with 3 collections

## member

member data is obtained after the user successfully login with Google Oauth

Example data structure:

```sh
{
    _id: ObjectId("61f6892a47cbdd690fc107be"),
    email: 'pembela.allah@gmail.com',
    name: 'budiyanto simo',
    picture: 'https://lh3.googleusercontent.com/a-/AOh14Gjif-lTqIRwpom74lk2uqWt-oKihr_WpCSsJNLCNA=s96-c',
    verified_email: true
}
```

## room

each member can create his own private room

Example data structure:

```sh
{
    _id: ObjectId("61f686d0062c44a112088ba6"),
    id_primary: 'GyKeWjvydfn',
    id_member_creator: 'budiyanto.simo@gmail.com',
    name: 'pemburu jamur barat',
    list_id_member: [ 'budiyanto.simo@gmail.com' ],
    list_id_member_moderator: [ 'budiyanto.simo@gmail.com' ],
    list_id_member_banned: [],
    list_id_member_enable_notification: [],
    date_created: '2022-01-30T19:38:40.696166798+07:00',
    date_last_activity: '2022-01-30T19:38:40.696166798+07:00',
    link_join: 'LJ8m287TdJys'
}
```


## room_activity

every log activity in the room will be stored in the database

Example data structure:

```sh
 {
    _id: ObjectId("61f647b224c1b250de71c826"),
    id_primary: 'RAvKCc61kmk2',
    id_parent: '',
    id_room: 'GyKeWjvydfn',
    id_member_actor: 'pembela.allah@gmail.com',
    id_member_target: '',
    type_activity: 'group_created',
    message: 'budiyanto simo membuat group baru',
    date_created: '2022-01-30T14:40:23.133120073+07:00',
    list_id_member_unread: []
}
```