package types

type BackendInterface interface {
	GetValue(key string) (string, error)
	PutValue(key, value string) error
}
