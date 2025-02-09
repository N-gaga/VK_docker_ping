package main

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"net"
	"os/exec"
	"time"
)

const (
	host     = "db"
	port     = 5432
	user     = "user"
	password = "password"
	dbname   = "db_ping"
)

var db *sql.DB

func connectToDB() error {
	psqlConn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	var err error
	db, err = sql.Open("postgres", psqlConn)
	if err != nil {
		return err
	}
	err = db.Ping()
	if err != nil {
		return err
	}
	return nil
}

func saveIpToDb(listContainer string, ip string, time time.Time) error {
	sqlStatement := `INSERT INTO ip_table (listContainer, ip_address, time) VALUES ($1, $2, $3)`
	_, err := db.Exec(sqlStatement, listContainer, ip, time)
	if err != nil {
		return err
	}
	return nil
}

func main() {
	err := connectToDB()
	if err != nil {
		panic(err)
	}
	defer db.Close()

	listContainers := []string{"backend", "db", "pinger"}

	for {
		for _, listContainer := range listContainers {
			ips, err := net.LookupIP(listContainer)
			if err != nil {
				fmt.Printf("Ошибка пинга к IP для %s: %v\n", listContainer, err)
				continue
			}

			ip := ips[0].String()
			now := time.Now()

			err = saveIpToDb(listContainer, ip, now)
			if err != nil {
				fmt.Printf("Ошибка при сохранении IP для %s (%s) в базу данных: %v\n", listContainer, ip, err)
			} else {
				fmt.Printf("IP для %s (%s) успешно сохранен в базу данных.\n", listContainer, ip)
			}

			// Выполняем команду ping
			cmd := exec.Command("ping", "-c", "1", ip)
			output, err := cmd.CombinedOutput()
			if err != nil {
				fmt.Printf("Ошибка при пинге %s (%s): %v\n", listContainer, ip, err)
			} else {
				fmt.Printf("Результат пинга %s (%s):\n%s\n", listContainer, ip, string(output))
			}
		}
		time.Sleep(time.Second * 10) // Пауза между пингами
	}
}
