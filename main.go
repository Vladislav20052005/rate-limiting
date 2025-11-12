package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
)

var ctx = context.Background()

func checkTime(rdb *redis.Client, ip string, delta int64) (bool, error) {
	currentTime := time.Now().UnixMilli()

	time, err := rdb.Get(ctx, ip).Result()
	if err != nil && err != redis.Nil {
		return false, err
	}
	var lastTime int64
	if time == "" {
		lastTime = 0
	} else {
		lastTime, err = strconv.ParseInt(time, 10, 64)
		if err != nil {
			return false, err
		}
	}

	if currentTime < lastTime+delta {
		return false, nil
	}

	err = rdb.Set(ctx, ip, currentTime, 0).Err()
	if err != nil {
		return false, err
	}
	return true, nil
}

func main() {
	if len(os.Args[1:]) < 3 {
		fmt.Printf("Not enough args")
		return
	}
	maxOverallRate, err := strconv.ParseInt(os.Args[1], 10, 64)
	if err != nil {
		fmt.Printf("Arg error: %s", err.Error())
		return
	}
	maxIPRate, err := strconv.ParseInt(os.Args[2], 10, 64)
	if err != nil {
		fmt.Printf("Arg error: %s", err.Error())
		return
	}
	address := os.Args[3]

	var (
		overallDelta int64 = 60 * 1000 / maxOverallRate
		ipDelta      int64 = 60 * 1000 / maxIPRate //delay in ms
	)

	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		ipPort := r.RemoteAddr
		ip := strings.Split(ipPort, ":")[0]

		isAllowedOverall, err := checkTime(rdb, "all", overallDelta)
		if err != nil {
			fmt.Printf("error: %s", err.Error())
			return
		}

		isAllowedIP, err := checkTime(rdb, ip, ipDelta)
		if err != nil {
			fmt.Printf("error: %s", err.Error())
			return
		}

		if isAllowedOverall && isAllowedIP {
			http.Redirect(w, r, address, http.StatusSeeOther) // Success
		} else {
			http.ServeFile(w, r, "dummy/index.html")
		}

	})

	fmt.Println("Server starting on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
