package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/l3lackShark/winservice/database"
	"github.com/l3lackShark/winservice/memory"
)

const updateTime int64 = 1000 //ms

func main() {
	exPath, err := os.Executable()
	if err != nil {
		panic(err)
	}

	db, err := database.New(exPath)
	if err != nil {
		panic(err)
	}

	memoryApi := memory.New()
	var prevProcs map[memory.UniqueProcess]memory.Process

	ticker := time.NewTicker(time.Duration(updateTime) * time.Millisecond)

	for ; true; <-ticker.C {
		iterationStartTime := time.Now()
		procs, changes, err := memoryApi.GetAllProcessesAndComputeDiff(prevProcs)
		if err != nil {
			log.Panicf("memoryApi.GetAllProcessesAndComputeDiff(): %e\n", err) //needs proper handling in production
		}

		//check if there is a difference
		if len(changes.Clsoed)+len(changes.New) > 0 {
			prevProcs = procs

			outJSON, err := json.Marshal(changes)
			if err != nil {
				log.Panicf("json.Marshal(): %e\n", err) //needs proper handling in production
			}

			fmt.Printf("New DIFF:\n %s", string(outJSON))
			//store the payload in the database (in goroutine to not cause a waitline for the next iteration)
			go func() {
				if err := db.SendPayload(outJSON); err != nil {
					log.Panicf("db.SendPayload(out): %e\n", err) //needs proper handling in production
				}
			}()
		}
		elapsed := time.Since(iterationStartTime).Milliseconds()
		log.Printf("Cycle took: %dms, len(procs): %d New tick in about %dms\n", elapsed, len(procs), updateTime-elapsed)
	}
}
