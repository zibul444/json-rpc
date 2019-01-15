package main

import (
	"encoding/hex"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/satori/go.uuid"
	"testing"
)

var c Client

func TestDoInsert(t *testing.T) {
	DoInsert(c)
	var received uuid.UUID

	db := InitConnDB(driverName)
	defer db.Close()
	rows, err := db.Query("SELECT uuid FROM client WHERE uuid = $1;", c.uuid)
	if err != nil {
		t.Error(err)
	}

	rows.Next()
	rows.Scan(&received)

	if c.uuid != received {
		t.Error("Expected", c.uuid, "got", received)
	}
}

func TestDoUpdate(t *testing.T) {
	rowsAffected := DoUpdate(toStringUUID(c.uuid), "TestDoUpdate")

	if rowsAffected != 1 {
		t.Error("Expected", 1, "got", rowsAffected)
	}
}

func TestDoSelect(t *testing.T) {
	rowsAffected, _ := DoSelect(fmt.Sprintf("SELECT * FROM client WHERE uuid = '%s'",
		c.uuid))

	if rowsAffected < 1  {
		t.Error("Expected >", 1, "got", rowsAffected)
	}
}


func toStringUUID(u uuid.UUID) string {
	buf := make([]byte, 36)

	hex.Encode(buf[0:8], u[0:4])
	buf[8] = '-'
	hex.Encode(buf[9:13], u[4:6])
	buf[13] = '-'
	hex.Encode(buf[14:18], u[6:8])
	buf[18] = '-'
	hex.Encode(buf[19:23], u[8:10])
	buf[23] = '-'
	hex.Encode(buf[24:], u[10:])

	return string(buf)
}

//func setup() {
//	fmt.Println("after test")
//	c = Client{login: "test_client", uuid: uuid.Must(uuid.NewV4())}
//}
//
//func shutdown() {
//	fmt.Println("before test")
//	db := InitConnDB(driverName)
//	db.Exec("truncate table public.client")
//
//}
//
//func TestMain(m *testing.M) {
//	setup()
//	code := m.Run()
//	shutdown()
//	os.Exit(code)
//}

