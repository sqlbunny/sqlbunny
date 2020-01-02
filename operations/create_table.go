package operations

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/sqlbunny/sqlschema/schema"
)

type Column struct {
	Name     string
	Type     string
	Default  string
	Nullable bool
}

func (o Column) Dump(buf *bytes.Buffer) {
	buf.WriteString(fmt.Sprintf("operations.Column{Name: %s, Type: %s, Default: %s, Nullable: %s}", esc(o.Name), esc(o.Type), esc(o.Default), dumpBool(o.Nullable)))
}

type CreateTable struct {
	Name    string
	Columns []Column
}

func (o CreateTable) GetSQL() string {
	var x []string
	for _, c := range o.Columns {
		var n string
		if !c.Nullable {
			n = " NOT NULL"
		}
		var d string
		if c.Default != "" {
			d = " DEFAULT " + c.Default
		}
		x = append(x, fmt.Sprintf("    \"%s\" %s%s%s", c.Name, c.Type, n, d))
	}
	return fmt.Sprintf("CREATE TABLE \"%s\" (\n%s\n)", o.Name, strings.Join(x, ",\n"))
}

func (o CreateTable) Dump(buf *bytes.Buffer) {
	buf.WriteString("operations.CreateTable {\n")
	buf.WriteString("Name: " + esc(o.Name) + ",\n")
	buf.WriteString("Columns: []operations.Column{\n")
	for _, c := range o.Columns {
		c.Dump(buf)
		buf.WriteString(",\n")
	}
	buf.WriteString("},\n")
	buf.WriteString("}")
}

func (o CreateTable) Apply(s *schema.Schema) error {
	if _, ok := s.Tables[o.Name]; ok {
		return fmt.Errorf("CreateTable on already-existing table: %s", o.Name)
	}

	t := schema.NewTable()
	for _, c := range o.Columns {
		t.Columns[c.Name] = &schema.Column{
			Nullable: c.Nullable,
			Type:     c.Type,
			Default:  c.Default,
		}
	}
	s.Tables[o.Name] = t

	return nil
}
