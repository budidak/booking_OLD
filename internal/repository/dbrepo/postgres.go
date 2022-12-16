package dbrepo

import (
	"context"
	"time"

	"github.com/budidak/booking/internal/models"
)

func (m *postgresDBRepo) AllUsers() bool {
	return true
}

// handlers.go içerisinde PostReservation() içinde güncelleme yaptık.
func (m *postgresDBRepo) InsertReservation(res models.Reservation) error {
	// web apps are unpredictable... What happens if;
	// - user lost the connection?
	// - user closes the browser while something processing in the background?
	// We use "contexts" for situations like these.
	// Bu fonksiyonun işlemi halletmesi 3 saniyeden fazla sürerse cancel edecek, arkaplanda context çalışacak.
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	query := `insert into reservations (first_name, last_name, email, phone, start_date, end_date, room_id, created_at, updated_at) values ($1, $2, $3, $4, $5, $6, $7, $8, $9)`

	// context kullanırken, Exec() yerine ExecContext() kullanmalıyız.
	_, err := m.DB.ExecContext(ctx, query,
		res.FirstName,
		res.LastName,
		res.Email,
		res.Phone,
		res.StartDate,
		res.EndDate,
		res.RoomID,
		time.Now(),
		time.Now(),
	)

	if err != nil {
		return err
	}

	return nil
}
