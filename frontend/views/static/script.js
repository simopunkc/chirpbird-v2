let socket
let endpointFE
let currentGroupID = ""
let userID = ""
let userName = ""
let userPicture = ""
let userVerified = false

const connectWebsocket = () => {
    console.log("Attempting websocket connection")
    endpointFE = "http://localhost"
    socket = new WebSocket("ws://be.localhost/ws")
}

connectWebsocket()

const inisialisasiMember = () => {
    let url = endpointFE + "/oauth/google/profile"
    var xhttp = new XMLHttpRequest();
    xhttp.onreadystatechange = function () {
        if (this.readyState == 4 && this.status == 200) {
            let raw = JSON.parse(this.responseText);
            userID = raw.message.email
            userName = raw.message.name
            userPicture = raw.message.picture
            userVerified = raw.message.verified_email
        }
    };
    xhttp.open("GET", url, true);
    xhttp.send();
}

socket.onopen = () => {
    console.log("Succesfully connected")
    inisialisasiMember()
    listGroup(1)
}

socket.onclose = (event) => {
    console.log("Socket closed connection:", event)
    setTimeout(function () {
        connectWebsocket();
    }, 1000);
}

socket.onmessage = (msg) => {
    if (msg.data) {
        let data = JSON.parse(msg.data);
        let url = endpointFE + "/activity/" + data.Id_room_activity
        var xhttp = new XMLHttpRequest();
        xhttp.onreadystatechange = function () {
            if (this.readyState == 4 && this.status == 200) {
                try {
                    let raw = JSON.parse(this.responseText);

                    const container = document.getElementById("messages_container");
                    var chatContainer = document.createElement("div");
                    chatContainer.classList.add("chat");
                    chatContainer.setAttribute("id", raw.message.id_primary);
                    var newChat = document.createElement("div");
                    newChat.classList.add("chatcontent");
                    var clearChat = document.createElement("div");
                    clearChat.classList.add("clear");
                    if (raw.message.type_activity != "new_chat") {
                        newChat.classList.add("activity");
                    } else if (raw.message.id_member_actor == userID) {
                        newChat.classList.add("currentuser");
                    } else {
                        newChat.classList.add("otheruser");
                    }

                    var headerChat = document.createElement("div");
                    headerChat.classList.add("headerChat");
                    var bodyChat = document.createElement("div");
                    bodyChat.classList.add("bodyChat");

                    var name = document.createElement("span");
                    name.innerHTML = raw.message.id_member_actor;
                    var tanggal = document.createElement("span");
                    tanggal.innerHTML = " " + raw.message.date_created;
                    headerChat.appendChild(name);
                    headerChat.appendChild(tanggal);

                    var pesan = document.createElement("p");
                    pesan.innerHTML = raw.message.message;
                    bodyChat.appendChild(pesan);

                    var btn1 = document.createElement("button");
                    btn1.innerHTML = "delete"
                    btn1.setAttribute("onClick", 'deleteChat("' + raw.message.id_primary + '")');
                    var btn2 = document.createElement("button");
                    btn2.innerHTML = "kick member"
                    btn2.setAttribute("onClick", 'kickMember("' + raw.message.id_room + '","' + raw.message.id_member_actor + '")');
                    var btn3 = document.createElement("button");
                    btn3.innerHTML = "ubah member jadi moderator"
                    btn3.setAttribute("onClick", 'memberToModerator("' + raw.message.id_room + '","' + raw.message.id_member_actor + '")');
                    var btn4 = document.createElement("button");
                    btn4.innerHTML = "cancel moderator"
                    btn4.setAttribute("onClick", 'moderatorToMember("' + raw.message.id_room + '","' + raw.message.id_member_actor + '")');

                    newChat.appendChild(headerChat);
                    newChat.appendChild(bodyChat);
                    newChat.appendChild(btn1);
                    newChat.appendChild(btn2);
                    newChat.appendChild(btn3);
                    newChat.appendChild(btn4);

                    chatContainer.appendChild(newChat);
                    chatContainer.appendChild(clearChat);

                    container.appendChild(chatContainer)
                    container.scrollTop = container.scrollHeight;
                } catch (e) {
                    console.log(e.message)
                }
            }
        };
        xhttp.open("GET", url, true);
        xhttp.send();
    }
}

