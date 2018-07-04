// UDPServer project main.go
package main

import (
	"fmt"
	"net"
	"strconv"
	"time"
	"github.com/salamander-mh/Hammer/udp/pack"
	"github.com/salamander-mh/Hammer/udp/utils"
)

const Port = 60000

func checkError(err error) {
	if err != nil {
		panic(err.Error())
	}
}

//保存所有的客户端连接
var addrs *utils.Set

func main() {

	addrs = utils.NewSet()

	ServerAddr, err := net.ResolveUDPAddr("udp", ":"+strconv.Itoa(Port))
	checkError(err)

	conn, err := net.ListenUDP("udp", ServerAddr)
	checkError(err)
	defer conn.Close()

	fmt.Println("服务端启动成功，绑定端口:", Port)

	go func() {
		var buf [400]byte
		for {
			time.Sleep(time.Second)
			n, RemoteAddr, err := conn.ReadFromUDP(buf[:])
			checkError(err)

			p := pack.Pack{}

			if !p.Decode(buf[:n]) {
				fmt.Println("收到来自", RemoteAddr, "的无效的数据包,", err.Error())
				continue
			}

			if p.Type == pack.Join {

				//有新用户加入，向其他所有用户发送加入请求包，请求包包含了新用户ip：port信息，用于所有用户对他进行打洞
				if addrs.Add(RemoteAddr.String()) {

					fmt.Println("新用户加入", RemoteAddr, "已广播给所有用户")
					hosts := pack.Hosts{}
					for _, v := range addrs.Elements() {
						if v.(string) == RemoteAddr.String() {
							continue
						}

						tAddr, err := net.ResolveUDPAddr("udp", v.(string))
						pk := pack.Pack{pack.Join, []byte(RemoteAddr.String())}
						bs := pk.Encode()
						_, err = conn.WriteToUDP(bs, tAddr)
						if err != nil {
							fmt.Println("send error:", err.Error())
						}

						//添加到主机列表，后面要发送给当前新加入用户，用于打洞
						hosts.Add(v.(string))
						fmt.Println(v)

					}

					//新用户添加成功之后，把当前的所有用户信息都告诉【新用户】，用于一一打洞
					pk := pack.Pack{pack.JoinReturn, hosts.Encode()}
					bs := pk.Encode()
					conn.WriteToUDP(bs, RemoteAddr)

				}
			}
		}
	}()

	for {
		time.Sleep(time.Second)
	}

}
