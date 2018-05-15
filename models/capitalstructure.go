package models

import (
	"database/sql"
	"fmt"
	"github.com/jmoiron/sqlx"
	"time"
	"github.com/shopspring/decimal"
)

func NewCapitalStructure(db *sqlx.DB) *CapitalStructure {
	cs := &CapitalStructure{}
	cs.db = db
	cs.table = "capital_structure"
	cs.hasID = true
	return cs
}

type CapitalStructure struct {
	Base
}

type CapitalizationStructure struct {
	ID            int64     `db:"id"`
	Investment_ID int64     `db:"investment_id"`
	ReportingDate time.Time `db:"ReportingDate"`
	//MaturityDate   time.Time `db:"maturity_date"`
	ClosingValue   decimal.Decimal `db:"ClosingValue"`
	YearEndValue   decimal.Decimal `db:"YearEndValue"`
	Capitalization string  `db:"Capitalization"`
}

func (i *CapitalStructure) userRowFromSqlResult(tx *sqlx.Tx, sqlResult sql.Result) (*CapitalizationStructure, error) {
	csId, err := sqlResult.LastInsertId()
	if err != nil {
		return nil, err
	}

	return i.GetById(tx, csId)
}

// AllUsers returns all user rows.
func (i *CapitalStructure) AllCapitalStructures(tx *sqlx.Tx) ([]*CapitalizationStructure, error) {
	css := []*CapitalizationStructure{}
	query := fmt.Sprintf("SELECT * FROM %v", i.table)
	err := i.db.Select(&css, query)

	return css, err
}

// GetById returns record by id.
func (i *CapitalStructure) GetById(tx *sqlx.Tx, id int64) (*CapitalizationStructure, error) {
	cs := &CapitalizationStructure{}
	query := fmt.Sprintf("SELECT * FROM %v WHERE id=?", i.table)
	err := i.db.Get(cs, query, id)

	return cs, err
}


// create a new record of user.
func (i *CapitalStructure) Create(tx *sqlx.Tx, m map[string]interface{}) (*CapitalizationStructure, error) {
	sqlResult, err := i.InsertIntoTable(tx, m)
	if err != nil {
		return nil, err
	}
	return i.userRowFromSqlResult(tx, sqlResult)
}

// Get All by Investment ID.
func (i *CapitalStructure) GetAllByInvestmentId(tx *sqlx.Tx, Investment_ID int64) ([]*CapitalizationStructure, error) {
	css := []*CapitalizationStructure{}
	query := fmt.Sprintf("SELECT * FROM %v WHERE Investment_ID=%v", i.table, Investment_ID)
	err := i.db.Select(&css, query)

	return css, err
}

// UpdateEmailAndPasswordById updates user email and password.
func (i *CapitalStructure) DeleteByID(tx *sqlx.Tx, csId int64) (sql.Result, error) {

	//calling base.go function
	sqlResult, err := i.DeleteById(tx, csId)
	if err != nil {
		return nil, err
	}

	return sqlResult, nil
}
