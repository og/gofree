package connectRDS

import f "github.com/og/gofree"

type RDS struct {
	Main f.Database
}
func (rds RDS) Close() error {
	return rds.Main.Close()
}
func NewRDS () (RDS, error) {
	db, err := f.NewDatabase(dataSourceName)
	if err != nil {return RDS{}, err}
	return RDS{
		Main: db,
	}, nil
}

