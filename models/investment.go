package models

import (
	"database/sql"
	"fmt"
	"github.com/jmoiron/sqlx"
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
	ID                               int64     `db:"id"`
	StartupName                      string    `db:"StartupName"`
	Industry                         string    `db:"Industry"`
	InvestmentDate                   time.Time `db:"InvestmentDate"`
	Headquarters                     string    `db:"Headquarters"`
	BoardRepresentation              string    `db:"BoardRepresentation"`
	BoardMembers                     string    `db:"BoardMembers"`
	CapTable                         string    `db:"CapTable"`
	InvestmentBackground             string    `db:"InvestmentBackground"`
	InvestmentThesis                 string    `db:"InvestmentThesis"`
	ExitValueAtClosing               float64       `db:"ExitValueAtClosing"`
	FundOwnershipPercentage          float64       `db:"FundOwnershipPercentage"`
	InvestorGroupPercentage 		 float64       `db:"InvestorGroupPercentage"`
	ManagementOwnership              float64       `db:"ManagementOwnership"`
	InvestmentCommittment            float64       `db:"InvestmentCommittment"`
	InvestedCapital                  float64       `db:"InvestedCapital"`
	RealizedProceeds                 float64       `db:"RealizedProceeds"`
	ReportedValue                    float64       `db:"ReportedValue"`
	InvestmentMultiple               float64       `db:"InvestmentMultiple"`
	GrossIRR                         float64       `db:"GrossIRR"`

}
type News struct {
	ID            int64     `db:"id"`
	Investment_ID int64     `db:"investment_id"`
	NewsDate      time.Time `db:"NewsDate"`
	News          string    `db:"News"`
}

type InvestmentStructure struct {
	ID               int64     `db:"id"`
	Investment_ID    int64     `db:"investment_id"`
	ReportingDate    time.Time `db:"ReportingDate"`
	Units            float64       `db:"Units"`
	TotalInvested    float64       `db:"TotalInvested"`
	ReportedValue    float64       `db:"ReportedValue"`
	RealizedProceeds float64       `db:"RealizedProceeds"`
	Structure        string    `db:"Structure"`
}

type CapitalizationStructure struct {
	ID             int64     `db:"id"`
	Investment_ID  int64     `db:"investment_id"`
	ReportingDate  time.Time `db:"ReportingDate"`
	//MaturityDate   time.Time `db:"maturity_date"`
	ClosingValue   float64       `db:"ClosingValue"`
	YearEndValue   float64       `db:"YearEndValue"`
	Capitalization string    `db:"Capitalization"`
	
}

type FinancialResults struct {
	ID                     int64     `db:"id"`
	Investment_ID          int64     `db:"investment_id"`
	ReportingDate          time.Time `db:"ReportingDate"`
	Revenue                float64       `db:"Revenue"`
	YoYGrowthPercentage1   float64       `db:"YoYGrowthPercentage1"`
	LTMEBITDA              float64       `db:"LTMEBITDA"`
	YoYGrowthPercentage2   float64       `db:"YoYGrowthPercentage2"`
	EBITDAMargin           float64       `db:"EBITDAMargin"`
	TotalExitValue         float64       `db:"TotalExitValue"`
	TotalExitValueMultiple float64       `db:"TotalExitValueMultiple"`
	TotalLeverage          float64       `db:"TotalLeverage"`
	TotalLeverageMultiple  float64       `db:"TotalLeverageMultiple"`
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

// create a new record of user.
// func (i *Investment) Create(tx *sqlx.Tx, name, industry string) (*InvestmentRow, error) {
// 	if name == "" {
// 		return nil, errors.New("Name cannot be blank.")
// 	}

// 	data := make(map[string]interface{})
// 	data["startup_name"] = name
// 	data["industry"] = industry

// 	sqlResult, err := i.InsertIntoTable(tx, data)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return i.userRowFromSqlResult(tx, sqlResult)
// }

// create a new record of user.
func (i *Investment) Create(tx *sqlx.Tx, m map[string]interface{}) (*InvestmentRow, error) {
	sqlResult, err := i.InsertIntoTable(tx, m)
	if err != nil {
		return nil, err
	}
	return i.userRowFromSqlResult(tx, sqlResult)
}

// UpdateEmailAndPasswordById updates user email and password.
func (i *Investment) UpdateById(tx *sqlx.Tx, investmentId int64, name, industry string) (*InvestmentRow, error) {
	data := make(map[string]interface{})

	if name != "" {
		data["startup_name"] = name
	}
	if industry != "" {
		data["industry"] = industry
	}

	if len(data) > 0 {
		_, err := i.UpdateByID(tx, data, investmentId)
		if err != nil {
			return nil, err
		}
	}

	return i.GetById(tx, investmentId)
}
