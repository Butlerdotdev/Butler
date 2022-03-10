package db

import "github.com/scylladb/gocqlx/v2/table"

var TableRules = table.New(RulesMetaData)

var RulesMetaData = table.Metadata{
	Name:    "rules",
	Columns: []string{"ruleId", "timestamp"},
	PartKey: []string{"ruleId"},
	SortKey: []string{},
}
