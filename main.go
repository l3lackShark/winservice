package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/go-test/deep"
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
	var prevProcs ([]memory.Process)

	for {
		iterationStartTime := time.Now()
		procs, err := memoryApi.GetAllProcesses()
		if err != nil {
			log.Printf("memoryApi.GetAllProcesses(): %e\n", err)
			continue
		}
		//if we need a diff in the database, this could see some logic improvements. But for the sake of simplicity, the whole new payload is uploaded everytime there is a change
		if diff := deep.Equal(prevProcs, procs); diff != nil { //there is a change
			fmt.Printf("NEW DIFF:\n %s\n", diff)
			prevProcs = procs

			out, err := json.Marshal(procs)
			if err != nil {
				log.Printf("json.Marshal(): %e\n", err)
				continue
			}

			//store the payload in the database (in goroutine to not cause a waitline for the next iteration)
			go func() {
				if err := db.SendPayload(out); err != nil {
					log.Printf("json.Marshal(): %e\n", err) //we just log an error in this case, needs proper handling in production
				}
			}()
		}
		elapsed := time.Since(iterationStartTime).Milliseconds()
		fmt.Printf("Cycle took: %dms, len(procs): %d Sleeping for ~~%dms\n", elapsed, len(procs), updateTime-elapsed)
		time.Sleep(time.Duration(updateTime-elapsed) * time.Millisecond)
	}
}
