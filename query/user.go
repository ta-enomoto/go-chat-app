//ユーザー情報を扱うクエリパッケージ
package query

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"
	"goserver/config"
)

// マスタからSELECTしたデータをマッピングする構造体
type User struct {
	UserId   string `db:"USER_ID"`  // ID
	Password []byte `db:"PASSWORD"` // パスワード
}

var ConStrUsr string

func init() {
	confDbUsr, err := config.ReadConfDbUsr()
	if err != nil {
		fmt.Println(err.Error())
	}
	ConStrUsr = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s", confDbUsr.User, confDbUsr.Pass, confDbUsr.Host, confDbUsr.Port, confDbUsr.DbName, confDbUsr.Charset)
}

// ユーザーをdbに登録する関数
func InsertUser(userId string, password []byte, db *sql.DB) bool {

	hashed_pass, err := bcrypt.GenerateFromPassword(password, 10)
	if err != nil {
		fmt.Println(err.Error())
	}

	stmt, err := db.Prepare("INSERT INTO USERS(USER_ID,PASSWORD) VALUES(?,?)")
	if err != nil {
		fmt.Println(err.Error())
	}
	defer stmt.Close()

	_, err = stmt.Exec(userId, hashed_pass)
	if err != nil {
		return false
	} else {
		return true
	}
}

//ユーザーIDと一致するユーザー情報をdbから取得する関数
func SelectUserById(userId string, db *sql.DB) (user User) {

	err := db.QueryRow("SELECT USER_ID,PASSWORD FROM USERS WHERE USER_ID = ?", userId).Scan(&user.UserId, &user.Password)
	if err != nil {
		return
	}
	return
}

//ユーザーIDに一致するユーザーをdbから削除する関数。ハンドラでチェックはしてるが関数内でもパスも一致させたほうがいいかも
func DeleteUserById(userId string, db *sql.DB) bool {

	stmt, err := db.Prepare("DELETE FROM USERS WHERE USER_ID = ?")
	if err != nil {
		fmt.Println(err.Error())
	}
	defer stmt.Close()

	_, err = stmt.Exec(userId)
	if err != nil {
		return false
	} else {
		return true
	}
}

//全ユーザー情報をスライスとしてdbから取得する関数
func SelectAllUser(db *sql.DB) (users []User) {

	rows, err := db.Query("SELECT * FROM USERS")
	if err != nil {
		fmt.Println(err.Error())
	}

	for rows.Next() {
		user := User{}
		err := rows.Scan(&user.UserId, &user.Password)
		if err != nil {
			fmt.Println(err.Error())
		}
		users = append(users, user)
	}
	return
}

//ユーザー名が重複していないか確認する関数
func ContainsUserName(s []User, e string) bool {
	for _, v := range s {
		if e == v.UserId {
			return true
		}
	}
	return false
}
