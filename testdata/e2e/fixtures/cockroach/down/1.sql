-- # 1 column
-- # row 1
-- ## 207
CREATE TABLE public.schema_migration (
	version VARCHAR(14) NOT NULL,
	CONSTRAINT "primary" PRIMARY KEY (version ASC),
	UNIQUE INDEX schema_migration_version_idx (version ASC),
	FAMILY "primary" (version)
);
-- # row 2
-- ## 247
CREATE TABLE public.e2e_users (
	id UUID NOT NULL,
	created_at TIMESTAMP NOT NULL,
	updated_at TIMESTAMP NOT NULL,
	username VARCHAR(255) NULL,
	CONSTRAINT "primary" PRIMARY KEY (id ASC),
	FAMILY "primary" (id, created_at, updated_at, username)
);
-- # 2 rows
