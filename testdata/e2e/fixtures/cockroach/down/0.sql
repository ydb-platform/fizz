-- # 1 column
-- # row 1
-- ## 207
CREATE TABLE public.schema_migration (
	version VARCHAR(14) NOT NULL,
	CONSTRAINT "primary" PRIMARY KEY (version ASC),
	UNIQUE INDEX schema_migration_version_idx (version ASC),
	FAMILY "primary" (version)
);
-- # 1 row
