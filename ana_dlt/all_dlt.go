package ana_dlt

import (
	"fmt"
	"sort"
	"strings"

	"github.com/before80/lot/dbop"
	"github.com/before80/lot/gen"
	"github.com/before80/lot/models"
)

// DltAllFrontBackCombToTable
//
//	@Description:
//	A 的个数有：  179721
//	B 的个数有： 4728859
//	C 的个数有：11898754
//	D 的个数有： 4463550
//	E 的个数有：  154828
func DltAllFrontBackCombToTable() {
	combs := gen.CrossComb(gen.AllDltFrontHms, gen.AllDltBackHms, 5, 2)
	var allDlts []models.AllDlt
	var hz int
	var aeHz, oe, qzh string
	for _, comb := range combs {
		newCombStr := strings.ReplaceAll(comb, "|", ",")
		newCombs := strings.Split(newCombStr, ",")
		hz = CalDltHz(newCombs)
		aeHz = DltHzABCDE(hz)
		oe = CalDltOe(newCombs)
		qzh = CalDltQzh(newCombs)
		allDlts = append(allDlts, models.AllDlt{
			Hm:   comb,
			Hz:   hz,
			Oe:   oe,
			AeHz: aeHz,
			Qzh:  qzh,
		})
	}

	insertRow, err := dbop.InsertAllDltBatch(allDlts, 10000)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("insert %d rows successfully\n", insertRow)
	//if i%100000 == 0 {
	//	fmt.Printf("已完成 %8d 还有%8d 个\n", i, 21425712-i)
	//}
	fmt.Printf("DltAllFrontBackCombToTable done\n")
}

// GenAllDltsInfoForGoCode 生成所有大乐透的相关信息的Go代码
//
//	@Description:
func GenAllDltsInfoForGoCode() {
	dlts, _ := dbop.ReadAllDlts(false)
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
	for _, dlt := range dlts {
		fbs := strings.Split(dlt.Hm, "|")
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

		oe2Count[dlt.Oe] = oe2Count[dlt.Oe] + 1
		hz2Count[fmt.Sprintf("%d", dlt.Hz)] = hz2Count[fmt.Sprintf("%d", dlt.Hz)] + 1
		aeHz2Count[dlt.AeHz] = aeHz2Count[dlt.AeHz] + 1
		qzh2Count[dlt.Qzh] = qzh2Count[dlt.Qzh] + 1
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

	fmt.Printf("// AllDltOnlyOneFront2Count 所有大乐透前区号码中的单独一个号码对应的个数\n")
	fmt.Printf("var AllDltOnlyOneFront2Count = map[string]int{\n")
	for _, v := range onlyOneFrontKLen {
		fmt.Printf("\t\"%s\": %d,\n", v.Key, v.Length)
	}
	fmt.Printf("} \n")

	fmt.Printf("// AllDltOnlyOneBack2Count 所有大乐透后区号码中的单独一个号码对应的个数\n")
	fmt.Printf("var AllDltOnlyOneBack2Count = map[string]int{\n")
	for _, v := range onlyOneBackKLen {
		fmt.Printf("\t\"%s\": %d,\n", v.Key, v.Length)
	}
	fmt.Printf("} \n")

	fmt.Printf("// AllDltOe2Count 所有大乐透可能号码中的奇偶对应的个数\n")
	fmt.Printf("var AllDltOe2Count = map[string]int{\n")
	for _, v := range oeKLen {
		fmt.Printf("\t\"%s\": %d,\n", v.Key, v.Length)
	}
	fmt.Printf("} \n")
	fmt.Printf("// AllDltHz2Count 所有大乐透可能号码中的和值对应的个数\n")
	fmt.Printf("var AllDltHz2Count = map[string]int{\n")
	for _, v := range hzKLen {
		fmt.Printf("\t\"%s\": %d,\n", v.Key, v.Length)
	}
	fmt.Printf("} \n")
	fmt.Printf("// AllDltAeHz2Count 所有大乐透可能号码中的ae和值对应的个数\n// A -> 18~52\n// B -> 53~86\n// C -> 87~120\n// D -> 121~154\n// E -> 155~188\n")
	fmt.Printf("var AllDltAeHz2Count = map[string]int{\n")
	for _, v := range aeHzKLen {
		fmt.Printf("\t\"%s\": %d,\n", v.Key, v.Length)
	}
	fmt.Printf("} \n")
	fmt.Printf("// AllDltQzh2Count 所有大乐透可能号码中的前中后对应的个数\n")
	fmt.Printf("var AllDltQzh2Count = map[string]int{\n")
	for _, v := range qzhKLen {
		fmt.Printf("\t\"%s\": %d,\n", v.Key, v.Length)
	}
	fmt.Printf("} \n")

}
