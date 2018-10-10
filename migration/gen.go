package migration

import (
	"bytes"
	"errors"
	"fmt"

	"github.com/kernelpayments/sqlbunny/common"
	"github.com/kernelpayments/sqlbunny/schema"
)

func Run(s2 *schema.Schema, mstore *MigrationStore) error {
	if mstore == nil {
		return errors.New("I don't have the existing migrations. Set them using migration.Set()")
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
	buf.WriteString("import \"github.com/kernelpayments/sqlbunny/migration\"\n")
	buf.WriteString(fmt.Sprintf("func init() {\nMigrations.Register(%d, ", migrationNumber))
	ops.Dump(&buf)
	buf.WriteString(")\n}")

	if err := common.WriteFile("migrations", migrationFile, &buf); err != nil {
		return err
	}
	return nil
}
