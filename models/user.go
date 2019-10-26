package models

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
	"golang.org/x/crypto/bcrypt"
)

func NewUser(db *sqlx.DB) *User {
	user := &User{}
	user.db = db
	user.table = "users"
	user.hasID = true
	return user
}

type UserRow struct {
	ID       int64  `db:"id"`
	Email    string `db:"email"`
	Phone    string `db:"phone"`
	Password string `db:"password"`
	Admin    bool   `db:"admin"`
	Dsc      bool
	FundOne  bool
	FundTwo  bool
}

type User struct {
	Base
}

func (u *User) userRowFromSqlResult(tx *sqlx.Tx, sqlResult sql.Result) (*UserRow, error) {
	userId, err := sqlResult.LastInsertId()
	if err != nil {
		return nil, err
	}

	return u.GetById(tx, userId)
}

// AllUsers returns all user rows.
func (u *User) AllUsers(tx *sqlx.Tx) ([]*UserRow, error) {
	users := []*UserRow{}
	query := fmt.Sprintf("SELECT * FROM %v", u.table)
	err := u.db.Select(&users, query)
	return users, err
}

// GetById returns record by id.
func (u *User) GetById(tx *sqlx.Tx, id int64) (*UserRow, error) {
	user := &UserRow{}
	query := fmt.Sprintf("SELECT * FROM %v WHERE id=?", u.table)
	err := u.db.Get(user, query, id)
	user.Dsc = isDsc(user.Email)
	user.FundOne = isFundI(user.ID, u.db)
	user.FundTwo = isFundII(user.ID, u.db)
	return user, err
}

// GetByEmail returns record by email.
func (u *User) GetByEmail(tx *sqlx.Tx, email string) (*UserRow, error) {
	user := &UserRow{}
	query := fmt.Sprintf("SELECT *  FROM %v WHERE email=?", u.table)
	err := u.db.Get(user, query, email)
	user.Dsc = isDsc(user.Email)
	user.FundOne = isFundI(user.ID, u.db)
	user.FundTwo = isFundII(user.ID, u.db)
	return user, err
}

// GetByEmail returns record by email but checks password first.
func (u *User) GetUserByEmailAndPassword(tx *sqlx.Tx, email, password string) (*UserRow, error) {
	user, err := u.GetByEmail(tx, email)
	if err != nil {
		return nil, err
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return nil, err
	}
	user.Dsc = isDsc(user.Email)
	user.FundOne = isFundI(user.ID, u.db)
	user.FundTwo = isFundII(user.ID, u.db)
	return user, err
}

// Signup create a new record of user.
func (u *User) Signup(tx *sqlx.Tx, email, password, passwordAgain string, phone string) (*UserRow, error) {
	if email == "" {
		return nil, errors.New("Email cannot be blank.")
	}
	if phone == "" {
		return nil, errors.New("Phone is invalid.")
	}
	if password == "" {
		return nil, errors.New("Password cannot be blank.")
	}
	if password != passwordAgain {
		return nil, errors.New("Password is invalid.")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 5)
	if err != nil {
		return nil, err
	}

	data := make(map[string]interface{})
	data["email"] = email
	data["password"] = hashedPassword
	data["phone"] = phone
	data["admin"] = 0

	sqlResult, err := u.InsertIntoTable(tx, data)
	if err != nil {
		return nil, err
	}

	return u.userRowFromSqlResult(tx, sqlResult)
}

// UpdateEmailAndPasswordById updates user email and password.
func (u *User) UpdateEmailAndPasswordById(tx *sqlx.Tx, userId int64, email, password, passwordAgain string, phone string) (*UserRow, error) {
	data := make(map[string]interface{})

	if email != "" {
		data["email"] = email
	}
	if phone != "" {
		data["phone"] = phone
	}

	if password != "" && passwordAgain != "" && password == passwordAgain {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 5)
		if err != nil {
			return nil, err
		}

		data["password"] = hashedPassword
	}

	if len(data) > 0 {
		_, err := u.UpdateByID(tx, data, userId)
		if err != nil {
			return nil, err
		}
	}

	return u.GetById(tx, userId)
}

// AllUsers returns all user rows.
func (u *User) AllEmails(tx *sqlx.Tx) ([]string, error) {
	emails := []string{}
	query := fmt.Sprintf("SELECT email FROM %v", u.table)
	err := u.db.Select(&emails, query)

	return emails, err
}

// UpdateEmailAndPasswordById updates user email and password.
func (i *User) DeleteByID(tx *sqlx.Tx, csId int64) (sql.Result, error) {

	//calling base.go function
	sqlResult, err := i.DeleteById(tx, csId)
	if err != nil {
		return nil, err
	}

	return sqlResult, nil
}
func isDsc(a string) bool {
	//hardcode roles temporarily
	dsc_team := []string{
		"rmaddhi@gmail.com",
		"roy.rajiv@gmail.com",
		"arun.taman@gmail.com",
		"sgosala99@gmail.com",
		"dsc@3lines.vc",
	}
	for _, b := range dsc_team {
		if b == a {
			return true
		}
	}
	return false
}


//big hack..******fix this crap

func SplitContributions(Contributions []*ContributionRow) ([]*ContributionRow, []*ContributionRow) {
	var fundone []*ContributionRow
	var fundtwo []*ContributionRow

	for _, contribution := range Contributions {
		switch fundName := contribution.FundLegalName; fundName {
		case FUNDI:
			fundone = append(fundone, contribution)
		case FUNDII:
			fundtwo = append(fundtwo, contribution)

		default:
			//fmt.Printf("%s. is unknown investor type", fundName)
		}
	}
	return fundone, fundtwo
}



func isFundI(User_ID int64, db *sqlx.DB) bool {
	contributions, _ := NewContribution(db).AllContributions(nil)
	fundone_insvestors, _ := SplitContributions(contributions)
	for _, b := range fundone_insvestors {
		if b.User_ID == User_ID {
			return true
		}
	}
	return false
}
func isFundII(User_ID int64, db *sqlx.DB) bool {
	contributions, _ := NewContribution(db).AllContributions(nil)
	_,fundtwo_insvestors := SplitContributions(contributions)
	for _, b := range fundtwo_insvestors {
		if b.User_ID == User_ID {
			return true
		}
	}
	return false
}
