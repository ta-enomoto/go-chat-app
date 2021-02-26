package query // 独自のクエリパッケージ

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"goserver/config"
)

// マスタからSELECTしたデータをマッピングする構造体
type User struct {
	UserId   string `db:"USER_ID"`  // ID
	Password string `db:"PASSWORD"` // パスワード
}

var confDbUsr *config.ConfigUsr
var ConStrUsr string

func init() {
	_confDbUsr, err := config.ReadConfDbUsr()
	if err != nil {
		fmt.Println(err.Error())
	}
	confDbUsr = _confDbUsr
	_conStrUsr := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s", confDbUsr.User, confDbUsr.Pass, confDbUsr.Host, confDbUsr.Port, confDbUsr.DbName, confDbUsr.Charset)
	ConStrUsr = _conStrUsr
}

// データ登録関数
func InsertUser(userId string, password string, db *sql.DB) bool {

	//ユーザー名・パスの両方が重複したときの処理書く！！
	// プリペアードステートメント
	stmt, err := db.Prepare("INSERT INTO USERS(USER_ID,PASSWORD) VALUES(?,?)")
	if err != nil {
		fmt.Println(err.Error())
	}
	defer stmt.Close()

	// クエリ実行
	insertedOrNot, err := stmt.Exec(userId, password)
	if err != nil {
		return false
	} else {
		_ = insertedOrNot
		return true
	}
}

// 単一行データ取得関数
func SelectUserById(userId string, db *sql.DB) (user User) {

	// プリペアードステートメント
	err := db.QueryRow("SELECT USER_ID,PASSWORD FROM USERS WHERE USER_ID = ?", userId).Scan(&user.UserId, &user.Password)
	if err != nil {
		return
	}
	return
}

//単一行データ削除関数
func DeleteUserById(userId string, db *sql.DB) bool {

	stmt, err := db.Prepare("DELETE FROM USERS WHERE USER_ID = ?")
	if err != nil {
		fmt.Println(err.Error())
	}
	defer stmt.Close()
	deletedOrNot, err := stmt.Exec(userId)
	if err != nil {
		return false
	} else {
		_ = deletedOrNot
		return true
	}
}

//全ユーザー取得関数
func SelectAllUser(db *sql.DB) (users []User) {

	// プリペアードステートメント
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

//ユーザー名重複確認関数
func ContainsUserName(s []User, e string) bool {
	for _, v := range s {
		if e == v.UserId {
			return true
		}
	}
	return false
}
