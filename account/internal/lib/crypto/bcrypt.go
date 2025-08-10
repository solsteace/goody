package crypto

import bc "golang.org/x/crypto/bcrypt"

type bcrypt struct {
	cost int
}

func NewBcrypt(cost int) bcrypt {
	return bcrypt{cost: cost}
}

func (b bcrypt) Generate(payload string) ([]byte, error) {
	return bc.GenerateFromPassword([]byte(payload), b.cost)
}

func (b bcrypt) Compare(digest, payload string) error {
	return bc.CompareHashAndPassword([]byte(digest), []byte(payload))
}
