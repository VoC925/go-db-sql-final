package internal

import (
	"fmt"
	"time"

	"github.com/Yandex-Practicum/go-db-sql-final/pkg/logging"
)

// Интрефейс обертка сервиса.
type Service interface {
	Register(int, string) (Parcel, error)
	PrintClientParcels(int) error
	NextStatus(int) error
	ChangeAddress(int, string) error
	Delete(int) error
}

var _ Service = &ParcelService{}

// Статусы заказа.
const (
	ParcelStatusRegistered = "registered"
	ParcelStatusSent       = "sent"
	ParcelStatusDelivered  = "delivered"
)

// Структура сервиса заказов.
type ParcelService struct {
	store  Storage
	logger *logging.Logger
}

// NewParcelService возвращает экземпляр интерфейса на основе переданного хранилища.
func NewParcelService(store Storage, logger *logging.Logger) Service {
	return &ParcelService{
		store:  store,
		logger: logger,
	}
}

// Register создает новый заказ на основе основе идентификатора клиента client и адреса address.
func (s ParcelService) Register(client int, address string) (Parcel, error) {
	// модель заказа
	parcel := Parcel{
		Client:    client,
		Status:    ParcelStatusRegistered,
		Address:   address,
		CreatedAt: time.Now().UTC().Format(time.RFC3339),
	}
	// добавление в ДБ нового заказа
	id, err := s.store.Add(parcel)
	if err != nil {
		s.logger.Error("ошибка добавления посылки: " + err.Error())
		return parcel, err
	}

	parcel.Number = id // номер посылки

	fmt.Printf("Новая посылка № %d на адрес %s от клиента с идентификатором %d зарегистрирована %s\n",
		parcel.Number, parcel.Address, parcel.Client, parcel.CreatedAt)
	s.logger.Infof(`Добавлен заказ № %d с клиентом %d по адресу "%s", статус: %s`, id, client, address, parcel.Status)
	return parcel, nil
}

// PrintClientParcels показывает все заказы с идентификатором client.
func (s ParcelService) PrintClientParcels(client int) error {
	// извлечение из ДБ всех заказов с идентификатором client
	parcels, err := s.store.GetByClient(client)
	if err != nil {
		s.logger.Errorf("ошибка поиска посылки с идентификатором клиента № %d : %s", client, err.Error())
		return err
	}

	fmt.Printf("Посылки клиента %d:\n", client)
	for _, parcel := range parcels {
		fmt.Printf("Посылка № %d на адрес %s от клиента с идентификатором %d зарегистрирована %s, статус %s\n",
			parcel.Number, parcel.Address, parcel.Client, parcel.CreatedAt, parcel.Status)
	}
	fmt.Println()
	s.logger.Debugf("Заказы с идентификатором клиента № %d найдены в ДБ", client)
	return nil
}

// NextStatus возвращает заказ с номером number и присвает ему новый статус.
func (s ParcelService) NextStatus(number int) error {
	parcel, err := s.store.Get(number)
	if err != nil {
		s.logger.Errorf("ошибка получения посылки № %d: %s", number, err.Error())
		return err
	}
	// проверка на существование заказа в БД
	if parcel == nil {
		s.logger.Infof("Посылка с №: %d отсутствует в БД", number)
		return nil
	}

	var nextStatus string
	switch parcel.Status {
	case ParcelStatusRegistered:
		nextStatus = ParcelStatusSent
	case ParcelStatusSent:
		nextStatus = ParcelStatusDelivered
	case ParcelStatusDelivered:
		return nil
	}

	fmt.Printf("У посылки № %d новый статус: %s\n", number, nextStatus)
	s.logger.Infof(`Статус заказа № %d изменен на "%s"`, number, nextStatus)
	return s.store.SetStatus(number, nextStatus)
}

// ChangeAddress изменяет адрес у заказа number на новый.
func (s ParcelService) ChangeAddress(number int, address string) error {
	parcel, err := s.store.Get(number)
	if err != nil {
		s.logger.Errorf("ошибка получения посылки № %d: %s", number, err.Error())
		return err
	}
	// проверка на существование заказа в БД
	if parcel == nil {
		s.logger.Infof("Посылка с №: %d отсутствует в БД", number)
		return nil
	}
	s.logger.Infof(`Адрес заказа № %d изменен на "%s"`, number, address)
	return s.store.SetAddress(number, address)
}

// Delete удаляет заказ с номером number.
func (s ParcelService) Delete(number int) error {
	parcel, err := s.store.Get(number)
	if err != nil {
		s.logger.Errorf("ошибка получения посылки № %d: %s", number, err.Error())
		return err
	}
	// проверка на существование заказа в БД
	if parcel == nil {
		s.logger.Infof("Посылка с №: %d отсутствует в БД", number)
		return nil
	}
	if parcel.Status != ParcelStatusRegistered {
		s.logger.Infof(`Заказ № %d имеет статус "%s" не может быть удален`, number, parcel.Status)
		return nil
	}
	s.logger.Infof("Заказ № %d удален из БД", number)
	return s.store.Delete(number)
}
