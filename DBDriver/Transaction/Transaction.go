package transaction

// 假设有两个表： accounts 表（包含字段 id 主键， balance 账户余额）和 transactions 表（包含字段 id 主键， from_account_id 转出账户ID， to_account_id 转入账户ID， amount 转账金额）。
// 要求 ：
// 编写一个事务，实现从账户 A 向账户 B 转账 100 元的操作。在事务中，需要先检查账户 A 的余额是否足够，如果足够则从账户 A 扣除 100 元，向账户 B 增加 100 元，并在 transactions 表中记录该笔转账信息。如果余额不足，则回滚事务。

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/logger"
)

type Account struct {
	gorm.Model
	Balance int64 `gorm:"not null"`
}

type Transfer struct {
	gorm.Model
	FromAccountID uint  `gorm:"not null;index"`
	ToAccountID   uint  `gorm:"not null;index"`
	Amount        int64 `gorm:"not null"`
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

func AutoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(&Account{}, &Transfer{})
}

func CreateAccounts(db *gorm.DB, balances ...int) error {
	var accounts []Account
	for _, b := range balances {
		accounts = append(accounts, Account{Balance: int64(b)})
	}
	return db.CreateInBatches(accounts, 1000).Error
}

func TransferAmount(db *gorm.DB, fromID, toID uint, amount int64) error {
	if amount <= 0 {
		return fmt.Errorf("amount must be positive")
	}
	if fromID == toID {
		return fmt.Errorf("from and to must differ")
	}

	return db.Transaction(func(tx *gorm.DB) error {
		var from, to Account
		// lock two row
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&from, fromID).Error; err != nil {
			return err
		}
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&to, toID).Error; err != nil {
			return err
		}

		if from.Balance < amount {
			return errors.New("insufficient balance")
		}

		// execute transaciton
		// minus from money
		// plus to money
		if err := tx.Model(&Account{}).
			Where("id = ? AND balance >= ?", fromID, amount).
			UpdateColumn("balance", gorm.Expr("balance - ?", amount)).Error; err != nil {
			return err
		}

		if err := tx.Model(&Account{}).
			Where("id = ?", toID).
			UpdateColumn("balance", gorm.Expr("balance + ?", amount)).Error; err != nil {
			return err
		}
		// record a transation
		rec := Transfer{FromAccountID: fromID, ToAccountID: toID, Amount: amount}
		if err := tx.Create(&rec).Error; err != nil {
			return err
		}

		return nil
	})
}

func StdDB(db *gorm.DB) (*sql.DB, error) {
	return db.DB()
}

func Close(db *DB) error {
	sdb, err := db.DB.DB()
	if err != nil {
		return err
	}
	return sdb.Close()
}
