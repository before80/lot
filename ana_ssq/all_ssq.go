package ana_ssq

import (
	"fmt"
	"sort"
	"strings"

	"github.com/before80/lot/dbop"
	"github.com/before80/lot/gen"
	"github.com/before80/lot/models"
)

// SsqAllFrontBackCombToTable
//
//	@Description:
//	A 的个数有：  179721
//	B 的个数有： 4728859
//	C 的个数有：11898754
//	D 的个数有： 4463550
//	E 的个数有：  154828
func SsqAllFrontBackCombToTable() {
	combs := gen.CrossComb(gen.AllSsqFrontHms, gen.AllSsqBackHms, 6, 1)
	var allSsqs []models.AllSsq
	var hz int
	var aeHz, oe, qzh string
	for _, comb := range combs {
		newCombStr := strings.ReplaceAll(comb, "|", ",")
		newCombs := strings.Split(newCombStr, ",")
		hz = CalSsqHz(newCombs)
		aeHz = SsqHzABCDE(hz)
		oe = CalSsqOe(newCombs)
		qzh = CalSsqQzh(newCombs)
		allSsqs = append(allSsqs, models.AllSsq{
			Hm:   comb,
			Hz:   hz,
			Oe:   oe,
			AeHz: aeHz,
			Qzh:  qzh,
		})
	}

	insertRow, err := dbop.InsertAllSsqBatch(allSsqs, 10000)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("insert %d rows successfully\n", insertRow)
	//if i%100000 == 0 {
	//	fmt.Printf("已完成 %8d 还有%8d 个\n", i, 21425712-i)
	//}
	fmt.Printf("SsqAllFrontBackCombToTable done\n")
}

