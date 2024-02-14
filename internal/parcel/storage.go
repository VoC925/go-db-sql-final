package parcel

// Интерфейс хранилища.
type Storage interface {
	Add(*Parcel) (int, error)
	Get(int) (*Parcel, error)
	GetByClient(int) ([]*Parcel, error)
	SetStatus(int, string) error
	SetAddress(int, string) error
	Delete(int) error
}
