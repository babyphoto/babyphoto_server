package main

import (
	"database/sql"
	"time"

	"github.com/babyphoto/babyphoto_server/service/apiserver"
	"github.com/babyphoto/babyphoto_server/service/database/babyphoto"
	"github.com/babyphoto/babyphoto_server/service/util"
	_ "github.com/go-sql-driver/mysql"
	"golang.org/x/sys/windows/svc"
)

type dummyService struct {
}

func (srv *dummyService) Execute(args []string, req <-chan svc.ChangeRequest, stat chan<- svc.Status) (svcSpecificEC bool, exitCode uint32) {
	stat <- svc.Status{State: svc.StartPending}

	// 실제 서비스 내용
	stopChan := make(chan bool, 1)
	go runServer(stopChan)

	stat <- svc.Status{State: svc.Running, Accepts: svc.AcceptStop | svc.AcceptShutdown}

LOOP:
	for {
		// 서비스 변경 요청에 대해 핸들링
		switch r := <-req; r.Cmd {
		case svc.Stop, svc.Shutdown:
			stopChan <- true
			break LOOP

		case svc.Interrogate:
			stat <- r.CurrentStatus
			time.Sleep(100 * time.Millisecond)
			stat <- r.CurrentStatus
		}
	}

	stat <- svc.Status{State: svc.StopPending}
	return
}

func runServer(stopChan chan bool) {
	for {
		select {
		case <-stopChan:
			return
		default:
			db, err := sql.Open("mysql", "sherwher:lskun@tcp(174.129.71.0:3306)/babyphoto")
			util.CheckError("main ::: db connectoin ::: ", err)
			defer db.Close()

			bybyPhotodb := &babyphoto.BabyPhotoDB{
				DB: db,
			}

			API := apiserver.NewAPIServer(bybyPhotodb)
			go API.Run(":38080")
			select {}
		}
	}
}

func main() {
	err := svc.Run("DummyService", &dummyService{})
	//err := debug.Run("DummyService", &dummyService{}) //콘솔출력 디버깅시
	if err != nil {
		panic(err)
	}
}
