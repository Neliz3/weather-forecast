package model

type Subscription struct {
	ID        int    `db:"id"`
	Email     string `db:"email"`
	City      string `db:"city"`
	Frequency string `db:"frequency"`
	Confirmed bool   `db:"confirmed"`
}
