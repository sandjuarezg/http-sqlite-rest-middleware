package functionality

import (
	"database/sql"
	"io/ioutil"
	"log"
	"os"
)

type User struct {
	Id          int    `json:"id"`
	Name        string `json:"name"`
	Username    string `json:"username"`
	Pass        string `json:"pass"`
	Description string `json:"description"`
}

func SqlMigration() {
	//Check migration.sql
	var _, err = os.Stat("./database/migration.sql")
	if os.IsNotExist(err) {
		log.Fatal(err)
	}

	//Get content
	file, _ := os.Open("./database/migration.sql")
	content, err := ioutil.ReadAll(file)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	//Check user.db
	_, err = os.Stat("./database/user.db")
	if os.IsNotExist(err) {
		var file, err = os.Create("./database/user.db")
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()
	}

	//Check table
	db, err := sql.Open("sqlite3", "./database/user.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	_, err = db.Query("SELECT * from users")
	if err != nil {
		_, err = db.Exec(string(content))
		if err != nil {
			log.Fatal(err)
		}
	}
}
