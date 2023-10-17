package user

import (
	"encoding/json"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

type Users struct {
	UserId   int    `db:"user_id"`
	Username string `db:"username"`
	Password string `db:"password"`
	Sex      string `db:"sex"`
	Email    string `db:"email"`
	Money    int    `db:"money"`
}

var db *sqlx.DB

func init() {
	database, err := sqlx.Open("mysql", "root:123456@tcp(127.0.0.1:3306)/mytest")
	if err != nil {
		fmt.Println(err)
	}
	db = database
}

func AddUser(username, password, sex, email string, money int) (int64, error) {
	sql := "insert into user(username,password,sex, email, money)values (?,?,?,?,?)"
	r, err := db.Exec(sql, username, password, sex, email, money)
	if err != nil {
		fmt.Println("exec failed", err)
		return 0, err
	}
	id, err := r.LastInsertId()
	if err != nil {
		fmt.Println("get last insert id failed", err)
		return 0, err
	}
	return id, nil
}

func UpdateUser(username, password, sex, email string, money int, id int) error {
	sql := "update user set username=?, password=?, sex=?, email=?, money=? where user_id=?"
	_, err := db.Exec(sql, username, password, sex, email, money, id)
	if err != nil {
		fmt.Println("exec failed", err)
		return err
	}
	return nil
}

func SelectUser(id int) (*Users, error) {
	sql := "select user_id, username, sex, email from user where user_id=?"
	var user Users
	err := db.Get(&user, sql, id)
	if err != nil {
		fmt.Println("get failed", err)
		return nil, err
	}
	return &user, nil
}

type listReq struct {
	Id       string
	Pagesize int
	Current  int
}

// 查询所有
func SelectAllUser(pagesize int, current int, username string) ([]Users, int, error) {
	// var users Users
	offset := (current - 1) * pagesize
	// err := db.Select(&users, "SELECT * FROM table LIMIT ? OFFSET ?", pagesize, offset)
	// if err != nil {
	// 	if err != nil {
	// 		fmt.Println("get failed", err)
	// 		return nil, err
	// 	}
	// }
	// query := fmt.Sprintf("SELECT * FROM user  LIMIT %d OFFSET %d", pagesize, offset)
	con := make([]string, 0)
	if username != "" {
		con = append(con, fmt.Sprintf("username LIKE '%%%s%%'", username))
	}
	var users []Users
	var total int
	sql := "SELECT * FROM user WHERE username REGEXP `?` LIMIT ? OFFSET ?"
	err := db.Select(&users, username, sql, pagesize, offset)
	db.Get(&total, "select count(*) from user")
	if err != nil {
		fmt.Println("get failed", err)
		return nil, 0, err
	}
	return users, total, nil
}

type BtnsType struct {
	Look     bool
	Withdrew bool
}

type PayApplyOrder struct {
	Id             int    `db:"id"`
	PayApplyNo     string `db:"payApplyNo"`
	ApplyTime      string `db:"applyTime"`
	PurchaseNo     string `db:"purchaseNo"`
	SupplierNo     string `db:"supplierNo"`
	SupplierName   string `db:"supplierName"`
	ApplyPayAmount string `db:"applyPayAmount"`
	DiscountAmount string `db:"discountAmount"`
	CashAmount     string `db:"cashAmount"`
	RealPayAmount  string `db:"realPayAmount"`
	Applicant      string `db:"applicant"`
	ApplyDept      string `db:"applyDept"`
	Status         string `db:"status"`
	EbsSyncTime    string `db:"ebsSyncTime"`
	Btns           string `db:"btns"`
}

type PayApplyOrder1 struct {
	Id             int
	PayApplyNo     string
	ApplyTime      string
	PurchaseNo     string
	SupplierNo     string
	SupplierName   string
	ApplyPayAmount string
	DiscountAmount string
	CashAmount     string
	RealPayAmount  string
	Applicant      string
	ApplyDept      string
	Status         string
	EbsSyncTime    string
	Btns           BtnsType
}

type SearchCondition struct {
	Supplier   string `json:"supplier"`
	ApplyNo    string `json:"applyNo"`
	PurchaseNo string `json:"purchaseNo"`
	SupplierNo string `json:"supplierNo"`
	ApplyUser  string `json:"applyUser"`
	ApplyDept  string `json:"applyDept"`
	StartTime  string `json:"startTime"`
	EndTime    string `json:"endTime"`
}

func PayApplyOrderList(pagesize int, current int, searchCondition SearchCondition) ([]PayApplyOrder1, int, error) {
	offset := current * pagesize
	fmt.Println(pagesize, offset)
	var total int
	var users []PayApplyOrder
	var result []PayApplyOrder1
	payApplyNo := searchCondition.ApplyNo
	sqlMap := "where 1=1"
	if searchCondition.ApplyNo != "" {
		sqlMap = sqlMap + " and `payApplyNo` LIKE '%" + payApplyNo + "%'"
	}
	if searchCondition.Supplier != "" {
		sqlMap = sqlMap + " and `supplierNo` LIKE '%" + searchCondition.Supplier + "%'"
	}
	if searchCondition.PurchaseNo != "" {
		sqlMap = sqlMap + " and `purchaseNo` LIKE '%" + searchCondition.PurchaseNo + "%'"
	}
	if searchCondition.SupplierNo != "" {
		sqlMap = sqlMap + " and `supplierNo` LIKE '%" + searchCondition.SupplierNo + "%'"
	}
	if searchCondition.ApplyUser != "" {
		sqlMap = sqlMap + " and `applicant` LIKE '%" + searchCondition.ApplyUser + "%'"
	}
	if searchCondition.StartTime != "" {
		sqlMap = sqlMap + " and `applyTime` >= '" + searchCondition.StartTime + "'"
	}
	if searchCondition.EndTime != "" {
		sqlMap = sqlMap + " and `applyTime` <= '" + searchCondition.EndTime + "'"
	}
	sql := "select id,payApplyNo,applyTime,purchaseNo,supplierNo,supplierName,applyPayAmount,discountAmount,cashAmount,realPayAmount,applicant,applyDept,status,ebsSyncTime,btns from paymentapplyorder " + sqlMap + " LIMIT ? OFFSET ? "
	fmt.Println(sql)
	err := db.Select(&users, sql, pagesize, offset)
	for i := 0; i < len(users); i++ {
		var Btns BtnsType
		err := json.Unmarshal([]byte(users[i].Btns), &Btns)
		if err != nil {
			fmt.Println("unmarshal failed", err)
			return nil, 0, err
		}
		result = append(result, PayApplyOrder1{
			Id:             users[i].Id,
			PayApplyNo:     users[i].PayApplyNo,
			ApplyTime:      users[i].ApplyTime,
			PurchaseNo:     users[i].PurchaseNo,
			SupplierNo:     users[i].SupplierNo,
			SupplierName:   users[i].SupplierName,
			ApplyPayAmount: users[i].ApplyPayAmount,
			DiscountAmount: users[i].DiscountAmount,
			CashAmount:     users[i].CashAmount,
			RealPayAmount:  users[i].RealPayAmount,
			Applicant:      users[i].Applicant,
			ApplyDept:      users[i].ApplyDept,
			Status:         users[i].Status,
			EbsSyncTime:    users[i].EbsSyncTime,
			Btns:           Btns,
		})
	}
	db.Get(&total, "select count(*) from paymentapplyorder "+sqlMap)
	if err != nil {
		fmt.Println("get failed", err)
		return nil, 0, err
	}
	if result == nil {
		return []PayApplyOrder1{}, total, nil
	}
	return result, total, nil
}

type PurchaseOrderType struct {
	Id              int    `db:"id"`
	PurchaseNo      string `db:"purchaseNo"`
	PurchaseStatus  string `db:"purchaseStatus"`
	DespatchStatus  string `db:"despatchStatus"`
	PayStatus       string `db:"payStatus"`
	Version         string `db:"version"`
	BusiType        string `db:"busiType"`
	ChannelClass    string `db:"channelClass"`
	SaleChannel     string `db:"saleChannel"`
	OrderType       string `db:"orderType"`
	OrderSource     string `db:"orderSource"`
	Project         string `db:"project"`
	OwnerPoOrder    string `db:"ownerPoOrder"`
	SupplierName    string `db:"supplierName"`
	Product         string `db:"product"`
	RealPurchaseNum string `db:"realPurchaseNum"`
	TotalAmount     string `db:"totalAmount"`
	DiscountAmount  string `db:"discountAmount"`
	Freight         string `db:"freight"`
	DeliverNum      string `db:"deliverNum"`
	IncomingNum     string `db:"incomingNum"`
	CancelReturnNum string `db:"cancelReturnNum"`
	ClearingForm    string `db:"clearingForm"`
	NoteTaker       string `db:"noteTaker"`
	CreateTime      string `db:"createTime"`
	ActiveTime      string `db:"activeTime"`
}

func PurchaseOrderList(pagesize int, current int) ([]PurchaseOrderType, int, error) {
	offset := current * pagesize
	fmt.Println(pagesize, offset)
	var total int
	var users []PurchaseOrderType
	sql := "select * from purchaseorderlist LIMIT ? OFFSET ?"
	err := db.Select(&users, sql, pagesize, offset)
	db.Get(&total, "select count(*) from purchaseorderlist")
	if err != nil {
		fmt.Println("get failed", err)
		return nil, 0, err
	}
	return users, total, nil
}
func DeleteUser(id int) error {
	sql := "delete from user where user_id=?"
	_, err := db.Exec(sql, id)
	if err != nil {
		fmt.Println("exec failed", err)
		return err
	}
	return nil
}
