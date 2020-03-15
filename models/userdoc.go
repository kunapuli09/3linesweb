package models

import (
	"database/sql"
	"fmt"
	"github.com/jmoiron/sqlx"
	"time"
)

func NewUserDoc(db *sqlx.DB) *UserDoc {
	n := &UserDoc{}
	n.db = db
	n.table = "userdocs"
	n.hasID = true
	return n
}

type UserDoc struct {
	Base
}

type UserDocRow struct {
	ID            int64     `db:"id"`
	User_ID 	  int64   `db:"user_id"`
	UploadDate    time.Time `db:"UploadDate"`
	Hash          string    `db:"Hash"`
	DocPath       string    `db:"DocPath"`
	DocName       string    `db:"DocName"`
}

func (n *UserDocRow) FormattedUploadDate() string {
	return n.UploadDate.Format("01/02/2006")
}

func (i *UserDoc) userRowFromSqlResult(tx *sqlx.Tx, sqlResult sql.Result) (*UserDocRow, error) {
	nId, err := sqlResult.LastInsertId()
	if err != nil {
		return nil, err
	}

	return i.GetById(tx, nId)
}

// AllUsers returns all user rows.
func (i *UserDoc) AllDocs(tx *sqlx.Tx) ([]*UserDocRow, error) {
	nrs := []*UserDocRow{}
	query := fmt.Sprintf("SELECT * FROM %v", i.table)
	err := i.db.Select(&nrs, query)

	return nrs, err
}

// GetById returns record by id.
func (i *UserDoc) GetById(tx *sqlx.Tx, id int64) (*UserDocRow, error) {
	n := &UserDocRow{}
	query := fmt.Sprintf("SELECT * FROM %v WHERE id=?", i.table)
	err := i.db.Get(n, query, id)

	return n, err
}

// GetByName returns record by name.
func (i *UserDoc) GetByName(tx *sqlx.Tx, name string) (*UserDocRow, error) {
	n := &UserDocRow{}
	query := fmt.Sprintf("SELECT * FROM %v WHERE name=?", i.table)
	err := i.db.Get(n, query, name)

	return n, err
}

// create a new record of user.
func (i *UserDoc) Create(tx *sqlx.Tx, m map[string]interface{}) (*UserDocRow, error) {
	sqlResult, err := i.InsertIntoTable(tx, m)
	if err != nil {
		return nil, err
	}
	return i.userRowFromSqlResult(tx, sqlResult)
}

func (i *UserDoc) BatchInsert(tx *sqlx.Tx, docs []*UserDocRow) (sql.Result, error) {
	sqlStr := "INSERT INTO userdocs(user_id, UploadDate, DocPath, Hash, DocName) VALUES "
	vals := []interface{}{}
	for _, doc := range docs {
		sqlStr += "(?, ?, ?, ?, ?),"
		vals = append(vals, doc.User_ID, doc.UploadDate, doc.DocPath, doc.Hash, doc.DocName)
	}
	//trim the last ,
	sqlStr = sqlStr[0 : len(sqlStr)-1]
	//prepare the statement
	stmt, err := i.db.Prepare(sqlStr)
	if err != nil {
		fmt.Println(err)
	}

	//format all vals at once
	sqlResult, err := stmt.Exec(vals...)

	return sqlResult, err
}

// UpdateEmailAndPasswordById updates user email and password.
func (i *UserDoc) UpdateById(tx *sqlx.Tx, nId int64, data map[string]interface{}) (*UserDocRow, error) {
	if len(data) > 0 {
		//calling base.go function
		_, err := i.UpdateByID(tx, data, nId)
		if err != nil {
			return nil, err
		}
	}

	return i.GetById(tx, nId)
}

// Get All by Investment ID.
func (i *UserDoc) GetAllByUserId(tx *sqlx.Tx, User_ID int64) ([]*UserDocRow, error) {
	css := []*UserDocRow{}
	query := fmt.Sprintf("SELECT * FROM %v WHERE user_id=%v", i.table, User_ID)
	err := i.db.Select(&css, query)

	return css, err
}

// UpdateEmailAndPasswordById updates user email and password.
func (i *UserDoc) DeleteByID(tx *sqlx.Tx, csId int64) (sql.Result, error) {

	//calling base.go function
	sqlResult, err := i.DeleteById(tx, csId)
	if err != nil {
		return nil, err
	}

	return sqlResult, nil
}
