package main

import (
	"database/sql"
	"fmt"
	"net/url"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	user := "xxxxx"
	password := "xxxxx"
	enPwd := url.QueryEscape(password)
	dsn := fmt.Sprintf("%s:%s@tcp(127.0.0.1:3306)/go_test?charset=utf8mb4&parseTime=True&loc=Local", user, enPwd)
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		fmt.Println(err)
		return
	}

	if err := db.Ping(); err != nil {
		fmt.Println(err)
		return
	}

	defer db.Close()
	//insertStu(db)
	//selectStu(db)
	//updateStu(db)
	//deleteStu(db)
	transaction(db, 2, 1, 100)
}

func insertStu(db *sql.DB) {
	var sql string = fmt.Sprintf("INSERT INTO tbl_student(name, age, grade) VALUES ('%s', %d, '%s')",
		"张三", 20, "三年级")
	res, err := db.Exec(sql)
	if err != nil {
		fmt.Println(err)
	}
	id, _ := res.LastInsertId()
	fmt.Println(id)
}

func selectStu(db *sql.DB) {
	var sql string = "SELECT * FROM tbl_student WHERE age > 18"
	rows, err := db.Query(sql)
	if err != nil {
		fmt.Println(err)
	}
	defer rows.Close()
	for rows.Next() {
		var id, age int
		var name, grade string
		if err := rows.Scan(&id, &name, &age, &grade); err != nil {
			fmt.Println(err)
		}
		fmt.Println(id, name, age, grade)
	}
}

func updateStu(db *sql.DB) {
	var sql string = fmt.Sprintf("UPDATE tbl_student SET grade = '%s' WHERE name = '%s' ", "四年级", "张三")
	res, err := db.Exec(sql)
	if err != nil {
		fmt.Println(err)
	}
	id, _ := res.LastInsertId()
	fmt.Println(id)
}

func deleteStu(db *sql.DB) {
	var sql string = fmt.Sprintf("DELETE FROM tbl_student WHERE age < %d ", 15)
	_, err := db.Exec(sql)
	if err != nil {
		fmt.Println(err)
	}
}

// 转账
func transaction(db *sql.DB, fromAccountId int, toAccountId int, amount float64) {
	tx, err := db.Begin()
	if err != nil {
		fmt.Println(err)
	}
	defer tx.Rollback() //事务回滚
	var balance float64
	err = tx.QueryRow("SELECT balance from tbl_account WHERE id = ? FOR UPDATE", fromAccountId).Scan(&balance)
	if err != nil {
		fmt.Printf("查询账号A的余额失败！err = %s\n", err.Error())
		return
	}

	if balance < amount {
		fmt.Println("账号A的余额小于转账金额！")
		return
	}

	// 扣除账户 A 余额
	_, err = tx.Exec("UPDATE tbl_account SET balance = balance - ? WHERE id = ?", amount, fromAccountId)
	if err != nil {
		fmt.Printf("执行扣除账号A的余额失败！ %s\n", err.Error())
		return
	}

	// 更新账户 B 余额
	_, err = tx.Exec("UPDATE tbl_account SET balance = balance + ? WHERE id = ?", amount, toAccountId)
	if err != nil {
		fmt.Printf("更新账号B的余额失败！ %s\n", err.Error())
		return
	}
	// 插入交易记录
	_, err = tx.Exec("INSERT INTO tbl_transaction (from_account_id, to_account_id, amount) VALUES (?, ?, ?)",
		fromAccountId, toAccountId, amount)
	if err != nil {
		fmt.Printf("插入交易记录失败！ %s\n", err.Error())
		return
	}

	// 提交事务
	if err := tx.Commit(); err != nil {
		fmt.Printf("提交事务失败！ %s", err.Error())
		return
	}
}
