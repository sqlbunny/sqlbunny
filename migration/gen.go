package migration

import (
	"bytes"
	"fmt"
	"go/build"
	"io/ioutil"
	"os"
	"os/exec"
	"plugin"

	"github.com/KernelPay/sqlboiler/common"
	"github.com/KernelPay/sqlboiler/schema"
)

func Run(s2 *schema.Schema) error {
	mstore, err := LoadMigrations()
	if err != nil {
		return err
	}
	s1 := schema.NewSchema()
	mstore.applyAll(s1)

	ops := Diff(nil, s1, s2)

	if len(ops) == 0 {
		return fmt.Errorf("No model changes found, doing nothing.")
	}
	migrationNumber := mstore.nextFree()
	migrationFile := fmt.Sprintf("migration_%05d.go", migrationNumber)

	var buf bytes.Buffer
	common.WritePackageName(&buf, "migrations")
	buf.WriteString("import \"github.com/KernelPay/sqlboiler/migration\"\n")
	buf.WriteString(fmt.Sprintf("func init() {\nMigrations.Register(%d, ", migrationNumber))
	ops.Dump(&buf)
	buf.WriteString(")\n}")

	if err := common.WriteFile("migrations", migrationFile, &buf); err != nil {
		return err
	}
	return nil
}

const loaderProgram = `package main

import "%s/migrations"

var Migrations = &migrations.Migrations`

func LoadMigrations() (*MigrationStore, error) {
	if err := os.RemoveAll("sqlboiler_tmp"); err != nil {
		return nil, err
	}
	if err := os.Mkdir("sqlboiler_tmp", 0777); err != nil {
		return nil, err
	}

	cwd, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("Error getting cwd: %v", err)
	}

	pkg, err := build.Import(".", cwd, 0)
	if err != nil {
		return nil, fmt.Errorf("Error finding package for cwd: %v", err)
	}

	if err := ioutil.WriteFile("./sqlboiler_tmp/main.go", []byte(fmt.Sprintf(loaderProgram, pkg.ImportPath)), 0666); err != nil {
		return nil, err
	}
	cmd := exec.Command("go", "build", "-buildmode=plugin", "-o=./sqlboiler_tmp/migrations.so", "./sqlboiler_tmp")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		return nil, fmt.Errorf("Error compiling migrations package: %v", err)
	}

	p, err := plugin.Open("./sqlboiler_tmp/migrations.so")
	if err != nil {
		return nil, fmt.Errorf("Error loading migrations package: %v", err)
	}
	m, err := p.Lookup("Migrations")
	if err != nil {
		return nil, fmt.Errorf("Error looking for 'Migrations' variable in migrations package: %v", err)
	}

	ops, ok := m.(**MigrationStore)
	if !ok {
		return nil, fmt.Errorf("'Migrations' variable is the wrong type, should be migration.MigrationStore")
	}

	os.RemoveAll("sqlboiler_tmp")

	return *ops, nil
}
