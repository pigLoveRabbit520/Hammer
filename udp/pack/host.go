package pack

import (
	"bytes"
	"encoding/gob"
)

type Hosts struct {
	Hosts []string
}

//添加主机到列表
func (this *Hosts) Add(host string) {
	this.Hosts = append(this.Hosts, host)
}

func (this *Hosts) Encode() []byte {
	buf := bytes.NewBuffer([]byte{})
	ge := gob.NewEncoder(buf)
	err := ge.Encode(this)
	if err != nil {
		panic(err)
	}
	return buf.Bytes()
}

func (this *Hosts) Decode(bs []byte) bool {
	b := bytes.NewBuffer(bs)
	gd := gob.NewDecoder(b)

	err := gd.Decode(this)

	return err == nil
}

func (this *Hosts) Elements() []string {
	return this.Hosts
}
