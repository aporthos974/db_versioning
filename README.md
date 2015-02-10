# db_versioning

DB Versioning permit to manage your database schema with a version system.
This allows you to manipulate the database schema with a tracking queries executed during update.<br/>

Installation
-----------
```bash
  cd db_versioning/
  go get
  go install
```
Usage
-----
```bash
Usage of ./db_versioning [option] <schema>
  -host="localhost": Database environment (not implemented)
  -i=false: Initialize versioning system for database schema
  -u=false: Upgrade database schema
  -v=false: Display database schema version
```

Directory
-----
```bash
db_versioning
schema_name/
├── 1.0.0
│   └── test.sql
├── 1.0.1
│   ├── test1.sql
│   └── test2.sql
```
