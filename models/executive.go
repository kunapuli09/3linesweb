package models

import (
	"database/sql"
	"fmt"
	"github.com/jmoiron/sqlx"
	"time"
)

func NewExecutive(db *sqlx.DB) *Executive {
	Executive := &Executive{}
	Executive.db = db
	Executive.table = "executives"
	Executive.hasID = true

	return Executive
}

type ExecutiveRow struct {
	ID              int64           `db:"id"`
	ApplicationDate time.Time       `db:"ApplicationDate"`
	Name       		string          `db:"Name"`
	Email           string          `db:"Email"`
	SocialMediaHandle        string `db:"SocialMediaHandle"`
	//Google Captcha Response Field Required for Form Parsing
	//Not storing this field in database
	rcres string
}

type Executive struct {
	Base
}

func (ar *ExecutiveRow) FormattedApplicationDate() string {
	return ar.ApplicationDate.Format("01/02/2006")
}

func (i *Executive) userRowFromSqlResult(tx *sqlx.Tx, sqlResult sql.Result) (*ExecutiveRow, error) {
	isId, err := sqlResult.LastInsertId()
	if err != nil {
		return nil, err
	}

	return i.GetById(tx, isId)
}

// GetById returns record by id.
func (i *Executive) GetById(tx *sqlx.Tx, id int64) (*ExecutiveRow, error) {
	isr := &ExecutiveRow{}
	query := fmt.Sprintf("SELECT * FROM %v WHERE id=?", i.table)
	err := i.db.Get(isr, query, id)
	return isr, err
}

// AllUsers returns all user rows.
func (i *Executive) AllExecutives(tx *sqlx.Tx) ([]*ExecutiveRow, error) {
	isrs := []*ExecutiveRow{}
	query := fmt.Sprintf("SELECT * FROM %v ORDER BY ApplicationDate DESC", i.table)
	err := i.db.Select(&isrs, query)
	if err != nil {
		fmt.Println(err)
	}
	return isrs, err
}

// GetByName returns record by name.
func (i *Executive) GetExisting(tx *sqlx.Tx, email string, website string, companyname string) bool {
	var count int
	query := fmt.Sprintf("SELECT COUNT(*) FROM %v WHERE Email=? OR SocialMediaHandle=? ORDER BY ApplicationDate DESC", i.table)
	err := i.db.Get(&count, query, email, website, companyname)
	if err != nil {
		fmt.Println("Existing Executive Search Error %v", err)
	}
	if count > 0 {
		return true
	}
	return false
}

// create a new record of user.
func (i *Executive) Create(tx *sqlx.Tx, m map[string]interface{}) (*ExecutiveRow, error) {
	sqlResult, err := i.InsertIntoTable(tx, m)
	if err != nil {
		return nil, err
	}
	return i.userRowFromSqlResult(tx, sqlResult)
}

// UpdateEmailAndPasswordById updates user email and password.
func (i *Executive) UpdateById(tx *sqlx.Tx, isId int64, data map[string]interface{}) (*ExecutiveRow, error) {
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
func (i *Executive) DeleteByID(tx *sqlx.Tx, csId int64) (sql.Result, error) {

	//calling base.go function
	sqlResult, err := i.DeleteById(tx, csId)
	if err != nil {
		return nil, err
	}

	return sqlResult, nil
}
