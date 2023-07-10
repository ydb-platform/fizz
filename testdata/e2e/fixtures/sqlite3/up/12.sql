CREATE TABLE IF NOT EXISTS "schema_migration" (
"version" TEXT PRIMARY KEY
);
CREATE UNIQUE INDEX "schema_migration_version_idx" ON "schema_migration" (version);
CREATE TABLE IF NOT EXISTS "e2e_authors" (
"id" TEXT PRIMARY KEY,
"created_at" DATETIME NOT NULL,
"updated_at" DATETIME NOT NULL
);
CREATE TABLE IF NOT EXISTS "e2e_user_posts" (
"id" TEXT PRIMARY KEY,
"author_id" char(36) NOT NULL,
"slug" TEXT NOT NULL,
"content" TEXT NOT NULL DEFAULT '',
FOREIGN KEY ("author_id") REFERENCES "e2e_authors" (id) ON UPDATE NO ACTION ON DELETE CASCADE
);
CREATE INDEX "e2e_user_notes_user_id_idx" ON "e2e_user_posts" ("author_id");
CREATE UNIQUE INDEX "e2e_user_notes_slug_idx" ON "e2e_user_posts" (slug);
