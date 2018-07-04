// UDPClient project main.go
package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"net"
	"strconv"
	"time"
	"github.com/salamander-mh/Hammer/udp/pack"
	"github.com/salamander-mh/Hammer/udp/utils"
)

const Server = "noxue.com:60000"

func checkError(err error) {
	if err != nil {
		panic(err.Error())
	}
}

var addrs *utils.Set

func main() {

	addrs = utils.NewSet()

	rand.Seed(time.Now().UnixNano())
	port := rand.Intn(1000)
	ClientAddr, err := net.ResolveUDPAddr("udp", ":"+strconv.Itoa(port+2000))
	checkError(err)

	conn, err := net.ListenUDP("udp", ClientAddr)
	checkError(err)
	defer conn.Close()

	fmt.Println("客户端启动成功，绑定端口：", port+2000)

	RemoteAddr, err := net.ResolveUDPAddr("udp", Server)
	checkError(err)
	//连接服务器，发送加入聊天的包
	pk := pack.Pack{pack.Join, []byte("hello")}
	bs := pk.Encode()

	_, err = conn.WriteToUDP(bs, RemoteAddr)
	checkError(err)

	//不断接收udp包
	go func() {
		var buf [400]byte

		for {
			time.Sleep(time.Second)
			n, RemoteAddr, err := conn.ReadFromUDP(buf[:])
			checkError(err)

			p := pack.Pack{}

			if !p.Decode(buf[:n]) {
				fmt.Println("收到来自", RemoteAddr, "的无效的数据包")
				continue
			}

			//如果有新人加入，就向他发个打洞的消息
			if p.Type == pack.Join {
				addr, err := net.ResolveUDPAddr("udp", string(p.Data[0:]))
				checkError(err)
				addrs.Add(string(p.Data[:]))
				pk := pack.Pack{pack.Hole, []byte("hello")}
				bs := pk.Encode()
				_, err = conn.WriteToUDP(bs, addr)
				checkError(err)
				fmt.Println("有新人加入，向新人", addr, "打了个洞")

			} else if p.Type == pack.JoinReturn { //加入后，会返回当前所有在线的主机，要一一对他们打洞

				hosts := pack.Hosts{}
				if !hosts.Decode(p.Data) {
					continue
				}

				for _, host := range hosts.Elements() {
					addr, err := net.ResolveUDPAddr("udp", host)
					if err != nil {
						continue
					}
					addrs.Add(host)
					pk = pack.Pack{pack.Hole, []byte("hello")}
					bs = pk.Encode()
					checkError(err)
					_, err = conn.WriteToUDP(bs, addr)
					checkError(err)
					fmt.Println("首次加入，向用户", addr, "打洞")
				}
			} else if p.Type == pack.Msg {
				fmt.Println(RemoteAddr, "说:", string(p.Data[:]))
			}
		}
	}()

	num := 1
	for {
		time.Sleep(time.Second * 3)

		buf := bytes.NewBuffer([]byte{})
		buf.WriteString("我是第 ")
		buf.WriteString(strconv.Itoa(num))
		buf.WriteString("条信息")

		pk := pack.Pack{pack.Msg, buf.Bytes()}

		for _, v := range addrs.Elements() {
			tAddr, err := net.ResolveUDPAddr("udp", v.(string))
			checkError(err)
			conn.WriteToUDP(pk.Encode(), tAddr)
		}
		num += 1
	}

	for {
		time.Sleep(time.Second)
	}

}