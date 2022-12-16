package repository

import "github.com/budidak/booking/internal/models"

// think interface as a contract. Önce burayı yazdık nasıl bir şey yapacağımızı tasarlarken.
type DatabaseRepo interface {
	AllUsers() bool
	InsertReservation(res models.Reservation) (int, error) // make-reservation sayfasındaki formu doldurup bilgileri db'ye post ederken kullanacağız.
	InsertRoomRestriction(r models.RoomRestriction) error
}
