package users

const (
	// https://dev.to/yogski/dealing-with-enum-type-in-postgresql-1j3g
	orgunitType     = "CREATE TYPE  orgunit_type AS ENUM ('Entreprise', 'Direction', 'Service', 'Office', 'Bureau', 'Unité', 'Division');"
	orgunitTypeList = "SELECT UNNEST(enum_range(null::orgunit_type)) AS orgunit;"
	orgUnitsCount   = "SELECT COUNT(*) FROM org_unit"
)
