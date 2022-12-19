package translators_test

import (
	"github.com/gobuffalo/fizz"
	"github.com/gobuffalo/fizz/translators"
)

var _ fizz.Translator = (*translators.YDB)(nil)
var ydbt = translators.NewYDB()

func (p *YDBSuite) Test_YDB_CreateTable() {
	r := p.Require()
	ddl := `CREATE TABLE users (
				id int32 NOT NULL,
				"first_name" string,
				"last_name" string,
				"email" string,
				"permissions" json,
				"age" int32,
				"raw" string NOT NULL,
				"created_at" timestamp NOT NULL,
				"updated_at" timestamp NOT NULL
				PRIMARY KEY("id")
			);`

	res, _ := fizz.AString(`
	create_table("users") {
		t.Column("id", "integer", {"primary": true})                                                                                                     
		t.Column("first_name", "string", {})
		t.Column("last_name", "string", {})
		t.Column("email", "string", {"size":20})
		t.Column("permissions", "jsonb", {"null": true})
		t.Column("age", "integer", {"null": true, "default": 40})
		t.Column("raw", "blob", {})
	}
	`, ydbt)
	r.Equal(ddl, res)
}

func (p *YDBSuite) Test_YDB_CreateTables_WithCompositePrimaryKey() {
	r := p.Require()
	ddl := `CREATE TABLE "user_profiles" (
			"user_id" INT NOT NULL,
			"profile_id" INT NOT NULL,
			"created_at" timestamp,
			"updated_at" timestamp,
			PRIMARY KEY("user_id", "profile_id")
			);`

	res, _ := fizz.AString(`
	create_table("user_profiles") {
		t.Column("user_id", "INT")
		t.Column("profile_id", "INT")
		t.PrimaryKey("user_id", "profile_id")
	}
	`, ydbt)
	r.Equal(ddl, res)
}

func (p *YDBSuite) Test_YDB_DropTable() {
	r := p.Require()

	ddl := `DROP TABLE "users";`

	res, _ := fizz.AString(`drop_table("users")`, ydbt)
	r.Equal(ddl, res)
}

func (p *YDBSuite) Test_YDB_RenameTable() {
	r := p.Require()

	ddl := `ALTER TABLE "users" RENAME TO "people";`

	res, _ := fizz.AString(`rename_table("users", "people")`, ydbt)
	r.Equal(ddl, res)
}

func (p *YDBSuite) Test_YDB_RenameTable_NotEnoughValues() {
	r := p.Require()

	_, err := ydbt.RenameTable([]fizz.Table{})
	r.Error(err)
}

func (p *YDBSuite) Test_YDB_ChangeColumn() {
	_, err := fizz.AString(`change_column("mytable", "mycolumn", "string", {"default": "foo", "size": 50})`, ydbt)

	p.Require().Error(err)
}

func (p *YDBSuite) Test_YDB_AddColumn() {
	ddl := `ALTER TABLE my_table ADD COLUMN my_column string;`
	res, _ := fizz.AString(`add_column("my_table", "my_column", "string"`, ydbt)
	p.Require().Equal(ddl, res)
}

func (p *YDBSuite) Test_YDB_DropColumn() {
	r := p.Require()
	ddl := `ALTER TABLE table_name DROP COLUMN column_name;`

	res, _ := fizz.AString(`drop_column("table_name", "column_name")`, ydbt)

	r.Equal(ddl, res)
}

func (p *YDBSuite) Test_YDB_RenameColumn() {
	_, err := fizz.AString(`rename_column("table_name", "old_column", "new_column")`, ydbt)
	p.Require().Error(err)
}

func (p *YDBSuite) Test_YDB_AddIndex() {
	ddl := `ALTER TABLE table_name ADD INDEX table_name_column_name_idx GLOBAL ON (column_name);`
	res, _ := fizz.AString(`add_index("table_name", "column_name", {})`, ydbt)
	p.Require().Equal(ddl, res)
}

func (p *YDBSuite) Test_YDB_AddIndex_MultiColumn() {
	ddl := `ALTER TABLE table_name ADD INDEX table_name_col1_col2_col3_idx GLOBAL ON (col1, col2, col3);`
	res, _ := fizz.AString(`add_index("table_name", ["col1", "col2", "col3"], {})`, ydbt)
	p.Require().Equal(ddl, res)
}

func (p *YDBSuite) Test_YDB_AddIndex_CustomName() {
	ddl := `ALTER TABLE table_name ADD INDEX custom_name GLOBAL ON (column_name);`
	res, _ := fizz.AString(`add_index("table_name", "column_name", {"name": "custom_name"})`, ydbt)
	p.Require().Equal(ddl, res)
}

func (p *YDBSuite) Test_YDB_DropIndex() {
	ddl := `ALTER TABLE users DROP INDEX my_idx;`
	res, _ := fizz.AString(`drop_index("users", "my_idx")`, ydbt)
	p.Require().Equal(ddl, res)
}

func (p *YDBSuite) Test_YDB_RenameIndex() {
	_, err := fizz.AString(`rename_index("table", "old_ix", "new_ix")`, ydbt)
	p.Require().Error(err)
}

func (p *YDBSuite) Test_YDB_AddForeignKey() {
	_, err := fizz.AString(`add_foreign_key("profiles", "user_id", {"users": ["id"]}, {})`, ydbt)
	p.Require().Error(err)
}

func (p *YDBSuite) Test_YDB_DropForeignKey() {
	_, err := fizz.AString(`drop_foreign_key("profiles", "profiles_users_id_fk", {})`, ydbt)
	p.Require().Error(err)
}
