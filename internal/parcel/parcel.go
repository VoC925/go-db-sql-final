package parcel

// Модель заказа
type Parcel struct {
	Number    int    // номер посылки
	Client    int    // идентификатор клиента
	Status    string // статус посылки
	Address   string // адрес посылки
	CreatedAt string // дата и время создания посылки
}
