package models

import (
	"database/sql"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/shopspring/decimal"
	"strconv"
	"time"
)

func NewAssessment(db *sqlx.DB) *Assessment {
	is := &Assessment{}
	is.db = db
	is.table = "assessments"
	is.hasID = true
	return is
}

type Assessment struct {
	Base
}

type AssessmentRow struct {
	ID                         int64           `db:"id"`
	Investment_ID              int64           `db:"investment_id"`
	ReviewDate                 time.Time       `db:"ReviewDate"`
	Status                     string          `db:"Status"`
	RevenueGrowth              string          `db:"RevenueGrowth"`
	Execution                  string          `db:"Execution"`
	Leadership                 string          `db:"Leadership"`
	RevenueBreakEvenPlan       string          `db:"RevenueBreakEvenPlan"`
	KeyGrowthEnablers          string          `db:"KeyGrowthEnablers"`
	PlaybookAdoption           string          `db:"PlaybookAdoption"`
	StartupName                string          `db:"StartupName"`
	MarketMultiple             decimal.Decimal `db:"MarketMultiple"`
	YearThreeForecastedRevenue decimal.Decimal `db:"YearThreeForecastedRevenue"`
	ThreelinesValueAtExit      decimal.Decimal `db:"ThreelinesValueAtExit"`
	YearThreeExitMultiple      decimal.Decimal `db:"YearThreeExitMultiple"`
}

func (i *AssessmentRow) FormattedReviewDate() string {
	return i.ReviewDate.Format("01/02/2006")
}

func (i *Assessment) userRowFromSqlResult(tx *sqlx.Tx, sqlResult sql.Result) (*AssessmentRow, error) {
	isId, err := sqlResult.LastInsertId()
	if err != nil {
		return nil, err
	}

	return i.GetById(tx, isId)
}

// AllUsers returns all user rows.
func (i *Assessment) AllAssessment(tx *sqlx.Tx) ([]*AssessmentRow, error) {
	isrs := []*AssessmentRow{}
	query := fmt.Sprintf("SELECT * FROM %v", i.table)
	err := i.db.Select(&isrs, query)

	return isrs, err
}

// AllUsers returns all user rows.
func (i *Assessment) AllAssessmentByInvestmentId(tx *sqlx.Tx, Investment_ID int64) ([]*AssessmentRow, error) {
	isrs := []*AssessmentRow{}
	query := fmt.Sprintf("SELECT * FROM %v WHERE investment_id=%v", i.table, Investment_ID)
	err := i.db.Select(&isrs, query)

	return isrs, err
}

// GetById returns record by id.
func (i *Assessment) GetById(tx *sqlx.Tx, id int64) (*AssessmentRow, error) {
	isr := &AssessmentRow{}
	query := fmt.Sprintf("SELECT * FROM %v WHERE id=?", i.table)
	err := i.db.Get(isr, query, id)

	return isr, err
}

// GetByName returns record by name.
func (i *Assessment) GetByName(tx *sqlx.Tx, name string) (*AssessmentRow, error) {
	isr := &AssessmentRow{}
	query := fmt.Sprintf("SELECT * FROM %v WHERE StartupName=?", i.table)
	err := i.db.Get(isr, query, name)

	return isr, err
}

// create a new record of user.
func (i *Assessment) Create(tx *sqlx.Tx, m map[string]interface{}) (*AssessmentRow, error) {
	sqlResult, err := i.InsertIntoTable(tx, m)
	if err != nil {
		return nil, err
	}
	return i.userRowFromSqlResult(tx, sqlResult)
}

// UpdateEmailAndPasswordById updates user email and password.
func (i *Assessment) UpdateById(tx *sqlx.Tx, isId int64, data map[string]interface{}) (*AssessmentRow, error) {
	if len(data) > 0 {
		//calling base.go function
		_, err := i.UpdateByID(tx, data, isId)
		if err != nil {
			return nil, err
		}
	}

	return i.GetById(tx, isId)
}

// Get All by Application ID.
func (i *Assessment) GetByInvestmentId(tx *sqlx.Tx, Assessment_ID int64, Investment_ID int64) (*AssessmentRow, error) {
	isr := &AssessmentRow{}
	var query string
	if Assessment_ID == 0 {
		query = fmt.Sprintf("SELECT * FROM %v WHERE Investment_ID=%v", i.table, Investment_ID)
		//fmt.Printf("Executing Query %v", query)
	} else {
		query = fmt.Sprintf("SELECT * FROM %v WHERE ID=%v AND Investment_ID=%v", i.table, Assessment_ID, Investment_ID)
		//fmt.Printf("Executing Query for Existing Notes %v", query)
	}
	err := i.db.Get(isr, query)
	return isr, err
}

// UpdateEmailAndPasswordById updates user email and password.
func (i *Assessment) DeleteByID(tx *sqlx.Tx, csId int64) (sql.Result, error) {

	//calling base.go function
	sqlResult, err := i.DeleteById(tx, csId)
	if err != nil {
		return nil, err
	}

	return sqlResult, nil
}

// Search By CompanyName or Location returns records query.
func (i *Assessment) GetAssessmentsForInvestmentIds(tx *sqlx.Tx, investmentids []int64) ([]*AssessmentRow, error) {
	var query string
	var err error
	isrs := []*AssessmentRow{}
	query = fmt.Sprintf(`SELECT a.* FROM %s a LEFT JOIN investments i ON a.investment_id=i.id WHERE a.investment_id in (`, i.table)
	last := len(investmentids) - 1
	//fmt.Printf("Number of InvestmentIds for Assessment Query %s", len(investmentids))
	for index, id := range investmentids {
		if index == last {
			query += strconv.FormatInt(id, 10) + `)`
		} else {
			query += strconv.FormatInt(id, 10) + `,`
		}
	}
	//fmt.Printf("input query %s", query)
	err = i.db.Select(&isrs, query)
	if err != nil {
		fmt.Println("Search1 Error ", err)
		return nil, err
	}
	//fmt.Printf("Number of Assessments Returned %s", len(isrs))
	return isrs, err
}
