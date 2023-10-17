package user

import (
	"encoding/json"
	"fmt"
	userController "golang-study/controller"
	utils "golang-study/utils"
	"log"
	"net"
	"net/http"
	"strconv"
)

const port = ":8080"

func InitService() {
	http.HandleFunc("/", home)
	http.HandleFunc("/list", list)
	http.HandleFunc("/update", update)
	http.HandleFunc("/add", add)
	http.HandleFunc("/delete", delete)
	http.HandleFunc("/test", test)
	http.HandleFunc("/payApplyList", payApplyList)
	http.HandleFunc("/purchaseOrderList", purchaseOrderList)
	ips := GetIps()
	fmt.Println(ips[2])
	fmt.Printf("服务已启动,域名为 " + "\x1b]8;;" + "http://" + ips[2] + port + "\x1b\\" + "http://" + ips[2] + port + "\x1b]8;;\x1b\\\n")
	log.Fatal(http.ListenAndServe(port, nil))
}

// 初始页面
func home(w http.ResponseWriter, req *http.Request) {
	w.Write([]byte("服务已经启动成功>>>>>>"))
}

type listReq struct {
	Id       string `json:"id"`
	Pagesize int    `json:"pagesize"`
	Current  int    `json:"current"`
	UserName string `json:"username"`
}
type DataType struct {
	total int
	list  []userController.Users
}
type payReq struct {
	Pagesize        int `json:"limit"`
	Current         int `json:"page"`
	SearchCondition userController.SearchCondition
}

func payApplyList(w http.ResponseWriter, req *http.Request) {
	decoder := json.NewDecoder(req.Body)
	var data payReq
	err := decoder.Decode(&data)
	if err != nil {
		fmt.Println(err)
	}
	pagesize := data.Pagesize
	current := data.Current
	searchCondition := data.SearchCondition
	fmt.Println(pagesize, current)
	rows, total, err := userController.PayApplyOrderList(pagesize, current, searchCondition)
	// var data DataType
	datasouce := map[string]interface{}{
		"content": rows,
		"total":   total,
	}
	if err != nil {
		w.WriteHeader(500)
		utils.HandleResponse(w, 500, nil, "查询失败"+err.Error())
	} else {
		w.Header().Set("Content-Type", "application/json")
		utils.HandleResponse(w, 0, datasouce, "查询成功")
	}
}

func purchaseOrderList(w http.ResponseWriter, req *http.Request) {
	decoder := json.NewDecoder(req.Body)
	var data payReq
	err := decoder.Decode(&data)
	if err != nil {
		fmt.Println(err)
	}
	pagesize := data.Pagesize
	current := data.Current
	fmt.Println(pagesize, current)
	rows, total, err := userController.PurchaseOrderList(pagesize, current)
	// var data DataType
	datasouce := map[string]interface{}{
		"content": rows,
		"total":   total,
	}
	if err != nil {
		w.WriteHeader(500)
		utils.HandleResponse(w, 500, nil, "查询失败"+err.Error())
	} else {
		w.Header().Set("Content-Type", "application/json")
		utils.HandleResponse(w, 0, datasouce, "查询成功")
	}
}

func list(w http.ResponseWriter, req *http.Request) {
	decoder := json.NewDecoder(req.Body)
	var data listReq
	err := decoder.Decode(&data)
	if err != nil {
		fmt.Println(err)
	}
	id := data.Id
	pagesize := data.Pagesize
	current := data.Current
	username := data.UserName
	fmt.Println(id, pagesize, current, username)
	//查询所有
	if id == "" {
		rows, total, err := userController.SelectAllUser(pagesize, current, username)
		// var data DataType
		data := map[string]interface{}{
			"total": total,
			"list":  rows,
		}
		if err != nil {
			w.WriteHeader(500)
			utils.HandleResponse(w, 500, nil, "查询失败"+err.Error())
		} else {
			w.Header().Set("Content-Type", "application/json")
			utils.HandleResponse(w, 0, data, "查询成功")
		}
	} else {
		//查询单个
		int, _ := strconv.Atoi(id)
		row, err := userController.SelectUser(int)
		if err != nil {
			w.WriteHeader(500)
			utils.HandleResponse(w, 500, nil, "查询失败")
		} else {
			w.Header().Set("Content-Type", "application/json")
			utils.HandleResponse(w, 0, row, "查询成功")
		}
	}
}

