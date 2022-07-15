package main

import (
	"fmt"
	"github.com/roelofruis/upkeep/internal/infra"
	"testing"
	"time"
)

func Test(t *testing.T) {

	clock := infra.FixedClock{Time: time.Date(2022, 03, 19, 20, 21, 5, 0, time.UTC)}
	io := infra.NewInMemoryIO()

	router := Bootstrap(clock, io)

	res, err := router.Handle(infra.ParseArgs([]string{"upkeep", "start"}))
	fmt.Printf("RES: %s\nERR: %v\n", res, err)

}
