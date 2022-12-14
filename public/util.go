package public

import (
	"crypto/sha256"
	"fmt"
)

func GenSaltPassword(salt, pasword string) string {
	s1 := sha256.New()
	s1.Write([]byte(pasword))
	str1 := fmt.Sprintf("%x", s1.Sum(nil))

	s2 := sha256.New()
	s2.Write([]byte(str1 + salt))
	return fmt.Sprintf("%x", s2.Sum(nil))
}
