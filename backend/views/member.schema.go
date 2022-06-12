package views

type DatabaseMember struct {
	Email          string `json:"email,omitempty"`
	Name           string `json:"name,omitempty"`
	Picture        string `json:"picture,omitempty"`
	Verified_email bool   `json:"verified_email,omitempty"`
}
