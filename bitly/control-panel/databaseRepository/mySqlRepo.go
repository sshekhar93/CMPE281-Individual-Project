package repository

import (
	"database/sql"
	"fmt"

	entity "github.com/cmpe281-sshekhar93/bitly/control-panel/linkEntity"

	//Import disappears if writing the package name
	_ "github.com/go-sql-driver/mysql"
)

type mysqlRepo struct {
	username string
	password string
	host     string
	port     string
	database string
}

var (
	uri string
	db  *sql.DB
)

// NewMysqlRepository returns mysql driver for database
func NewMysqlRepository(uname, pass, host, port, dbase string) Repository {
	// mysqlRepo.username	= uname
	// mysqlRepo.password	= pass
	// mysqlRepo.host		= host
	if port == "" {
		port = "3306"
	}
	// mysqlRepo.database	= dbase
	uri = uname + ":" + pass + "@tcp(" + host + ":" + port + ")/" + dbase
	db, _ = connect()
	return &mysqlRepo{uname, pass, host, port, dbase}
}

func connect() (*sql.DB, error) {
	// fmt.Println("Inside repo connect func")
	db, err := sql.Open("mysql", uri)
	if err != nil {
		//panic(err.Error())
		db = nil
		// defer db.Close()
	}
	return db, err
}

//INSERT data into database
func (*mysqlRepo) INSERT(linkData *entity.LinkData) error {
	// fmt.Println("Inside repo INSERT func")
	if db == nil {
		db, err := connect()
		if err != nil {
			defer db.Close()
			return err
		}
	}

	insert, err := db.Query("INSERT INTO shortLinks (shortLink, uri) VALUES ( '" + linkData.ShotLink + "', '" + linkData.Uri + "' );")
	if err != nil {
		defer insert.Close()
		return err
	}
	defer insert.Close()

	fetch, err := db.Query("SELECT * FROM shortLinks WHERE shortLinks.shortLink = '" + linkData.ShotLink + "';")
	if err != nil {
		defer fetch.Close()
		return err
	}
	defer fetch.Close()

	for fetch.Next() {
		err = fetch.Scan(&linkData.Id, &linkData.ShotLink, &linkData.Uri)
		if err != nil {
			defer fetch.Close()
			return err
		}
	}
	fmt.Println("INSERT(): shortLink ID: " + string(linkData.Id))
	return err
}

func (*mysqlRepo) FETCH(linkData *entity.LinkData, id string) error {
	if db == nil {
		db, err := connect()
		if err != nil {
			defer db.Close()
			return err
		}
	}

	fetch, err := db.Query("SELECT * FROM shortLinks WHERE shortLinks.ID = '" + id + "';")
	if err != nil {
		defer fetch.Close()
		return err
	}
	defer fetch.Close()

	for fetch.Next() {
		err = fetch.Scan(&linkData.Id, &linkData.ShotLink, &linkData.Uri)
		if err != nil {
			defer fetch.Close()
			return err
		}
	}
	defer fetch.Close()
	return err
}

/*func (*mysqlRepo) FETCHREDIRECT(linkData *entity.LinkData, shortLink string) error {
	fmt.Println("FETCHREDIRECT(): ", linkData, shortLink)
	db, err := connect()
	if err != nil {
		defer db.Close()
		return err
	}
	fetch, err := db.Query("SELECT * FROM shortLinks WHERE shortLinks.shortLink = '" + shortLink + "';")
	if err != nil {
		defer fetch.Close()
		return err
	}
	for fetch.Next() {
		err = fetch.Scan(&linkData.Id, &linkData.ShotLink, &linkData.Uri)
		if err != nil {
			defer fetch.Close()
			return err
		}
	}
	return err
}*/

func (*mysqlRepo) FETCHALL() ([]entity.LinkData, error) {
	if db == nil {
		db, err := connect()
		if err != nil {
			defer db.Close()
			return nil, err
		}
	}

	fetch, err := db.Query("SELECT * FROM shortLinks;")
	if err != nil {
		defer fetch.Close()
		return nil, err
	}
	defer fetch.Close()

	var linkDatas []entity.LinkData
	linkData := entity.LinkData{}
	for fetch.Next() {
		err = fetch.Scan(&linkData.Id, &linkData.ShotLink, &linkData.Uri)
		if err != nil {
			defer fetch.Close()
			return nil, err
		}
		linkDatas = append(linkDatas, linkData)
	}
	return linkDatas, err
}

func (*mysqlRepo) UPDATE(linkData *entity.LinkData) error {
	var err error
	if db == nil {
		db, err = connect()
		if err != nil {
			defer db.Close()
			return err
		}
	}
	if linkData.Id != 0 {
		update, err := db.Query("UPDATE shortLinks SET shortLinks.uri = '" + linkData.Uri + "' WHERE shortLinks.ID = '" + string(linkData.Id) + "';")
		if err != nil {
			defer update.Close()
			return err
		}
		defer update.Close()
		for update.Next() {
			err = update.Scan(&linkData.Id, &linkData.ShotLink, &linkData.Uri)
			if err != nil {
				defer update.Close()
				return err
			}
		}
	} else if linkData.ShotLink != "" {
		update, err := db.Query("UPDATE shortLinks SET shortLinks.uri = '" + linkData.Uri + "' WHERE shortLinks.shortLink = '" + linkData.ShotLink + "';")
		for update.Next() {
			err = update.Scan(&linkData.Id, &linkData.ShotLink, &linkData.Uri)
			if err != nil {
				defer update.Close()
				return err
			}
		}
	}
	if err != nil {
		defer db.Close()
		return err
	}
	return err
}

func (*mysqlRepo) DELETE(id string) error {
	if db == nil {
		db, err := connect()
		if err != nil {
			defer db.Close()
			return err
		}
	}

	delete, err := db.Query("DELETE FROM shortLinks WHERE shortLinks.ID = '" + id + "';")
	if err != nil {
		defer delete.Close()
		return err
	}
	defer delete.Close()
	if err != nil {
		defer delete.Close()
	}
	return err
}
