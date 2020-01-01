package models

import (
	"database/sql"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/shopspring/decimal"
	"time"
)

func NewInvestmentStructure(db *sqlx.DB) *InvestmentStructure {
	is := &InvestmentStructure{}
	is.db = db
	is.table = "investment_structure"
	is.hasID = true
	return is
}

type InvestmentStructure struct {
	Base
}

type InvestmentStructureRow struct {
	ID               int64           `db:"id"`
	Investment_ID    int64           `db:"investment_id"`
	ReportingDate    time.Time       `db:"ReportingDate"`
	Units            decimal.Decimal `db:"Units"`
	TotalInvested    decimal.Decimal `db:"TotalInvested"`
	ReportedValue    decimal.Decimal `db:"ReportedValue"`
	RealizedProceeds decimal.Decimal `db:"RealizedProceeds"`
	Structure        string          `db:"Structure"`
}

func (i *InvestmentStructureRow) FormattedReportingDate() string {
	return i.ReportingDate.Format("01/02/2006")
}

func (i *InvestmentStructure) userRowFromSqlResult(tx *sqlx.Tx, sqlResult sql.Result) (*InvestmentStructureRow, error) {
	isId, err := sqlResult.LastInsertId()
	if err != nil {
		return nil, err
	}

	return i.GetById(tx, isId)
}

// AllUsers returns all user rows.
func (i *InvestmentStructure) AllInvestmentStructures(tx *sqlx.Tx) ([]*InvestmentStructureRow, error) {
	isrs := []*InvestmentStructureRow{}
	query := fmt.Sprintf("SELECT * FROM %v", i.table)
	err := i.db.Select(&isrs, query)

	return isrs, err
}

// GetById returns record by id.
func (i *InvestmentStructure) GetById(tx *sqlx.Tx, id int64) (*InvestmentStructureRow, error) {
	isr := &InvestmentStructureRow{}
	if id == 0 {
		isr.ReportingDate = time.Now().AddDate(0, 0, -3)
		return isr, nil
	}
	query := fmt.Sprintf("SELECT * FROM %v WHERE id=?", i.table)
	err := i.db.Get(isr, query, id)

	return isr, err
}

// GetByName returns record by name.
func (i *InvestmentStructure) GetByName(tx *sqlx.Tx, name string) (*InvestmentStructureRow, error) {
	isr := &InvestmentStructureRow{}
	query := fmt.Sprintf("SELECT * FROM %v WHERE name=?", i.table)
	err := i.db.Get(isr, query, name)

	return isr, err
}

// create a new record of user.
func (i *InvestmentStructure) Create(tx *sqlx.Tx, m map[string]interface{}) (*InvestmentStructureRow, error) {
	sqlResult, err := i.InsertIntoTable(tx, m)
	if err != nil {
		return nil, err
	}
	return i.userRowFromSqlResult(tx, sqlResult)
}

// UpdateEmailAndPasswordById updates user email and password.
func (i *InvestmentStructure) UpdateById(tx *sqlx.Tx, isId int64, data map[string]interface{}) (*InvestmentStructureRow, error) {
	if len(data) > 0 {
		//calling base.go function
		_, err := i.UpdateByID(tx, data, isId)
		if err != nil {
			return nil, err
		}
	}

	return i.GetById(tx, isId)
}

// Get All by Investment ID.
func (i *InvestmentStructure) GetAllByInvestmentId(tx *sqlx.Tx, Investment_ID int64) ([]*InvestmentStructureRow, error) {
	css := []*InvestmentStructureRow{}
	query := fmt.Sprintf("SELECT * FROM %v WHERE Investment_ID=%v ORDER BY ReportingDate DESC", i.table, Investment_ID)
	err := i.db.Select(&css, query)

	return css, err
}

// UpdateEmailAndPasswordById updates user email and password.
func (i *InvestmentStructure) DeleteByID(tx *sqlx.Tx, csId int64) (sql.Result, error) {

	//calling base.go function
	sqlResult, err := i.DeleteById(tx, csId)
	if err != nil {
		return nil, err
	}

	return sqlResult, nil
}
