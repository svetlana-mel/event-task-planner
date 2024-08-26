package migrations

const CreateBlankUserAndEvent = `
ALTER SEQUENCE user_user_id_seq RESTART WITH 1;
ALTER SEQUENCE event_event_id_seq RESTART WITH 1;
ALTER SEQUENCE task_task_id_seq RESTART WITH 1;

insert into "user" (user_id, name, email, pass_hash, created_date_time) 
values (0, '', '', '', to_timestamp(0));

insert into "event" (event_id, name, description, date_time, fk_user_id)
values
(0, '', '', CURRENT_TIMESTAMP, 0);
`
