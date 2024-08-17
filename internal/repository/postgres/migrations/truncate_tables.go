package migrations

const TruncateTablesStmt = `
TRUNCATE TABLE "user", "event", "task" RESTART IDENTITY;
`
