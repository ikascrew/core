package core

type Effect interface {
	Set(interface{}) error
	Convert(Video) (Video, error)
}
