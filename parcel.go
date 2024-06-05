package main

import (
	"database/sql"
	"fmt"
)

type ParcelStore struct {
	db *sql.DB
}

func NewParcelStore(db *sql.DB) ParcelStore {
	return ParcelStore{db: db}
}

func (s ParcelStore) Add(p Parcel) (int, error) {

	res, err := s.db.Exec("INSERT INTO parcel (client, status, address, created_at) VALUES (:client, :status, :address, :created_at)",
		sql.Named("client", p.Client),
		sql.Named("status", p.Status),
		sql.Named("address", p.Address),
		sql.Named("created_at", p.CreatedAt))

	if err != nil {
		fmt.Println(err)
		return 0, nil
	}
	lastInserted, err := res.LastInsertId()
	if err != nil {
		fmt.Println(err)
		return 0, nil
	}

	return int(lastInserted), nil
}

func (s ParcelStore) Get(number int) (Parcel, error) {

	p := Parcel{}
	row := s.db.QueryRow("SELECT * FROM parcel WHERE number = :number", sql.Named("number", number))
	err := row.Scan(&p.Number, &p.Client, &p.Status, &p.Address, &p.CreatedAt)

	if err != nil {
		return p, err
	}

	return p, nil
}

func (s ParcelStore) GetByClient(client int) ([]Parcel, error) {

	var res []Parcel

	rows, err := s.db.Query("SELECT * FROM parcel WHERE client = :client", sql.Named("client", client))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		p := Parcel{}

		err = rows.Scan(&p.Number, &p.Client, &p.Status, &p.Address, &p.CreatedAt)
		if err != nil {
			return nil, err
		}

		res = append(res, p)
	}

	return res, nil
}

func (s ParcelStore) SetStatus(number int, status string) error {

	_, err := s.db.Exec("UPDATE parcel SET status = :status WHERE number = :number",
		sql.Named("status", status),
		sql.Named("number", number))

	if err != nil {
		return err
	}

	return nil
}

func (s ParcelStore) SetAddress(number int, address string) error {

	_, err := s.db.Exec("UPDATE parcel SET address = :address WHERE number = :number AND status LIKE :status",
		sql.Named("address", address),
		sql.Named("number", number),
		sql.Named("status", ParcelStatusRegistered))

	if err != nil {
		return err
	}

	return nil
}

func (s ParcelStore) Delete(number int) error {

	_, err := s.db.Exec("DELETE FROM parcel WHERE number = :number AND status LIKE :status",
		sql.Named("number", number),
		sql.Named("status", ParcelStatusRegistered))

	if err != nil {
		return err
	}

	return nil
}
