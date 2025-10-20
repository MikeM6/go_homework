// Package main
package main

import (
	"context"
	"fmt"
	"log"

	advancegorm "advancegorm"
	basecrud "basecrud"
	querywithsqlx "querywithsqlx"
	transaction "transaction"
)

func main() {

	dsn := "root:root@tcp(127.0.0.1:3306)/gorm?charset=utf8mb4&parseTime=True&loc=Local"

	baseCRUD(dsn)
	transactionExam(dsn)

	sqlxQueryEmp(dsn)
	sqlxQueryBooks(dsn)

	blogAssociationQuery(dsn)
	hookExam(dsn)

}

func hookExam(dsn string) {

	db, err := advancegorm.NewMySQL(dsn)
	if err != nil {
		log.Fatalf("connect db: %v", err)
	}
	log.Println("connected to MySQL with sqlx")

	sqlDB, err := advancegorm.StdDB(db)
	if err != nil {
		log.Fatal(err)
	}
	defer sqlDB.Close()

	err = advancegorm.CreateOnePost(db)
	if err != nil {
		log.Fatal(err)
	}

	advancegorm.DeleteOneComment(db)
}

func blogAssociationQuery(dsn string) {

	db, err := advancegorm.NewMySQL(dsn)
	if err != nil {
		log.Fatalf("connect db: %v", err)
	}
	log.Println("connected to MySQL with sqlx")

	sqlDB, err := advancegorm.StdDB(db)
	if err != nil {
		log.Fatal(err)
	}
	defer sqlDB.Close()

	err = advancegorm.AutoMigrate(db)
	if err != nil {
		log.Fatal("Create Table has Fail: ", err)
	}

	err = advancegorm.BatchInsertBlogData(db)
	if err != nil {
		log.Fatal("Insert Blog Data Fail: ", err)
	}

	posts, err := advancegorm.GetUserPostsWithComments(db, 1)
	if err != nil {
		log.Fatalf("query user posts: %v", err)
	}

	log.Printf("user %d posts: %d\n", 1, len(posts))
	for _, p := range posts {
		log.Printf("Post #%d: %s | comments: %d\n", p.ID, p.Title, len(p.Comments))
	}

	post, cnt, err := advancegorm.GetPostWithMostComments(db)
	if err != nil {
		log.Fatalf("query top post: %v", err)
	}
	if post.ID == 0 {
		log.Println("no posts found")
		return
	}
	log.Printf("top post: id=%d title=%q comments=%d\n", post.ID, post.Title, cnt)
}

func sqlxQueryBooks(dsn string) {
	ctx := context.Background()
	db, err := querywithsqlx.NewDB(ctx, dsn)
	if err != nil {
		log.Fatalf("connect db: %v", err)
	}
	log.Println("connected to MySQL with sqlx")
	defer db.Close()

	err = querywithsqlx.CreateBooksTable(ctx, db)
	if err != nil {
		log.Fatalf("Create Book Table Issue: %v", err)
	}

	querywithsqlx.DoBatchInsertBooks(ctx, db)

	books, err := querywithsqlx.QueryExpensiveBooks(ctx, db, 100)
	if err != nil {
		log.Fatalf("Get Expensive Book Issue: %v", err)
	}
	fmt.Println(books)
}

func sqlxQueryEmp(dsn string) {
	ctx := context.Background()
	db, err := querywithsqlx.NewDB(ctx, dsn)
	if err != nil {
		log.Fatalf("connect db: %v", err)
	}
	log.Println("connected to MySQL with sqlx")
	defer db.Close()

	err = querywithsqlx.CreateEmployeesTable(ctx, db)
	if err != nil {
		log.Fatalf("Create Employee Table Issue: %v", err)
	}

	querywithsqlx.DoBatchInsertEmployees(ctx, db)

	emps, err := querywithsqlx.QueryTechDeptEmployees(ctx, db)
	if err != nil {
		log.Fatalf("Query Tech Dept Issue: %v", err)
	}
	fmt.Println("Query Tech Dept Result : ", len(emps))

	emp, err := querywithsqlx.GetTopPaidEmployee(ctx, db)
	if err != nil {
		log.Fatalf("Get Top Paid Emp Issue: %v", err)
	}
	fmt.Println("Top Paid Emp: ", emp)

}

func transactionExam(dsn string) {
	var (
		err error
	)
	// Transaction
	// connection db
	db, err := transaction.NewMySQL(dsn)
	if err != nil {
		log.Fatalf("connect db: %v", err)
	}
	log.Println("connected to MySQL")

	// // defer close db
	sqlDB, err := transaction.StdDB(db.DB)
	if err != nil {
		log.Fatal(err)
	}
	defer sqlDB.Close()

	err = transaction.AutoMigrate(db.DB)
	if err != nil {
		log.Fatalf("Create Table Has Issue", err)
	}

	err = transaction.CreateAccounts(db.DB, 1000, 0)
	if err != nil {
		log.Fatalf("Create Accounts Has Issue", err)
	}

	err = transaction.TransferAmount(db.DB, 1, 2, 1000)
	if err != nil {
		log.Fatalf("Transfer Issue", err)
	}

}

func baseCRUD(dsn string) {
	var (
		rowsAffected int64
		err          error
	)
	// get connection
	db, err := basecrud.NewMySQL(dsn)
	if err != nil {
		log.Fatalf("connect db: %v", err)
	}
	log.Println("connected to MySQL")

	// defer close db
	sqlDB, err := db.StdDB()
	if err != nil {
		log.Fatal(err)
	}
	defer sqlDB.Close()
	db.AutoMigrateModels()

	db.InsertStudent("Mike", 18, "1")

	students, err := db.QueryStudentsAgeGreaterThan(1)
	if err != nil {
		log.Fatalf("There are issue in QueryStudentsAgeGreaterThan", err)
	}
	log.Println(students)

	rowsAffected, err = db.UpdateStudentGradeByName("Mike", "2")
	if err != nil {
		log.Fatalf("There are issue in UpdateStudentGradeByName", err)
	}
	log.Println("Update rowNum:", rowsAffected)

	rowsAffected, err = db.DeleteStudentsAgeLessThan(18)
	if err != nil {
		log.Fatalf("There are issue in UpdateStudentGradeByName", err)
	}
	log.Println("Update rowNum:", rowsAffected)
}
