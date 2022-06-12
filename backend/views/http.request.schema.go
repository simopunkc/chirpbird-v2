package views

type BodyRequestVerifyLogin struct {
	State string `json:"state"`
	Code  string `json:"code"`
}

type BodyRequestRefreshLogin struct {
	Ref_token string `json:"ref_token"`
}

type BodyRequestCreateRoom struct {
	Name string `json:"name"`
}

type BodyRequestCreateRoomActivity struct {
	Id_parent string `json:"id_parent"`
	Message   string `json:"message"`
}

type BodyRequestJoinRoom struct {
	Token string `json:"token"`
}

type BodyRequestManipulateRoomMember struct {
	Id_target string `json:"id_target"`
}
