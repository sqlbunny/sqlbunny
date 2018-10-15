# geo

Go geometry/geography data types.

Features:
- Point, LineString, Polygon
- MultiPoint, MultiLineString, MultiPolygon
- `*Z`, `*M`, `*ZM` variants, including the Z and/or M coordinate.
- `*S` variants, including an SRID.
- `null.*` variants, to handle nullable SQL columns nicely.
- EWKB marshal/unmarshal.
- sql.Scanner, sql.Valuer implementation, tested with Postgres PostGIS.
- JSON marshal/unmarshal works nicely.
