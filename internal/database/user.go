package database

import "log"

type Customer struct {
	Id               int
	Email            string
	VeriffSessionId  string
	VeriffSessionUrl string
}

func CreateUserTable() {
	_, err := CreateTable(`
		create table if not exists customer (
			id serial primary key,
			email text not null,
			veriffSessionId text not null,
			veriffSessionUrl text not null
		);
	`)

	if err != nil {
		log.Fatalf("Couldn't create customer table! %v", err)
	}
}

func CreateUser(email string, veriffSessionId string, veriffSessionUrl string) error {
	return Exec(`
		insert into customer (email, veriffSessionId, veriffSessionUrl)
		values ($1, $2, $3)
	`,
		email,
		veriffSessionId,
		veriffSessionUrl,
	)
}
