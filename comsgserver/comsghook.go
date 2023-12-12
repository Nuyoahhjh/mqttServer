package msgserver

import (
	"bytes"
	"codeskserver/log"
	"fmt"

	mqtt "github.com/mochi-mqtt/server/v2"
	"github.com/mochi-mqtt/server/v2/packets"
)

type CoMsgHook struct {
	mqtt.HookBase
}

func (h *CoMsgHook) ID() string {
	return "codesk"
}

func (h *CoMsgHook) Provides(b byte) bool {
	return bytes.Contains([]byte{
		mqtt.OnConnectAuthenticate, //自己来鉴权
		mqtt.OnACLCheck,            //自己来鉴权
		mqtt.OnConnect,
		mqtt.OnDisconnect,
		mqtt.OnSubscribed,
		mqtt.OnUnsubscribed,
		mqtt.OnPublished,
		mqtt.OnPublish,
		mqtt.OnSessionEstablished,
	}, []byte{b})
}

func (h *CoMsgHook) Init(config any) error {
	log.Print("CoMsgHook initialised")
	return nil
}

func (h *CoMsgHook) OnConnect(cl *mqtt.Client, pk packets.Packet) error {
	log.Print("client connected lient %s from %s", cl.ID, cl.Net.Remote)
	//这个地方有新设备连上了,但是还没有经过授权同意，所以不能再这个地方添加设备

	return nil
}

func (h *CoMsgHook) OnDisconnect(cl *mqtt.Client, err error, expire bool) {
	if err != nil {
		log.Print("client disconnected + client %s expire %t, erro:%s", cl.ID, expire, err.Error())
	} else {
		log.Print("client disconnected + client %s expire %t", cl.ID, expire)
	}

}

func (h *CoMsgHook) OnSubscribed(cl *mqtt.Client, pk packets.Packet, reasonCodes []byte) {
	log.Print(fmt.Sprintf("subscribed qos=%v", reasonCodes), "client", cl.ID, "filters", pk.Filters)
}

func (h *CoMsgHook) OnUnsubscribed(cl *mqtt.Client, pk packets.Packet) {
	//log.Print("unsubscribed", "client", cl.ID, "filters", pk.Filters)
}

func (h *CoMsgHook) OnPublish(cl *mqtt.Client, pk packets.Packet) (packets.Packet, error) {
	//log.Print("received from client client %s payload %s", cl.ID, string(pk.Payload))

	/*pkx := pk
	if string(pk.Payload) == "hello" {
		pkx.Payload = []byte("hello world")
		log.Print("received modified packet from client", "client", cl.ID, "payload", string(pkx.Payload))
	}*/

	return pk, nil
}

func (h *CoMsgHook) OnPublished(cl *mqtt.Client, pk packets.Packet) {
	//log.Print("published to client client %s payload %s", cl.ID, string(pk.Payload))
}

func (h *CoMsgHook) OnConnectAuthenticate(cl *mqtt.Client, pk packets.Packet) bool {
	// 在这里进行用户名和密码的验证
	// 如果验证通过，返回 true；否则，返回 false
	// if _, ok := h.ledger.AuthOk(cl, pk); ok {
	// 	return true
	// }
	log.Print("OnConnectAuthenticate %s password:%s", string(pk.Connect.Username), string(pk.Connect.Password))

	if string(pk.Connect.Username) == "test" && string(pk.Connect.Password) == "123456" {
		return true
	}
	return false
	//h.Log.Info("client failed authentication check",
	//	"username", string(pk.Connect.Username),
	//	"remote", cl.Net.Remote)
}

func (h *CoMsgHook) OnACLCheck(cl *mqtt.Client, topic string, write bool) bool {
	// 在这里进行 ACL（Access Control List） 检查
	// 根据传入的 topic 和 write 参数，决定是否允许客户端执行相应的操作
	// 如果允许，返回 true；否则，返回 false
	//if _, ok := h.ledger.AuthOk(cl, pk); ok {
	//	return true
	//}
	//log.Print("client ACL check done")

	return true
}

func (h *CoMsgHook) OnSessionEstablished(cl *mqtt.Client, pk packets.Packet) {

	//这个地方添加设备

	log.Print("OnSessionEstablished  id:%s", cl.ID)
}
