package migrateAction

import f "github.com/og/gofree"

func (MasterMigrate) Migrate20201205133727UserAddFieldMobileAndPassword(mi f.Migrate) {
	mi.Exec(`ALTER TABLE user add COLUMN mobile CHAR(11) NOT NULL DEFAULT "";`)
}
