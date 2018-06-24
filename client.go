package main
import (
	"flag"
	"fmt"
	"net"
	"os"
	"time"
)

const BUFFER_SIZE = 256
var host = flag.String("host", "localhost", "host")
var port = flag.String("port", "37", "port")
//go run client.go -host time.nist.gov

func main() {
	flag.Parse()
	addr, err := net.ResolveUDPAddr("udp", *host+":"+*port)
	if err != nil {
		fmt.Println("Can't resolve address: ", err)
		os.Exit(1)
	}
	conn, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		fmt.Println("Can't dial: ", err)
		os.Exit(1)
	}
	defer conn.Close()
	_, err = conn.Write([]byte("hello salamander"))
	if err != nil {
		fmt.Println("failed:", err)
		os.Exit(1)
	}
	data := make([]byte, BUFFER_SIZE)
	_, err = conn.Read(data)
	if err != nil {
		fmt.Println("failed to read UDP msg because of ", err)
		os.Exit(1)
	}
	fmt.Println(string(data))
	fmt.Println(time.Now().Format("2006-01-02 15:04:05"))
	os.Exit(0)
}