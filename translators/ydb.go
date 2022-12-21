package translators

import (
	"errors"
	"fmt"
	"strings"

	"github.com/gobuffalo/fizz"
)

type YDB struct {
}

func NewYDB() *YDB {
	return &YDB{}
}

func (*YDB) Name() string {
	return "ydb"
}

func (p *YDB) CreateTable(table fizz.Table) (string, error) {
	var cols []string
	var primaryColumn string
	for _, column := range table.Columns {
		if column.Primary {
			switch strings.ToLower(column.ColType) {
			case "int64", "uint8", "uint32", "uint64", "dynumber", "utf8", "date", "datetime": // make sure that we don'table fall into default
			case "integer", "int32":
				column.ColType = "int32"
			case "bool", "boolean":
				column.ColType = "bool"
			case "time", "timestamp":
				column.ColType = "timestamp"
			case "string", "text":
				column.ColType = "string"
			default:
				return "", fmt.Errorf("can not use %s as a primary key", column.ColType)
			}
			primaryColumn = column.Name
		}
		col, err := p.buildAddColumn(column)
		if err != nil {
			return "", err
		}
		cols = append(cols, col)
	}

	primaryKeys := table.PrimaryKeys()
	if len(primaryKeys) == 0 {
		if primaryColumn != "" {
			primaryKeys = append(primaryKeys, primaryColumn)
		} else {
			return "", errors.New("need to specify primary key")
		}
	}
	cols = append(cols, fmt.Sprintf("PRIMARY KEY(%s)", strings.Join(primaryKeys, ", ")))
	sql := []string{fmt.Sprintf("CREATE TABLE %s (\n%s\n);", table.Name, strings.Join(cols, ",\n"))}

	for _, i := range table.Indexes {
		s, err := p.AddIndex(fizz.Table{
			Name:    table.Name,
			Indexes: []fizz.Index{i},
		})
		if err != nil {
			return "", err
		}
		sql = append(sql, s)
	}

	return strings.Join(sql, "\n"), nil
}

func (p *YDB) DropTable(table fizz.Table) (string, error) {
	return fmt.Sprintf("DROP TABLE %s;", table.Name), nil
}

func (p *YDB) RenameTable(t []fizz.Table) (string, error) {
	if len(t) < 2 {
		return "", fmt.Errorf("not enough table names supplied")
	}
	return fmt.Sprintf("ALTER TABLE %s RENAME TO %s;", t[0].Name, t[1].Name), nil
}

func (p *YDB) ChangeColumn(fizz.Table) (string, error) {
	return "", errors.New("changing column is unsupported in YDB")
}

func (p *YDB) AddColumn(table fizz.Table) (string, error) {
	if len(table.Columns) == 0 {
		return "", fmt.Errorf("not enough columns supplied")
	}
	c := table.Columns[0]
	s := fmt.Sprintf("ALTER TABLE %s ADD COLUMN %s %s;", table.Name, c.Name, p.colType(c))
	return s, nil
}

func (p *YDB) DropColumn(table fizz.Table) (string, error) {
	if len(table.Columns) == 0 {
		return "", fmt.Errorf("not enough columns supplied")
	}
	column := table.Columns[0]
	if column.Primary {
		return "", fmt.Errorf("YBD doesn't support drop of primary columns")
	}
	return fmt.Sprintf("ALTER TABLE %s DROP COLUMN %s;", table.Name, column.Name), nil
}

func (p *YDB) RenameColumn(fizz.Table) (string, error) {
	return "", errors.New("renaming column is unsupported in YDB")
}

func (p *YDB) AddIndex(table fizz.Table) (string, error) {
	if len(table.Indexes) == 0 {
		return "", fmt.Errorf("not enough indexes supplied")
	}
	index := table.Indexes[0]
	s := fmt.Sprintf("ALTER TABLE %s ADD INDEX %s GLOBAL ON (%s);", table.Name, index.Name, strings.Join(index.Columns, ", "))
	if index.Options["Cover"] != nil {
		s = fmt.Sprintf("%s COVER (%s)", s, strings.Join(index.Options["Cover"].([]string), ","))
	}
	return s, nil
}

func (p *YDB) DropIndex(table fizz.Table) (string, error) {
	if len(table.Indexes) == 0 {
		return "", fmt.Errorf("not enough indexes supplied")
	}
	index := table.Indexes[0]
	return fmt.Sprintf("ALTER TABLE %s DROP INDEX %s;", table.Name, index.Name), nil
}

func (p *YDB) RenameIndex(fizz.Table) (string, error) {
	return "", errors.New("renaming index is unsupported in YDB")
}

func (p *YDB) AddForeignKey(fizz.Table) (string, error) {
	return "", errors.New("foreign key is unsupported in YDB")
}

func (p *YDB) DropForeignKey(fizz.Table) (string, error) {
	return "", errors.New("foreign key is unsupported in YDB")
}

func (p *YDB) buildAddColumn(c fizz.Column) (string, error) {
	s := fmt.Sprintf("%s %s", c.Name, p.colType(c))
	if c.Options["null"] == nil && c.Primary {
		s = fmt.Sprintf("%s NOT NULL", s)
	}
	return s, nil
}

func (p *YDB) colType(c fizz.Column) string {
	switch c.ColType {
	case "bool", "boolean":
		return "bool"
	case "int", "integer", "int32":
		return "int32"
	case "decimal", "numeric":
		return "decimal"
	case "time", "timestamp":
		return "timestamp"
	case "blob", "[]byte", "string", "text":
		return "string"
	default:
		return c.ColType
	}
}
