package main

import (
	"database/sql"
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	randSource = rand.NewSource(time.Now().UnixNano())
	randRange  = rand.New(randSource)
)

func getTestParcel() Parcel {
	return Parcel{
		Client:    1000,
		Status:    ParcelStatusRegistered,
		Address:   "test",
		CreatedAt: time.Now().UTC().Format(time.RFC3339),
	}
}

func TestAddGetDelete(t *testing.T) {
	db, err := sql.Open(DBDriver, DBFile)
	if err != nil {
		require.NoError(t, err)
	}
	defer db.Close()

	store := NewParcelStore(db)
	parcel := getTestParcel()

	id, err := store.Add(parcel)
	parcel.Number = id

	require.NoError(t, err)
	require.NotEmpty(t, id)

	getParcel, err := store.Get(id)
	require.NoError(t, err)
	assert.Equal(t, parcel, getParcel)

	err = store.Delete(id)
	require.NoError(t, err)

	_, err = store.Get(id)
	assert.ErrorIs(t, sql.ErrNoRows, err)
}

func TestSetAddress(t *testing.T) {
	db, err := sql.Open(DBDriver, DBFile)
	if err != nil {
		require.NoError(t, err)
	}
	defer db.Close()

	store := NewParcelStore(db)
	parcel := getTestParcel()

	id, err := store.Add(parcel)
	require.NoError(t, err)
	require.NotEmpty(t, id)

	newAddress := "new test address"
	err = store.SetAddress(id, newAddress)
	require.NoError(t, err)

	getParcel, err := store.Get(id)
	require.NoError(t, err)
	require.Equal(t, newAddress, getParcel.Address)
}

func TestSetStatus(t *testing.T) {
	db, err := sql.Open(DBDriver, DBFile)
	if err != nil {
		require.NoError(t, err)
	}
	defer db.Close()

	store := NewParcelStore(db)
	parcel := getTestParcel()

	id, err := store.Add(parcel)
	require.NoError(t, err)
	require.NotEmpty(t, id)

	err = store.SetStatus(id, ParcelStatusSent)
	require.NoError(t, err)

	getParcel, err := store.Get(id)
	require.NoError(t, err)
	require.Equal(t, ParcelStatusSent, getParcel.Status)
}

func TestGetByClient(t *testing.T) {
	db, err := sql.Open(DBDriver, DBFile)
	if err != nil {
		require.NoError(t, err)
	}
	defer db.Close()

	store := NewParcelStore(db)

	parcels := []Parcel{
		getTestParcel(),
		getTestParcel(),
		getTestParcel(),
	}

	client := randRange.Intn(10_000_000)
	for indx, parsc := range parcels {
		parsc.Client = client

		id, err := store.Add(parsc)
		require.NoError(t, err)

		parcels[indx].Client = client
		parcels[indx].Number = id

	}

	storedParcels, err := store.GetByClient(client)
	require.NoError(t, err)
	assert.ElementsMatch(t, parcels, storedParcels)
}
