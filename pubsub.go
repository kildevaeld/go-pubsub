// Package pubsub implement the Publish/Subscribe messaging paradigm
// where (citing Wikipedia) senders (publishers) are not programmed to
// send their messages to specific receivers (subscribers).
package pubsub

import (
	"path/filepath"
)

// Pubsub implement the Publish/Subscribe messaging paradigm.
type Pubsub struct {
	channels map[string][]chan interface{}
	patterns map[string][]chan interface{}
}

// Create a Pubsub.
func New() *Pubsub {
	return &Pubsub{
		channels: make(map[string][]chan interface{}),
		patterns: make(map[string][]chan interface{}),
	}
}

// Subscribe the message with specified name and send to channel c.
func (p *Pubsub) Subscribe(name string, c chan interface{}) {
	if c == nil {
		return
	}
	chans, ok := p.channels[name]
	if !ok {
		chans = []chan interface{}{c}
	} else {
		for _, ch := range chans {
			if ch == c {
				return
			}
		}
		chans = append(chans, c)
	}
	p.channels[name] = chans
}

// Unsubscribe the channel c with specified name.
func (p *Pubsub) Unsubscribe(name string, c chan interface{}) {
	if c == nil {
		return
	}
	chans, ok := p.channels[name]
	if !ok {
		return
	}
	for i := len(chans) - 1; i >= 0; i-- {
		if chans[i] == c {
			chans = append(chans[:i], chans[i+1:]...)
		}
	}
	if len(chans) == 0 {
		delete(p.channels, name)
	} else {
		p.channels[name] = chans
	}
}

// Subscribes the message with the specified pattern and send to channel c.
// Pattern supported glob-style patterns:
//
//  - h?llo matches hello, hallo and hxllo
//  - h*llo matches hllo and heeeello
//  - h[ae]llo matches hello and hallo, but not hillo
func (p *Pubsub) PSubscribe(pattern string, c chan interface{}) {
	if c == nil {
		return
	}
	chans, ok := p.patterns[pattern]
	if !ok {
		chans = []chan interface{}{c}
	} else {
		for _, ch := range chans {
			if ch == c {
				return
			}
		}
		chans = append(chans, c)
	}
	p.patterns[pattern] = chans
}

// Unsubscribes the channel c with the specified pattern.
func (p *Pubsub) PUnsubscribe(pattern string, c chan interface{}) {
	if c == nil {
		return
	}
	chans, ok := p.patterns[pattern]
	if !ok {
		return
	}
	for i := len(chans) - 1; i >= 0; i-- {
		if chans[i] == c {
			chans = append(chans[:i], chans[i+1:]...)
		}
	}
	if len(chans) == 0 {
		delete(p.patterns, pattern)
	} else {
		p.patterns[pattern] = chans
	}
}

// Publish a message with specifid name. Publish won't be blocked by channel receiving,
// if a channel doesn't ready when publish, it will be ignored.
func (p *Pubsub) Publish(name string, message interface{}) {
	if chans, ok := p.channels[name]; ok {
		for _, c := range chans {
			select {
			case c <- message:
			default:
			}
		}
	}
	for pattern, chans := range p.patterns {
		if ok, err := filepath.Match(pattern, name); err == nil && ok {
			for _, c := range chans {
				select {
				case c <- message:
				default:
				}
			}
		}
	}
}
