# geo

Go geometry/geography data types.

Features:
- Point, LineString, Polygon
- MultiPoint, MultiLineString, MultiPolygon
- `*S` variants of the above including an SRID.
- EWKB marshal/unmarshal.
- sql.Scanner, sql.Valuer implementation, tested with Postgres PostGIS.
- JSON marshal/unmarshal works nicely.

To Do:
- `*M`, `*Z`, `*ZM` variants, including the Z and/or M coordinate.
