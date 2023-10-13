package auth

type Config interface {
	ReadPrivateKey() (any, error)
	ReadPublicKey() (any, error)
}
