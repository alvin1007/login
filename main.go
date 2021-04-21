package main

import (
	"crypto/sha1"
	"database/sql"
	"encoding/hex"
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	var choice int
	var log_id int
	var log_password int
	var id string
	var password string
	fmt.Print("회원가입(0)/로그인(1) : ")
	fmt.Scan(&choice)
	if choice == 0 {
		fmt.Print("아이디 입력 : ")
		fmt.Scan(&id)
		nRow1 := id_sql(id)
		fmt.Print("비밀번호 입력 : ")
		fmt.Scan(&password)
		nRow2 := password_sql(id, password)
		if nRow1 == 1 && nRow2 == 1 {
			fmt.Println("회원가입 되셨습니다.")
		} else {
			fmt.Println("중간에 오류가 났습니다.")
		}
	} else if choice == 1 {
		fmt.Print("아이디 입력 : ")
		fmt.Scan(&id)
		log_id = login_id(id)
		fmt.Print("비밀번호 입력 : ")
		fmt.Scan(&password)
		log_password = login_password(id, password)
		if log_id == 1 && log_password == 1 {
			fmt.Println("로그인 되었습니다.")
		} else {
			fmt.Println("로그인에 실패하셨습니다.")
		}
	}
}

func id_sql(id string) int {
	conn, err := sql.Open("mysql", "alvin:alvin1007@tcp(127.0.0.1:3306)/alvin")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	ins, err := conn.Exec("insert into login (id, password) value (?, 'null')", id)
	if err != nil {
		return 0
	}
	nRow, err := ins.RowsAffected()
	if err != nil {
		return 0
	}
	conn.Close()
	return int(nRow)
}

func password_sql(id string, password string) int {
	conn, err := sql.Open("mysql", "alvin:alvin1007@tcp(127.0.0.1:3306)/alvin")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	hash := sha1.New()
	hash.Write([]byte(password))
	bs := hash.Sum(nil)
	bs_str := hex.EncodeToString(bs)
	ins, err := conn.Exec("update login set password = ? where id = ?", bs_str, id)
	if err != nil {
		return 0
	}
	nRow, err := ins.RowsAffected()
	if err != nil {
		return 0
	}
	conn.Close()
	return int(nRow)
}

func login_id(id string) int {
	conn, err := sql.Open("mysql", "alvin:alvin1007@tcp(127.0.0.1:3306)/alvin")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	var sql_id string
	err = conn.QueryRow("select id from login where id = ?", id).Scan(&sql_id)
	if err != nil {
		conn.Close()
		return 0
	}
	if id == sql_id {
		conn.Close()
		return 1
	}
	conn.Close()
	return 0
}

func login_password(id string, password string) int {
	conn, err := sql.Open("mysql", "alvin:alvin1007@tcp(127.0.0.1:3306)/alvin")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	var sql_password string
	err = conn.QueryRow("select password from login where id = ?", id).Scan(&sql_password)
	if err != nil {
		conn.Close()
		return 0
	}
	hash := sha1.New()
	hash.Write([]byte(password))
	bs := hash.Sum(nil)
	bs_str := hex.EncodeToString(bs)
	if sql_password == bs_str {
		conn.Close()
		return 1
	}
	conn.Close()
	return 0
}
