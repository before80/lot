package main

import (
	"fmt"
	"sync"

	"github.com/before80/lot/ana_dlt"
	"github.com/before80/lot/lg"
)

func main() {
	//threadNum := arg.GetThreadNum()
	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		if i == 0 {
			go ana_dlt.NewAbcdeQsWithWg("77777", &wg)
		}

		if i == 1 {
			go ana_dlt.NewAbcdeQsWithWg("215432", &wg)
		}
		if i == 2 {
			go ana_dlt.NewAbcdeQsWithWg("224432", &wg)
		}

		if i == 3 {
			go ana_dlt.NewAbcdeQsWithWg("224441", &wg)
		}
		if i == 4 {
			go ana_dlt.NewAbcdeQsWithWg("253322", &wg)
		}

		if i == 5 {
			go ana_dlt.NewAbcdeQsWithWg("272222", &wg)
		}

		if i == 6 {
			go ana_dlt.NewAbcdeQsWithWg("116666", &wg)
		}

		if i == 7 {
			go ana_dlt.NewAbcdeQsWithWg("155555", &wg)
		}

		if i == 8 {
			go ana_dlt.NewAbcdeQsWithWg("194444", &wg)
		}

		if i == 9 {
			go ana_dlt.NewAbcdeQsWithWg("233333", &wg)
		}
	}
	lg.InfoToFileAndStdOut(fmt.Sprintln("已经开始运行..."))
	wg.Wait()
	lg.InfoToFileAndStdOut(fmt.Sprintln("已完成"))
}
