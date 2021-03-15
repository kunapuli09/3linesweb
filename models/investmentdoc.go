package models

import (
	"database/sql"
	"fmt"
	"github.com/jmoiron/sqlx"
	"time"
)

func NewInvestmentDoc(db *sqlx.DB) *InvestmentDoc {
	n := &InvestmentDoc{}
	n.db = db
	n.table = "docs"
	n.hasID = true
	return n
}

type InvestmentDoc struct {
	Base
}

type InvestmentDocRow struct {
	ID            int64     `db:"id"`
	Investment_ID int64     `db:"investment_id"`
	UploadDate    time.Time `db:"UploadDate"`
	DocPath       string    `db:"DocPath"`
	Hash          string    `db:"Hash"`
	DocName       string    `db:"DocName"`
}

func (n *InvestmentDocRow) FormattedUploadDate() string {
	return n.UploadDate.Format("01/02/2006")
}

func (i *InvestmentDoc) userRowFromSqlResult(tx *sqlx.Tx, sqlResult sql.Result) (*InvestmentDocRow, error) {
	nId, err := sqlResult.LastInsertId()
	if err != nil {
		return nil, err
	}

	return i.GetById(tx, nId)
}

// AllUsers returns all user rows.
func (i *InvestmentDoc) AllDocs(tx *sqlx.Tx) ([]*InvestmentDocRow, error) {
	nrs := []*InvestmentDocRow{}
	query := fmt.Sprintf("SELECT * FROM %v", i.table)
	err := i.db.Select(&nrs, query)

	return nrs, err
}

// GetById returns record by id.
func (i *InvestmentDoc) GetById(tx *sqlx.Tx, id int64) (*InvestmentDocRow, error) {
	n := &InvestmentDocRow{}
	query := fmt.Sprintf("SELECT * FROM %v WHERE id=?", i.table)
	err := i.db.Get(n, query, id)

	return n, err
}

// GetByName returns record by name.
func (i *InvestmentDoc) GetByName(tx *sqlx.Tx, name string) (*InvestmentDocRow, error) {
	n := &InvestmentDocRow{}
	query := fmt.Sprintf("SELECT * FROM %v WHERE name=?", i.table)
	err := i.db.Get(n, query, name)

	return n, err
}

// create a new record of user.
func (i *InvestmentDoc) Create(tx *sqlx.Tx, m map[string]interface{}) (*InvestmentDocRow, error) {
	sqlResult, err := i.InsertIntoTable(tx, m)
	if err != nil {
		return nil, err
	}
	return i.userRowFromSqlResult(tx, sqlResult)
}

func (i *InvestmentDoc) BatchInsert(tx *sqlx.Tx, docs []*InvestmentDocRow) (sql.Result, error) {
	sqlStr := "INSERT INTO docs(investment_id,  UploadDate, DocPath, Hash, DocName) VALUES "
	vals := []interface{}{}
	for _, doc := range docs {
		sqlStr += "(?, ?, ?, ?, ?, ?),"
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
func (i *InvestmentDoc) UpdateById(tx *sqlx.Tx, nId int64, data map[string]interface{}) (*InvestmentDocRow, error) {
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
func (i *InvestmentDoc) GetAllByInvestmentId(tx *sqlx.Tx, Investment_ID int64) ([]*InvestmentDocRow, error) {
	css := []*InvestmentDocRow{}
	query := fmt.Sprintf("SELECT * FROM %v WHERE Investment_ID=%v", i.table, Investment_ID)
	err := i.db.Select(&css, query)

	return css, err
}

// UpdateEmailAndPasswordById updates user email and password.
func (i *InvestmentDoc) DeleteByID(tx *sqlx.Tx, csId int64) (sql.Result, error) {

	//calling base.go function
	sqlResult, err := i.DeleteById(tx, csId)
	if err != nil {
		return nil, err
	}

	return sqlResult, nil
}
