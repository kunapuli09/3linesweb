package models

import (
	"database/sql"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/shopspring/decimal"
	"time"
)

const FUNDI = "3Lines 2016 Discretionary Fund, LLC"
const FUNDII = "3Lines Rocket Fund, L.P"

func NewContribution(db *sqlx.DB) *Contribution {
	contribution := &Contribution{}
	contribution.db = db
	contribution.table = "contributions"
	contribution.hasID = true
	return contribution
}

type Contribution struct {
	Base
}

type ContributionRow struct {
	ID                  int64           `db:"id"`
	User_ID             int64           `db:"user_id"`
	FundLegalName       string          `db:"FundLegalName"`
	InvestorLegalName   string          `db:"InvestorLegalName"`
	InvestorAddress     string          `db:"InvestorAddress"`
	InvestorType        string          `db:"InvestorType"`
	GroupContact        string          `db:"GroupContact"`
	InvestmentGroupName string          `db:"InvestmentGroupName"`
	CommitmentDate      time.Time       `db:"CommitmentDate"`
	OwnershipPercentage decimal.Decimal `db:"OwnershipPercentage"`
	InvestmentAmount    decimal.Decimal `db:"InvestmentAmount"`
	Comments            string          `db:"Comments"`
	Status              string          `db:"Status"`
}

type SearchContribution struct {
	InvestorLegalName string
	FundLegalNames    []string
}

func (i *ContributionRow) FormattedCommitmentDate() string {
	return i.CommitmentDate.Format("01/02/2006")
}

func (i *Contribution) userRowFromSqlResult(tx *sqlx.Tx, sqlResult sql.Result) (*ContributionRow, error) {
	contributionId, err := sqlResult.LastInsertId()
	if err != nil {
		return nil, err
	}

	return i.GetById(tx, contributionId)
}

// AllUsers returns all user rows.
func (i *Contribution) AllContributions(tx *sqlx.Tx) ([]*ContributionRow, error) {
	contributions := []*ContributionRow{}
	query := fmt.Sprintf("SELECT * FROM %v", i.table)
	err := i.db.Select(&contributions, query)

	return contributions, err
}

// GetById returns record by id.
func (i *Contribution) GetById(tx *sqlx.Tx, id int64) (*ContributionRow, error) {
	contribution := &ContributionRow{}
	if id == 0 {
		contribution.CommitmentDate = time.Now().AddDate(0, 0, -3)
		return contribution, nil
	}
	query := fmt.Sprintf("SELECT * FROM %v WHERE id=?", i.table)
	err := i.db.Get(contribution, query, id)

	return contribution, err
}

// Get All by Application ID.
func (i *Contribution) GetAllByContributionIdAndUserId(tx *sqlx.Tx, Contribution_ID int64, User_ID int64) (*ContributionRow, error) {
	isr := &ContributionRow{}
	var query string
	if Contribution_ID == 0 {
		query = fmt.Sprintf("SELECT * FROM %v WHERE User_ID=%v", i.table, User_ID)
		//fmt.Printf("Executing Query %v", query)
	} else {
		query = fmt.Sprintf("SELECT * FROM %v WHERE ID=%v AND User_ID=%v", i.table, Contribution_ID, User_ID)
		//fmt.Printf("Executing Query for Existing Notes %v", query)
	}
	err := i.db.Get(isr, query)
	return isr, err
}

// Get All by Application ID.
func (i *Contribution) GetAllByFundNameAndUserId(tx *sqlx.Tx, FundLegalName string, User_ID int64) ([]*ContributionRow, error) {
	var query string
	contributions := []*ContributionRow{}
	if len(FundLegalName) == 0 || FundLegalName == "" {
		query = fmt.Sprintf("SELECT * FROM %v WHERE User_ID=%v", i.table, User_ID)
		//fmt.Printf("Executing Query %v", query)
	} else {
		query = fmt.Sprintf("SELECT * FROM %v WHERE FundLegalName=%v AND User_ID=%v", i.table, FundLegalName, User_ID)
		//fmt.Printf("Executing Query for Existing Notes %v", query)
	}
	err := i.db.Select(&contributions, query)
	return contributions, err
}

// GetByName returns record by name.
func (i *Contribution) GetByName(tx *sqlx.Tx, FundLegalName string) (*ContributionRow, error) {
	contribution := &ContributionRow{}
	query := fmt.Sprintf("SELECT * FROM %v WHERE FundLegalName=?", i.table)
	err := i.db.Get(contribution, query, FundLegalName)

	return contribution, err
}

// GetByName returns record by name.
func (i *Contribution) GetInvestorNames(tx *sqlx.Tx) ([]*ContributionRow, error) {
	contributions := []*ContributionRow{}
	query := fmt.Sprintf("SELECT * FROM %v", i.table)
	err := i.db.Select(&contributions, query)
	return contributions, err
}

// create a new record of user.
func (i *Contribution) Create(tx *sqlx.Tx, m map[string]interface{}) (*ContributionRow, error) {
	sqlResult, err := i.InsertIntoTable(tx, m)
	if err != nil {
		return nil, err
	}
	return i.userRowFromSqlResult(tx, sqlResult)
}

// UpdateEmailAndPasswordById updates user email and password.
func (i *Contribution) UpdateById(tx *sqlx.Tx, ContributionId int64, data map[string]interface{}) (*ContributionRow, error) {
	if len(data) > 0 {
		//calling base.go function
		_, err := i.UpdateByID(tx, data, ContributionId)
		if err != nil {
			return nil, err
		}
	}

	return i.GetById(tx, ContributionId)
}

// Search By CompanyName or Location returns records query.
func (i *Contribution) SearchContributions(tx *sqlx.Tx, data SearchContribution) ([]*ContributionRow, error) {
	var query string
	var err error
	isrs := []*ContributionRow{}
	investorName := data.InvestorLegalName
	if len(data.FundLegalNames) > 0 {
		if len(investorName) > 0 {
			query = fmt.Sprintf(`SELECT a.* FROM %s a LEFT JOIN users u ON a.user_id=u.id WHERE a.InvestorLegalName Like '%%%s%%' AND a.FundLegalName in (`, i.table, investorName)
		} else {
			query = fmt.Sprintf(`SELECT a.* FROM %s a LEFT JOIN users u ON a.user_id=u.id WHERE a.FundLegalName in (`, i.table)
		}
		last := len(data.FundLegalNames) - 1
		for index, fundName := range data.FundLegalNames {
			if index == last {
				query += `'` + fundName + `')`
			} else {
				query += `'` + fundName + `',`
			}
		}
		fmt.Printf("input query %s", query)
		err = i.db.Select(&isrs, query)
		if err != nil {
			fmt.Println("Contribution Search Error ", err)
			return nil, err
		}
		return isrs, err
	} else {
		if len(investorName) > 0 {
			query = fmt.Sprintf(`SELECT * FROM %s WHERE InvestorLegalName Like '%%%s%%'`, i.table, investorName)
			//			fmt.Printf("input query %s for investor %s", query, investorName)
			err = i.db.Select(&isrs, query)
			if err != nil {
				fmt.Println("Search2 Error", err)
				return nil, err
			}
			return isrs, err
		}

	}
	return i.AllContributions(tx)

}
