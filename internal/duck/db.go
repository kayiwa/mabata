package duck

import "database/sql"

func Init(db *sql.DB) error {
	_, err := db.Exec(`
		create table if not exists sample_events (
			id integer,
			category varchar,
			value integer
		);
		delete from sample_events;
		insert into sample_events values
			(1, 'alpha', 10),
			(2, 'alpha', 12),
			(3, 'beta', 7),
			(4, 'gamma', 19);
	`)
	return err
}
