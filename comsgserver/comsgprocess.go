package msgserver

import (
	"codeskserver/log"
	"sync"
)

type ProcMsgType int

const (
	RecvMsg int = iota
	RecvEventTmr
	RecvEventQuit
)

type ProcMsg struct {
	Type     int //body类型
	ClientID string
	Topic    string
	Body     []byte
}

type CoMsgProcess struct {
	isLoop int
	job    sync.WaitGroup
	msgQ   chan *ProcMsg
}

// 初始化
func (p *CoMsgProcess) Init() {
	p.msgQ = make(chan *ProcMsg, 128)
	p.isLoop = 0
}

func (p *CoMsgProcess) onNewMsg(msg *ProcMsg) {
	//这个地方我开始处理消息

	log.Debug("onNewMsg from: %s, topic:%s, msg:%s", msg.ClientID, msg.Topic, string(msg.Body))
}

// 外部调用这个发送消息出去
func (p *CoMsgProcess) PushRecvMsg(clientid string, topic string, body []byte) {
	msg := &ProcMsg{
		Type:     RecvMsg,
		ClientID: clientid,
		Topic:    topic,
		Body:     body,
	}

	p.msgQ <- msg
}

func (p *CoMsgProcess) run() {
	p.job.Add(1)
	for {
		if p.isLoop == 0 {
			break
		}

		select {
		case msg := <-p.msgQ:
			//处理消息
			switch msg.Type {
			case RecvMsg:
				p.onNewMsg(msg)
			case RecvEventTmr:
				//处理事件
			case RecvEventQuit:
				//处理事件
			}
		}
	}
	p.job.Done()
}

// 启动
func (p *CoMsgProcess) Start() {
	if p.isLoop == 1 {
		return
	}
	p.isLoop = 1
	go p.run()
}

// 停止
func (p *CoMsgProcess) Stop() {
	p.isLoop = 0
	msg := &ProcMsg{
		Type: RecvEventQuit,
	}
	p.msgQ <- msg
	p.job.Wait()
}
