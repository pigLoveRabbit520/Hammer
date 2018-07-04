package utils

import (
	"bytes"
	"fmt"
	"sync"
)

type Set struct {
	mutex sync.Mutex
	m     map[interface{}]bool
}

func NewSet() *Set {

	return &Set{mutex: sync.Mutex{}, m: make(map[interface{}]bool)}
}

//添加    true 添加成功 false 添加失败
func (set *Set) Add(e interface{}) (b bool) {

	set.mutex.Lock()
	defer set.mutex.Unlock()

	if !set.m[e] {
		set.m[e] = true
		return true
	}
	return false
}

//删除
func (set *Set) Remove(e interface{}) {
	set.mutex.Lock()
	defer set.mutex.Unlock()
	delete(set.m, e)
}

//清除
func (set *Set) Clear() {
	set.mutex.Lock()
	defer set.mutex.Unlock()
	set.m = make(map[interface{}]bool)
}

//是否包含
func (set *Set) Contains(e interface{}) bool {
	set.mutex.Lock()
	defer set.mutex.Unlock()
	return set.m[e]
}

//获取元素数量
func (set *Set) Len() int {
	set.mutex.Lock()
	defer set.mutex.Unlock()
	return len(set.m)
}

//判断两个set时候相同
//true 相同 false 不相同
func (set *Set) Same(other *Set) bool {
	set.mutex.Lock()
	defer set.mutex.Unlock()
	if other == nil {
		return false
	}

	if set.Len() != other.Len() {
		return false
	}

	for k, _ := range set.m {
		if !other.Contains(k) {
			return false
		}
	}

	return true
}

//迭代
func (set *Set) Elements() []interface{} {
	set.mutex.Lock()
	defer set.mutex.Unlock()
	initlen := len(set.m)

	snaphot := make([]interface{}, initlen)

	actuallen := 0

	for k, _ := range set.m {
		if actuallen < initlen {
			snaphot[actuallen] = k
		} else {
			snaphot = append(snaphot, k)
		}
		actuallen++
	}

	if actuallen < initlen {
		snaphot = snaphot[:actuallen]
	}

	return snaphot
}

//获取自身字符串
func (set *Set) String() string {
	set.mutex.Lock()
	defer set.mutex.Unlock()

	var buf bytes.Buffer

	buf.WriteString("set{")

	first := true

	for k, _ := range set.m {
		if first {
			first = false
		} else {
			buf.WriteString(" ")
		}

		buf.WriteString(fmt.Sprintf("%v", k))
	}

	buf.WriteString("}")

	return buf.String()
}
