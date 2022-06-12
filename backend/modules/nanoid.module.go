package modules

import (
	"github.com/aidarkhanov/nanoid"
)

func GenerateUniqueID(prefix string) string {
	defaultAlphabet := "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	defaultSize := 10
	id, err := nanoid.Generate(defaultAlphabet, defaultSize)
	if err != nil {
		panic(err)
	}
	return prefix + id
}
