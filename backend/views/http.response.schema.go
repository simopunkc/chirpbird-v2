package views

type HttpErrorMessage struct {
	Status  int
	Message string
}

type HttpSuccessMessage struct {
	Status  int
	Message interface{}
}

type BodyPostResponseVerifyLogin struct {
	Access_token  string `json:"access_token"`
	Id_token      string `json:"id_token"`
	Refresh_token string `json:"refresh_token"`
}

type BodyPostResponseRefreshLogin struct {
	Access_token string `json:"access_token"`
	Id_token     string `json:"id_token"`
}
