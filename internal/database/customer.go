package database

import "log"

type Customer struct {
	Id               int
	Email            string
	VeriffIdvSessionId  string
	VeriffIdvSessionUrl string
	VeriffPoaSessionId  string
	VeriffPoaSessionUrl string
}

func CreateCustomerTable() {
	_, err := CreateTable(`
		create table if not exists customer (
			id serial primary key,
			email text not null,
			veriffIdvSessionId text not null,
			veriffIdvSessionUrl text not null,
			veriffPoaSessionId text not null,
			veriffPoaSessionUrl text not null
		);
	`)

	if err != nil {
		log.Fatalf("Couldn't create customer table! %v", err)
	}
}

func CreateCustomer(
	email string, 
	veriffIdvSessionId string, 
	veriffIdvSessionUrl string,
	veriffPoaSessionId string,
	veriffPoaSessionUrl string,
) error {
	return Exec(`
		insert into customer (email, veriffIdvSessionId, veriffIdvSessionUrl, veriffPoaSessionId, veriffPoaSessionUrl)
		values ($1, $2, $3, $4, $5)
	`,
		email,
		veriffIdvSessionId,
		veriffIdvSessionUrl,
		veriffPoaSessionId,
		veriffPoaSessionUrl,
	)
}
