// Package querywithsqlx with type mapping
package querywithsqlx

import (
	"context"
	"log"
	"time"

	"github.com/jmoiron/sqlx"
	"gorm.io/gorm"
)

// 假设有一个 books 表，包含字段 id 、 title 、 author 、 price 。
// 要求 ：
// 定义一个 Book 结构体，包含与 books 表对应的字段。
// 编写Go代码，使用Sqlx执行一个复杂的查询，例如查询价格大于 50 元的书籍，并将结果映射到 Book 结构体切片中，确保类型安全。

// Book represents a row in the books table.
// Field tags ensure sqlx maps columns correctly and type-safely.
type Book struct {
	gorm.Model
	Title  string  `db:"title"`
	Author string  `db:"author"`
	Price  float64 `db:"price"`
}

// CreateBooksTable creates the books table if it does not exist.
func CreateBooksTable(ctx context.Context, db *sqlx.DB) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	const schema = `
CREATE TABLE IF NOT EXISTS books (
  id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  title VARCHAR(255) NOT NULL,
  author VARCHAR(255) NOT NULL,
  price DECIMAL(12,2) NOT NULL,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP NULL DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP,
  deleted_at datetime(3) DEFAULT NULL,
  PRIMARY KEY (id),
  KEY idx_price (price)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
`
	_, err := db.ExecContext(ctx, schema)
	return err
}

func DoBatchInsertBooks(ctx context.Context, db *sqlx.DB) {
	// Prepare data
	books := []Book{
		{Title: "book1", Author: "author1", Price: 100},
		{Title: "book2", Author: "author2", Price: 200},
		{Title: "book3", Author: "author3", Price: 300},
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := BatchInsertBooks(ctx, db, books); err != nil {
		log.Fatalf("batch insert: %v", err)
	}
	log.Println("batch insert ok")
}

func BatchInsertBooks(ctx context.Context, db *sqlx.DB, books []Book) error {
	if len(books) == 0 {
		return nil
	}
	tx, err := db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	stmt, err := tx.PrepareNamedContext(ctx, `
		INSERT INTO books (title, author, price)
		VALUES (:title, :author, :price)
	`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, e := range books {
		if _, err := stmt.Exec(e); err != nil {
			return err
		}
	}
	return tx.Commit()
}

// QueryExpensiveBooks returns all books with price greater than minPrice.
func QueryExpensiveBooks(ctx context.Context, db *sqlx.DB, minPrice float64) ([]Book, error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	const q = `
	SELECT id, title, author, price
	FROM books
	WHERE price > ?
	ORDER BY price DESC, id ASC
`

	var books []Book
	if err := db.SelectContext(ctx, &books, q, minPrice); err != nil {
		return nil, err
	}
	return books, nil
}
