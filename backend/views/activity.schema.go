package views

type DatabaseRoomActivity struct {
	Id_primary            string   `json:"id_primary,omitempty"`
	Id_parent             string   `json:"id_parent,omitempty"`
	Id_room               string   `json:"id_room,omitempty"`
	Id_member_actor       string   `json:"id_member_actor,omitempty"`
	Id_member_target      string   `json:"id_member_target,omitempty"`
	Type_activity         string   `json:"type_activity,omitempty"`
	Message               string   `json:"message,omitempty"`
	Date_created          string   `json:"date_created,omitempty"`
	List_id_member_unread []string `json:"list_id_member_unread,omitempty"`
}
