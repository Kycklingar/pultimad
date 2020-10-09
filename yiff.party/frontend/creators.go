package main

type creator struct {
	ID   int
	Name string
}

func archivedCreators() ([]creator, error) {
	rows, err := db.Query(`
		SELECT id, name
		FROM creators
		WHERE download
		`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var creators []creator

	for rows.Next() {
		var c creator
		err = rows.Scan(&c.ID, &c.Name)
		if err != nil {
			return nil, err
		}
		creators = append(creators, c)
	}

	return creators, nil
}

