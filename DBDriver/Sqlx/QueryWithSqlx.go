package querywithsqlx

import (
	"context"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"gorm.io/gorm"
)

// 题目1：使用SQL扩展库进行查询
// 假设你已经使用Sqlx连接到一个数据库，并且有一个 employees 表，包含字段 id 、 name 、 department 、 salary 。
// 要求 ：
// 编写Go代码，使用Sqlx查询 employees 表中所有部门为 "技术部" 的员工信息，并将结果映射到一个自定义的 Employee 结构体切片中。
// 编写Go代码，使用Sqlx查询 employees 表中工资最高的员工信息，并将结果映射到一个 Employee 结构体中。

type Employee struct {
	gorm.Model
	Name       string  `db:"name"`
	Department string  `db:"department"`
	Salary     float64 `db:"salary"`
}

func NewDB(ctx context.Context, dsn string) (*sqlx.DB, error) {
	// Open does not establish connections immediately
	db, err := sqlx.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	// Pooling
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxIdleTime(5 * time.Minute)
	db.SetConnMaxLifetime(30 * time.Minute)

	// Verify connectivity with timeout
	ctxPing, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	if err := db.PingContext(ctxPing); err != nil {
		_ = db.Close()
		return nil, err
	}
	return db, nil
}

func CreateEmployeesTable(ctx context.Context, db *sqlx.DB) error {
	// Keep DDL fast with a short timeout
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	const schema = `
CREATE TABLE IF NOT EXISTS employees (
  id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  name VARCHAR(255) NOT NULL,
  department VARCHAR(255) NOT NULL,
  salary DECIMAL(12,2) NOT NULL,
  PRIMARY KEY (id),
  KEY idx_department (department),
  KEY idx_salary (salary)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
`
	_, err := db.ExecContext(ctx, schema)
	return err
}

func DoBatchInsertEmployees(ctx context.Context, db *sqlx.DB) {
	// Prepare data
	emps := []Employee{
		{Name: "Alice", Department: "技术部", Salary: 120000},
		{Name: "Bob", Department: "技术部", Salary: 115000},
		{Name: "Carol", Department: "HR", Salary: 90000},
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := BatchInsertEmployees(ctx, db, emps); err != nil {
		log.Fatalf("batch insert: %v", err)
	}
	log.Println("batch insert ok")
}

func BatchInsertEmployees(ctx context.Context, db *sqlx.DB, emps []Employee) error {
	if len(emps) == 0 {
		return nil
	}
	tx, err := db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	stmt, err := tx.PrepareNamedContext(ctx, `
		INSERT INTO employees (name, department, salary)
		VALUES (:name, :department, :salary)
	`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, e := range emps {
		if _, err := stmt.Exec(e); err != nil {
			return err
		}
	}
	return tx.Commit()
}

func QueryTechDeptEmployees(ctx context.Context, db *sqlx.DB) ([]Employee, error) {
	// Keep queries bounded with a timeout
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	const sql = `
		SELECT id, name, department, salary
		FROM employees
		WHERE department = ?
		ORDER BY id
	`
	var emps []Employee
	if err := db.SelectContext(ctx, &emps, sql, "技术部"); err != nil {
		return nil, err
	}
	return emps, nil
}

func GetTopPaidEmployee(ctx context.Context, db *sqlx.DB) (Employee, error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	const sql = `
		SELECT id, name, department, salary
		FROM employees
		ORDER BY salary DESC, id ASC
		LIMIT 1
	`

	var emp Employee
	if err := db.GetContext(ctx, &emp, sql); err != nil {
		// 可能返回 sql.ErrNoRows（表为空）
		return Employee{}, err
	}
	return emp, nil
}
