package tms

type Manager interface {
	Upload() error
	Translate() (string, error)
}
