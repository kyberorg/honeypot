package util

import "github.com/kyberorg/honeypot/cmd/honeypot/dto"

// taken from here https://stackoverflow.com/questions/36417199/how-to-broadcast-message-using-channel

type ConnectionBroker struct {
	stopCh    chan struct{}
	publishCh chan *dto.CollectedData
	subCh     chan chan *dto.CollectedData
	unsubCh   chan chan *dto.CollectedData
}

func NewBroker() *ConnectionBroker {
	return &ConnectionBroker{
		stopCh:    make(chan struct{}),
		publishCh: make(chan *dto.CollectedData, 1),
		subCh:     make(chan chan *dto.CollectedData, 1),
		unsubCh:   make(chan chan *dto.CollectedData, 1),
	}
}

func (b *ConnectionBroker) Start() {
	subs := map[chan *dto.CollectedData]struct{}{}
	for {
		select {
		case <-b.stopCh:
			for msgCh := range subs {
				close(msgCh)
			}
		case msgCh := <-b.subCh:
			subs[msgCh] = struct{}{}
		case msgCh := <-b.unsubCh:
			delete(subs, msgCh)
		case msg := <-b.publishCh:
			for msgCh := range subs {
				// msgCh is buffered, use non-blocking send to protect the broker:
				select {
				case msgCh <- msg:
				default:
				}
			}
		}
	}
}

func (b *ConnectionBroker) Stop() {
	close(b.stopCh)
}

func (b *ConnectionBroker) Subscribe() chan *dto.CollectedData {
	msgCh := make(chan *dto.CollectedData, 5)
	b.subCh <- msgCh
	return msgCh
}

func (b *ConnectionBroker) Unsubscribe(msgCh chan *dto.CollectedData) {
	b.unsubCh <- msgCh
}

func (b *ConnectionBroker) Publish(msg *dto.CollectedData) {
	b.publishCh <- msg
}
