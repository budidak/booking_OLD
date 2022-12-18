package repository

import (
	"time"

	"github.com/budidak/booking/internal/models"
)

// think interface as a contract. Önce burayı yazdık nasıl bir şey yapacağımızı tasarlarken.
type DatabaseRepo interface {
	AllUsers() bool
	InsertReservation(res models.Reservation) (int, error) // make-reservation sayfasındaki formu doldurup bilgileri db'ye post ederken kullanacağız.
	InsertRoomRestriction(r models.RoomRestriction) error  // belirli tarihler için oda rezerve edildiğinde o odaya restriction uygulayacağız, başkası seçemeyecek o tarihlerde.
	SearchAvailabilityForDatesByRoomID(start, end time.Time, roomID int) (bool, error)
}
