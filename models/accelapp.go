package models

import (
	"database/sql"	
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/shopspring/decimal"
)

func NewAppl(db *sqlx.DB) *Appl {
	Appl := &Appl{}
	Appl.db = db
	Appl.table = "applications"
	Appl.hasID = true

	return Appl
}

type ApplRow struct {
	ID       int64  `db:"id"`
	Email    string `db:"Email"`
	FirstName    string `db:"FirstName"`
	LastName string `db:"LastName"`
	CompanyName    string `db:"CompanyName"`
	Phone    string `db:"Phone"`
	Website string `db:"Website"`
	Title    string `db:"Title"`
	Industries string `db:"Industries"`
	Locations    string `db:"Locations"`
	Comments string `db:"Comments"`
	CapitalRaised  decimal.Decimal `db:"CapitalRaised"`
}

type Appl struct {
	Base
}

func (i *Appl) userRowFromSqlResult(tx *sqlx.Tx, sqlResult sql.Result) (*ApplRow, error) {
	isId, err := sqlResult.LastInsertId()
	if err != nil {
		return nil, err
	}

	return i.GetById(tx, isId)
}

// AllUsers returns all user rows.
func (i *Appl) AllAppls(tx *sqlx.Tx) ([]*ApplRow, error) {
	isrs := []*ApplRow{}
	query := fmt.Sprintf("SELECT * FROM %v", i.table)
	err := i.db.Select(&isrs, query)

	return isrs, err
}

// GetById returns record by id.
func (i *Appl) GetById(tx *sqlx.Tx, id int64) (*ApplRow, error) {
	isr := &ApplRow{}
	query := fmt.Sprintf("SELECT * FROM %v WHERE id=?", i.table)
	err := i.db.Get(isr, query, id)

	return isr, err
}

// GetByName returns record by name.
func (i *Appl) GetByLastName(tx *sqlx.Tx, name string) (*ApplRow, error) {
	isr := &ApplRow{}
	query := fmt.Sprintf("SELECT * FROM %v WHERE LastName=?", i.table)
	err := i.db.Get(isr, query, name)

	return isr, err
}

// create a new record of user.
func (i *Appl) Create(tx *sqlx.Tx, m map[string]interface{}) (*ApplRow, error) {
	sqlResult, err := i.InsertIntoTable(tx, m)
	if err != nil {
		return nil, err
	}
	return i.userRowFromSqlResult(tx, sqlResult)
}

// UpdateEmailAndPasswordById updates user email and password.
func (i *Appl) UpdateById(tx *sqlx.Tx, isId int64, data map[string]interface{}) (*ApplRow, error) {
	if len(data) > 0 {
		//calling base.go function
		_, err := i.UpdateByID(tx, data, isId)
		if err != nil {
			return nil, err
		}
	}

	return i.GetById(tx, isId)
}

// UpdateEmailAndPasswordById updates user email and password.
func (i *Appl) DeleteByID(tx *sqlx.Tx, csId int64) (sql.Result, error) {

	//calling base.go function
	sqlResult, err := i.DeleteById(tx, csId)
	if err != nil {
		return nil, err
	}

	return sqlResult, nil
}
