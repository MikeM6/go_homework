module github.com/MikeM6/go_homework/DBDriver

go 1.25.0

replace basecrud => ./BaseCRUD

replace transaction => ./Transaction

replace querywithsqlx => ./Sqlx

replace advancegorm => ./AdvanceGorm

require (
	advancegorm v0.0.0-00010101000000-000000000000
	basecrud v0.0.0-00010101000000-000000000000
	querywithsqlx v0.0.0-00010101000000-000000000000
	transaction v0.0.0-00010101000000-000000000000
)

require (
	filippo.io/edwards25519 v1.1.0 // indirect
	github.com/go-sql-driver/mysql v1.9.3 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	github.com/jmoiron/sqlx v1.4.0 // indirect
	golang.org/x/text v0.30.0 // indirect
	gorm.io/driver/mysql v1.6.0 // indirect
	gorm.io/gorm v1.31.0 // indirect
)
