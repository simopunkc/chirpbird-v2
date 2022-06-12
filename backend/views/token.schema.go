package views

type Csrftoken struct {
	Random  string
	Expired int64
}

type RawData struct {
	Plaintext  string
	Ciphertext string
}

type TempRefreshToken struct {
	Value   string
	Expired int64
}