socket.onerror = (error) => {
    console.log("Socket Error:", error)
}

const sendmessage = () => {
    const input = document.getElementById("message_input")
    const msg = input.value.trim()
    if (userID == "" || userName == "" || userPicture == "") {
        alert("user id tidak valid / belum login")
    } else if (currentGroupID == "") {
        alert("group belum diload")
    } else if (msg.length) {
        let url = endpointFE + "/room/" + currentGroupID + "/newChat"
        let id_parent = ""
        let formdata = "id_parent=" + id_parent + "&message=" + msg
        var xhttp = new XMLHttpRequest();
        xhttp.onreadystatechange = function () {
            if (this.readyState == 4 && this.status == 200) {
                try {
                    let raw = JSON.parse(this.responseText);
                    socket.send(JSON.stringify({
                        Id_room_activity: raw.message.id_room_activity,
                        Id_member_actor: raw.message.id_member_actor,
                        Id_room: raw.message.id_room
                    }))
                    input.value = ""
                } catch (e) {
                    console.log(e.message)
                }
            }
        };
        xhttp.open("POST", url, true);
        xhttp.setRequestHeader("Content-type", "application/x-www-form-urlencoded");
        xhttp.send(formdata);
    } else {
        console.log("bad request")
    }
}

const closeSession = () => {
    socket.close()
    window.location.href = endpointFE + "/oauth/logout";
}

var modal1 = document.getElementById("myModal1");
var modal2 = document.getElementById("myModal2");
const openModal1 = () => {
    modal1.style.display = "block";
}
const openModal2 = () => {
    modal2.style.display = "block";
}
const closeModal1 = () => {
    modal1.style.display = "none";
}
const closeModal2 = () => {
    modal2.style.display = "none";
}

const listGroup = (page) => {
    const container = document.getElementById("group_container");
    container.innerHTML = ""
    let pages = page > 1 ? page : 1;
    let url = endpointFE + "/member/room/page" + pages
    var xhttp = new XMLHttpRequest();
    xhttp.onreadystatechange = function () {
        if (this.readyState == 4 && this.status == 200) {
            try {
                let raw = JSON.parse(this.responseText);
                raw.message.map((data) => {
                    var newGroup = document.createElement("div");
                    newGroup.classList.add("group");
                    newGroup.setAttribute("id", data.id_primary);

                    var newText = document.createElement("div");
                    newGroup.classList.add("center");
                    newText.innerHTML = data.name + " (join: " + data.link_join + ")"
                    var newBtn1 = document.createElement("button");
                    newBtn1.setAttribute("onClick", 'loadGroup("' + data.id_primary + '", 1)');
                    newBtn1.innerHTML = "load chat";
                    var newBtn2 = document.createElement("button");
                    newBtn2.setAttribute("onClick", 'exitGroup("' + data.id_primary + '")');
                    newBtn2.innerHTML = "exit";
                    var newBtn3 = document.createElement("button");
                    newBtn3.setAttribute("onClick", 'deleteGroup("' + data.id_primary + '")');
                    newBtn3.innerHTML = "delete";
                    var newBtn4 = document.createElement("button");
                    newBtn4.setAttribute("onClick", 'renameGroup("' + data.id_primary + '")');
                    newBtn4.innerHTML = "rename";
                    var newBtn5 = document.createElement("button");
                    newBtn5.setAttribute("onClick", 'addMember("' + data.id_primary + '")');
                    newBtn5.innerHTML = "add member";
                    var newBtn6 = document.createElement("button");
                    newBtn6.setAttribute("onClick", 'enableNotif("' + data.id_primary + '")');
                    newBtn6.innerHTML = "enable notif";
                    var newBtn7 = document.createElement("button");
                    newBtn7.setAttribute("onClick", 'disableNotif("' + data.id_primary + '")');
                    newBtn7.innerHTML = "disable notif";

                    newGroup.appendChild(newText);
                    newGroup.appendChild(newBtn1);
                    newGroup.appendChild(newBtn2);
                    newGroup.appendChild(newBtn3);
                    newGroup.appendChild(newBtn4);
                    newGroup.appendChild(newBtn5);
                    newGroup.appendChild(newBtn6);
                    newGroup.appendChild(newBtn7);
                    container.appendChild(newGroup);
                });
                container.scrollTop = container.scrollHeight;
            } catch (e) {
                console.log(e.message)
            }
        }
    };
    xhttp.open("GET", url, true);
    xhttp.send();
}

