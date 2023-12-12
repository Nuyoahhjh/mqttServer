package msgserver

import (
	"codeskserver/log"
	"sync"
)

const (
	ClientMsgQuit int = iota
	ClientMsgAdd
	ClientMsgRemove
	ClientMsgUpdate
)

type CoClientMsg struct {
	Type int //body类型
	item *CoClientItem
}

type CoClientItem struct {
	ClientID   string
	ClientType int
	Token      string
	AccessID   string
	Addr       string
	Os         string
	Info       string
}

type CoClients struct {
	isLoop int
	job    sync.WaitGroup
	msgQ   chan *CoClientMsg
	lock   sync.Mutex
	items  map[string]*CoClientMsg
}

func (c *CoClients) Init() {
	c.isLoop = 0
	c.msgQ = make(chan *CoClientMsg, 128)
	c.items = make(map[string]*CoClientMsg)
}

func (c *CoClients) AddClient() {

}

func (c *CoClients) RemvoeClient(clientid string) {
	msg := &CoClientMsg{
		Type: ClientMsgQuit,
	}

	msg.item = &CoClientItem{
		ClientID: clientid,
	}

	c.msgQ <- msg
}

func (c *CoClients) Start() {
	if c.isLoop == 1 {
		return
	}
	c.isLoop = 1
	go c.run()
}

func (c *CoClients) Stop() {
	c.isLoop = 0
	msg := &CoClientMsg{
		Type: ClientMsgQuit,
	}
	c.msgQ <- msg
	c.job.Wait()
}

func (c *CoClients) run() {

	c.job.Add(1)
	for {
		if c.isLoop == 0 {
			break
		}

		msg := <-c.msgQ
		switch msg.Type {
		case ClientMsgQuit:
			c.isLoop = 0
		case ClientMsgAdd:
			c.lock.Lock()

			log.Debug("try to add client id:%s", msg.item.ClientID)

			c.lock.Unlock()
		case ClientMsgRemove:
			c.lock.Lock()

			log.Debug("try to delete client id:%s", msg.item.ClientID)

			delete(c.items, msg.item.ClientID)

			c.lock.Unlock()
		}
	}
	c.job.Done()
}
