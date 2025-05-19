package db

import (
	"database/sql"
	"fmt"
	"weather-forecast/internal/config"
	"weather-forecast/internal/model"

	_ "github.com/lib/pq"
)

func Connect() (*sql.DB, error) {
	cfg := config.Load()

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.DB.Host, cfg.DB.Port, cfg.DB.User, cfg.DB.Password, cfg.DB.DBName)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	return db, db.Ping()
}

func InsertSubscription(db *sql.DB, sub model.Subscription) error {
	query := `INSERT INTO subscriptions (email, city, frequency, confirmed) VALUES ($1, $2, $3, $4)`
	_, err := db.Exec(query, sub.Email, sub.City, sub.Frequency, sub.Confirmed)
	return err
}

func GetSubscriptionByEmailCity(db *sql.DB, email, city string) (*model.Subscription, error) {
	query := `SELECT id, email, city, frequency, confirmed FROM subscriptions WHERE email=$1 AND city=$2`
	row := db.QueryRow(query, email, city)

	var sub model.Subscription
	err := row.Scan(&sub.ID, &sub.Email, &sub.City, &sub.Frequency, &sub.Confirmed)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &sub, nil
}

func MarkSubscriptionConfirmed(db *sql.DB, email string) error {
	query := `UPDATE subscriptions SET confirmed=true WHERE email=$1`
	_, err := db.Exec(query, email)
	return err
}

func DeleteSubscription(db *sql.DB, email string) error {
	query := `DELETE FROM subscriptions WHERE email=$1`
	_, err := db.Exec(query, email)
	return err
}
