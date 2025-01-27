package crypto

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

func main() {
	password := "password"

	hashedpassword, _ := bcrypt.GenerateFromPassword([]byte(password), 12)

	fmt.Println(string(hashedpassword))
}
