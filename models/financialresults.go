package models

import (
	"database/sql"
	"fmt"
	"github.com/jmoiron/sqlx"
	"time"
)

func NewFinancialResults(db *sqlx.DB) *FinancialResults {
	fr := &FinancialResults{}
	fr.db = db
	fr.table = "financial_results"
	fr.hasID = true
	return fr
}

type FinancialResults struct {
	Base
}

type FinancialResultsRow struct {
	ID                     int64     `db:"id"`
	Investment_ID          int64     `db:"investment_id"`
	ReportingDate          time.Time `db:"ReportingDate"`
	Revenue                float64   `db:"Revenue"`
	YoYGrowthPercentage1   float64   `db:"YoYGrowthPercentage1"`
	LTMEBITDA              float64   `db:"LTMEBITDA"`
	YoYGrowthPercentage2   float64   `db:"YoYGrowthPercentage2"`
	EBITDAMargin           float64   `db:"EBITDAMargin"`
	TotalExitValue         float64   `db:"TotalExitValue"`
	TotalExitValueMultiple float64   `db:"TotalExitValueMultiple"`
	TotalLeverage          float64   `db:"TotalLeverage"`
	TotalLeverageMultiple  float64   `db:"TotalLeverageMultiple"`
	Assessment             string    `db:"Assessment"`
}

func (i *FinancialResults) userRowFromSqlResult(tx *sqlx.Tx, sqlResult sql.Result) (*FinancialResultsRow, error) {
	frId, err := sqlResult.LastInsertId()
	if err != nil {
		return nil, err
	}

	return i.GetById(tx, frId)
}

// AllUsers returns all user rows.
func (i *FinancialResults) AllFinancialResultss(tx *sqlx.Tx) ([]*FinancialResultsRow, error) {
	frs := []*FinancialResultsRow{}
	query := fmt.Sprintf("SELECT * FROM %v", i.table)
	err := i.db.Select(&frs, query)

	return frs, err
}

// GetById returns record by id.
func (i *FinancialResults) GetById(tx *sqlx.Tx, id int64) (*FinancialResultsRow, error) {
	fr := &FinancialResultsRow{}
	query := fmt.Sprintf("SELECT * FROM %v WHERE id=?", i.table)
	err := i.db.Get(fr, query, id)

	return fr, err
}

// GetByName returns record by name.
func (i *FinancialResults) GetByName(tx *sqlx.Tx, name string) (*FinancialResultsRow, error) {
	fr := &FinancialResultsRow{}
	query := fmt.Sprintf("SELECT * FROM %v WHERE name=?", i.table)
	err := i.db.Get(fr, query, name)

	return fr, err
}

// create a new record of user.
func (i *FinancialResults) Create(tx *sqlx.Tx, m map[string]interface{}) (*FinancialResultsRow, error) {
	sqlResult, err := i.InsertIntoTable(tx, m)
	if err != nil {
		return nil, err
	}
	return i.userRowFromSqlResult(tx, sqlResult)
}

// UpdateEmailAndPasswordById updates user email and password.
func (i *FinancialResults) UpdateById(tx *sqlx.Tx, frId int64, data map[string]interface{}) (*FinancialResultsRow, error) {
	if len(data) > 0 {
		//calling base.go function
		_, err := i.UpdateByID(tx, data, frId)
		if err != nil {
			return nil, err
		}
	}

	return i.GetById(tx, frId)
}

// Get All by Investment ID.
func (i *FinancialResults) GetAllByInvestmentId(tx *sqlx.Tx, Investment_ID int64) ([]*FinancialResultsRow, error) {
	css := []*FinancialResultsRow{}
	query := fmt.Sprintf("SELECT * FROM %v WHERE Investment_ID=%v", i.table, Investment_ID)
	err := i.db.Select(&css, query)

	return css, err
}

// UpdateEmailAndPasswordById updates user email and password.
func (i *FinancialResults) DeleteByID(tx *sqlx.Tx, csId int64) (sql.Result, error) {

	//calling base.go function
	sqlResult, err := i.DeleteById(tx, csId)
	if err != nil {
		return nil, err
	}

	return sqlResult, nil
}

