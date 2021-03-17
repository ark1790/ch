package localqueue

type LocalQueue struct {
	ch chan string
}

func NewLocalQueue() *LocalQueue {
	return &LocalQueue{
		ch: make(chan string, 1024),
	}
}

func (l *LocalQueue) Push(msg string) error {
	l.ch <- msg
	return nil
}

func (l *LocalQueue) Pop(ack <-chan bool) (*string, error) {
	select {
	case msg := <-l.ch:
		return &msg, nil
	default:
		return nil, nil
	}
}
