package duck

import (
	"database/sql"
	"fmt"
)

type QueryDef struct {
	Name        string
	Description string
	SQL         string
}

var Registry = []QueryDef{
	{
		Name:        "event_counts",
		Description: "Count rows by category",
		SQL:         `select category, count(*) as n from sample_events group by category order by category`,
	},
	{
		Name:        "event_totals",
		Description: "Sum values by category",
		SQL:         `select category, sum(value) as total from sample_events group by category order by total desc`,
	},
}

func Names() []string {
	out := make([]string, 0, len(Registry))
	for _, q := range Registry {
		out = append(out, q.Name)
	}
	return out
}

func Run(db *sql.DB, name string) ([]string, [][]string, error) {
	var found *QueryDef
	for i := range Registry {
		if Registry[i].Name == name {
			found = &Registry[i]
			break
		}
	}
	if found == nil {
		return nil, nil, fmt.Errorf("unknown query: %s", name)
	}

	rows, err := db.Query(found.SQL)
	if err != nil {
		return nil, nil, err
	}
	defer rows.Close()

	cols, err := rows.Columns()
	if err != nil {
		return nil, nil, err
	}
	data := [][]string{}
	for rows.Next() {
		vals := make([]any, len(cols))
		ptrs := make([]any, len(cols))
		for i := range vals {
			ptrs[i] = &vals[i]
		}
		if err := rows.Scan(ptrs...); err != nil {
			return nil, nil, err
		}
		row := make([]string, len(cols))
		for i, v := range vals {
			row[i] = fmt.Sprint(v)
		}
		data = append(data, row)
	}
	return cols, data, rows.Err()
}
