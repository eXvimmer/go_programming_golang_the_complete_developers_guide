package mdb

import (
	"database/sql"
	"github.com/mattn/go-sqlite3"
	"log"
	"time"
)

type EmailEntry struct {
	Id          int64
	Email       string
	ConfirmedAt *time.Time
	OptOut      bool
}

type GetEmailBatchQueryparams struct {
	Page  int
	Count int
}

func TryCreate(db *sql.DB) {
	_, err := db.Exec(`
    CREATE TABLE emails (
      id           INTEGER PRIMARY KEY,
      email        TEXT UNIQUE,
      confirmed_at INTEGER,
      opt_out      INTEGER
    );
  `)

	// NOTE: casting err type to sqlite3.Error type
	if sqlError, ok := err.(sqlite3.Error); ok {
		// NOTE: 1 means table already exists
		if sqlError.Code != 1 {
			log.Fatal(sqlError)
		}
	} else {
		log.Fatal(err)
	}
}

func EmailEntryFromRow(row *sql.Rows) (*EmailEntry, error) {
	var (
		id          int64
		email       string
		confirmedAt int64
		optOut      bool
	)

	err := row.Scan(&id, &email, &confirmedAt, &optOut)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	t := time.Unix(confirmedAt, 0)

	return &EmailEntry{
		Id:          id,
		Email:       email,
		ConfirmedAt: &t,
		OptOut:      optOut,
	}, nil
}

func CreateEmail(db *sql.DB, email string) error {
	_, err := db.Exec(
		`
    INSERT INTO emails (email, confirmed_at, opt_out)
    VALUES (?, 0, false);
    `,
		email,
	)

	if err != nil {
		log.Panicln(err)
		return err
	}

	return nil
}

func GetEmail(db *sql.DB, email string) (*EmailEntry, error) {
	rows, err := db.Query(`
    SELECT id, email, confirmed_at, opt_out
    FROM emails
    WHERE email = ?
  `, email)

	defer rows.Close()

	if err != nil {
		log.Println(err)
		return nil, err
	}

	for rows.Next() {
		return EmailEntryFromRow(rows)
	}

	return nil, nil
}

func UpdateEmail(db *sql.DB, entry EmailEntry) error {
	t := entry.ConfirmedAt.Unix()

	// NOTE: up-sert operation
	_, err := db.Exec(`
    INSERT INTO emails (email, confirmed_at, opt_out)
    VALUES (?, ?, ?)
    ON CONFLICT (email) DO UPDATE SET
      confirmed_at = ?
      opt_out = ?
  `, entry.Email, t, entry.OptOut, t, entry.OptOut)

	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

func DeleteEmail(db *sql.DB, email string) error {
	_, err := db.Exec(`
    UPDATE emails
    SET opt_out = true
    WHERE email = ?
  `, email)

	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

func GetEmailBatch(db *sql.DB, params GetEmailBatchQueryparams) ([]EmailEntry, error) {
	var empty []EmailEntry

	rows, err := db.Query(`
    SELECT id, email, confirmed_at, opt_out
    FROM emails
    WHERE opt_out = false
    ORDER BY id ASC
    LIMIT ?
    OFFSET ?
  `,
		params.Count, (params.Page-1)*params.Count,
	)

	if err != nil {
		log.Println(err)
		return empty, err
	}
	defer rows.Close()

	emails := make([]EmailEntry, 0, params.Count)

	for rows.Next() {
		email, err := EmailEntryFromRow(rows)
		if err != nil {
			return nil, err
		}
		emails = append(emails, *email)
	}

	return emails, nil
}