const loadGroup = (id, page) => {
    currentGroupID = id
    const container = document.getElementById("messages_container");
    let pages = page > 1 ? page : 1;
    let url = endpointFE + "/messenger/" + id + "/page" + pages
    container.innerHTML = ""
    var xhttp = new XMLHttpRequest();
    xhttp.onreadystatechange = function () {
        if (this.readyState == 4 && this.status == 200) {
            try {
                let raw = JSON.parse(this.responseText);
                raw.message.map((data) => {
                    var chatContainer = document.createElement("div");
                    chatContainer.classList.add("chat");
                    chatContainer.setAttribute("id", data.id_primary);
                    var newChat = document.createElement("div");
                    newChat.classList.add("chatcontent");
                    var clearChat = document.createElement("div");
                    clearChat.classList.add("clear");
                    if (data.type_activity != "new_chat") {
                        newChat.classList.add("activity");
                    } else if (data.id_member_actor == userID) {
                        newChat.classList.add("currentuser");
                    } else {
                        newChat.classList.add("otheruser");
                    }

                    var headerChat = document.createElement("div");
                    headerChat.classList.add("headerChat");
                    var bodyChat = document.createElement("div");
                    bodyChat.classList.add("bodyChat");

                    var name = document.createElement("span");
                    name.innerHTML = data.id_member_actor;
                    var tanggal = document.createElement("span");
                    tanggal.innerHTML = " " + data.date_created;
                    headerChat.appendChild(name);
                    headerChat.appendChild(tanggal);

                    var pesan = document.createElement("p");
                    pesan.innerHTML = data.message;
                    bodyChat.appendChild(pesan);

                    var btn1 = document.createElement("button");
                    btn1.innerHTML = "delete"
                    btn1.setAttribute("onClick", 'deleteChat("' + data.id_primary + '")');
                    var btn2 = document.createElement("button");
                    btn2.innerHTML = "kick member"
                    btn2.setAttribute("onClick", 'kickMember("' + data.id_room + '","' + data.id_member_actor + '")');
                    var btn3 = document.createElement("button");
                    btn3.innerHTML = "become moderator"
                    btn3.setAttribute("onClick", 'memberToModerator("' + data.id_room + '","' + data.id_member_actor + '")');
                    var btn4 = document.createElement("button");
                    btn4.innerHTML = "cancel moderator"
                    btn4.setAttribute("onClick", 'moderatorToMember("' + data.id_room + '","' + data.id_member_actor + '")');

                    newChat.appendChild(headerChat);
                    newChat.appendChild(bodyChat);
                    newChat.appendChild(btn1);
                    newChat.appendChild(btn2);
                    newChat.appendChild(btn3);
                    newChat.appendChild(btn4);
                    chatContainer.appendChild(newChat);
                    chatContainer.appendChild(clearChat);

                    container.appendChild(chatContainer);
                    container.scrollTop = container.scrollHeight;
                });
                const buttons = document.querySelectorAll('.group');
                buttons.forEach(button => {
                    if (button.id == id) {
                        button.classList.add('groupaktif')
                    } else {
                        button.classList.remove('groupaktif')
                    }
                })
            } catch (e) {
                console.log(e.message)
            }
        }
    };
    xhttp.open("GET", url, true);
    xhttp.send();
}

const createGroup = () => {
    const input = document.getElementById("group_input")
    const msg = input.value.trim()
    input.value = ""
    if (msg.length) {
        let url = endpointFE + "/member/room/create"
        let formdata = "name=" + msg
        var xhttp = new XMLHttpRequest();
        xhttp.onreadystatechange = function () {
            if (this.readyState == 4 && this.status == 200) {
                closeModal1()
                listGroup(1)
            }
        };
        xhttp.open("POST", url, true);
        xhttp.setRequestHeader("Content-type", "application/x-www-form-urlencoded");
        xhttp.send(formdata);
    }
}

