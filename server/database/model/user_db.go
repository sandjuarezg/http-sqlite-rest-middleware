package model

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
	"github.com/sandjuarezg/http-sqlite-rest-middleware/server/functionality"
)

func AddUser(db *sql.DB, user functionality.User) (err error) {
	smt, err := db.Prepare("INSERT INTO users (name, username, password, description) VALUES (?, ?, ?, ?)")
	if err != nil {
		return
	}
	defer smt.Close()

	_, err = smt.Exec(user.Name, user.Username, user.Pass, user.Description)
	if err != nil {
		return
	}

	return
}

func ShowUser(db *sql.DB) (users []functionality.User, err error) {
	rows, err := db.Query("SELECT id, name, username, password, description FROM users")
	if err != nil {
		return
	}
	defer rows.Close()

	var user functionality.User

	for rows.Next() {
		err = rows.Scan(&user.Id, &user.Name, &user.Username, &user.Pass, &user.Description)
		if err != nil {
			return
		}
		users = append(users, user)
	}

	return
}

func SearchUser(db *sql.DB, id int) (user functionality.User, err error) {
	var row *sql.Row = db.QueryRow("SELECT id, name, username, password, description FROM users WHERE id = ?", id)
	err = row.Scan(&user.Id, &user.Name, &user.Username, &user.Pass, &user.Description)
	if err != nil {
		return
	}

	return
}
