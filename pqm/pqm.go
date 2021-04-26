package pqm

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"

	_ "github.com/lib/pq"
)

type Table struct {
	Title  string
	Column map[string]*column
	Keys   map[string]*key
}
type column struct {
	Type      string
	IsNotNull bool
	Default   interface{}
	Length    int64
}
type key struct {
	FromColumns     []string
	ToColumns       []string
	ToTableTitle    string
	IsUnicue        bool
	IsReferences    bool
	IsUpdateCascade bool
}
type tableInfo struct {
	Column     string
	ColumnType string
	Default    string
	Length     int64
	IsNotNull  string
	Key        string
	KeyType    string
	KeyColumn  string
	KeyTable   string
}

func InitTable(db *sql.DB, table *Table) error {
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("start tx is failed: %v", err)
	}
	t, err := scanInfo(table.Title, tx)
	if err != nil {
		return fmt.Errorf("scan table %v is failed: %v", table.Title, err)
	}
	qry := `create table if not exists ` + table.Title + ` (id bigserial primary key);`

	for k, v := range table.Column {
		if tt, ok := t.Column[k]; ok {
			if tt.Type != v.Type {
				deleteColumn(&qry, table.Title, k)
				addColumn(&qry, table.Title, k, v)
				continue
			}
			if tt.Type == "character varying" && tt.Length != v.Length {
				setLengthColumn(&qry, table.Title, k, v.Type, v.Length)
			}
			if tt.IsNotNull != v.IsNotNull {
				setNullColumn(&qry, table.Title, k, v.IsNotNull)
			}
			if tt.Default != v.Default {
				setDefaultColumn(&qry, table.Title, k, v.Type, v.Default)
			}
		} else {
			addColumn(&qry, table.Title, k, v)
		}
	}

	for k, v := range table.Keys {
		if kk, ok := t.Keys[k]; ok {
			if kk.IsReferences != v.IsReferences ||
				kk.IsUnicue != v.IsUnicue ||
				(v.ToTableTitle != "" && kk.ToTableTitle != v.ToTableTitle) ||
				!equalsArray(v.FromColumns, kk.FromColumns) ||
				!equalsArray(v.ToColumns, kk.ToColumns) {
				deleteKey(&qry, table.Title, k)
				addKey(&qry, table.Title, k, v)
			}
		} else {
			addKey(&qry, table.Title, k, v)
		}
	}

	if _, err = tx.Exec(qry); err != nil {
		return fmt.Errorf("migration is failed: %v", err)
	}
	return tx.Commit()
}

func Integer(def int32, isNotNull bool) *column {
	return &column{
		Type:      "integer",
		Default:   def,
		IsNotNull: isNotNull,
		Length:    0,
	}
}
func Bigint(def int64, isNotNull bool) *column {
	return &column{
		Type:      "bigint",
		Default:   def,
		IsNotNull: isNotNull,
		Length:    0,
	}
}
func DPrecision(def float64, isNotNull bool) *column {
	return &column{
		Type:      "double precision",
		Default:   def,
		IsNotNull: isNotNull,
		Length:    0,
	}
}
func VarChar(def string, length int64, isNotNull bool) *column {
	return &column{
		Type:      "character varying",
		Default:   def,
		IsNotNull: isNotNull,
		Length:    length,
	}
}
func Text(def string, isNotNull bool) *column {
	return &column{
		Type:      "text",
		Default:   def,
		IsNotNull: isNotNull,
		Length:    0,
	}
}
func Bytea(def []byte, isNotNull bool) *column {
	return &column{
		Type:      "bytea",
		Default:   def,
		IsNotNull: isNotNull,
		Length:    0,
	}
}
func Array(def []interface{}, isNotNull bool) *column {
	return &column{
		Type:      "array",
		Default:   def,
		IsNotNull: isNotNull,
		Length:    0,
	}
}
func JsonB(def json.RawMessage, isNotNull bool) *column {
	return &column{
		Type:      "jsonb",
		Default:   def,
		IsNotNull: isNotNull,
		Length:    0,
	}
}

func Unique(fromColumn []string) *key {
	return &key{
		FromColumns: fromColumn,
		IsUnicue:    true,
	}
}
func Reference(fromColumn, toTable, toColumn string) *key {
	return &key{
		FromColumns:  []string{fromColumn},
		ToColumns:    []string{toColumn},
		ToTableTitle: toTable,
		IsReferences: true,
	}
}

