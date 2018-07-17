package models

import (
	"database/sql"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/shopspring/decimal"
	"time"
)

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
	query := fmt.Sprintf("SELECT * FROM %v WHERE id=?", i.table)
	err := i.db.Get(contribution, query, id)

	return contribution, err
}

// GetByName returns record by name.
func (i *Contribution) GetByName(tx *sqlx.Tx, name string) (*ContributionRow, error) {
	contribution := &ContributionRow{}
	query := fmt.Sprintf("SELECT * FROM %v WHERE name=?", i.table)
	err := i.db.Get(contribution, query, name)

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
