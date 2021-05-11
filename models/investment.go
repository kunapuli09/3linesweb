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
	FundLegalName           string          `db:"FundLegalName"`
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
	TotalCapitalRaised      decimal.Decimal `db:"TotalCapitalRaised"`
	RealizedProceeds        decimal.Decimal `db:"RealizedProceeds"`
	ReportedValue           decimal.Decimal `db:"ReportedValue"`
	InvestmentMultiple      decimal.Decimal `db:"InvestmentMultiple"`
	GrossIRR                decimal.Decimal `db:"GrossIRR"`
	Status                  string          `db:"Status"`
}

type RevenueSummary struct {
	ID                 int64           `db:"id"`
	StartupName        string          `db:"StartupName"`
	TotalCapitalRaised decimal.Decimal `db:"TotalCapitalRaised"`
	InvestmentMultiple decimal.Decimal `db:"InvestmentMultiple"`
	ReportingDate      time.Time       `db:"ReportingDate"`
	Revenue            decimal.Decimal `db:"Revenue"`
	EBIDTA             decimal.Decimal `db:"LTMEBITDA"`
}

type RevenueDisplay struct {
	ID                                   int64
	StartupName                          string
	TotalCapitalRaised                   decimal.Decimal
	InvestmentMultiple                   decimal.Decimal
	LastYearEBIDTA                       decimal.Decimal
	ForecastedEBIDTA                     decimal.Decimal
	LastYearRevenue                      decimal.Decimal
	ForecastedRevenue                    decimal.Decimal
	LastYearRevenueToCapital             decimal.Decimal
	ForecastedRevenueToCapital           decimal.Decimal
	IsLastYearEBIDTANegative             bool
	IsLastYearRevenueNegative            bool
	IsForecastedEBIDTANegative           bool
	IsForecastedRevenueNegative          bool
	IsLastYearRevenueToCapitalNegative   bool
	IsForecastedRevenueToCapitalNegative bool
}

func (i *InvestmentRow) FormattedInvestmentDate() string {
	return i.InvestmentDate.Format("01/02/2006")
}

func (r *RevenueSummary) FormattedReportingDate() string {
	return r.ReportingDate.Format("01/02/2006")
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
	if id == 0 {
		investment.InvestmentDate = time.Now().AddDate(0, 0, -3)
		return investment, nil
	}
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
func (i *Investment) GetRevenueSummary(tx *sqlx.Tx) ([]*RevenueSummary, error) {
	revenues := []*RevenueSummary{}
	query := "SELECT i.id, i.StartupName,i.TotalCapitalRaised, i.InvestmentMultiple, fr.Revenue, fr.ReportingDate,fr.LTMEBITDA FROM investments AS i INNER JOIN financial_results AS fr ON i.ID = fr.Investment_ID"
	err := i.db.Select(&revenues, query)
	return revenues, err
}

// GetByName returns record by name.
func (i *Investment) GetAllInvestmentsWithoutSyndicates(tx *sqlx.Tx) ([]*InvestmentRow, error) {
	investments := []*InvestmentRow{}
	query := fmt.Sprintf("SELECT * FROM %v WHERE StartupName NOT LIKE '%%%s%%'", i.table, "LLC")
	err := i.db.Select(&investments, query)
	return investments, err
}

// GetByName returns record by name.
func (i *Investment) GetStartupNames(tx *sqlx.Tx) ([]*InvestmentRow, error) {
	investments := []*InvestmentRow{}
	query := fmt.Sprintf("SELECT * FROM %v", i.table)
	err := i.db.Select(&investments, query)
	return investments, err
}

// GetByName returns record by name.
func (i *Investment) GetUserInvestments(tx *sqlx.Tx, partcipatedFundNames []string) ([]*InvestmentRow, error) {
	investments := []*InvestmentRow{}
	query := `SELECT * FROM investments WHERE FundLegalName in (`
	last := len(partcipatedFundNames) - 1
	for index, fundName := range partcipatedFundNames {
		if index == last {
			query += `'` + fundName + `')`
		} else {
			query += `'` + fundName + `',`
		}
	}
	query += " ORDER BY Status ASC"
	//fmt.Printf("input query %v", query)
	err := i.db.Select(&investments, query)
	if err != nil {
		fmt.Println("GetUserInvestments Query Error %v", err)
		return nil, err
	}
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
