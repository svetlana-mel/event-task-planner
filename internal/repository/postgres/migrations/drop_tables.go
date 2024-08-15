package migrations

const DropTablesStmt = `
DROP TABLE IF EXISTS "task";
DROP TABLE IF EXISTS "event";
DROP TABLE IF EXISTS "user";
`
