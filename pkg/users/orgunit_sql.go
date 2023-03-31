package users

const (
	// https://dev.to/yogski/dealing-with-enum-type-in-postgresql-1j3g
	orgunitType      = "CREATE TYPE  orgunit_type AS ENUM ('Entreprise', 'Direction', 'Service', 'Office', 'Bureau', 'Unité', 'Division');"
	orgunitTypeExist = "select exists (select 1 from pg_type where typname = 'orgunit_type');"
	orgunitTypeList  = "SELECT UNNEST(enum_range(null::orgunit_type)) AS orgunit_type;"
	orgUnitsCount    = "SELECT COUNT(*) FROM go_orgunit"
)
