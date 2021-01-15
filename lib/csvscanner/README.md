This library provides a convenience wrapper for fetching remote data formatted as csv.

Library converts CSV rows to json rows and provides Bytes and Row functions, which can be used to map object into desired type.

Bytes return marshalled type `map[string]string`.

Row returns raw data in format `[]string`

Warning: doesn't support files with BOM.
