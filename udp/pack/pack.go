// Pack project Pack.go
package pack

import (
	"bytes"
	"encoding/gob"
)

const (
	//新用户加入
	Join = iota
	//加入后接收到的返回包
	JoinReturn
	//打洞消息
	Hole

	//聊天内容
	Msg
)

type Pack struct {
	Type byte
	Data []byte
}

/**
序列化pack，返回序列化后的字节数组
**/
func (this *Pack) Encode() []byte {
	buf := bytes.NewBuffer([]byte{})
	ge := gob.NewEncoder(buf)
	err := ge.Encode(this)
	if err != nil {
		panic(err)
	}
	return buf.Bytes()

}

func (this *Pack) Decode(bs []byte) bool {
	b := bytes.NewBuffer(bs)
	gd := gob.NewDecoder(b)

	err := gd.Decode(this)

	return err == nil
}
