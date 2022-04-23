package mysql

//// MySQL specific prefixes, sql
//// templates, functions and other helpers
//
//import (
//	"context"
//	"database/sql"
//	"fmt"
//
//	"github.com/cortezaproject/corteza-server/store/adapters/rdbms"
//	"github.com/cortezaproject/corteza-server/store/adapters/rdbms/schema"
//	_ "github.com/go-sql-driver/mysql"
//	"go.uber.org/zap"
//)
//
//type (
//	upgrader struct {
//		log     *zap.Logger
//		s       *rdbms.Store
//		dialect *schema.CommonDialect
//	}
//)
//
//// NewUpgrader returns MySQL schema upgrader
//func NewUpgrader(store *rdbms.Store) *upgrader {
//	var g = &upgrader{
//		log:     zap.NewNop(),
//		s:       store,
//		dialect: schema.NewCommonDialect(),
//	}
//
//	// All modifications we need for the DDL generator
//	// to properly support MySQL dialect:
//	g.dialect.AddTemplate("create-table-suffix", "ENGINE=InnoDB DEFAULT CHARSET=utf8")
//
//	// Sadly, MySQL does not support partial indexes
//	//
//	// To work around this, we'll ignore partial indexes
//	// and solve this on application level
//	g.dialect.AddTemplate(
//		"create-index",
//		`{{ if not .Condition }}CREATE {{ if .Unique }}UNIQUE {{ end }}INDEX {{ template "index-name" . }} ON {{ .Table }} {{ template "index-fields" .Fields }}{{ else }}SELECT 1 -- dummy sql, just to prevent "empty query" errors...{{ end }}`,
//	)
//
//	g.dialect.AddTemplate("if-not-exists-clause", "")
//	g.dialect.AddTemplate("index-fields", `
//({{ range $n, $f := . -}}
//	{{ if $n }}, {{ end }}
//	{{- if .Expr}}({{ end }}
//	{{- .Field }}
//	{{- if .Length}}({{ .Length }}){{ end }}
//	{{- if .Expr}}){{ end }}
//	{{- if .Desc }} DESC{{ end }}
//{{- end }})
//`)
//
//	// Cover mysql exceptions
//	g.dialect.AddTemplateFunc("columnType", func(ct *schema.ColumnType) string {
//		switch ct.Type {
//		case schema.ColumnTypeIdentifier:
//			return "BIGINT UNSIGNED"
//		case schema.ColumnTypeText:
//			if y, has := ct.Flags["mysqlLongText"].(bool); has && y {
//				return "LONGTEXT"
//			}
//
//			return "TEXT"
//		case schema.ColumnTypeBinary:
//			return "BLOB"
//		case schema.ColumnTypeTimestamp:
//			return "DATETIME"
//		case schema.ColumnTypeBoolean:
//			return "TINYINT(1)"
//		default:
//			return schema.GenColumnType(ct)
//		}
//	})
//
//	return g
//}
//
//func (u *upgrader) SetLogger(l *zap.Logger) {
//	u.log = l
//}
//
//func (u *upgrader) Store() *rdbms.Store {
//	return u.s
//}
//
//// Before runs before all tables are upgraded
//func (u *upgrader) Before(ctx context.Context) error {
//	tt := []func() error{
//		func() error {
//			const migrations = "migrations"
//			if exists, err := u.TableExists(ctx, migrations); err != nil || !exists {
//				return err
//			}
//
//			if _, err := u.s.DB().ExecContext(ctx, fmt.Sprintf(`DROP TABLE "%s"`, migrations)); err != nil {
//				return err
//			}
//
//			u.log.Debug(fmt.Sprintf("%s table removed", migrations))
//
//			return nil
//		},
//		func() error {
//			// in the first versions we created and used foreign keys
//			// this is now obsolete with future plans with the (composite) store architecture
//			var (
//				// find and remove all foreign keys
//				find = `SELECT CONSTRAINT_NAME, TABLE_NAME FROM information_schema.TABLE_CONSTRAINTS where CONSTRAINT_SCHEMA = '%s' AND CONSTRAINT_TYPE = 'FOREIGN KEY';`
//				drop = `ALTER TABLE %s DROP FOREIGN KEY %s`
//
//				table, constraint string
//			)
//
//			if rows, err := u.s.DB().QueryContext(ctx, fmt.Sprintf(find, u.s.Config().DBName)); err != nil {
//				return err
//			} else {
//				for rows.Next() {
//					if err = rows.Scan(&constraint, &table); err != nil {
//						return err
//					}
//
//					u.log.Debug(fmt.Sprintf("removing foreign key %s from table %s", constraint, table))
//					if _, err = u.s.DB().ExecContext(ctx, fmt.Sprintf(drop, table, constraint)); err != nil {
//						return err
//					}
//				}
//			}
//
//			return nil
//		},
//	}
//
//	for _, t := range tt {
//		if err := t(); err != nil {
//			return err
//		}
//	}
//
//	return schema.CommonUpgrades(u.log, u).Before(ctx)
//}
//
//// After runs after all tables are upgraded
//func (u *upgrader) After(ctx context.Context) error {
//	return schema.CommonUpgrades(u.log, u).After(ctx)
//}
//
//// CreateTable is triggered for every table defined in the rdbms package
////
//// It checks if table is missing and creates it, otherwise
//// it runs
//func (u upgrader) CreateTable(ctx context.Context, t *schema.Table) (err error) {
//	var exists bool
//	if exists, err = u.TableExists(ctx, t.Name); err != nil {
//		return
//	}
//
//	if !exists {
//		if err = u.Exec(ctx, u.dialect.CreateTable(t)); err != nil {
//			return err
//		}
//
//		for _, i := range t.Indexes {
//			if err = u.Exec(ctx, u.dialect.CreateIndex(i)); err != nil {
//				return fmt.Errorf("could not create index %s on table %s: %w", i.Name, i.Table, err)
//			}
//		}
//	}
//
//	if err = u.upgradeTable(ctx, t); err != nil {
//		return
//	}
//
//	return nil
//}
//
//func (u upgrader) DropTable(ctx context.Context, table string) (dropped bool, err error) {
//	var exists bool
//	exists, err = u.TableExists(ctx, table)
//	if err != nil || !exists {
//		return false, err
//	}
//
//	err = u.Exec(ctx, fmt.Sprintf(`DROP TABLE "%s"`, table))
//	if err != nil {
//		return false, err
//	}
//
//	return true, nil
//}
//
//func (u upgrader) Exec(ctx context.Context, sql string, aa ...interface{}) error {
//	_, err := u.s.DB().ExecContext(ctx, sql, aa...)
//	return err
//}
//
//// upgradeTable applies any necessary changes connected to that specific table
//func (u *upgrader) upgradeTable(ctx context.Context, t *schema.Table) error {
//	g := schema.CommonUpgrades(u.log, u)
//
//	switch t.Name {
//	default:
//		return g.Upgrade(ctx, t)
//	}
//}
//
//func (u upgrader) TableExists(ctx context.Context, table string) (bool, error) {
//	var tmp interface{}
//	if err := u.s.DB().GetContext(ctx, &tmp, fmt.Sprintf(`SHOW TABLES LIKE '%s'`, table)); err == sql.ErrNoRows {
//		return false, nil
//	} else if err != nil {
//		return false, fmt.Errorf("could not check if table exists: %w", err)
//	}
//
//	return true, nil
//}
//
//func (u upgrader) TableSchema(ctx context.Context, table string) (schema.Columns, error) {
//	return nil, fmt.Errorf("pending implementation")
//}
//
//// AddColumn adds column to table
//func (u upgrader) AddColumn(ctx context.Context, table string, col *schema.Column) (added bool, err error) {
//	err = func() error {
//		var columns schema.Columns
//		if columns, err = u.getColumns(ctx, table); err != nil {
//			return err
//		}
//
//		if columns.Get(col.Name) != nil {
//			return nil
//		}
//
//		if col.Type.Type == schema.ColumnTypeText || col.Type.Type == schema.ColumnTypeJson {
//			col.DefaultValue = ""
//		}
//
//		if err = u.Exec(ctx, u.dialect.AddColumn(table, col)); err != nil {
//			return err
//		}
//
//		added = true
//		return nil
//	}()
//
//	if err != nil {
//		return false, fmt.Errorf("could not add column %q to %q: %w", col.Name, table, err)
//	}
//
//	return
//}
//
//// DropColumn drops column from table
//func (u upgrader) DropColumn(ctx context.Context, table, column string) (dropped bool, err error) {
//	err = func() error {
//		var columns schema.Columns
//		if columns, err = u.getColumns(ctx, table); err != nil {
//			return err
//		}
//
//		if columns.Get(column) == nil {
//			return nil
//		}
//
//		if err = u.Exec(ctx, u.dialect.DropColumn(table, column)); err != nil {
//			return err
//		}
//
//		dropped = true
//		return nil
//	}()
//
//	if err != nil {
//		return false, fmt.Errorf("could not drop column %q from %q: %w", column, table, err)
//	}
//
//	return
//}
//
//// RenameColumn renames column on a table
//func (u upgrader) RenameColumn(ctx context.Context, table, oldName, newName string) (changed bool, err error) {
//	err = func() error {
//		if oldName == newName {
//			return nil
//		}
//
//		var columns schema.Columns
//		if columns, err = u.getColumns(ctx, table); err != nil {
//			return err
//		}
//
//		if columns.Get(oldName) == nil {
//			// Old column does not exist anymore
//
//			if columns.Get(newName) == nil {
//				return fmt.Errorf("old and new columns are missing")
//			}
//
//			return nil
//		}
//
//		if columns.Get(newName) != nil {
//			return fmt.Errorf("new column already exists")
//
//		}
//
//		if err = u.Exec(ctx, u.dialect.RenameColumn(table, oldName, newName)); err != nil {
//			return err
//		}
//
//		changed = true
//		return nil
//	}()
//
//	if err != nil {
//		return false, fmt.Errorf("could not rename column %q on table %q to %q: %w", oldName, table, newName, err)
//	}
//
//	return
//}
//
//func (u upgrader) AddPrimaryKey(ctx context.Context, table string, ind *schema.Index) (added bool, err error) {
//	if err = u.Exec(ctx, u.dialect.AddPrimaryKey(table, ind)); err != nil {
//		return false, fmt.Errorf("could not add primary key to table %s: %w", table, err)
//	}
//
//	return true, nil
//}
//
//func (u upgrader) CreateIndex(ctx context.Context, ind *schema.Index) (added bool, err error) {
//	if added, err = u.hasIndex(ctx, ind.Table, ind.Name); added || err != nil {
//		return
//	}
//
//	if err = u.Exec(ctx, u.dialect.CreateIndex(ind)); err != nil {
//		return false, fmt.Errorf("could not create index on table %s: %w", ind.Table, err)
//	}
//
//	return true, nil
//}
//
//func (u upgrader) hasIndex(ctx context.Context, table, name string) (has bool, err error) {
//	var (
//		lookup = "SELECT COUNT(*) > 0 FROM information_schema.statistics where table_schema = ? AND table_name = ? AND index_name = ?"
//	)
//
//	return has, u.s.DB().GetContext(ctx, &has, lookup, u.s.Config().DBName, table, table+"_"+name)
//}
//
//// loads and returns all tables columns
//func (u upgrader) getColumns(ctx context.Context, table string) (out schema.Columns, err error) {
//	type (
//		col struct {
//			Name       string `db:"COLUMN_NAME"`
//			IsNullable bool   `db:"IS_NULLABLE"`
//			DataType   string `db:"DATA_TYPE"`
//		}
//	)
//
//	var (
//		lookup = `SELECT COLUMN_NAME,
//                         IS_NULLABLE = 'YES' AS IS_NULLABLE,
//                         DATA_TYPE
//                    FROM INFORMATION_SCHEMA.COLUMNS
//                   WHERE TABLE_SCHEMA = ?
//                     AND TABLE_NAME = ?`
//
//		cols []*col
//	)
//
//	if err = u.s.DB().SelectContext(ctx, &cols, lookup, u.s.Config().DBName, table); err != nil {
//		return nil, err
//	}
//
//	out = make([]*schema.Column, len(cols))
//	for i := range cols {
//		out[i] = &schema.Column{
//			Name: cols[i].Name,
//			//Type:         schema.ColumnType{},
//			IsNull: cols[i].IsNullable,
//			//DefaultValue: "",
//		}
//	}
//
//	return out, nil
//}