func scanInfo(title string, tx *sql.Tx) (*Table, error) {
	t := &Table{
		Title:  title,
		Column: map[string]*column{},
		Keys:   map[string]*key{},
	}
	res, err := tx.Query(`
	select
		c.column_name,
		c.data_type,
		case when c.column_default is not null then c.column_default else '' end,
		case when c.character_maximum_length is not null then c.character_maximum_length else 0 end,
		c.is_nullable,
		case when kcu.constraint_name is not null then kcu.constraint_name else '' end,
		case when tc.constraint_type is not null then tc.constraint_type else '' end,
		case when ccu.column_name is not null then ccu.column_name else '' end,
		case when ccu.table_name is not null then ccu.table_name else '' end
	from
		information_schema."columns" c
	left join information_schema.key_column_usage kcu on
		kcu.column_name = c.column_name
	left join information_schema.constraint_column_usage ccu on
		ccu.constraint_name = kcu.constraint_name
	left join information_schema.table_constraints tc on
		tc.constraint_name = ccu.constraint_name
	where
		c.table_name = $1
		and c.column_name <> 'id'
	`, title)
	if err != nil {
		return t, fmt.Errorf("table info not found: %v", err)
	}
	defer res.Close()

	for res.Next() {
		ti := &tableInfo{}
		if err = res.Scan(&ti.Column, &ti.ColumnType, &ti.Default, &ti.Length, &ti.IsNotNull, &ti.Key, &ti.KeyType, &ti.KeyColumn, &ti.KeyTable); err != nil {
			return t, fmt.Errorf("Table scan is failed: %v", err)
		}
		if _, ok := t.Column[ti.Column]; !ok {
			t.Column[ti.Column] = &column{
				Type:    ti.ColumnType,
				Default: ti.Default,
				Length:  ti.Length,
			}
			if ti.IsNotNull == "NO" {
				t.Column[ti.Column].IsNotNull = true
			}
		}
		if ti.Key != "" {
			if k, ok := t.Keys[ti.Key]; ok {
			fcolumns:
				for {
					for _, c := range k.FromColumns {
						if c == ti.Column {
							break fcolumns
						}
					}
					k.FromColumns = append(k.FromColumns, ti.Column)
					break
				}
				if ti.KeyTable != title {
				tcolumns:
					for {
						for _, c := range k.ToColumns {
							if c == ti.KeyColumn {
								break tcolumns
							}
						}
						k.ToColumns = append(k.ToColumns, ti.KeyColumn)
						break
					}
				}
			} else {
				k = &key{
					FromColumns:  []string{ti.Column},
					ToColumns:    []string{},
					ToTableTitle: ti.KeyTable,
				}
				if ti.KeyTable != title {
					k.ToColumns = []string{ti.KeyColumn}
				}
				switch ti.KeyType {
				case "UNIQUE":
					k.IsUnicue = true
				case "FOREIGN KEY":
					k.IsReferences = true
				}
				t.Keys[ti.Key] = k
			}
		}
	}

	return t, nil
}

func addColumn(qry *string, title, key string, c *column) {
	*qry += fmt.Sprintf("\nalter table %v add %v %v", title, key, c.Type)
	if c.Type == "character varying" && c.Length > 0 {
		*qry += fmt.Sprintf("(%v)", c.Length)
	}
	if c.Default != nil {
		if v, ok := c.Default.(string); ok {
			if v != "" {
				*qry += fmt.Sprintf(" default '%v'::%v", v, c.Type)
			}
		} else if v, ok := c.Default.(json.RawMessage); ok {
			if v != nil {
				*qry += fmt.Sprintf(" default '%v'::%v", string(v), c.Type)
			}
		} else {
			*qry += fmt.Sprintf(" default %v::%v", c.Default, c.Type)
		}
	}
	if c.IsNotNull {
		*qry += " not null"
	}
	*qry += ";"
}
func setLengthColumn(qry *string, title, key, typ string, length int64) {
	*qry += fmt.Sprintf("\nalter table %v alter column %v type %v", title, key, typ)
	if length > 0 {
		*qry += fmt.Sprintf("(%v)", length)
	}
	*qry += fmt.Sprintf(" using %v::%v;", key, typ)
}
func setNullColumn(qry *string, title, key string, isNotNull bool) {
	*qry += fmt.Sprintf("\nalter table %v alter column %v", title, key)
	if isNotNull {
		*qry += " set not null;"
	} else {
		*qry += " drop not null;"
	}
}
func setDefaultColumn(qry *string, title, key, typ string, def interface{}) {
	*qry += fmt.Sprintf("\nalter table %v alter column %v", title, key)
	if def != nil {
		if v, ok := def.(string); ok {
			if v != "" {
				*qry += fmt.Sprintf(" set default '%v'::%v;", v, typ)
			} else {
				*qry += " drop default;"
			}
		} else if v, ok := def.(json.RawMessage); ok {
			if v != nil {
				*qry += fmt.Sprintf(" set default '%v'::%v;", string(v), typ)
			} else {
				*qry += " drop default;"
			}
		} else {
			*qry += fmt.Sprintf(" set default %v::%v;", def, typ)
		}
	} else {
		*qry += " drop default;"
	}
}
func deleteColumn(qry *string, title, key string) {
	*qry += fmt.Sprintf("\nalter table %v drop column %v;", title, key)
}

func addKey(qry *string, title, key string, k *key) {
	if k.IsUnicue {
		if len(k.FromColumns) > 0 {
			*qry += fmt.Sprintf("\nalter table %v add constraint %v unique(", title, key)
			*qry += strings.Join(k.FromColumns, ",")
			*qry += ");"
		}
	} else if k.IsReferences {
		if len(k.FromColumns) == 1 && len(k.ToColumns) == 1 {
			*qry += fmt.Sprintf("\nalter table %v add constraint %v foreign key (%v) references %v(%v) on delete cascade;", title, key, k.FromColumns[0], k.ToTableTitle, k.ToColumns[0])
		}
	}
}
func deleteKey(qry *string, title, key string) {
	*qry += fmt.Sprintf("\nalter table %v drop constraint %v;", title, key)
}

func equalsArray(from, to []string) bool {
	flag := false
	if len(from) == 0 && len(to) == 0 {
		flag = true
	}
loop:
	for _, f := range from {
		for _, t := range to {
			if f == t {
				flag = true
				continue loop
			}
		}
		flag = false
	}
	return flag
}
