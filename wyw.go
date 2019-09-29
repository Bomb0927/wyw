package main

import (
	"github.com/gorilla/mux"
	"bufio"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
	"os/exec"
	"time"
)

//设置资产的结构
type Asset struct {
	ID     string   `json:"id,omitempty"`				//ID
	Name string   `json:"name,omitempty"`				//名字
	Value  string   `json:"value,omitempty"`			//价值
	Longitude string `json:"longitude,omitempty"`		//横坐标
	Latitude string `json:"latitude,omitempty"`			//纵坐标
	Priority string `json:"priority,omitempty"`			//优先级
}

//客户端选择
func select_client(){
	for true {
		fmt.Println("please input the client(web portal/client/command line/read csv)(退出为exit)")
		inputReader := bufio.NewReader(os.Stdin)
		COMMAND, err := inputReader.ReadString('\n')							 //输入客户端想要的形式：portal web/client/command line
		if err == nil {
			if COMMAND == "exit\n" {
				break
			}

			//web portal 的形式
			if COMMAND == "web portal\n" {
				router := mux.NewRouter()
				router.HandleFunc("/Assets", getassets).Methods("GET") //如果是get型method且路径为/Assets，下面几个命令同理
				router.HandleFunc("/Assets/{id}", getasset).Methods("GET")
				// Post handle function
				router.HandleFunc("/Assets/{id}", postasset).Methods("POST")
				// Delete handle function:
				router.HandleFunc("/Assets/{id}", deleteasset).Methods("DELETE")
				// 启动 API端口0618
				log.Fatal(http.ListenAndServe(":0618", router))
			}

			//client的形式
			if COMMAND == "client\n" {
				client()
			}

			//command line的形式
			if COMMAND == "command line\n" {
				checkexe2()
			}
			if COMMAND == "read csv\n" {
				csv_read()
			}
		}
	}
}

//选择所需要的client服务
func client(){
	var command string
	var temp Asset
	for true {
		fmt.Println("please input the method(get/post/delete)(退出为exit)")
		fmt.Scan(&command)
		//退出程序
		if command == "exit" {
			return
		}

		//获得资产信息
		if command == "get" {
			for true {
				fmt.Println("please input the ID of zhe asset which you want to get(全部为all,退出为exit)")
				fmt.Scan(&temp.ID) //输入id/all
				if temp.ID == "exit" {
					break
				}
				if temp.ID == "all" {
					get_all()
				} else{
					get_id(temp.ID)
				}
			}
		}

		//插入
		if command == "post" {
			for true {
				fmt.Println("please input the ID of the asset which you want to post(退出为exit)")

				fmt.Scan(&temp.ID)						//输入要插入的资产的id
				if temp.ID== "exit" {
					break
				}
				post(temp)
			}
		}

		//删除
		if command == "delete" {
			for true {
				fmt.Println("please input the ID of the asset which you want to delete(退出为exit)")
				fmt.Scan(&temp.ID)
				if temp.ID == "exit" {
					break
				}
				delete(temp.ID)
			}
		}
	}
}

//获得所有资产的信息
func get_all(){
	url := "mongodb://localhost"
	session, err := mgo.Dial(url)
	if err != nil {
		panic(err)
	}
	defer session.Close()
	session.SetMode(mgo.Monotonic, true)
	c := session.DB("test").C("Assets")
	var assets []Asset
	c.Find(nil).All(&assets)
	fmt.Println(assets)
}

//获得指定id的资产的资产信息
func get_id(id string) {
	url := "mongodb://localhost"
	session, err := mgo.Dial(url)
	if err != nil {
		panic(err)
	}
	defer session.Close()
	session.SetMode(mgo.Monotonic, true)
	c := session.DB("test").C("Assets")
	var assets []Asset
	c.Find(nil).All(&assets)
	flag := 0
	for _, item := range assets {
		if item.ID == id {
			flag = 1
			var ass []Asset
			err = c.Find(bson.M{"id":id}).All(&ass)
			fmt.Println(ass)
		}
	}
	if flag == 0 {
		fmt.Println("This ID does not exist!")
	}
	return
}

//插入资产信息
func post(temp Asset){
	url := "mongodb://localhost"
	session, err := mgo.Dial(url)
	if err != nil {
		panic(err)
	}
	defer session.Close()
	session.SetMode(mgo.Monotonic, true)
	c := session.DB("test").C("Assets")
	var assets []Asset
	flag := 0
	c.Find(nil).All(&assets)

	for _, item := range assets {			//检查是否有重复的id
		if item.ID == temp.ID {
			flag = 1
			fmt.Println("this ID is repeated!")
		}
	}
	if flag == 1 {
		return
	}

	//如果没有重复id，则开始输入其他值
	fmt.Println("please input the name of it")
	fmt.Scan(&temp.Name)
	fmt.Println("please input the value of it")
	fmt.Scan(&temp.Value)
	fmt.Println("please enter the longitude of it")
	fmt.Scan(&temp.Longitude)
	fmt.Println("please enter the latitude of it")
	fmt.Scan( &temp.Latitude)
	fmt.Println("please enter the priority of it")
	fmt.Scan( &temp.Priority)
	c.Insert(temp)
}

