package main

import (
	"fmt"
	"net"
	"time"
)

func main(){
	//用于检查服务的端口是否可用
	timeout := time.Duration(5 * time.Second)
	t1 := time.Now()
	_, err := net.DialTimeout("tcp","www.baidu.com:443", timeout)
	fmt.Println("waist time :", time.Now().Sub(t1))
	if err != nil {
		fmt.Println("Site unreachable, error: ", err)
		return
	}
	fmt.Println("tcp server is ok")
}
