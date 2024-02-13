package db

import (
	"database/sql"
	"math/rand"
	"testing"
	"time"

	"github.com/Yandex-Practicum/go-db-sql-final/internal"
	"github.com/stretchr/testify/require"

	_ "modernc.org/sqlite"
)

const pathToDB = "tracker_copy_test.db"

var (
	// randSource источник псевдо случайных чисел.
	// Для повышения уникальности в качестве seed
	// используется текущее время в unix формате (в виде числа)
	randSource = rand.NewSource(time.Now().UnixNano())
	// randRange использует randSource для генерации случайных чисел
	randRange = rand.New(randSource)
)

// getTestParcel возвращает тестовую посылку
func getTestParcel() internal.Parcel {
	return internal.Parcel{
		Client:    1000,
		Status:    internal.ParcelStatusRegistered,
		Address:   "test",
		CreatedAt: time.Now().UTC().Format(time.RFC3339),
	}
}

// TestAddGetDelete проверяет добавление, получение и удаление посылки
func TestAddGetDelete(t *testing.T) {
	db, err := sql.Open("sqlite", pathToDB)
	if err != nil {
		require.NoError(t, err)
	}
	defer db.Close()

	store := NewParcelStore(db)
	parcel := getTestParcel()

	// add
	id, err := store.Add(parcel)
	parcel.Number = id
	require.NoError(t, err)
	// Проверка id > 0
	require.Greater(t, id, 0)

	// get
	parcelActual, err := store.Get(id)
	// Добавил проверку на не nil (исправлено)
	require.NotNil(t, parcelActual)
	require.NoError(t, err)
	// Исправлена проверка равенства структур (исправлено)
	require.Equal(t, parcel, *parcelActual)

	// delete
	err = store.Delete(id)
	require.NoError(t, err)

	// get удаленной посылки
	parcelDeletedActual, err := store.Get(id)
	// Переделал проверку (исправлено)
	require.Nil(t, parcelDeletedActual)
	require.ErrorIs(t, err, internal.ErrEmptyData)
}

// TestSetAddress проверяет обновление адреса
func TestSetAddress(t *testing.T) {
	db, err := sql.Open("sqlite", pathToDB)
	if err != nil {
		require.NoError(t, err)
	}
	defer db.Close()

	store := NewParcelStore(db)
	parcel := getTestParcel()

	// add
	id, err := store.Add(parcel)
	require.NoError(t, err)
	// Проверка id > 0
	require.Greater(t, id, 0)

	// set address
	newAddress := "new test address"
	err = store.SetAddress(id, newAddress)
	require.NoError(t, err)

	// check
	parcelActual, err := store.Get(id)
	// Переделал проверку (исправлено)
	require.NotNil(t, parcelActual)
	require.NoError(t, err)
	require.Equal(t, newAddress, parcelActual.Address)
}

// TestSetStatus проверяет обновление статуса
func TestSetStatus(t *testing.T) {
	db, err := sql.Open("sqlite", pathToDB)
	if err != nil {
		require.NoError(t, err)
	}
	defer db.Close()

	store := NewParcelStore(db)
	parcel := getTestParcel()

	// add
	id, err := store.Add(parcel)
	require.NoError(t, err)
	// Проверка id > 0
	require.Greater(t, id, 0)

	// set status
	newStatus := internal.ParcelStatusSent
	err = store.SetStatus(id, newStatus)
	require.NoError(t, err)

	// check
	parcelActual, err := store.Get(id)
	// Переделал проверку (исправлено)
	require.NotNil(t, parcelActual)
	require.NoError(t, err)
	require.Equal(t, newStatus, parcelActual.Status)
}

// TestGetByClient проверяет получение посылок по идентификатору клиента
func TestGetByClient(t *testing.T) {
	db, err := sql.Open("sqlite", pathToDB)
	if err != nil {
		require.NoError(t, err)
	}
	defer db.Close()

	store := NewParcelStore(db)

	// слайс посылок
	parcels := []internal.Parcel{
		getTestParcel(),
		getTestParcel(),
		getTestParcel(),
	}

	// задаём всем посылкам один и тот же идентификатор клиента
	client := randRange.Intn(10_000_000) // случайное число из диапазон [0,9999999)
	parcels[0].Client = client
	parcels[1].Client = client
	parcels[2].Client = client

	// add
	for i := 0; i < len(parcels); i++ {
		id, err := store.Add(parcels[i])
		require.NoError(t, err)
		require.Greater(t, id, 0)
		// обновляем идентификатор добавленной у посылки
		parcels[i].Number = id
	}
	// get by client
	storedParcels, err := store.GetByClient(client)
	require.NotNil(t, storedParcels)
	require.NoError(t, err)
	// Изменена проверка равенства слайсов
	require.Equal(t, storedParcels, parcels)
}
