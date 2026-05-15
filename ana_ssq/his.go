package ana_ssq

import (
	"sync"

	"github.com/before80/lot/dbop"
	"github.com/before80/lot/models"
)

type SsqHis struct {
	Typ      string // 类型
	Cs       int    // 历史总的出现期数
	AllCount int    // 所有可能
	Last10   int    // 最近10期出现期数
	Last20   int    // 最近20期出现期数
	Last30   int
	Last50   int
	Last100  int
	Last200  int
	Last500  int
	Last1000 int
	Last1500 int
	Last2000 int
	Last2500 int
	Last3500 int // 最近3500期出现期数
}

// DxSsqs 倒序
var DxSsqs []models.Ssq

// ZxSsqs 正序
var ZxSsqs []models.Ssq

// ZxAllSsqs 正序
//var ZxAllSsqs = make([]models.AllSsq, 0, 21425712)

func InitSsqs() {
	var wg sync.WaitGroup
	//startTime0 := time.Now()
	for i := 0; i < 2; i++ {
		wg.Add(1)
		if i == 0 {
			go func() {
				defer wg.Done()
				//startTime := time.Now()
				ZxSsqs, _ = dbop.ReadAllSsq(false)
				//fmt.Println(time.Now().Sub(startTime).Round(time.Second), time.Now().Sub(startTime0).Round(time.Second))
			}()
		}
		if i == 1 {
			go func() {
				defer wg.Done()
				//startTime := time.Now()
				DxSsqs, _ = dbop.ReadAllSsq(true)
				//fmt.Println(time.Now().Sub(startTime).Round(time.Second), time.Now().Sub(startTime0).Round(time.Second))
			}()
		}
		//if i == 2 {
		//	go func() {
		//		defer wg.Done()
		//		startTime := time.Now()
		//		ZxAllSsqs, _ = dbop.ReadAllSsqs(false)
		//		fmt.Println(time.Now().Sub(startTime).Round(time.Second), time.Now().Sub(startTime0).Round(time.Second))
		//	}()
		//}
	}
	wg.Wait()
}
