CREATE TABLE users
(
	id serial NOT NULL,
	username character varying(30) NOT NULL,
	password character varying(255),
	email character varying(60),
	first_name character varying(60),
	last_name character varying(60),
	status integer DEFAULT 0,
	created_at timestamp with time zone,
	updated_at timestamp with time zone,
	deleted_at timestamp with time zone,
	CONSTRAINT users_pkey PRIMARY KEY (id),
	CONSTRAINT users_username_lowercase_ck CHECK (username::text = lower(username::text))
)
WITH (
	OIDS=FALSE
);

CREATE UNIQUE INDEX username_unique_idx
	ON users
	USING btree
	(username COLLATE pg_catalog."default");
