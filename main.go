package main

import (
	"bufio"
	"fmt"
	"github.com/dutchcoders/goftp"
	"github.com/fsnotify/fsnotify"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path"
	"strings"
	"time"
)

func main() {
	configMap := InitConfig("./configuration")
	fmt.Println(configMap)
	var  endflag = configMap["endflag"]
	var server = configMap["server"]
	var port = configMap["port"]
	//创建一个监控对象
	watch, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watch.Close()
	//添加要监控的对象，文件或文件夹
	err = watch.Add("./FTPdemo")
	if err != nil {
		log.Fatal(err)
	}
	//我们另启一个goroutine来处理监控对象的事件
	go func() {
	   for {
		select {
		case ev := <-watch.Events:
		   {
		      //判断事件发生的类型，如下5种
		      // Create 创建
		      // Write 写入
		      // Remove 删除
		      // Rename 重命名
		      // Chmod 修改权限
			if ev.Op&fsnotify.Create == fsnotify.Create {
				log.Println("创建文件 : ", ev.Name)
				go func() {
					filecontext, _ := ioutil.ReadFile(ev.Name)
					//fmt.Println(len(filecontext))
					if endflag == string(filecontext[len(filecontext)-5:len(filecontext)-1]) {
						fmt.Println(string(filecontext))
						for i:=0;i<10;i++{
							result:=checkFTP(server,port)
							if result{
								uploadFile(server,port,ev.Name)
								break
							}
							time.Sleep(10*time.Second)
						}
						return
					}
				}()
			}
			if ev.Op&fsnotify.Write == fsnotify.Write {
			        log.Println("写入文件 : ", ev.Name)
			}
			if ev.Op&fsnotify.Remove == fsnotify.Remove {
				log.Println("删除文件 : ", ev.Name)
			}
			if ev.Op&fsnotify.Rename == fsnotify.Rename {
				log.Println("重命名文件 : ", ev.Name)
			}
			if ev.Op&fsnotify.Chmod == fsnotify.Chmod {
				log.Println("修改权限 : ", ev.Name)
			}
		     }
		  case err := <-watch.Errors:
		     {
			log.Println("error : ", err)
			return
		     }
		  }
	    }
	}()
	// 阻塞:阻塞的方法是用channel，当然都是可以的，不过用select{}更加简洁
	select {}
}
func InitConfig(path string) map[string]string {
	//初始化
	myMap := make(map[string]string)

	//打开文件指定目录，返回一个文件f和错误信息
	f, err := os.Open(path)
	defer f.Close()
	//异常处理 以及确保函数结尾关闭文件流
	if err != nil {
		panic(err)
	}
	//创建一个输出流向该文件的缓冲流*Reader
	r := bufio.NewReader(f)
	for {
		//读取，返回[]byte 单行切片给b
		b, _, err := r.ReadLine()
		if err != nil {
			if err == io.EOF {
				break
			}
			panic(err)
		}
		//去除单行属性两端的空格
		s := strings.TrimSpace(string(b))
		//fmt.Println(s)

		//判断等号=在该行的位置
		index := strings.Index(s, "=")
		if index < 0 {
			continue
		}
		//取得等号左边的key值，判断是否为空
		key := strings.TrimSpace(s[:index])
		if len(key) == 0 {
			continue
		}
		//取得等号右边的value值，判断是否为空
		value := strings.TrimSpace(s[index+1:])
		if len(value) == 0 {
			continue
		}
		//这样就成功吧配置文件里的属性key=value对，成功载入到内存中c对象里
		myMap[key] = value
	}
	return myMap
}

//检查FTP是否可连接
func checkFTP(server,port string)bool{
	//连接FTP
	if _,err:=goftp.Connect(server+":"+port);err!=nil{
		fmt.Println(err)
		return false
	}
	return true
}
//上传文件到FTP
func uploadFile(server,port,filename string)bool {
	var err error
	var ftp *goftp.FTP
	//连接FTP
	if ftp,err=goftp.Connect(server+":"+port);err!=nil{
		fmt.Println(err)
		return false
	}
	defer ftp.Close()
	//匿名登录FTP
	if err = ftp.Login("anonymous","anonymous");err!=nil{
		fmt.Println(err)
		return false
	}
	//上传文件到FTP
	var file *os.File
	if file, err = os.Open(filename); err != nil {
		fmt.Println(err)
		return false
	}
	if err := ftp.Stor(filename, file); err != nil {
		fmt.Println(err)
		return false
	}
    return true
}
//获取文件名称
func getfilenameonly(fullFilename string)string{
	var filenameWithSuffix string
	filenameWithSuffix = path.Base(fullFilename) //获取文件名带后缀
	//fmt.Println("filenameWithSuffix =", filenameWithSuffix)
	var fileSuffix string
	fileSuffix = path.Ext(filenameWithSuffix) //获取文件后缀
	//fmt.Println("fileSuffix =", fileSuffix)
	var filenameOnly string
	filenameOnly = strings.TrimSuffix(filenameWithSuffix, fileSuffix)//获取文件名
	return filenameOnly
}
