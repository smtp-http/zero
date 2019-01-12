package zero

import (
	"fmt"
	"net"
	"testing"
	"time"
	"container/list"
)

var SessionList *list.List

func TestService(t *testing.T) {
	host := "127.0.0.1:18787"

	SessionList = list.New()

	fmt.Println(SessionList)

	ss, err := NewSocketService(host)
	if err != nil {
		return
	}

	

	// ss.SetHeartBeat(5*time.Second, 30*time.Second)

	ss.RegMessageHandler(HandleMessage)
	ss.RegConnectHandler(HandleConnect)
	ss.RegDisconnectHandler(HandleDisconnect)

	go NewClientConnect()

	timer := time.NewTimer(time.Second * 1000)
	go func() {
		<-timer.C
		ss.Stop("stop service")
		t.Log("service stoped")
	}()

	t.Log("service running on " + host)
	ss.Serv()

}

func HandleMessage(s *Session, msg *Message) {
	fmt.Println("receive msgID:", msg)
	fmt.Println("receive data:", string(msg.GetData()))
	for e := SessionList.Front(); e != nil; e = e.Next() {
		//fmt.Print(e.Value) //输出list的值,01234
		session := e.Value.(*Session)
		fmt.Println("++++++++++++++++++++++ 2")
		fmt.Println(session)
		fmt.Println("---------------------- 2")
		conn := session.GetConn()
		str := "shit ===== shit ++++"
		var data []byte = []byte(str)
		msg := NewMessage(12,data)
		conn.SendMessage(msg)
	}
}

func HandleDisconnect(s *Session, err error) {
	fmt.Println(s.GetConn().GetName() + " lost.")
}

func HandleConnect(s *Session) {
	fmt.Println("++++++++++++++++++++++ 1")
	fmt.Println(s)
	fmt.Println("---------------------- 1")
	fmt.Println(s.GetConn().GetName() + " connected." )
	SessionList.PushBack(s)
}

func NewClientConnect() {
	host := "127.0.0.1:18787"
	tcpAddr, err := net.ResolveTCPAddr("tcp", host)
	if err != nil {
		return
	}

	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		return
	}

	msg := NewMessage(1, []byte("Hello Zero!"))
	data, err := Encode(msg)
	if err != nil {
		return
	}
	conn.Write(data)

	

}
