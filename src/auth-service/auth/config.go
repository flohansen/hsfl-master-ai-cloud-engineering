package auth

type Config interface {
	GetPrivateKey() (any, error)
}
