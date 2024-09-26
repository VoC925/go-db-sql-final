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
	require.NoError(t, err)
	require.NotEqual(t, id, -1)

	// get
	parcelActual, err := store.Get(id)
	require.NoError(t, err)
	require.Equal(t, parcel.Address, parcelActual.Address)
	require.Equal(t, parcel.Client, parcelActual.Client)
	require.Equal(t, parcel.CreatedAt, parcelActual.CreatedAt)
	require.Equal(t, parcel.Status, parcelActual.Status)

	// delete
	err = store.Delete(id)
	require.NoError(t, err)

	// get удаленной посылки
	var parcelDeletedExpected *internal.Parcel // nil значение
	parcelDeletedActual, err := store.Get(id)
	require.Equal(t, parcelDeletedExpected, parcelDeletedActual)
	require.NoError(t, err)
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
	require.NotEqual(t, id, -1)

	// set address
	newAddress := "new test address"
	err = store.SetAddress(id, newAddress)
	require.NoError(t, err)

	// check
	parcelActual, err := store.Get(id)
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
	require.NotEqual(t, id, -1)

	// set status
	newStatus := internal.ParcelStatusSent
	err = store.SetStatus(id, newStatus)
	require.NoError(t, err)

	// check
	parcelActual, err := store.Get(id)
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
	// мапа [id посылки] посылка
	parcelMap := map[int]internal.Parcel{}

	// задаём всем посылкам один и тот же идентификатор клиента
	client := randRange.Intn(10_000_000) // случайное число из диапазон [0,9999999)
	parcels[0].Client = client
	parcels[1].Client = client
	parcels[2].Client = client

	// add
	for i := 0; i < len(parcels); i++ {
		id, err := store.Add(parcels[i])
		require.NoError(t, err)
		require.NotEqual(t, id, -1)

		// обновляем идентификатор добавленной у посылки
		parcels[i].Number = id

		// сохраняем добавленную посылку в структуру map, чтобы её можно было легко достать по идентификатору посылки
		parcelMap[id] = parcels[i]
	}

	// get by client
	storedParcels, err := store.GetByClient(client)
	require.NoError(t, err)
	// проверка количества посылок
	require.Equal(t, len(parcels), len(storedParcels))

	// check
	for _, parcel := range storedParcels {
		// проверка наличия в мапе
		if _, ok := parcelMap[parcel.Number]; !ok {
			require.True(t, ok)
		}
		require.Equal(t, parcelMap[parcel.Number].Address, parcel.Address)
		require.Equal(t, parcelMap[parcel.Number].Client, parcel.Client)
		require.Equal(t, parcelMap[parcel.Number].CreatedAt, parcel.CreatedAt)
		require.Equal(t, parcelMap[parcel.Number].Number, parcel.Number)
		require.Equal(t, parcelMap[parcel.Number].Status, parcel.Status)
	}
}
