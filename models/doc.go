package models

import (
	"database/sql"
	"fmt"
	"github.com/jmoiron/sqlx"
	"time"
)

func NewDoc(db *sqlx.DB) *Doc {
	n := &Doc{}
	n.db = db
	n.table = "docs"
	n.hasID = true
	return n
}

type Doc struct {
	Base
}

type DocRow struct {
	ID            int64     `db:"id"`
	Investment_ID int64     `db:"investment_id"`
	UploadDate    time.Time `db:"UploadDate"`
	DocPath       string    `db:"DocPath"`
	Hash          string    `db:"Hash"`
	DocName       string    `db:"DocName"`
}

func (n *DocRow) FormattedUploadDate() string {
	return n.UploadDate.Format("01/02/2006")
}

func (i *Doc) userRowFromSqlResult(tx *sqlx.Tx, sqlResult sql.Result) (*DocRow, error) {
	nId, err := sqlResult.LastInsertId()
	if err != nil {
		return nil, err
	}

	return i.GetById(tx, nId)
}

// AllUsers returns all user rows.
func (i *Doc) AllDocs(tx *sqlx.Tx) ([]*DocRow, error) {
	nrs := []*DocRow{}
	query := fmt.Sprintf("SELECT * FROM %v", i.table)
	err := i.db.Select(&nrs, query)

	return nrs, err
}

// GetById returns record by id.
func (i *Doc) GetById(tx *sqlx.Tx, id int64) (*DocRow, error) {
	n := &DocRow{}
	query := fmt.Sprintf("SELECT * FROM %v WHERE id=?", i.table)
	err := i.db.Get(n, query, id)

	return n, err
}

// GetByName returns record by name.
func (i *Doc) GetByName(tx *sqlx.Tx, name string) (*DocRow, error) {
	n := &DocRow{}
	query := fmt.Sprintf("SELECT * FROM %v WHERE name=?", i.table)
	err := i.db.Get(n, query, name)

	return n, err
}

// create a new record of user.
func (i *Doc) Create(tx *sqlx.Tx, m map[string]interface{}) (*DocRow, error) {
	sqlResult, err := i.InsertIntoTable(tx, m)
	if err != nil {
		return nil, err
	}
	return i.userRowFromSqlResult(tx, sqlResult)
}

func (i *Doc) BatchInsert(tx *sqlx.Tx, docs []*DocRow) (sql.Result, error) {
	sqlStr := "INSERT INTO docs(investment_id,  UploadDate, DocPath, Hash, DocName) VALUES "
	vals := []interface{}{}
	for _, doc := range docs {
		sqlStr += "(?, ?, ?, ?, ?),"
		vals = append(vals, doc.Investment_ID, time.Now(), doc.DocPath, doc.Hash, doc.DocName)
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
func (i *Doc) UpdateById(tx *sqlx.Tx, nId int64, data map[string]interface{}) (*DocRow, error) {
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
func (i *Doc) GetAllByInvestmentId(tx *sqlx.Tx, Investment_ID int64) ([]*DocRow, error) {
	css := []*DocRow{}
	query := fmt.Sprintf("SELECT * FROM %v WHERE Investment_ID=%v", i.table, Investment_ID)
	err := i.db.Select(&css, query)

	return css, err
}

// UpdateEmailAndPasswordById updates user email and password.
func (i *Doc) DeleteByID(tx *sqlx.Tx, csId int64) (sql.Result, error) {

	//calling base.go function
	sqlResult, err := i.DeleteById(tx, csId)
	if err != nil {
		return nil, err
	}

	return sqlResult, nil
}
