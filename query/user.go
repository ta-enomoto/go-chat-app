package query // 独自のクエリパッケージ

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

// マスタからSELECTしたデータをマッピングする構造体
type USER struct {
	Id       string `db:"ID"`       // ID
	Password string `db:"PASSWORD"` // パスワード
}

// データ登録関数
func InsertUser(loginID, password string, db *sql.DB) error {

	// プリペアードステートメント
	stmt, err := db.Prepare("INSERT INTO USERS(ID,PASSWORD) VALUES(?,?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	// クエリ実行
	_, err = stmt.Exec(loginID, password)
	if err != nil {
		return err
	}
	return nil
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