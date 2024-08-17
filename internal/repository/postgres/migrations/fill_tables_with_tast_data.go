package migrations

const FillTablesWithTestDataStmt = `
BEGIN;

-- create test user with user_id=1
insert into "user" (name, email, password, created_date_time, updated_date_time, last_login) 
values ('Ivan', 'sample@mail.com', 'dfasfsdf', to_timestamp(0), to_timestamp(0), to_timestamp(0));

insert into "event" (name, description, date_time, fk_user_id)
values
('Interviev', '', CURRENT_TIMESTAMP, 1),
('apply for a building pass', 'need a pass tothe interviev', CURRENT_TIMESTAMP, 1);

insert into "task" 
(name, description, list, start_date_time, fk_event_id, fk_user_id)
values
('read textbook', '', NULL, CURRENT_TIMESTAMP, 1, 1),
('finish my pet prodject', 'api for postman', NULL, CURRENT_TIMESTAMP, 1, 1),
('write email', 'use form from site', NULL, CURRENT_TIMESTAMP, 2, 1),
('by a new jeans', 'black and stretchy', NULL, CURRENT_TIMESTAMP, NULL, 1);

END;
`
