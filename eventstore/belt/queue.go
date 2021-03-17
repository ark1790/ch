package belt

type Queue interface {
	Push(string) error
	Pop(<-chan bool) (*string, error)
}
