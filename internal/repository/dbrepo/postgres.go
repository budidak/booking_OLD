package dbrepo

import (
	"context"
	"time"

	"github.com/budidak/booking/internal/models"
)

func (m *postgresDBRepo) AllUsers() bool {
	return true
}

// handlers.go içerisinde PostReservation() içinde güncelleme yaptık. Bu ilk yazdığımız fonksiyon. Bunda room id return etmiyoruz.
/*
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
*/

// Bu fonksiyon room id, error return edecek.
// Dikkat et, bu fonksiyon rezervasyon için yer ayırıp database'e yazacak ama reservations kısmında foreign key olarak roomID bağlıydı.
// Bu bağlılık nedeniyle, rooms table boş iken reservations table'a  veri eklemeye çalışırsak hata verir.
// Dolayısıyla database'de rooms table açıp önce bir satır değer ekledik.
func (m *postgresDBRepo) InsertReservation(res models.Reservation) (int, error) {
	// web apps are unpredictable... What happens if;
	// - user lost the connection?
	// - user closes the browser while something processing in the background?
	// We use "contexts" for situations like these.
	// Bu fonksiyonun işlemi halletmesi 3 saniyeden fazla sürerse cancel edecek, arkaplanda context çalışacak.
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `insert into reservations (first_name, last_name, email, phone, start_date, end_date, room_id, created_at, updated_at) values ($1, $2, $3, $4, $5, $6, $7, $8, $9) returning id`
	// Burada sql query'imizi değiştirdik, sonunda "returning id" kısmı var. Query'den sonra id değerini return ediyor. Bunu da aşağıda .Scan() ile istediğimiz değişkene atıyoruz.

	var newID int // rezervasyon yapılan odaya newID atayacağız.

	// query sonucunda bir değişkene değer atamak için .QueryRowContext().Scan() kullanabiliriz.
	err := m.DB.QueryRowContext(ctx, query,
		res.FirstName,
		res.LastName,
		res.Email,
		res.Phone,
		res.StartDate,
		res.EndDate,
		res.RoomID,
		time.Now(),
		time.Now(),
	).Scan(&newID) // newID değişkenine query'deki return değer olan ID değerini atadık.
	if err != nil {
		return 0, err
	}
	return newID, nil
}

// Bunu çalıştırınca RoomRestrictions table'a veri ekleyecek. Ama o tabloda foreign keys var. Bunların karşılığı olmalı, yani Rooms, Reservations, Restrictions tabloları buraya veri eklemeden önce boş olmamalı.
func (m *postgresDBRepo) InsertRoomRestriction(r models.RoomRestriction) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `inser into room_restrictions (start_date, end_date, room_id, reservation_id, created_at, updated_at, restriction_id) values ($1, $2, $3, $4, $5, $6, $7)`

	_, err := m.DB.ExecContext(ctx, query,
		r.StartDate,
		r.EndDate,
		r.RoomID,
		r.ReservationID,
		time.Now(),
		time.Now(),
		r.RestrictionID,
	)

	if err != nil {
		return err
	}

	return nil
}
