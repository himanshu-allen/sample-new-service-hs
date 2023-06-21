package main

import (
	"flag"
	"fmt"
	"github.com/Allen-Career-Institute/go-kratos-sample/internal/conf"
	"github.com/Allen-Career-Institute/go-kratos-sample/internal/data/entity"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ssm"
	"github.com/go-kratos/kratos/v2/config"
	"github.com/go-kratos/kratos/v2/config/file"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
)

var (
	// flagconf is the config flag.
	flagconf string
)

func init() {
	flag.StringVar(&flagconf, "conf", "../../../configs/config.yaml", "config path, eg: -conf config.yaml")
}

func main() {
	flag.Parse()

	// Load configurations.
	c := config.New(
		config.WithSource(
			file.NewSource(flagconf),
		),
	)
	defer c.Close()

	if err := c.Load(); err != nil {
		panic(err)
	}

	var bc conf.Bootstrap
	if err := c.Scan(&bc); err != nil {
		panic(err)
	}

	// TODO: Fetch the config from Common library
	sess, _ := session.NewSession(&aws.Config{
		Region: aws.String("ap-south-1")},
	)

	svc := ssm.New(sess)

	dbParameter, _ := svc.GetParameter(&ssm.GetParameterInput{Name: aws.String("testUserDb"),
		WithDecryption: aws.Bool(true)})

	dbConn := dbParameter.Parameter.Value

	dsn := fmt.Sprintf("%s?charset=utf8mb4&parseTime=True&loc=Local", *dbConn)

	// Open connection to the database
	db, err := gorm.Open(mysql.Open(dsn))
	if err != nil {
		log.Fatalf("Unable to open db connection: %v", err)
	}

	log.Printf("DB Connection opened...")
	if err := db.Set("gorm:table_options", "ENGINE=InnoDB").Table("users").AutoMigrate(&entity.UserEntity{}); err != nil {
		if err.Error() != "record not found" {
			log.Fatalf("Error occurred during AutoMigrate: %v", err)
		}
	}
	log.Printf("Migrations ran successfully...")
	db.Migrator().HasTable("users")
	log.Printf("Table Verified...SUCCESS!")
}
