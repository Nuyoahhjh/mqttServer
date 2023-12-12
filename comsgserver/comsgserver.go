package msgserver

import (
	"codeskserver/log"
	"crypto/tls"
	"fmt"
	"sync"
	"time"

	mqtt "github.com/mochi-mqtt/server/v2"
	"github.com/mochi-mqtt/server/v2/listeners"
	"github.com/mochi-mqtt/server/v2/packets"
)

type MsgBodyType int

const (
	PublishMsg MsgBodyType = iota
	TimerEvent
)

type PushMessage struct {
	Type  int //body类型
	Topic string
	Body  string
	Qos   int
}

type CoMsgServer struct {
	tcpPort int //tcp端口
	tlsPort int
	wsPort  int
	wssPort int

	isLoop int
	job    sync.WaitGroup

	pushMsgChan chan *PushMessage

	procMsg *CoMsgProcess
	clients *CoClients
}

var CoMsgServerInstance *CoMsgServer

func Instance() *CoMsgServer {
	if CoMsgServerInstance == nil {
		CoMsgServerInstance = CreateCoMsgServer(1883, 1884, 8083, 8084)
	}
	return CoMsgServerInstance
}

func CreateCoMsgServer(tcpPort int, tlsPort int, wsPort int, wssPort int) *CoMsgServer {
	svr := &CoMsgServer{
		tcpPort:     tcpPort,
		tlsPort:     tlsPort,
		wsPort:      wsPort,
		wssPort:     wssPort,
		pushMsgChan: make(chan *PushMessage, 128),
		procMsg:     &CoMsgProcess{},
	}

	svr.Init()
	svr.procMsg.Init()
	svr.procMsg.Start()

	svr.clients = &CoClients{}
	svr.clients.Init()
	svr.clients.Start()

	return svr
}

func (c *CoMsgServer) Init() {

}

func (c *CoMsgServer) Start() {
	if c.isLoop == 1 {
		return
	}
	c.isLoop = 1
	go c.run()
}

func (c *CoMsgServer) PublishMsg(topic string, body string, qos int) {
	//直接发送消息
	msg := &PushMessage{
		Type:  int(PublishMsg),
		Body:  body,
		Qos:   qos,
		Topic: topic,
	}
	c.pushMsgChan <- msg
}

func (c *CoMsgServer) Stop() {
	if c.isLoop == 0 {
		return
	}

	if c.procMsg != nil {
		c.procMsg.Stop()
	}

	if c.clients != nil {
		c.clients.Stop()
	}

	c.isLoop = 0
	c.job.Wait()
}

func (c *CoMsgServer) run() {

	bInitInlineClient := true

	c.job.Add(1)

	//cert, err := tls.X509KeyPair(testCertificate, testPrivateKey)
	cert, err := tls.LoadX509KeyPair("tool/test.pem", "tool/test.key")
	if err != nil {
		log.Error("LoadX509KeyPair failed: %s", err)
	}

	// Basic TLS Config
	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
	}

	server := mqtt.New(&mqtt.Options{
		InlineClient: true,
	})

	// Allow all connections.
	//_ = server.AddHook(new(auth.AllowHook), nil)

	// 在标准端口上创建TCP侦听器。
	tcp := listeners.NewTCP("tcp", fmt.Sprintf(":%d", c.tcpPort), nil)
	err = server.AddListener(tcp)
	if err != nil {
		log.Error("AddListener failed: %s", err)
	}

	tls := listeners.NewTCP("tls", fmt.Sprintf(":%d", c.tlsPort), &listeners.Config{
		TLSConfig: tlsConfig,
	})
	err = server.AddListener(tls)
	if err != nil {
		log.Error("AddListener failed: %s", err)
	}

	//在标准端口上创建WS侦听器
	ws := listeners.NewWebsocket("ws", fmt.Sprintf(":%d", c.wsPort), nil)
	err = server.AddListener(ws)
	if err != nil {
		log.Error("AddListener failed: %s", err)
	}

	err = server.AddHook(new(CoMsgHook), map[string]any{})
	if err != nil {
		log.Error("AddHook failed: %s", err)
	}

	c.job.Add(1)
	go func() {
		err := server.Serve()
		if err != nil {
			log.Error("Serve failed: %s", err)
		}
		c.job.Done()
	}()

	go func() {
		//开启定时器
		for {
			if c.isLoop == 0 {
				break
			}

			time.Sleep(time.Millisecond * 500)
			msg := &PushMessage{
				Type: int(TimerEvent),
			}
			c.pushMsgChan <- msg
		}
	}()

	//利用定时控制循环延迟
	for {
		if c.isLoop == 0 {
			break
		}

		select {
		case msg := <-c.pushMsgChan:
			if msg.Type == int(PublishMsg) {
				server.Publish(msg.Topic, []byte(msg.Body), false, byte(msg.Qos))
			} else {
				//log.Debug("publish type: %d", msg.Type)

			}
		}

		if bInitInlineClient {
			callbackFn := func(cl *mqtt.Client, sub packets.Subscription, pk packets.Packet) {
				//server.Log.Info("inline client received message from subscription", "client", cl.ID, "subscriptionId", sub.Identifier, "topic", pk.TopicName, "payload", string(pk.Payload))
				c.procMsg.PushRecvMsg(cl.ID, pk.TopicName, pk.Payload)
			}

			//订阅所有的管理消息
			server.Subscribe("codesk/admin/#", 1, callbackFn)
			server.Subscribe("codesk/report/#", 2, callbackFn)

			bInitInlineClient = false
		}
	}

	// Cleanup
	server.Close()

	c.job.Done()
}
