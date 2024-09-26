package db

import (
	"database/sql"
	"fmt"

	"github.com/Yandex-Practicum/go-db-sql-final/internal"

	_ "modernc.org/sqlite"
)

var _ internal.Storage = &ParcelStore{}

// Структура ДБ.
type ParcelStore struct {
	db *sql.DB // БД
}

// NewParcelStore возвращает экземпляр интерфейса Storage.
func NewParcelStore(db *sql.DB) internal.Storage {
	return &ParcelStore{
		db: db,
	}
}

// Add добавляет новый заказ в таблицу parcel.
func (s ParcelStore) Add(p internal.Parcel) (int, error) {
	q := `INSERT INTO parcel (client, status, address, created_at) VALUES (:clientVal, :statusVal, :addressVal, :dateVal)`
	res, err := s.db.Exec(q,
		sql.Named("clientVal", p.Client),
		sql.Named("statusVal", p.Status),
		sql.Named("addressVal", p.Address),
		sql.Named("dateVal", p.CreatedAt))
	if err != nil {
		return 0, err
	}
	id, err := res.LastInsertId()
	if err != nil {

		return 0, err
	}
	return int(id), nil
}

// Get возвращает структуру заказа, если его нет в БД возвращает nil.
func (s ParcelStore) Get(number int) (*internal.Parcel, error) {
	q := `SELECT client, status, address, created_at FROM parcel WHERE number=:numberVal`
	row := s.db.QueryRow(q, sql.Named("numberVal", number))

	p := internal.Parcel{}

	err := row.Scan(&p.Client, &p.Status, &p.Address, &p.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (s ParcelStore) GetByClient(client int) ([]internal.Parcel, error) {
	q := `SELECT number, status, address, created_at FROM parcel WHERE client=:clientVal`
	rows, err := s.db.Query(q, sql.Named("clientVal", client))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var res []internal.Parcel

	for rows.Next() {
		p := internal.Parcel{Client: client}
		err = rows.Scan(&p.Number, &p.Status, &p.Address, &p.CreatedAt)
		if err == sql.ErrNoRows {
			return []internal.Parcel{}, nil
		}
		if err != nil {
			return nil, err
		}
		res = append(res, p)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return res, nil
}

func (s ParcelStore) SetStatus(number int, status string) error {
	q := `UPDATE parcel SET status=:statusVal WHERE number=:numberVal`
	_, err := s.db.Exec(q,
		sql.Named("statusVal", status),
		sql.Named("numberVal", number))
	if err != nil {
		return err
	}
	return nil
}

func (s ParcelStore) SetAddress(number int, address string) error {
	parcel, err := s.Get(number)
	if err != nil {
		return err
	}
	if parcel.Status != internal.ParcelStatusRegistered {
		return fmt.Errorf("нельзя изменить адрес для статуса заказа " + parcel.Status)
	}
	q := `UPDATE parcel SET address=:addressVal WHERE number=:numberVal`
	_, err = s.db.Exec(q,
		sql.Named("addressVal", address),
		sql.Named("numberVal", number))
	if err != nil {
		return err
	}
	return nil
}

func (s ParcelStore) Delete(number int) error {
	q := `DELETE FROM parcel WHERE number=:numberVal`
	_, err := s.db.Exec(q,
		sql.Named("numberVal", number))
	if err != nil {
		return err
	}
	return nil
}