// GenAllSsqsInfoForGoCode 生成所有双色球的相关信息的Go代码
//
//	@Description:
func GenAllSsqsInfoForGoCode() {
	ssqs, _ := dbop.ReadAllSsqs(false)
	// 奇偶
	// 和值
	// ae和值
	// 前中后
	onlyOneFront2Count := make(map[string]int)
	onlyOneBack2Count := make(map[string]int)
	oe2Count := make(map[string]int)
	hz2Count := make(map[string]int)
	aeHz2Count := make(map[string]int)
	qzh2Count := make(map[string]int)
	for _, ssq := range ssqs {
		fbs := strings.Split(ssq.Hm, "|")
		frontHmSlice := strings.Split(fbs[0], ",")
		backHmSlice := strings.Split(fbs[1], ",")
		for _, v := range frontHmSlice {
			if _, ok := onlyOneFront2Count[v]; !ok {
				onlyOneFront2Count[v] = 1
			} else {
				onlyOneFront2Count[v] += 1
			}
		}

		for _, v := range backHmSlice {
			if _, ok := onlyOneBack2Count[v]; !ok {
				onlyOneBack2Count[v] = 1
			} else {
				onlyOneBack2Count[v] += 1
			}
		}

		oe2Count[ssq.Oe] = oe2Count[ssq.Oe] + 1
		hz2Count[fmt.Sprintf("%d", ssq.Hz)] = hz2Count[fmt.Sprintf("%d", ssq.Hz)] + 1
		aeHz2Count[ssq.AeHz] = aeHz2Count[ssq.AeHz] + 1
		qzh2Count[ssq.Qzh] = qzh2Count[ssq.Qzh] + 1
	}

	onlyOneFrontKLen := make([]KeyWithLength, 0, len(onlyOneFront2Count))
	onlyOneBackKLen := make([]KeyWithLength, 0, len(onlyOneBack2Count))
	oeKLen := make([]KeyWithLength, 0, len(oe2Count))
	hzKLen := make([]KeyWithLength, 0, len(hz2Count))
	aeHzKLen := make([]KeyWithLength, 0, len(aeHz2Count))
	qzhKLen := make([]KeyWithLength, 0, len(qzh2Count))
	for k, v := range onlyOneFront2Count {
		onlyOneFrontKLen = append(onlyOneFrontKLen, KeyWithLength{k, v})
	}

	for k, v := range onlyOneBack2Count {
		onlyOneBackKLen = append(onlyOneBackKLen, KeyWithLength{k, v})
	}
	for k, v := range oe2Count {
		oeKLen = append(oeKLen, KeyWithLength{k, v})
	}

	for k, v := range hz2Count {
		hzKLen = append(hzKLen, KeyWithLength{k, v})
	}
	for k, v := range aeHz2Count {
		aeHzKLen = append(aeHzKLen, KeyWithLength{k, v})
	}
	for k, v := range qzh2Count {
		qzhKLen = append(qzhKLen, KeyWithLength{k, v})
	}

	sort.Slice(onlyOneFrontKLen, func(i, j int) bool { return onlyOneFrontKLen[i].Length > onlyOneFrontKLen[j].Length })
	sort.Slice(onlyOneBackKLen, func(i, j int) bool { return onlyOneBackKLen[i].Length > onlyOneBackKLen[j].Length })
	sort.Slice(oeKLen, func(i, j int) bool { return oeKLen[i].Length > oeKLen[j].Length })
	sort.Slice(hzKLen, func(i, j int) bool { return hzKLen[i].Length > hzKLen[j].Length })
	sort.Slice(aeHzKLen, func(i, j int) bool { return aeHzKLen[i].Length > aeHzKLen[j].Length })
	sort.Slice(qzhKLen, func(i, j int) bool { return qzhKLen[i].Length > qzhKLen[j].Length })

	fmt.Printf("// AllSsqOnlyOneFront2Count 所有双色球前区号码中的单独一个号码对应的个数\n")
	fmt.Printf("var AllSsqOnlyOneFront2Count = map[string]int{\n")
	for _, v := range onlyOneFrontKLen {
		fmt.Printf("\t\"%s\": %d,\n", v.Key, v.Length)
	}
	fmt.Printf("} \n")

	fmt.Printf("// AllSsqOnlyOneBack2Count 所有双色球后区号码中的单独一个号码对应的个数\n")
	fmt.Printf("var AllSsqOnlyOneBack2Count = map[string]int{\n")
	for _, v := range onlyOneBackKLen {
		fmt.Printf("\t\"%s\": %d,\n", v.Key, v.Length)
	}
	fmt.Printf("} \n")

	fmt.Printf("// AllSsqOe2Count 所有双色球可能号码中的奇偶对应的个数\n")
	fmt.Printf("var AllSsqOe2Count = map[string]int{\n")
	for _, v := range oeKLen {
		fmt.Printf("\t\"%s\": %d,\n", v.Key, v.Length)
	}
	fmt.Printf("} \n")
	fmt.Printf("// AllSsqHz2Count 所有双色球可能号码中的和值对应的个数\n")
	fmt.Printf("var AllSsqHz2Count = map[string]int{\n")
	for _, v := range hzKLen {
		fmt.Printf("\t\"%s\": %d,\n", v.Key, v.Length)
	}
	fmt.Printf("} \n")
	fmt.Printf("// AllSsqAeHz2Count 所有双色球可能号码中的ae和值对应的个数\n// A -> 22~57\n// B -> 58~93\n// C -> 94~129\n// D -> 130~165\n// E -> 166~199\n")
	fmt.Printf("var AllSsqAeHz2Count = map[string]int{\n")
	for _, v := range aeHzKLen {
		fmt.Printf("\t\"%s\": %d,\n", v.Key, v.Length)
	}
	fmt.Printf("} \n")
	fmt.Printf("// AllSsqQzh2Count 所有双色球可能号码中的前中后对应的个数\n")
	fmt.Printf("var AllSsqQzh2Count = map[string]int{\n")
	for _, v := range qzhKLen {
		fmt.Printf("\t\"%s\": %d,\n", v.Key, v.Length)
	}
	fmt.Printf("} \n")

}
