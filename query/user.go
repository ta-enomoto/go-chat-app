package query // 独自のクエリパッケージ

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"goserver/config"
)

// マスタからSELECTしたデータをマッピングする構造体
type USER struct {
	Id       string `db:"ID"`       // ID
	Password string `db:"PASSWORD"` // パスワード
}

var confDB *config.Config
var ConStr string

func init() {
	_confDB, err := config.ReadConfDB()
	if err != nil {
		fmt.Println(err.Error())
	}
	confDB = _confDB
	_conStr := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s", confDB.User, confDB.Pass, confDB.Host, confDB.Port, confDB.DbName, confDB.Charset)
	ConStr = _conStr
}

// データ登録関数
func InsertUser(id string, password string, db *sql.DB) bool {

	// プリペアードステートメント
	stmt, err := db.Prepare("INSERT INTO USERS(ID,PASSWORD) VALUES(?,?)")
	if err != nil {
		fmt.Println(err.Error())
	}
	defer stmt.Close()

	// クエリ実行
	insertedOrNot, err := stmt.Exec(id, password)
	if err != nil {
		return false
	} else {
		_ = insertedOrNot
		return true
	}
}

// 単一行データ取得関数
func SelectUserById(id string, db *sql.DB) (userinfo USER) {

	// 構造体USER型の変数userを宣言
	user := USER{}

	// プリペアードステートメント
	err := db.QueryRow("SELECT ID,PASSWORD FROM USERS WHERE ID = ?", id).Scan(&user.Id, &user.Password)
	if err != nil {
		return
	}
	return user
}

//単一行データ削除関数
func DeleteUserById(id string, db *sql.DB) bool {

	stmt, err := db.Prepare("DELETE FROM USERS WHERE ID = ?")
	if err != nil {
		fmt.Println(err.Error())
	}
	defer stmt.Close()
	deletedOrNot, err := stmt.Exec(id)
	if err != nil {
		return false
	} else {
		_ = deletedOrNot
		return true
	}
}
