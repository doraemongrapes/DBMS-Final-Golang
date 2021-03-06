package queryMySQL

import (
	daoMySQL "dbmsfinal/DAOMySQL"
	dataMySQL "dbmsfinal/dataMySQL"
	"os"
	"time"

	gocsv "github.com/gocarina/gocsv"
)

func getAllCustomers() []*dataMySQL.Customer {
	customersFile, err := os.OpenFile("customers.csv", os.O_RDWR|os.O_CREATE, os.ModePerm)
	if err != nil {
		panic(err)
	}
	defer customersFile.Close()

	customers := []*dataMySQL.Customer{}

	if err := gocsv.UnmarshalFile(customersFile, &customers); err != nil { // Load clients from file
		panic(err)
	}

	return customers
}

func insertCustomerRowByRow() int64 {
	customers := getAllCustomers()

	start := time.Now()
	sqlStr := "INSERT INTO Customers(customer_id, customer_name, customer_email, customer_address) VALUES (?, ?, ?, ?)"
	stmt, _ := daoMySQL.DB.Prepare(sqlStr)
	for _, customer := range customers {
		_, err := stmt.Exec(customer.CustomerID, customer.CustomerName, customer.CustomerEmail, customer.CustomerAddress)
		if err != nil {
			panic(err)
		}
	}
	elapsed := time.Since(start).Nanoseconds()

	return elapsed
}

func insertCustomersBulk() int64 {
	customers := getAllCustomers()

	start := time.Now()

	sqlStr := "INSERT INTO Customers(customer_id, customer_name, customer_email, customer_address) VALUES "
	vals := []interface{}{}
	for _, customer := range customers {
		sqlStr += "(?, ?, ?, ?),"
		vals = append(vals, customer.CustomerID, customer.CustomerName, customer.CustomerEmail, customer.CustomerAddress)
	}

	sqlStr = sqlStr[0 : len(sqlStr)-1]
	stmt, _ := daoMySQL.DB.Prepare(sqlStr)
	_, err := stmt.Exec(vals...)
	if err != nil {
		panic(err)
	}

	elapsed := time.Since(start).Nanoseconds()

	return elapsed
}
