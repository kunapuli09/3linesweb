package models

import (
	"database/sql"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/shopspring/decimal"
	"strings"
	"time"
)

func NewAppl(db *sqlx.DB) *Appl {
	Appl := &Appl{}
	Appl.db = db
	Appl.table = "applications"
	Appl.hasID = true

	return Appl
}

type ApplRow struct {
	ID              int64           `db:"id"`
	ApplicationDate time.Time       `db:"ApplicationDate"`
	Email           string          `db:"Email"`
	FirstName       string          `db:"FirstName"`
	LastName        string          `db:"LastName"`
	CompanyName     string          `db:"CompanyName"`
	Phone           string          `db:"Phone"`
	Website         string          `db:"Website"`
	Title           sql.NullString  `db:"Title"`
	Referrer        string          `db:"Referrer"`
	Industries      string          `db:"Industries"`
	Locations       string          `db:"Locations"`
	Revenue         string          `db:"Revenue"`
	Comments        string          `db:"Comments"`
	ElevatorPitch   string          `db:"ElevatorPitch"`
	CapitalRaised   decimal.Decimal `db:"CapitalRaised"`
}

type Appl struct {
	Base
}

type Search struct {
	CompanyName string
	Location    string
	Status      []string
}

func (ar *ApplRow) FormattedApplicationDate() string {
	return ar.ApplicationDate.Format("01/02/2006")
}

func (i *Appl) userRowFromSqlResult(tx *sqlx.Tx, sqlResult sql.Result) (*ApplRow, error) {
	isId, err := sqlResult.LastInsertId()
	if err != nil {
		return nil, err
	}

	return i.GetById(tx, isId)
}

// GetById returns record by id.
func (i *Appl) GetById(tx *sqlx.Tx, id int64) (*ApplRow, error) {
	isr := &ApplRow{}
	query := fmt.Sprintf("SELECT * FROM %v WHERE id=?", i.table)
	err := i.db.Get(isr, query, id)
	return isr, err
}

// AllUsers returns all user rows.
func (i *Appl) AllAppls(tx *sqlx.Tx) ([]*ApplRow, error) {
	isrs := []*ApplRow{}
	query := fmt.Sprintf("SELECT * FROM %v ORDER BY ApplicationDate DESC", i.table)
	err := i.db.Select(&isrs, query)
	if err != nil {
		fmt.Println(err)
	}
	return isrs, err
}

// Search By CompanyName or Location returns records query.
func (i *Appl) Search(tx *sqlx.Tx, data Search) ([]*ApplRow, error) {
	var query string
	var err error
	isrs := []*ApplRow{}
	companyName := data.CompanyName
	location := data.Location
	statuses := strings.Join(data.Status, ",")

	if len(data.Status) > 0 {
		query = fmt.Sprintf("SELECT a.* FROM %v a LEFT JOIN screeningnotes s ON a.id=s.application_id WHERE a.CompanyName Like ? AND a.Locations Like ? AND s.Status in (?) ORDER BY ApplicationDate DESC", i.table)
		err = i.db.Select(&isrs, query, location+"%", companyName+"%", statuses)
		if err != nil {
			fmt.Println("Search1 Error %v", err)
			return nil, err
		}
		return isrs, err
	} else {
		query = fmt.Sprintf("SELECT * FROM %v WHERE Locations Like ? AND CompanyName Like ? ORDER BY ApplicationDate DESC", i.table)
		err = i.db.Select(&isrs, query, location+"%", companyName+"%")
		if err != nil {
			fmt.Println("Search2 Error %v", err)
			return nil, err
		}
		return isrs, err
	}
	return i.AllAppls(tx)

}

// GetByName returns record by name.
func (i *Appl) GetExisting(tx *sqlx.Tx, email string, website string, companyname string) (bool) {
	var count int
	query := fmt.Sprintf("SELECT COUNT(*) FROM %v WHERE Email=? OR Website=? OR CompanyName=? ORDER BY ApplicationDate DESC", i.table)
	err := i.db.Get(&count, query, email, website, companyname)
	if err != nil {
		fmt.Println("Existing Application Search Error %v", err)
	}
	if count > 0 {
		return true
	}
	return false
}

// create a new record of user.
func (i *Appl) Create(tx *sqlx.Tx, m map[string]interface{}) (*ApplRow, error) {
	sqlResult, err := i.InsertIntoTable(tx, m)
	if err != nil {
		return nil, err
	}
	return i.userRowFromSqlResult(tx, sqlResult)
}

// UpdateEmailAndPasswordById updates user email and password.
func (i *Appl) UpdateById(tx *sqlx.Tx, isId int64, data map[string]interface{}) (*ApplRow, error) {
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
func (i *Appl) DeleteByID(tx *sqlx.Tx, csId int64) (sql.Result, error) {

	//calling base.go function
	sqlResult, err := i.DeleteById(tx, csId)
	if err != nil {
		return nil, err
	}

	return sqlResult, nil
}
