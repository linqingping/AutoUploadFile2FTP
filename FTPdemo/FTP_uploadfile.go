package main

import (
	"fmt"
	"github.com/dutchcoders/goftp"
	"os"
)

const (
	filename ="./tmp/test.txt"
	server = "192.168.122.151"
	port = ":21"
)
func main()  {
	var err error
	var ftp *goftp.FTP
	//连接FTP
	if ftp,err=goftp.Connect(server+port);err!=nil{
		panic(err)
	}
	defer ftp.Close()
	fmt.Println("Successfully connected !!")
	//匿名登录FTP
	if err = ftp.Login("anonymous","anonymous");err!=nil{
		panic(err)
	}
	//切换FTP到指定目录
	/*------------------------------------
	if err=ftp.Cwd("/Documents");err!=nil{
	   panic(err)
	}
	--------------------------------------*/
	//显示FTP当前目录
	/*-------------------------------------
	var curpath string
	if curpath, err = ftp.Pwd(); err != nil {
	    panic(err)
	}
	fmt.Println(curpath)
	---------------------------------------*/
	//显示FTP目录下文件列表
	/*--------------------------------------
	var files []string
	if files, err = ftp.List("./"); err != nil {
	   panic(err)
	}
	fmt.Println("Directory listing:/n", files)
	------------------------------------------*/
  	//上传文件到FTP
	var file *os.File
	if file, err = os.Open(filename); err != nil {
		panic(err)
	}
	if err := ftp.Stor(filename, file); err != nil {
		panic(err)
	}
	//result:=checkFTP(server,port)
	//fmt.Print(result)
}

//func checkFTP(server,port string)bool{
//	//连接FTP
//	if _,err:=goftp.Connect(server+port);err!=nil{
//		    fmt.Println(err)
//			return false
//	}
//	return true
//}
