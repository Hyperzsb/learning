package demo

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
)

func printRows(db *sql.DB) error {
	queryStmt := `select * from demo`
	rows, err := db.Query(queryStmt)
	if err != nil {
		return err
	}
	defer rows.Close()

	rowCount := 0
	for rows.Next() {
		rowCount++

		var id int
		var number string

		err = rows.Scan(&id, &number)
		if err != nil {
			return err
		}

		fmt.Println(id, number)
	}
	fmt.Println("Row count: ", rowCount)

	return nil
}

func PostgresDemo() error {
	const (
		host     = "localhost"
		port     = 5432
		user     = "postgres"
		password = "postgres"
		dbname   = "mydb"
	)

	conn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

	// Connect data source
	db, err := sql.Open("postgres", conn)
	if err != nil {
		return err
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		return err
	}

	// Insert records
	insertStmt := `insert into demo (id, number) values (1, '123456')`
	_, err = db.Exec(insertStmt)
	if err != nil {
		return err
	}

	insertStmt = `insert into demo (id, number) values ($1, $2)`
	_, err = db.Exec(insertStmt, 2, "234567")
	if err != nil {
		return err
	}

	// Query or select records
	err = printRows(db)
	if err != nil {
		return err
	}

	// Update records
	updateStmt := `update demo set number=$1 where id=$2`
	_, err = db.Exec(updateStmt, "234567", 1)
	if err != nil {
		return err
	}

	// Delete records
	deleteStmt := `delete from demo where id=$1 or id=$2`
	_, err = db.Exec(deleteStmt, 1, 2)
	if err != nil {
		return err
	}

	err = printRows(db)
	if err != nil {
		return err
	}

	return nil
}