const renameGroup = (id) => {
    let newname = prompt("Please enter new group name:");
    if (newname != null && newname != "") {
        let url = endpointFE + "/room/" + id + "/rename"
        let formdata = "name=" + newname
        var xhttp = new XMLHttpRequest();
        xhttp.onreadystatechange = function () {
            if (this.readyState == 4 && this.status == 200) {
                try {
                    let raw = JSON.parse(this.responseText);
                    socket.send(JSON.stringify({
                        Id_room_activity: raw.message.id_room_activity,
                        Id_member_actor: raw.message.id_member_actor,
                        Id_room: raw.message.id_room
                    }))
                    listGroup(1)
                } catch (e) {
                    console.log(e.message)
                }
            }
        };
        xhttp.open("PUT", url, true);
        xhttp.setRequestHeader("Content-type", "application/x-www-form-urlencoded");
        xhttp.send(formdata);
    }
}

const addMember = (id) => {
    let newemail = prompt("insert email member");
    if (newemail != null && newemail != "") {
        let url = endpointFE + "/room/" + id + "/addMember"
        let formdata = "id_target=" + newemail
        var xhttp = new XMLHttpRequest();
        xhttp.onreadystatechange = function () {
            if (this.readyState == 4 && this.status == 200) {
                try {
                    let raw = JSON.parse(this.responseText);
                    socket.send(JSON.stringify({
                        Id_room_activity: raw.message.id_room_activity,
                        Id_member_actor: raw.message.id_member_actor,
                        Id_room: raw.message.id_room
                    }))
                } catch (e) {
                    console.log(e.message)
                }
            }
        };
        xhttp.open("PUT", url, true);
        xhttp.setRequestHeader("Content-type", "application/x-www-form-urlencoded");
        xhttp.send(formdata);
    }
}

const joinGroup = () => {
    const input = document.getElementById("join_input")
    const msg = input.value.trim()
    input.value = ""
    if (msg.length) {
        let url = endpointFE + "/member/room/join"
        let formdata = "token=" + msg
        var xhttp = new XMLHttpRequest();
        xhttp.onreadystatechange = function () {
            if (this.readyState == 4 && this.status == 200) {
                try {
                    let raw = JSON.parse(this.responseText);
                    closeModal2()
                    socket.send(JSON.stringify({
                        Id_room_activity: raw.message.id_room_activity,
                        Id_member_actor: raw.message.id_member_actor,
                        Id_room: raw.message.id_room
                    }))
                    listGroup(1)
                } catch (e) {
                    console.log(e.message)
                }
            }
        };
        xhttp.open("PUT", url, true);
        xhttp.setRequestHeader("Content-type", "application/x-www-form-urlencoded");
        xhttp.send(formdata);
    }
}

const exitGroup = (id) => {
    let text = "exit from group?";
    if (confirm(text)) {
        let url = endpointFE + "/room/" + id + "/exit"
        var xhttp = new XMLHttpRequest();
        xhttp.onreadystatechange = function () {
            if (this.readyState == 4 && this.status == 200) {
                try {
                    let raw = JSON.parse(this.responseText);
                    var myobj = document.getElementById(id);
                    myobj.remove();
                    socket.send(JSON.stringify({
                        Id_room_activity: raw.message.id_room_activity,
                        Id_member_actor: raw.message.id_member_actor,
                        Id_room: raw.message.id_room
                    }))
                } catch (e) {
                    console.log(e.message)
                }
            }
        };
        xhttp.open("PUT", url, true);
        xhttp.setRequestHeader("Content-type", "application/x-www-form-urlencoded");
        xhttp.send();
    }
}

const enableNotif = (id) => {
    let text = "enable group notification?";
    if (confirm(text)) {
        let url = endpointFE + "/room/" + id + "/enableNotif"
        var xhttp = new XMLHttpRequest();
        xhttp.onreadystatechange = function () {
            if (this.readyState == 4 && this.status == 200) {
                console.log(this.responseText)
            }
        };
        xhttp.open("PUT", url, true);
        xhttp.setRequestHeader("Content-type", "application/x-www-form-urlencoded");
        xhttp.send();
    }
}

const disableNotif = (id) => {
    let text = "disable group notification?";
    if (confirm(text)) {
        let url = endpointFE + "/room/" + id + "/disableNotif"
        var xhttp = new XMLHttpRequest();
        xhttp.onreadystatechange = function () {
            if (this.readyState == 4 && this.status == 200) {
                console.log(this.responseText)
            }
        };
        xhttp.open("PUT", url, true);
        xhttp.setRequestHeader("Content-type", "application/x-www-form-urlencoded");
        xhttp.send();
    }
}

