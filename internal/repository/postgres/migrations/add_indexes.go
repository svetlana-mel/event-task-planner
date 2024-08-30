package migrations

const AddIndexesStmt = `
CREATE INDEX IF NOT EXISTS idx_tasks_event ON task(fk_event_id);
CREATE INDEX IF NOT EXISTS idx_tasks_user ON task(fk_user_id);
CREATE INDEX IF NOT EXISTS idx_events_user ON event(fk_user_id);
CREATE INDEX IF NOT EXISTS idx_user_email ON "user"(email);`
