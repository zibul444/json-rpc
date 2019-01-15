package main

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/satori/go.uuid"
	"log"
	"math/rand"
	"net"
	"net/rpc"
	"net/rpc/jsonrpc"
	"time"
)

var connStr = "host=192.168.1.68 user=postgres password=`12`12 dbname=RPC_DB sslmode=disable"
var separator = "------------------------------------------------------------"
var driverName = "postgres"

type Client struct {
	login        string    `json:"login"`
	uuid         uuid.UUID `json:"uuid"`
	registration time.Time `json:"registration"`
}

func (this *Client) RPCInsert(i int64, reply *int64) error {
	fmt.Print("Insert ")
	u := uuid.Must(uuid.NewV4())
	c := Client{login: fmt.Sprintf("client_%d", i), uuid: u}

	r := DoInsert(c)
	*reply = r

	return nil
}

func (this *Client) RPCUpdate(c [2]string, reply *int64) error {
	fmt.Print("Update ")
	r := DoUpdate(c[0], c[1])
	*reply = r

	return nil
}

func (this *Client) RPCSelect(s string, reply *int64) error {
	fmt.Println("Select: ")

	r, clientsList := DoSelect(fmt.Sprintf("select %s from client", s))
	*reply = r

	for _, c := range clientsList {
		fmt.Printf("login: %s, uuid: %s, registration: %s \n", c.login, c.uuid, c.registration)
	}

	return nil
}

func server() {
	server := rpc.NewServer()
	server.Register(new(Client))

	ln, err := net.Listen("tcp", ":8800")
	if err != nil {
		log.Fatal(err)
		return
	}
	for {
		c, err := ln.Accept()
		if err != nil {
			log.Fatal(err)
			continue
		}
		go server.ServeCodec(jsonrpc.NewServerCodec(c))
	}
}

func client() {
	conn, err := net.Dial("tcp", ":8800")
	if err != nil {
		log.Fatal(err)
		return
	}
	c := jsonrpc.NewClient(conn)

	rand.Seed(time.Now().UnixNano())
	i := rand.Int31n(1000)
	var result int64
	fmt.Printf("\033[33m%s\033[m\n", separator)

	err = c.Call("Client.RPCInsert", int64(i), &result)
	if err != nil {
		log.Fatal(err)
	} else {
		fmt.Println("rowsAffected: ", result)
	}

	fmt.Printf("\033[34m%s\033[m\n", separator)

	err = c.Call("Client.RPCUpdate",
		[2]string{"ff5a62e5-101e-48aa-939c-ea9a8aadbd67", fmt.Sprint("client_", i-1)},
		//Client{uuid: uuid, login: fmt.Sprint("client_", i-1)},
		&result)
	if err != nil {
		log.Fatal(err)
	} else {
		fmt.Println("rowsAffected: ", result)
	}

	fmt.Printf("\033[34m%s\033[m\n", separator)

	err = c.Call("Client.RPCSelect", "*",
		&result)
	if err != nil {
		log.Fatal(err)
	} else {
		fmt.Println("Count rows = ", result)
	}

	fmt.Printf("\033[31m%s\033[m\n", separator)

}

func InitConnDB(driverName string) (db *sql.DB) {
	db, err := sql.Open(driverName, connStr)
	if err != nil {
		log.Fatal(err)
	}
	return
}

func DoInsert(client Client) (rowsAffected int64) {
	db := InitConnDB(driverName)
	defer db.Close()

	result, err := db.Exec("Insert into Client (login, uuid, registration) values ($1, $2, $3)",
		client.login, client.uuid, time.Now())
	if err != nil {
		log.Fatal(err)
	}

	rowsAffected, err = result.RowsAffected()
	return
}

func DoUpdate(uuid string, newLogin string) (rowsAffected int64) {
	db := InitConnDB(driverName)
	defer db.Close()

	result, err := db.Exec("UPDATE Client SET login = $2 WHERE uuid = $1;",
		uuid, newLogin)
	if err != nil {
		log.Fatal(err)
	}

	rowsAffected, err = result.RowsAffected()

	return
}

func DoSelect(query string) (count int64, clients []Client) {
	db := InitConnDB(driverName)
	defer db.Close()

	rows, err := db.Query(query)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	for rows.Next() {
		c := Client{}
		err := rows.Scan(&c.login, &c.uuid, &c.registration)
		if err != nil {
			log.Fatal(err)
			continue
		}
		clients = append(clients, c)
	}

	count = int64(len(clients))
	return
}

func main() {
	go server()
	go client()

	var input string
	fmt.Scanln(&input)
}
