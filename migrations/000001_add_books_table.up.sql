CREATE TABLE IF NOT EXISTS books (
	id bigserial NOT NULL,
	created_at timestamptz NULL,
	updated_at timestamptz NULL,
	deleted_at timestamptz NULL,
	title varchar(255) NOT NULL,
	author varchar(100) NOT NULL,
	price numeric(10, 2) NULL,
	description text NULL,
	CONSTRAINT books_pkey PRIMARY KEY (id)
);
CREATE INDEX idx_books_deleted_at ON books USING btree (deleted_at);