//删除资产信息
func delete(delete_id string){
	url := "mongodb://localhost"
	session, err := mgo.Dial(url)
	if err != nil {
		panic(err)
	}
	defer session.Close()
	session.SetMode(mgo.Monotonic, true)
	c := session.DB("test").C("Assets")
	var assets []Asset
	c.Find(nil).All(&assets)
	fmt.Println("this asset has been delete")
	_, err = c.RemoveAll(bson.M{"id": delete_id})
}

//打开命令行
func checkexe2() {
	checkEXE2("C:/go/src/project/mongo.exe")			//读取mongodb的命令行输入exe文件的位置
	time.Sleep(100 * time.Millisecond)       					//为了输出美观加了个停顿
	fmt.Println("退出时输入exit")
	var flag string
	for true {
		fmt.Scan(&flag)
		if flag == "exit" {
			return
		}
	}
}

func checkEXE2(exeAdress string) {									//启动命令行的函数
	cmd := exec.Command("cmd.exe", "/c", "start "+exeAdress)				//命令行程序地址
	err := cmd.Run()
	if err != nil {
		log.Println("启动失败:", err)
	} else {
		log.Println("启动成功!")
	}
}

//读取csv文件
func csv_read(){
	for true {
		fmt.Println("please enter the address of the csv file(退出为exit)")
		inputReader := bufio.NewReader(os.Stdin)
		way, err := inputReader.ReadString('\n')							 //输入客户端想要的形式：portal web/client/command line
		if err != nil {
			fmt.Println("!!")
		}
		if way == "exit\n" {
			break
		}
		file, err := os.Open(way[:len(way)-1])				//打开文件的位置
		if err != nil {							//检查地址是否有效，无效则退出
			fmt.Println("Error:", err)
			continue
		}
		defer file.Close()
		reader := csv.NewReader(file)
		url := "mongodb://localhost"
		session, err := mgo.Dial(url)
		if err != nil {
			panic(err)
		}
		defer session.Close()
		session.SetMode(mgo.Monotonic, true)
		c := session.DB("test").C("Assets")
		flag_read := 1							//设置一个变量，检查该文件是否被读取正常
		for {
			record, err := reader.Read()
			if err == io.EOF {
				break
			} else if err != nil {
				flag_read=0
				fmt.Println("记录集错误:", err)
				break
			}
			var temp Asset
			var assets []Asset
			c.Find(nil).All(&assets)
			flag_repeated := 0					//设置一个变量，检查该资产ID是否在数据库中存在
			for _, ass := range assets {
				if ass.ID == record[0] {
					flag_repeated = 1
				}
			}
			if flag_repeated == 1 {
				fmt.Println("the id \"" + record[0] + "\" is repeated")		//如果该资产ID已存在，则不存储并报错
			} else {																//如果该资产ID不存在，将该资产信息存入数据库
				temp.ID = record[0]
				temp.Name = record[1]
				temp.Value = record[2]
				temp.Longitude = record[3]
				temp.Latitude = record[4]
				temp.Priority = record[5]
				c.Insert(temp)
			}
		}
		if flag_read == 1{															//如果该文件读取正常，显示文件读取完毕
			fmt.Println("this csv file has being read")
		}
	}
}

//获得数据库所有资产的信息
func getassets(w http.ResponseWriter, req *http.Request) {
	url:="mongodb://localhost"								//打开数据库
	session,err:=mgo.Dial(url)
	if err!=nil{
		panic(err)
	}
	defer session.Close()
	session.SetMode(mgo.Monotonic,true)
	c:=session.DB("test").C("Assets")
	var assets []Asset
	c.Find(nil).All(&assets)							//获得所有的资产的信息
	json.NewEncoder(w).Encode(assets)						//将所有的资产信息转码为json数据在网页中输出
}

