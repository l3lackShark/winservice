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
	var prevProcs (map[memory.UniqueProcess]memory.Process)

	for {
		iterationStartTime := time.Now()
		procs, changes, err := memoryApi.GetAllProcessesAndComputeDiff(prevProcs)
		if err != nil {
			log.Printf("memoryApi.GetAllProcessesAndComputeDiff(): %e\n", err)
			continue
		}
		prevProcs = procs

		//check if there is a difference
		if len(changes.Clsoed) > 0 || len(changes.New) > 0 {

			outJSON, err := json.Marshal(changes)
			if err != nil {
				log.Printf("json.Marshal(): %e\n", err)
				continue
			}

			fmt.Printf("New DIFF:\n %s", string(outJSON))
			//store the payload in the database (in goroutine to not cause a waitline for the next iteration)
			go func() {
				if err := db.SendPayload(outJSON); err != nil {
					log.Printf("db.SendPayload(out): %e\n", err) //we just log an error in this case, needs proper handling in production
				}
			}()
		}
		elapsed := time.Since(iterationStartTime).Milliseconds()
		fmt.Printf("Cycle took: %dms, len(procs): %d Sleeping for ~~%dms\n", elapsed, len(procs), updateTime-elapsed)
		time.Sleep(time.Duration(updateTime-elapsed) * time.Millisecond)
	}
}
