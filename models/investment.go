package models

import (
	"database/sql"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/shopspring/decimal"
	"time"
)

func NewInvestment(db *sqlx.DB) *Investment {
	investment := &Investment{}
	investment.db = db
	investment.table = "investments"
	investment.hasID = true
	return investment
}

type Investment struct {
	Base
}

type InvestmentRow struct {
	ID                      int64           `db:"id"`
	StartupName             string          `db:"StartupName"`
	LogoPath                string          `db:"LogoPath"`
	Website                 string          `db:"Website"`
	Description             string          `db:"Description"`
	Team                    string          `db:"Team"`
	Industry                string          `db:"Industry"`
	InvestmentDate          time.Time       `db:"InvestmentDate"`
	Headquarters            string          `db:"Headquarters"`
	BoardRepresentation     string          `db:"BoardRepresentation"`
	BoardMembers            string          `db:"BoardMembers"`
	CapTable                string          `db:"CapTable"`
	InvestmentBackground    string          `db:"InvestmentBackground"`
	InvestmentThesis        string          `db:"InvestmentThesis"`
	ValuationMethodology    string          `db:"ValuationMethodology"`
	RiskAssessment          string          `db:"RiskAssessment"`
	ExitValueAtClosing      decimal.Decimal `db:"ExitValueAtClosing"`
	FundOwnershipPercentage decimal.Decimal `db:"FundOwnershipPercentage"`
	InvestorGroupPercentage decimal.Decimal `db:"InvestorGroupPercentage"`
	ManagementOwnership     decimal.Decimal `db:"ManagementOwnership"`
	InvestmentCommittment   decimal.Decimal `db:"InvestmentCommittment"`
	InvestedCapital         decimal.Decimal `db:"InvestedCapital"`
	RealizedProceeds        decimal.Decimal `db:"RealizedProceeds"`
	ReportedValue           decimal.Decimal `db:"ReportedValue"`
	InvestmentMultiple      decimal.Decimal `db:"InvestmentMultiple"`
	GrossIRR                decimal.Decimal `db:"GrossIRR"`
}

func (i *InvestmentRow) FormattedInvestmentDate() string {
	return i.InvestmentDate.Format("01/02/2006")
}

func (i *Investment) userRowFromSqlResult(tx *sqlx.Tx, sqlResult sql.Result) (*InvestmentRow, error) {
	investmentId, err := sqlResult.LastInsertId()
	if err != nil {
		return nil, err
	}

	return i.GetById(tx, investmentId)
}

// AllUsers returns all user rows.
func (i *Investment) AllInvestments(tx *sqlx.Tx) ([]*InvestmentRow, error) {
	investments := []*InvestmentRow{}
	query := fmt.Sprintf("SELECT * FROM %v", i.table)
	err := i.db.Select(&investments, query)

	return investments, err
}

// GetById returns record by id.
func (i *Investment) GetById(tx *sqlx.Tx, id int64) (*InvestmentRow, error) {
	investment := &InvestmentRow{}
	query := fmt.Sprintf("SELECT * FROM %v WHERE id=?", i.table)
	err := i.db.Get(investment, query, id)

	return investment, err
}

// GetByName returns record by name.
func (i *Investment) GetByName(tx *sqlx.Tx, name string) (*InvestmentRow, error) {
	investment := &InvestmentRow{}
	query := fmt.Sprintf("SELECT * FROM %v WHERE name=?", i.table)
	err := i.db.Get(investment, query, name)

	return investment, err
}

// GetByName returns record by name.
func (i *Investment) GetStartupNames(tx *sqlx.Tx) ([]*InvestmentRow, error) {
	investments := []*InvestmentRow{}
	query := fmt.Sprintf("SELECT * FROM %v", i.table)
	err := i.db.Select(&investments, query)
	return investments, err
}

// create a new record of user.
func (i *Investment) Create(tx *sqlx.Tx, m map[string]interface{}) (*InvestmentRow, error) {
	sqlResult, err := i.InsertIntoTable(tx, m)
	if err != nil {
		return nil, err
	}
	return i.userRowFromSqlResult(tx, sqlResult)
}

// UpdateEmailAndPasswordById updates user email and password.
func (i *Investment) UpdateById(tx *sqlx.Tx, investmentId int64, data map[string]interface{}) (*InvestmentRow, error) {
	if len(data) > 0 {
		//calling base.go function
		_, err := i.UpdateByID(tx, data, investmentId)
		if err != nil {
			return nil, err
		}
	}

	return i.GetById(tx, investmentId)
}
