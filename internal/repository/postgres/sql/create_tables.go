package sql

const CreateTablesStmt = `
CREATE TABLE IF NOT EXISTS public.event
(
    event_id bigserial NOT NULL,
    name character varying(256) NOT NULL,
    date_time timestamp with time zone NOT NULL,
    canceled boolean NOT NULL DEFAULT false,
    deleted boolean NOT NULL DEFAULT false,
    fk_user_id bigint NOT NULL,
    PRIMARY KEY (event_id)
);

CREATE TABLE IF NOT EXISTS public.task
(
    task_id bigserial NOT NULL,
    name character varying(256) NOT NULL,
    description text,
    list text[],
    start_date_time timestamp with time zone,
    end_date_time timestamp with time zone,
    fk_event_id bigint,
    completed boolean NOT NULL DEFAULT false,
    deleted boolean NOT NULL DEFAULT false,
    fk_user_id bigint NOT NULL,
    PRIMARY KEY (task_id)
);

CREATE TABLE IF NOT EXISTS public.user
(
    user_id bigserial NOT NULL,
    name character varying(255) NOT NULL,
    email character varying(255) NOT NULL,
    password character varying(255) NOT NULL,
    created_date_time timestamp with time zone NOT NULL,
    updated_date_time timestamp with time zone NOT NULL,
    last_login timestamp with time zone NOT NULL,
    refresh_token character varying(512),
    refresh_token_expiry timestamp with time zone,
    PRIMARY KEY (user_id)
);

ALTER TABLE IF EXISTS public.event
    ADD FOREIGN KEY (fk_user_id)
    REFERENCES public.user (user_id) MATCH SIMPLE
    ON UPDATE NO ACTION
    ON DELETE NO ACTION
    NOT VALID;


ALTER TABLE IF EXISTS public.task
    ADD FOREIGN KEY (fk_event_id)
    REFERENCES public.event (event_id) MATCH SIMPLE
    ON UPDATE NO ACTION
    ON DELETE NO ACTION
    NOT VALID;


ALTER TABLE IF EXISTS public.task
    ADD FOREIGN KEY (fk_user_id)
    REFERENCES public.user (user_id) MATCH SIMPLE
    ON UPDATE NO ACTION
    ON DELETE NO ACTION
    NOT VALID;
`
