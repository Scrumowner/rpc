package scanner

//
//import (
//	"reflect"
//	"strings"
//)
//
//const (
//	AllFields = "all"
//	Create    = "create"
//	Update    = "update"
//	Upsert    = "upsert"
//	Conflict  = "conflict"
//)
//
//type Scanner interface {
//	RegisterTable(entities ...Tabler)
//	OperationFields(tableName, operation string) []string
//	Table(tableName string) Table
//	Tables() map[string]Table
//}
//
//type Tabler interface {
//	TableName() string
//}
//
//type Table struct {
//	Name            string
//	Fields          []Field
//	FieldsMap       map[string]Field
//	Constraint      []Constraint
//	OperationFields map[string][]string
//	Entity          Tabler
//}
//
//type Field struct {
//	Name       string
//	Type       string
//	Default    string
//	Constraint Constraint
//	Table      *Table
//}
//type Constraint struct {
//	Index  bool
//	Unique bool
//	Field  *Field
//}
//
//type TableScanner struct {
//	tables map[string]Table
//}
//
//func NewTableScanner() Scanner {
//	return &TableScanner{}
//}
//
//func (t *TableScanner) RegisterTable(entities ...Tabler) {
//	tableEntites := make(map[string]Tabler, len(entities))
//	t.tables = make(map[string]Table, len(entities))
//	for i := range entities {
//		tableEntites[entities[i].TableName()] = entities[i]
//	}
//
//	for name, entity := range tableEntites {
//		table := Table{
//			Name:            name,
//			FieldsMap:       map[string]Field{},
//			OperationFields: map[string][]string{},
//			Entity:          entity,
//		}
//		reflected := reflect.TypeOf(entity).Elem()
//		for i := 0; i < reflected.NumField(); i++ {
//			structFild := reflected.Field(i)
//			fieldName := structFild.Tag.Get("db")
//			if fieldName == "" || fieldName == "-" {
//				continue
//			}
//			table.OperationFields[AllFields] = append(table.OperationFields[AllFields], fieldName)
//
//			field := Field{
//				Name:    fieldName,
//				Type:    structFild.Tag.Get("db_type"),
//				Default: structFild.Tag.Get("db_default"),
//				Table:   &table,
//			}
//			contraintRaw := structFild.Tag.Get("db_index")
//			constraintPieces := strings.Split(contraintRaw, ",")
//			if len(constraintPieces) > 0 {
//				for j := range constraintPieces {
//					switch constraintPieces[j] {
//					case "index":
//						field.Constraint.Index = true
//					case "unique":
//						field.Constraint.Unique = true
//
//					}
//				}
//			}
//			if field.Constraint.Index {
//
//			}
//
//		}
//	}
//
//}
