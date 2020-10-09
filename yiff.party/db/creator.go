package db

import (
	yp "github.com/kycklingar/pultimad/yiff.party/parser"
)

func (db *DB) GetCreators(limit int) ([]*yp.Creator, error) {
	rows, err := db.Query(
		"SELECT id, name FROM creators WHERE (downloaded IS NULL OR downloaded < NOW() - INTERVAL '7 days') AND download LIMIT $1",
		limit,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var creators []*yp.Creator

	for rows.Next() {
		var c = new(yp.Creator)
		if err := rows.Scan(&c.ID, &c.Name); err != nil {
			return nil, err
		}

		creators = append(creators, c)
	}

	return creators, rows.Err()
}

func (db *DB) CheckCreator(creator *yp.Creator) error {
	_, err := db.Exec(
		"UPDATE creators SET downloaded = CURRENT_TIMESTAMP WHERE id = $1",
		creator.ID,
	)
	return err
}

func (db *DB) StoreCreator(creator yp.Creator) error {
	_, err := db.Exec(
		"INSERT INTO creators(id, name) VALUES($1, $2) ON CONFLICT DO NOTHING",
		creator.ID,
		creator.Name,
	)
	return err
}

func (db *DB) LoadCreators() ([]*yp.Creator, error) {
	rows, err := db.Query("SELECT id, name FROM creators")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var creators []*yp.Creator

	for rows.Next() {
		var c = new(yp.Creator)
		err = rows.Scan(&c.ID, &c.Name)
		if err != nil {
			return nil, err
		}
		creators = append(creators, c)
	}

	return creators, rows.Err()
}