// 根据id获取对应固定资产
func getasset(w http.ResponseWriter, req *http.Request) {
	url:="mongodb://localhost"
	session,err:=mgo.Dial(url)
	if err!=nil{
		panic(err)
	}
	defer session.Close()
	session.SetMode(mgo.Monotonic,true)
	c:=session.DB("test").C("Assets")
	params := mux.Vars(req)									//获得所请求的id
	var assets []Asset
	c.Find(nil).All(&assets)							//获得全部资产

	for _, item := range assets {							//遍历全部资产 找到对应id的资产的信息并将其输出
		if item.ID == params["id"] {
			var ass []Asset
			err=c.Find(bson.M{"id":"0001"}).All(&ass)
			json.NewEncoder(w).Encode(item)
			return
		}
	}
	json.NewEncoder(w).Encode(assets)
}

// 向数据库中添加固定资产
func postasset(w http.ResponseWriter, req *http.Request) {
	url:="mongodb://localhost"
	session,err:=mgo.Dial(url)
	if err!=nil{
		panic(err)
	}
	defer session.Close()
	session.SetMode(mgo.Monotonic,true)
	c:=session.DB("test").C("Assets")
	var assets []Asset
	c.Find(nil).All(&assets)
	params := mux.Vars(req)
	for _, item := range assets {								//检查数据库中是否已经有该id的资产
		if item.ID == params["id"] {
			json.NewEncoder(w).Encode("This id is repeated!")
			return
		}
	}
	var Asset_1 Asset
	_ = json.NewDecoder(req.Body).Decode(&Asset_1)				//将网页中的json数据解码以便存入数据库
	Asset_1.ID = params["id"]
	c.Insert(Asset_1)
	c.Find(nil).All(&assets)
	json.NewEncoder(w).Encode(assets)
}


// 根据id进行删除操作
func deleteasset(w http.ResponseWriter, req *http.Request) {
	url:="mongodb://localhost"
	session,err:=mgo.Dial(url)
	if err!=nil{
		panic(err)
	}
	defer session.Close()
	//打开默认数据库test
	session.SetMode(mgo.Monotonic,true)
	c:=session.DB("test").C("Assets")
	params := mux.Vars(req)
	_, err = c.RemoveAll(bson.M{"id":params["id"]})					//从数据库中删除这个id的资产
	var assets []Asset
	c.Find(nil).All(&assets)
	json.NewEncoder(w).Encode(assets)
}

func random_assets() {
	url := "mongodb://localhost"
	session, err := mgo.Dial(url)
	if err != nil {
		panic(err)
	}
	defer session.Close()
	session.SetMode(mgo.Monotonic, true)
	c := session.DB("test").C("Assets")
	var assets,assets_A,assets_B,assets_C,assets_D,asset_A,asset_B,asset_C,asset_D []Asset
	c.Find(nil).All(&assets)
	c.Find(bson.M{"priority":"A"}).All(&assets_A)
	c.Find(bson.M{"priority":"B"}).All(&assets_B)
	c.Find(bson.M{"priority":"C"}).All(&assets_C)
	c.Find(bson.M{"priority":"D"}).All(&assets_D)
	fmt.Println("please input the the percent")
	var percentage int
	fmt.Scan(&percentage)

	number_A := len(assets_A)*percentage/100
	for i:=0;i<number_A;i++{
		temp:=rand.Intn(len(assets_A))
		asset_A=append(asset_A,assets_A[temp])
		assets_A = append(assets_A[:temp], assets_A[temp+1:]...)
	}

	number_B := len(assets_B)*percentage/100
	for i:=0;i<number_B;i++{
		temp:=rand.Intn(len(assets_B))
		asset_B=append(asset_B,assets_B[temp])
		assets_B = append(assets_B[:temp], assets_B[temp+1:]...)
	}

	number_C := len(assets_C)*percentage/100
	for i:=0;i<number_C;i++{
		temp:=rand.Intn(len(assets_C))
		asset_C=append(asset_C,assets_C[temp])
		assets_C= append(assets_C[:temp], assets_C[temp+1:]...)
	}

	number_D := len(assets_D)*percentage/100
	for i:=0;i<number_D;i++{
		temp:=rand.Intn(len(assets_D))
		asset_D=append(asset_D,assets_D[temp])
		assets_D = append(assets_D[:temp], assets_D[temp+1:]...)
	}

	fmt.Println("Priority A:",asset_A)
	fmt.Println("Priority B:",asset_B)
	fmt.Println("Priority C:",asset_C)
	fmt.Println("Priority D:",asset_D)
}

func main() {
	for true {
		// 输入需要选择的功能（生成随机资产组/查看更改数据库内容）

		fmt.Println("please choose the features which do you want to use")
		fmt.Println("select_client please enter 1, get random assets please enter 2, exit please enter 0")
		var choice int
		fmt.Scan(&choice)
		if choice == 0{
			break
		}else  if choice == 1{
			select_client()			//进行数据库查看修改功能的客户端选择
		}else if choice == 2{
			random_assets()			//生成随机资产组
		}
	}
}