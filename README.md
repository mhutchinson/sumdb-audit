# Auditor / Cloner for SumDB

The clone tool downloads all entries from the [Go SumDB](https://blog.golang.org/module-mirror-launch) into a local SQLite database, and verifies that the downloaded data matches the log commitment.

## Running

The following command will download all entries and store them in the database file provided:
```bash
go run ./cli/clone/clone.go -db ~/sum.db
```

The number of leaves downloaded can be queried:
```bash
sqlite3 ~/sum.db 'SELECT COUNT(*) FROM leaves;'
```

And the tile hashes at different levels inspected:
```bash
sqlite3 ~/sum.db 'SELECT level, COUNT(*) FROM tiles GROUP BY level;'
```

And the processed leaf data can be inspected to ensure that the same module+version does not appear twice:
```bash
sqlite3 ~/sum.db 'SELECT module, version, COUNT(*) cnt FROM leafMetadata GROUP BY module, version HAVING cnt > 1;'
```

## TODO
* This only downloads complete tiles, which means that at any point there could be up to 255 leaves missing from the database. These stragglers should be stored if the root hash checks out.
* The verified Checkpoint should be stored locally.
* Parse the downloaded data to key by module & version, and check no module & version appears twice in the log.
