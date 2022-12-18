package models

import "time"

// Database için kullanacağımız modelleri burada tanımlıyoruz. (tables)
type User struct {
	ID          int
	FirstName   string
	LastName    string
	Email       string
	Password    string
	AccessLevel int
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type Room struct {
	ID        int
	RoomName  string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Restriction struct {
	ID              int
	RestrictionName string
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

type Reservation struct {
	ID        int
	FirstName string
	LastName  string
	Email     string
	Phone     string
	StartDate time.Time
	EndDate   time.Time
	RoomID    int
	Room      Room // illa tabloyu olduğu gibi aktarmak zorunda değiliz, kendi oluşturduğumuz tipleri de ekstra bilgi için ekleyebiliriz.
	CreatedAt time.Time
	UpdatedAt time.Time
}

type RoomRestriction struct {
	ID            int
	StartDate     time.Time
	EndDate       time.Time
	RoomID        int
	Room          Room
	ReservationID int
	Reservation   Reservation
	RestrictionID int
	Restriction   Restriction
	CreatedAt     time.Time
	UpdatedAt     time.Time
}
