package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"http-server/danilkovalev/internal/models"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"time"

	"github.com/urfave/cli/v2"
	"github.com/beanstalkd/go-beanstalk"
	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	db            *gorm.DB
	beanstalkConn *beanstalk.Conn
)

func main() {
	app := &cli.App{
		Name:  "workerApp",
		Usage: "CLI for managing workers and jobs",
		Flags: []cli.Flag{
			&cli.IntFlag{
				Name:     "workers",
				Aliases:  []string{"w"},
				Value:    1,
				Usage:    "Number of workers to start",
				Required: true,
			},
		},
		Action: func(c *cli.Context) error {
			loadEnvironment()
			db = mustInitDB()
			beanstalkConn = mustInitBeanstalkClient()
			defer beanstalkConn.Close()

			numWorkers := c.Int("workers")
			startWorkers(numWorkers)

			select {}
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func loadEnvironment() {
	if err := godotenv.Load("../.env"); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
}

func mustInitDB() *gorm.DB {
	db, err := initDB()
	if err != nil {
		log.Fatalf("Failed to initialize MySQL: %v", err)
	}
	return db
}

func mustInitBeanstalkClient() *beanstalk.Conn {
	conn, err := initBeanstalkClient()
	if err != nil {
		log.Fatalf("Failed to create Beanstalk client: %v", err)
	}
	return conn
}

func startWorkers(numWorkers int) {
	for i := 0; i < numWorkers; i++ {
		cmd := exec.Command(os.Args[0], "run")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		if err := cmd.Start(); err != nil {
			log.Printf("Error starting worker %d: %v", i+1, err)
			continue
		}

		fmt.Printf("Started worker %d with PID %d\n", i+1, cmd.Process.Pid)
	}
}

func initDB() (*gorm.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		"localhost",
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"))

	return gorm.Open(mysql.Open(dsn), &gorm.Config{})
}

func initBeanstalkClient() (*beanstalk.Conn, error) {
	return beanstalk.Dial("tcp", "localhost:11300")
}

func runWorker() {
	db = mustInitDB()
	beanstalkConn = mustInitBeanstalkClient()
	defer beanstalkConn.Close()

	processJobs()
}

func processJobs() {
	for {
		id, body, err := beanstalkConn.Reserve(5 * time.Second)
		if err != nil {
			handleReservationError(err)
			continue
		}

		if err := processJob(id, body); err != nil {
			log.Printf("Error processing job: %v", err)
			beanstalkConn.Release(id, 1, 0)
			continue
		}

		if err := beanstalkConn.Delete(id); err != nil {
			log.Printf("Failed to delete job: %v", err)
		}
	}
}

func handleReservationError(err error) {
	if err == beanstalk.ErrTimeout {
		log.Println("No jobs available, waiting...")
	} else {
		log.Printf("Error reserving job: %v", err)
	}
}

func processJob(id uint64, body []byte) error {
	if contactChange, err := unmarshalContactChange(body); err == nil {
		return handleContactChange(contactChange)
	} else if task, err := unmarshalTask(body); err == nil {
		return handleTask(task)
	}
	log.Printf("Unknown task type or failed to unmarshal")
	return nil
}

func unmarshalContactChange(body []byte) (models.ContactChange, error) {
	var contactChange models.ContactChange
	err := json.Unmarshal(body, &contactChange)
	return contactChange, err
}

func unmarshalTask(body []byte) (models.Task, error) {
	var task models.Task
	err := json.Unmarshal(body, &task)
	return task, err
}

func handleContactChange(contactChange models.ContactChange) error {
	contact := models.Contact{
		Name:      contactChange.Name,
		Email:     contactChange.Email,
		AccountID: contactChange.AccountID,
		ID:        contactChange.ID,
	}

	switch contactChange.TypeChange {
	case "update":
		return UpdateContact(contact)
	case "add":
		return AddContact(contact)
	case "delete":
		return DeleteContact(contact)
	default:
		log.Printf("Unknown contact change type: %s", contactChange.TypeChange)
		return nil
	}
}

func handleTask(task models.Task) error {
	if err := AddUnisenderKey(task.AccountID, task.UnisenderKey); err != nil {
		return fmt.Errorf("failed to add Unisender key: %w", err)
	}

	return ToUnisender(task.AccountID, task.UnisenderKey)
}

func AddContact(contact models.Contact) error {
	return db.Create(&contact).Error
}

func DeleteContact(contact models.Contact) error {
	return db.Delete(&contact).Error
}

func UpdateContact(contact models.Contact) error {
	return db.Save(&contact).Error
}

func ToUnisender(accountID int, key string) error {
	emails, err := GetEmails(accountID)
	if err != nil {
		return fmt.Errorf("failed to get emails: %w", err)
	}

	apiURL := "https://api.unisender.com/ru/api/importContacts"
	params := createUnisenderParams(key, emails)

	resp, err := http.Post(apiURL+"?"+params.Encode(), "application/x-www-form-urlencoded", bytes.NewBuffer(nil))
	if err != nil {
		return fmt.Errorf("failed to send request to Unisender: %w", err)
	}
	defer resp.Body.Close()

	return nil
}

func createUnisenderParams(key string, emails []string) url.Values {
	params := url.Values{
		"format":  {"json"},
		"api_key": {key},
	}

	for i, email := range emails {
		params.Set(fmt.Sprintf("data[%d][0]", i), email)
	}
	return params
}

func GetEmails(accountID int) ([]string, error) {
	var contacts []models.Contact

	err := db.Where("account_id = ?", accountID).Find(&contacts).Error
	if err != nil {
		return nil, fmt.Errorf("failed to query emails: %w", err)
	}

	emails := make([]string, len(contacts))
	for i, contact := range contacts {
		emails[i] = contact.Email
	}
	return emails, nil
}

func AddUnisenderKey(accountID int, key string) error {
	var account models.Account

	if err := db.Where("account_id = ?", accountID).First(&account).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("account not found")
		}
		return err
	}

	account.UnisenderKey = key
	return db.Save(&account).Error
}

func init() {
	if len(os.Args) > 1 && os.Args[1] == "run" {
		runWorker()
		os.Exit(0)
	}
}
