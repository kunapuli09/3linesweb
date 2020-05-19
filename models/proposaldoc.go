package models

import (
	"database/sql"
	"fmt"
	"github.com/jmoiron/sqlx"
	"time"
)

func NewProposalDoc(db *sqlx.DB) *ProposalDoc {
	n := &ProposalDoc{}
	n.db = db
	n.table = "proposaldocs"
	n.hasID = true
	return n
}

type ProposalDoc struct {
	Base
}

type ProposalDocRow struct {
	ID            int64     `db:"id"`
	Investment_ID int64     `db:"investment_id"`
	Email         string    `db:"Email"`
	Phone         string    `db:"Phone"`
	FullName      string    `db:"FullName"`
	CompanyName   string    `db:"CompanyName"`
	UploadDate    time.Time `db:"UploadDate"`
	Hash          string    `db:"Hash"`
	DocPath       string    `db:"DocPath"`
	DocName       string    `db:"DocName"`
}

func (n *ProposalDocRow) FormattedUploadDate() string {
	return n.UploadDate.Format("01/02/2006")
}

func (i *ProposalDoc) userRowFromSqlResult(tx *sqlx.Tx, sqlResult sql.Result) (*ProposalDocRow, error) {
	nId, err := sqlResult.LastInsertId()
	if err != nil {
		return nil, err
	}

	return i.GetById(tx, nId)
}

// AllUsers returns all user rows.
func (i *ProposalDoc) AllDocs(tx *sqlx.Tx) ([]*ProposalDocRow, error) {
	nrs := []*ProposalDocRow{}
	query := fmt.Sprintf("SELECT * FROM %v", i.table)
	err := i.db.Select(&nrs, query)

	return nrs, err
}

// GetById returns record by id.
func (i *ProposalDoc) GetById(tx *sqlx.Tx, id int64) (*ProposalDocRow, error) {
	n := &ProposalDocRow{}
	query := fmt.Sprintf("SELECT * FROM %v WHERE id=?", i.table)
	err := i.db.Get(n, query, id)

	return n, err
}

// GetByName returns record by name.
func (i *ProposalDoc) GetByName(tx *sqlx.Tx, name string) (*ProposalDocRow, error) {
	n := &ProposalDocRow{}
	query := fmt.Sprintf("SELECT * FROM %v WHERE name=?", i.table)
	err := i.db.Get(n, query, name)

	return n, err
}

// create a new record of user.
func (i *ProposalDoc) Create(tx *sqlx.Tx, m map[string]interface{}) (*ProposalDocRow, error) {
	sqlResult, err := i.InsertIntoTable(tx, m)
	if err != nil {
		return nil, err
	}
	return i.userRowFromSqlResult(tx, sqlResult)
}

func (i *ProposalDoc) BatchInsert(tx *sqlx.Tx, docs []*ProposalDocRow) (sql.Result, error) {
	sqlStr := "INSERT INTO proposaldocs(Investment_ID, UploadDate, DocPath, Hash, DocName, Email, Phone, FullName, CompanyName) VALUES "
	vals := []interface{}{}
	for _, doc := range docs {
		sqlStr += "(?, ?, ?, ?, ?,?,?,?,?),"
		vals = append(vals, doc.Investment_ID, doc.UploadDate, doc.DocPath, doc.Hash, doc.DocName, doc.Email, doc.Phone, doc.FullName, doc.CompanyName)
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
func (i *ProposalDoc) UpdateById(tx *sqlx.Tx, nId int64, data map[string]interface{}) (*ProposalDocRow, error) {
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
func (i *ProposalDoc) GetAllByUserId(tx *sqlx.Tx, User_ID int64) ([]*ProposalDocRow, error) {
	css := []*ProposalDocRow{}
	query := fmt.Sprintf("SELECT * FROM %v WHERE user_id=%v", i.table, User_ID)
	err := i.db.Select(&css, query)

	return css, err
}

// UpdateEmailAndPasswordById updates user email and password.
func (i *ProposalDoc) DeleteByID(tx *sqlx.Tx, csId int64) (sql.Result, error) {

	//calling base.go function
	sqlResult, err := i.DeleteById(tx, csId)
	if err != nil {
		return nil, err
	}

	return sqlResult, nil
}
