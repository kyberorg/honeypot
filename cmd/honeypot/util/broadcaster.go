package util

import "github.com/kyberorg/honeypot/cmd/honeypot/dto"

// taken from here https://stackoverflow.com/questions/36417199/how-to-broadcast-message-using-channel

type LoginAttemptBroadcaster struct {
	stopCh    chan struct{}
	publishCh chan *dto.LoginAttempt
	subCh     chan chan *dto.LoginAttempt
	unsubCh   chan chan *dto.LoginAttempt
}

func NewLoginAttemptBroadcaster() *LoginAttemptBroadcaster {
	return &LoginAttemptBroadcaster{
		stopCh:    make(chan struct{}),
		publishCh: make(chan *dto.LoginAttempt, 1),
		subCh:     make(chan chan *dto.LoginAttempt, 1),
		unsubCh:   make(chan chan *dto.LoginAttempt, 1),
	}
}

func (b *LoginAttemptBroadcaster) Start() {
	subs := map[chan *dto.LoginAttempt]struct{}{}
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

func (b *LoginAttemptBroadcaster) Stop() {
	close(b.stopCh)
}

func (b *LoginAttemptBroadcaster) Subscribe() chan *dto.LoginAttempt {
	msgCh := make(chan *dto.LoginAttempt, 5)
	b.subCh <- msgCh
	return msgCh
}

func (b *LoginAttemptBroadcaster) Unsubscribe(msgCh chan *dto.LoginAttempt) {
	b.unsubCh <- msgCh
}

func (b *LoginAttemptBroadcaster) Send(msg *dto.LoginAttempt) {
	b.publishCh <- msg
}
