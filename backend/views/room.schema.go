package views

type DatabaseRoom struct {
	Id_primary                         string   `json:"id_primary,omitempty"`
	Id_member_creator                  string   `json:"id_member_creator,omitempty"`
	Name                               string   `json:"name,omitempty"`
	List_id_member                     []string `json:"list_id_member,omitempty"`
	List_id_member_moderator           []string `json:"list_id_member_moderator,omitempty"`
	List_id_member_banned              []string `json:"list_id_member_banned,omitempty"`
	List_id_member_enable_notification []string `json:"list_id_member_enable_notification,omitempty"`
	Date_created                       string   `json:"date_created,omitempty"`
	Date_last_activity                 string   `json:"date_last_activity,omitempty"`
	Link_join                          string   `json:"link_join,omitempty"`
}
