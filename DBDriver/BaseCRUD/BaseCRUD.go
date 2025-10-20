package basecrud

import (
	"database/sql"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// 假设有一个名为 students 的表，包含字段 id （主键，自增）、 name （学生姓名，字符串类型）、 age （学生年龄，整数类型）、 grade （学生年级，字符串类型）。
// 要求 ：
// 编写SQL语句向 students 表中插入一条新记录，学生姓名为 "张三"，年龄为 20，年级为 "三年级"。
// 编写SQL语句查询 students 表中所有年龄大于 18 岁的学生信息。
// 编写SQL语句将 students 表中姓名为 "张三" 的学生年级更新为 "四年级"。
// 编写SQL语句删除 students 表中年龄小于 15 岁的学生记录。

type Student struct {
	gorm.Model
	ID    uint   `gorm:"primaryKey"`
	Name  string `gorm:"size:100;not null"`
	Age   int    `gorm:"not null"`
	Grade string `gorm:"size:50;not null"`
}

type DB struct {
	*gorm.DB
}

func NewMySQL(dsn string) (*DB, error) {
	// set logger level, output info level log.
	gLogger := logger.Default.LogMode(logger.Info)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: gLogger,
	})
	if err != nil {
		return nil, err
	}

	// db.DB() gives access to that low-level object so you can tune the connection pool.
	// set connect pool
	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}
	sqlDB.SetMaxOpenConns(25)
	sqlDB.SetMaxIdleConns(25)
	sqlDB.SetConnMaxLifetime(5 * time.Minute)

	return &DB{db}, nil
}

func (db *DB) AutoMigrateModels() error {
	return db.AutoMigrate(&Student{})
}

func (db *DB) InsertStudent(name string, age int, grade string) error {
	return db.Create(&Student{Name: name, Age: age, Grade: grade}).Error
}

func (db *DB) QueryStudentsAgeGreaterThan(age int) ([]Student, error) {
	var out []Student
	err := db.Where("age > ?", age).Find(&out).Error
	return out, err
}

func (db *DB) UpdateStudentGradeByName(name, newGrade string) (int64, error) {
	res := db.Model(&Student{}).Where("name = ?", name).Update("grade", newGrade)
	return res.RowsAffected, res.Error
}

func (db *DB) DeleteStudentsAgeLessThan(age int) (int64, error) {
	res := db.Where("age < ?", age).Delete(&Student{})
	return res.RowsAffected, res.Error
}

func (db *DB) StdDB() (*sql.DB, error) {
	return db.DB.DB()
}

func (db *DB) Close() error {
	sdb, err := db.DB.DB()
	if err != nil {
		return err
	}
	return sdb.Close()
}
