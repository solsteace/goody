package crypto

type Cryptor interface {
	Generate(payload string) ([]byte, error)
	Compare(digest, payload string) error
}
