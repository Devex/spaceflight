// Package surgeon contains the types for schema 'public'.
package surgeon

import "fmt"

// GENERATED BY XO. DO NOT EDIT.

// GetProcesses runs a custom query, returning results as ProcessList.
func GetProcesses(db XODB, options map[string]bool) (fmt.Stringer, error) {
	var err error

	// sql query
	const sqlstr = `SELECT` +
		`      pid,` +
		`      usename,` +
		`      datname,` +
		`      client_addr,` +
		`      CASE` +
		`      WHEN state = 'active'::text THEN` +
		`        EXTRACT(EPOCH FROM now() - query_start)` +
		`      WHEN state = 'idle in transcation'::text THEN` +
		`        EXTRACT(EPOCH FROM now() - xact_start)` +
		`      ELSE` +
		`        EXTRACT(EPOCH FROM now() - state_change)` +
		`      END AS duration,` +
		`      CASE` +
		`      WHEN state = 'active'::text THEN` +
		`        CASE` +
		`        WHEN query IS NOT NULL THEN` +
		`          btrim(` +
		`            regexp_replace(` +
		`              regexp_replace(` +
		`                substring(query,0,80),` +
		`                E'[\n\r\u2028]+', '', 'g'` +
		`              ),` +
		`              E'[ ]+', ' ', 'g'` +
		`            ),` +
		`          ' ')` +
		`        ELSE` +
		`          'No query'` +
		`        END` +
		`      ELSE` +
		`        'Inactive'` +
		`      END AS query,` +
		`      CASE` +
		`      WHEN waiting THEN` +
		`        'true'` +
		`      ELSE` +
		`        'false'` +
		`      END` +
		`    FROM pg_stat_activity` +
		`    WHERE query NOT ILIKE '%pg_stat_activity%'`

	res := ProcessList{
		Processes: nil,
		Inactive:  options["showInactive"],
	}
	// run query
	XOLog(sqlstr)
	q, err := db.Query(sqlstr)
	if err != nil {
		return res, err
	}
	defer q.Close()

	// load results
	for q.Next() {
		pl := Process{}

		// scan
		err = q.Scan(
			&pl.Pid,
			&pl.Usename,
			&pl.Datname,
			&pl.Client,
			&pl.Duration,
			&pl.Query,
			&pl.Waiting,
		)
		if err != nil {
			return res, err
		}

		res.Processes = append(res.Processes, pl)
	}

	return res, nil
}