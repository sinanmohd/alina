BEGIN;

CREATE TABLE user_agents (
	id bigint GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
	user_agent text UNIQUE NOT NULL,
	created_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE uploads (
	id bigint GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
	created_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
	ip_addr inet NOT NULL,
	user_agent bigint REFERENCES user_agents(id) NOT NULL,
	file bigint REFERENCES files(id) NOT NULL,
	name text NOT NULL
);

DO $$ BEGIN
	IF EXISTS (SELECT FROM files) THEN
		INSERT INTO user_agents (user_agent)
		VALUES ('unknown: db migrated from old schema revision');

		INSERT INTO uploads(created_at, ip_addr, user_agent, name, file)
		SELECT files.created_at, files.ip_addr, 1, files.name, files.id
		FROM files;
	END IF;
END $$;

ALTER TABLE files
DROP COLUMN created_at,
DROP COLUMN ip_addr,
DROP COLUMN name,
ALTER last_access DROP NOT NULL,
ALTER COLUMN hash TYPE text,
ALTER COLUMN mime_type TYPE text;

ALTER TABLE chunked
ALTER COLUMN name TYPE text;
ALTER TABLE chunked
ADD user_agent bigint REFERENCES user_agents(id) NOT NULL;

COMMIT;