type useReq struct {
	Username string `json:"Username"`
	Password string `json:"Password"`
	Sex      string `json:"Sex"`
	Email    string `json:"Email"`
	Money    int    `json:"Money"`
	UserId   int    `json:"id"`
}

func update(w http.ResponseWriter, req *http.Request) {
	decoder := json.NewDecoder(req.Body)
	var data useReq
	err := decoder.Decode(&data)
	if err != nil {
		fmt.Println(err)
	}
	username := data.Username
	password := data.Password
	sex := data.Sex
	email := data.Email
	money := data.Money
	id := data.UserId
	fmt.Println(username, password, sex, email, money, id)
	if username != "" && password != "" && sex != "" && email != "" && id != 0 {
		err := userController.UpdateUser(username, password, sex, email, money, id)
		if err != nil {
			w.WriteHeader(500)
			utils.HandleResponse(w, 500, nil, "修改失败"+err.Error())
		} else {
			utils.HandleResponse(w, 0, nil, "修改成功")
		}
	} else {
		w.WriteHeader(400)
		utils.HandleResponse(w, 500, nil, "参数有误")
	}
}

func add(w http.ResponseWriter, req *http.Request) {
	decoder := json.NewDecoder(req.Body)
	var data useReq
	err := decoder.Decode(&data)
	if err != nil {
		fmt.Println(err)
	}
	username := data.Username
	password := data.Password
	sex := data.Sex
	email := data.Email
	money := data.Money
	if username != "" && password != "" && sex != "" && email != "" && money != 0 {
		_, err := userController.AddUser(username, password, sex, email, money)
		if err != nil {
			w.WriteHeader(500)
			utils.HandleResponse(w, 500, nil, "添加失败"+err.Error())
		} else {
			utils.HandleResponse(w, 0, nil, "添加成功")
		}
	} else {
		w.WriteHeader(400)
		utils.HandleResponse(w, 500, nil, "参数有误")
	}
}

func delete(w http.ResponseWriter, req *http.Request) {
	id := req.URL.Query().Get("id")
	if id != "" {
		id, _ := strconv.Atoi(id)
		err := userController.DeleteUser(id)
		if err != nil {
			w.WriteHeader(500)
			utils.HandleResponse(w, 500, nil, "删除失败"+err.Error())
		} else {
			utils.HandleResponse(w, 0, nil, "删除成功")
		}
	} else {
		w.WriteHeader(400)
		utils.HandleResponse(w, 500, nil, "参数有误")
	}
}

func test(w http.ResponseWriter, req *http.Request) {
	req.ParseForm()
	for k, v := range req.URL.Query() {
		fmt.Print("get")
		fmt.Println("key:", k, ", value:", v[0])
	}

	for k, v := range req.PostForm {
		fmt.Print("post")
		fmt.Println("key:", k, ", value:", v[0])
	}

	fmt.Fprintln(w, "这是一个开始")
}

func GetIps() (ips []string) {
	interfaceAddr, err := net.InterfaceAddrs()
	if err != nil {
		fmt.Printf("fail to get net interfaces ipAddress: %v\n", err)
		return ips
	}

	for _, address := range interfaceAddr {
		ipNet, isVailIpNet := address.(*net.IPNet)
		// 检查ip地址判断是否回环地址
		if isVailIpNet && !ipNet.IP.IsLoopback() {
			if ipNet.IP.To4() != nil {
				ips = append(ips, ipNet.IP.String())
			}
		}
	}
	return ips
}
