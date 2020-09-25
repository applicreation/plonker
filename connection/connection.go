package connection

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"math"
	"os"
	"os/signal"
	"time"

	"github.com/applicreation/plonker/config"
)

type Connection struct {
	Config config.Connection
	Db     *sql.DB
}

func (c *Connection) Count(table config.Table) (int, error) {
	query := fmt.Sprintf("SELECT COUNT(%s) AS count FROM `%s`", table.PrimaryKey, table.Name)

	if table.Timeframe.Period != "" && table.Timeframe.Column != "" {
		now := time.Now()
		yesterday := now.AddDate(0, 0, -1)

		query += fmt.Sprintf(" WHERE `%s` >= '%s'", table.Timeframe.Column, yesterday.Format("2006-01-02 15:04:05"))
	} else if table.Records > 0 {
		query += fmt.Sprintf(" LIMIT %d", table.Records)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	log.Printf("Count: %s", query)

	var count int
	err := c.Db.QueryRowContext(ctx, query).Scan(&count)
	if err != nil {
		return 0, err
	}

	if table.Percentage > 0 && table.Percentage < 100 {
		return int(math.Round(((float64(count) / float64(100)) * float64(table.Percentage)) + 0.5)), nil
	}

	return count, nil
}

func (c *Connection) RelationshipCount(table config.Table, relationship config.Relation, parentCount int) (int, error) {
	query := fmt.Sprintf("SELECT COUNT(DISTINCT `%s`) FROM `%s`", relationship.Key, table.Name)

	if table.Timeframe.Period != "" && table.Timeframe.Column != "" {
		now := time.Now()
		yesterday := now.AddDate(0, 0, -1)

		query += fmt.Sprintf(" WHERE `%s` >= '%s'", table.Timeframe.Column, yesterday.Format("2006-01-02 15:04:05"))
	} else if table.Records > 0 {
		query += fmt.Sprintf(" LIMIT %d", table.Records)
	} else if table.Percentage > 0 {
		query += fmt.Sprintf(" LIMIT %d", parentCount)
	} else {
		// error
	}

	if table.Order.Column != "" && table.Order.Direction != "" {
		query += fmt.Sprintf(" ORDER BY %s %s", table.Order.Column, table.Order.Direction)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	log.Printf("RelationshipCount: %s", query)

	var count int
	err := c.Db.QueryRowContext(ctx, query).Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (c *Connection) RelationshipKeys(table config.Table, relationship config.Relation) ([]string, error) {
	query := fmt.Sprintf("SELECT DISTINCT `%s` FROM `%s`", relationship.Key, table.Name)

	if table.Timeframe.Period != "" && table.Timeframe.Column != "" {
		now := time.Now()
		yesterday := now.AddDate(0, 0, -1)

		query += fmt.Sprintf(" WHERE `%s` >= '%s'", table.Timeframe.Column, yesterday.Format("2006-01-02 15:04:05"))
	} else if table.Records > 0 {
		query += fmt.Sprintf(" LIMIT %d", table.Records)
	} else if table.Percentage > 0 {
		parentCount, _ := c.Count(table)
		query += fmt.Sprintf(" LIMIT %d", parentCount)
	} else {
		// error
	}

	if table.Order.Column != "" && table.Order.Direction != "" {
		query += fmt.Sprintf(" ORDER BY %s %s", table.Order.Column, table.Order.Direction)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	log.Printf("RelationshipKeys: %s", query)

	var count int
	err := c.Db.QueryRowContext(ctx, query).Scan(&count)
	if err != nil {
		return nil, err
	}



	return nil, nil
}

func (c *Connection) FindAll(table config.Table) ([]interface{}, error) {
	return nil, nil
}

func GetConnection(c *config.Connection) *sql.DB {
	log.Printf("Connecting to %s:%d/%s with username %s", c.Host, c.Port, c.Name, c.Username)

	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s",
		c.Username, c.Password, c.Host, c.Port, c.Name,
	)

	db, err := sql.Open(c.Engine, dsn)
	if err != nil {
		log.Fatal("Unable to use data source name", err)
	}

	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)

	ctx, stop := context.WithCancel(context.Background())
	defer stop()

	appSignal := make(chan os.Signal, 3)
	signal.Notify(appSignal, os.Interrupt)

	go func() {
		select {
		case <-appSignal:
			stop()
		}
	}()

	ctx, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()

	err = db.PingContext(ctx)
	if err != nil {
		log.Fatalf("unable to connect to database: %v", err)
	}

	return db
}
