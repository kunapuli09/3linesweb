package models

import (
	"database/sql"
	"fmt"
	"github.com/jmoiron/sqlx"
	"time"
)

func NewNews(db *sqlx.DB) *News {
	n := &News{}
	n.db = db
	n.table = "news"
	n.hasID = true
	return n
}

type News struct {
	Base
}

type NewsRow struct {
	ID            int64     `db:"id"`
	Investment_ID int64     `db:"investment_id"`
	NewsDate      time.Time `db:"NewsDate"`
	News          string    `db:"News"`
}

type Notification struct {
	ID            int64     `db:"id"`
	Investment_ID int64     `db:"investment_id"`
	StartupName   string    `db:"StartupName"`
	Industry      string    `db:"Industry"`
	NewsDate      time.Time `db:"NewsDate"`
	News          string    `db:"News"`
}

func (n *Notification) FormattedNewsDate() string {
	return n.NewsDate.Format("01/02/2006")
}

func (n *NewsRow) FormattedNewsDate() string {
	return n.NewsDate.Format("01/02/2006")
}

func (i *News) userRowFromSqlResult(tx *sqlx.Tx, sqlResult sql.Result) (*NewsRow, error) {
	nId, err := sqlResult.LastInsertId()
	if err != nil {
		return nil, err
	}

	return i.GetById(tx, nId)
}

// AllUsers returns all user rows.
func (i *News) AllNews(tx *sqlx.Tx) ([]*NewsRow, error) {
	nrs := []*NewsRow{}
	query := fmt.Sprintf("SELECT * FROM %v", i.table)
	err := i.db.Select(&nrs, query)

	return nrs, err
}

func (n *News) AllNotifications(tx *sqlx.Tx) ([]*Notification, error) {
	notifications := []*Notification{}
	q := `SELECT investments.id, 
		investments.StartupName, 
		investments.Industry, 
		news.investment_id, news.NewsDate, news.News
  		FROM %v
  		INNER JOIN investments
    	ON LOWER(investments.id) = LOWER(news.investment_id)`
	query := fmt.Sprintf(q, n.table)
	err := n.db.Select(&notifications, query)
	return notifications, err
}

// GetById returns record by id.
func (i *News) GetById(tx *sqlx.Tx, id int64) (*NewsRow, error) {
	n := &NewsRow{}
	query := fmt.Sprintf("SELECT * FROM %v WHERE id=?", i.table)
	err := i.db.Get(n, query, id)

	return n, err
}

// GetByName returns record by name.
func (i *News) GetByName(tx *sqlx.Tx, name string) (*NewsRow, error) {
	n := &NewsRow{}
	query := fmt.Sprintf("SELECT * FROM %v WHERE name=?", i.table)
	err := i.db.Get(n, query, name)

	return n, err
}

// create a new record of user.
func (i *News) Create(tx *sqlx.Tx, m map[string]interface{}) (*NewsRow, error) {
	sqlResult, err := i.InsertIntoTable(tx, m)
	if err != nil {
		return nil, err
	}
	return i.userRowFromSqlResult(tx, sqlResult)
}

// UpdateEmailAndPasswordById updates user email and password.
func (i *News) UpdateById(tx *sqlx.Tx, nId int64, data map[string]interface{}) (*NewsRow, error) {
	if len(data) > 0 {
		//calling base.go function
		_, err := i.UpdateByID(tx, data, nId)
		if err != nil {
			return nil, err
		}
	}

	return i.GetById(tx, nId)
}

// Get All by Investment ID.
func (i *News) GetAllByInvestmentId(tx *sqlx.Tx, Investment_ID int64) ([]*NewsRow, error) {
	css := []*NewsRow{}
	query := fmt.Sprintf("SELECT * FROM %v WHERE Investment_ID=%v", i.table, Investment_ID)
	err := i.db.Select(&css, query)
	return css, err
}

// UpdateEmailAndPasswordById updates user email and password.
func (i *News) DeleteByID(tx *sqlx.Tx, csId int64) (sql.Result, error) {

	//calling base.go function
	sqlResult, err := i.DeleteById(tx, csId)
	if err != nil {
		return nil, err
	}

	return sqlResult, nil
}