const deleteGroup = (id) => {
    let text = "Delete group permanently?";
    if (confirm(text)) {
        let url = endpointFE + "/room/" + id + "/deleteRoom"
        var xhttp = new XMLHttpRequest();
        xhttp.onreadystatechange = function () {
            if (this.readyState == 4 && this.status == 200) {
                let raw = JSON.parse(this.responseText);
                if (raw.status == 200) {
                    var myobj = document.getElementById(id);
                    myobj.remove();
                    const container = document.getElementById("group_container");
                    container.innerHTML = ""
                }
            }
        };
        xhttp.open("DELETE", url, true);
        xhttp.send();
    }
}

const deleteChat = (id) => {
    let text = "Delete activity group permanently?";
    if (confirm(text)) {
        let url = endpointFE + "/activity/" + id + "/deleteChat"
        var xhttp = new XMLHttpRequest();
        xhttp.onreadystatechange = function () {
            if (this.readyState == 4 && this.status == 200) {
                let raw = JSON.parse(this.responseText);
                if (raw.status == 200) {
                    var myobj = document.getElementById(id);
                    myobj.remove();
                }
            }
        };
        xhttp.open("DELETE", url, true);
        xhttp.send();
    }
}

const kickMember = (id_room, id_target) => {
    let text = "Banned member from group?";
    if (confirm(text)) {
        let url = endpointFE + "/room/" + id_room + "/kickMember"
        let formdata = "id_target=" + id_target
        var xhttp = new XMLHttpRequest();
        xhttp.onreadystatechange = function () {
            if (this.readyState == 4 && this.status == 200) {
                try {
                    let raw = JSON.parse(this.responseText);
                    if (raw.status == 200) {
                        socket.send(JSON.stringify({
                            Id_room_activity: raw.message.id_room_activity,
                            Id_member_actor: raw.message.id_member_actor,
                            Id_room: raw.message.id_room
                        }))
                    }
                } catch (e) {
                    console.log(e.message)
                }
            }
        };
        xhttp.open("PUT", url, true);
        xhttp.setRequestHeader("Content-type", "application/x-www-form-urlencoded");
        xhttp.send(formdata);
    }
}

const memberToModerator = (id_room, id_target) => {
    let text = "become moderator?";
    if (confirm(text)) {
        let url = endpointFE + "/room/" + id_room + "/memberToModerator"
        let formdata = "id_target=" + id_target
        var xhttp = new XMLHttpRequest();
        xhttp.onreadystatechange = function () {
            if (this.readyState == 4 && this.status == 200) {
                try {
                    let raw = JSON.parse(this.responseText);
                    if (raw.status == 200) {
                        socket.send(JSON.stringify({
                            Id_room_activity: raw.message.id_room_activity,
                            Id_member_actor: raw.message.id_member_actor,
                            Id_room: raw.message.id_room
                        }))
                    }
                } catch (e) {
                    console.log(e.message)
                }
            }
        };
        xhttp.open("PUT", url, true);
        xhttp.setRequestHeader("Content-type", "application/x-www-form-urlencoded");
        xhttp.send(formdata);
    }
}

const moderatorToMember = (id_room, id_target) => {
    let text = "cancel moderator?";
    if (confirm(text)) {
        let url = endpointFE + "/room/" + id_room + "/ModeratorToMember"
        let formdata = "id_target=" + id_target
        var xhttp = new XMLHttpRequest();
        xhttp.onreadystatechange = function () {
            if (this.readyState == 4 && this.status == 200) {
                try {
                    let raw = JSON.parse(this.responseText);
                    if (raw.status == 200) {
                        socket.send(JSON.stringify({
                            Id_room_activity: raw.message.id_room_activity,
                            Id_member_actor: raw.message.id_member_actor,
                            Id_room: raw.message.id_room
                        }))
                    }
                } catch (e) {
                    console.log(e.message)
                }
            }
        };
        xhttp.open("PUT", url, true);
        xhttp.setRequestHeader("Content-type", "application/x-www-form-urlencoded");
        xhttp.send(formdata);
    }
}