-- public.roles definition

-- Drop table

-- DROP TABLE public.roles;

CREATE TABLE IF NOT EXISTS roles (
	id bigserial NOT NULL,
	created_at timestamptz NULL,
	updated_at timestamptz NULL,
	deleted_at timestamptz NULL,
	"name" varchar(255) NOT NULL,
	wx_user_id varchar(255) NOT NULL,
	avatar_url varchar(255) NOT NULL,
	role_data jsonb NOT NULL,
	description varchar(255) NOT NULL,
	CONSTRAINT roles_pkey PRIMARY KEY (id)
);
CREATE INDEX idx_roles_deleted_at ON roles USING btree (deleted_at);