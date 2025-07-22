package main

import (
	"fmt"
	"go_sqlx/model"
	"net/url"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

func main() {
	user := "xxxxx"
	password := "xxxxx"
	enPwd := url.QueryEscape(password)
	dsn := fmt.Sprintf("%s:%s@tcp(127.0.0.1:3306)/go_test?charset=utf8mb4&parseTime=True&loc=Local", user, enPwd)

	db, err := sqlx.Connect("mysql", dsn)
	if err != nil {
		fmt.Println(err)
		return
	}

	if err := db.Ping(); err != nil {
		fmt.Println(err)
		return
	}

	defer db.Close()
	//list, err := SelectByDeptName(db, "技术部")
	//if err != nil {
	//	fmt.Printf("SelectByDeptName failed. err: %s", err.Error())
	//	return
	//}
	//fmt.Println(list)

	employee, err := SelectMaxSalary(db)
	if err != nil {
		fmt.Printf("SelectMaxSalary failed. err: %s", err.Error())
		return
	}
	fmt.Println(employee)
}

func SelectByDeptName(db *sqlx.DB, deptName string) (list []model.Employee, err error) {
	err = db.Select(&list, "SELECT * FROM tbl_employee WHERE department = ?", deptName)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return list, nil
}

func SelectMaxSalary(db *sqlx.DB) (emplyoee *model.Employee, err error) {
	var e model.Employee
	err = db.Get(&e, "SELECT * FROM tbl_employee ORDER BY salary DESC LIMIT 1")
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return &e, nil
}
