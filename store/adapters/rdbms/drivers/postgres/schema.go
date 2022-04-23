package postgres

import (
	"context"
	"fmt"

	"github.com/cortezaproject/corteza-server/store/adapters/rdbms/ddl"
)

// PostgreSQL specific prefixes, sql
// templates, functions and other helpers

//import (
//	"context"
//	"fmt"
//
//	"github.com/cortezaproject/corteza-server/store/adapters/rdbms"
//	. "github.com/cortezaproject/corteza-server/store/adapters/rdbms/ddl"
//	"go.uber.org/zap"
//)
//
//type (
//	upgrader struct {
//		log     *zap.Logger
//		s       *rdbms.Store
//		dialect *CommonDialect
//	}
//)

type (
	SchemaModifier struct {
		*ddl.GenericSchemaModifier
	}
)

func (s *SchemaModifier) TableExists(ctx context.Context, db ddl.SqlExecer, table string) (bool, error) {
	var exists bool

	if err := db.GetContext(ctx, &exists, "SELECT TO_REGCLASS($1) IS NOT NULL", "public."+table); err != nil {
		return false, fmt.Errorf("could not check if table exists: %w", err)
	}

	return exists, nil
}

//
//var (
//	_ Upgrader = &upgrader{}
//)
//
//// NewUpgrader returns PostgreSQL schema upgrader
//func NewUpgrader(store *rdbms.Store) *upgrader {
//	var g = &upgrader{
//		log:     zap.NewNop(),
//		s:       store,
//		dialect: NewCommonDialect(),
//	}
//
//	// Modification to general DDL dialect to support Postgres requirements
//	g.dialect.AddTemplate("create-table-suffix", "WITHOUT OIDS")
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
//	return CommonUpgrades(u.log, u).Before(ctx)
//}
//
//// After runs after all tables are upgraded
//func (u *upgrader) After(ctx context.Context) error {
//	return CommonUpgrades(u.log, u).After(ctx)
//}
//
//// CreateTable is triggered for every table defined in the rdbms package
////
//// It checks if table is missing and creates it, otherwise
//// it runs
//func (u upgrader) CreateTable(ctx context.Context, t *Table) (err error) {
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
//	} else {
//		return u.upgradeTable(ctx, t)
//	}
//
//	return nil
//}
//
//func (u upgrader) Exec(ctx context.Context, sql string, aa ...interface{}) error {
//	_, err := u.s.DB().ExecContext(ctx, sql, aa...)
//	return err
//}
//
//// upgradeTable applies any necessary changes connected to that specific table
//func (u *upgrader) upgradeTable(ctx context.Context, t *Table) error {
//	g := CommonUpgrades(u.log, u)
//
//	switch t.Name {
//	default:
//		return g.Upgrade(ctx, t)
//	}
//}
//
//func (u upgrader) TableExists(ctx context.Context, table string) (bool, error) {
//	var exists bool
//
//	if err := u.s.DB().GetContext(ctx, &exists, "SELECT TO_REGCLASS($1) IS NOT NULL", "public."+table); err != nil {
//		return false, fmt.Errorf("could not check if table exists: %w", err)
//	}
//
//	return exists, nil
//}
//
//func (u upgrader) DropTable(ctx context.Context, table string) (dropped bool, err error) {
//	var exists bool
//	exists, err = u.TableExists(ctx, table)
//	if err != nil || !exists {
//		return false, err
//	}
//
//	err = u.Exec(ctx, fmt.Sprintf(`DROP TABLE %s`, table))
//	if err != nil {
//		return false, err
//	}
//
//	return true, nil
//}
//
//func (u upgrader) TableSchema(ctx context.Context, table string) (Columns, error) {
//	return nil, fmt.Errorf("pending implementation")
//}
//
//// AddColumn adds column to table
//func (u upgrader) AddColumn(ctx context.Context, table string, col *Column) (added bool, err error) {
//	err = func() error {
//		var columns Columns
//		if columns, err = u.getColumns(ctx, table); err != nil {
//			return err
//		}
//
//		if columns.Get(col.Name) != nil {
//			return nil
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
//		var columns Columns
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
//		var columns Columns
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
//func (u upgrader) AddPrimaryKey(ctx context.Context, table string, ind *Index) (added bool, err error) {
//	if err = u.Exec(ctx, u.dialect.AddPrimaryKey(table, ind)); err != nil {
//		return false, fmt.Errorf("could not add primary key to table %s: %w", table, err)
//	}
//
//	return true, nil
//}
//
//func (u upgrader) CreateIndex(ctx context.Context, ind *Index) (added bool, err error) {
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
//		lookup = "SELECT COUNT(*) > 0 FROM pg_indexes WHERE tablename = $1 AND indexname = $2"
//	)
//
//	return has, u.s.DB().GetContext(ctx, &has, lookup, table, table+"_"+name)
//}
//
//// loads and returns all tables columns
//func (u upgrader) getColumns(ctx context.Context, table string) (out Columns, err error) {
//	type (
//		col struct {
//			Name       string `db:"column_name"`
//			IsNullable bool   `db:"is_nullable"`
//			DataType   string `db:"data_type"`
//		}
//	)
//
//	var (
//		lookup = `SELECT column_name,
//                         is_nullable = 'YES' AS is_nullable,
//                         data_type
//                    FROM information_columns
//                   WHERE table_catalog = $1
//                     AND table_name = $2`
//
//		cols []*col
//	)
//
//	if err = u.s.DB().SelectContext(ctx, &cols, lookup, u.s.Config().DBName, table); err != nil {
//		return nil, err
//	}
//
//	out = make([]*Column, len(cols))
//	for i := range cols {
//		out[i] = &Column{
//			Name: cols[i].Name,
//			//Type:         ColumnType{},
//			IsNull: cols[i].IsNullable,
//			//DefaultValue: "",
//		}
//	}
//
//	return out, nil
//}