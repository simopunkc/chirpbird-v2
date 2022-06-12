package views

func GetActivity(type_activity string, actor_name string, target_name string) string {
	var tipe string
	switch type_activity {
	case "group_created":
		tipe = actor_name + " create a new group"
	case "join_group":
		tipe = actor_name + " join to the group using link join"
	case "exit_group":
		tipe = actor_name + " exit from the group"
	case "add_member":
		tipe = actor_name + " add " + target_name + " to the group"
	case "kick_member":
		tipe = actor_name + " removed " + target_name + " from the group"
	case "member_to_moderator":
		tipe = actor_name + " update " + target_name + " to become moderator"
	case "moderator_to_member":
		tipe = actor_name + " cancel " + target_name + " as moderator"
	case "rename_group":
		tipe = actor_name + " rename the group"
	}
	return tipe
}
