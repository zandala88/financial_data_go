package main

import (
	"context"
	"encoding/csv"
	_ "financia/public/db/connector"
	"financia/public/db/dao"
	"financia/public/db/model"
	"financia/router"
	"financia/util"
	"github.com/spf13/cast"
	"go.uber.org/zap"
	"log"
	"os"
	"strings"
	"time"
)

func processRecord(record []string) *model.StockData {
	parts := strings.Split(record[1], " ")
	return &model.StockData{
		TsCode:    record[0],
		TradeDate: util.ConvertDateStrToTime(parts[0], time.DateOnly),
		Open:      cast.ToFloat64(record[2]),
		High:      cast.ToFloat64(record[3]),
		Low:       cast.ToFloat64(record[4]),
		Close:     cast.ToFloat64(record[5]),
		PreClose:  cast.ToFloat64(record[6]),
		Change:    cast.ToFloat64(record[7]),
		PctChg:    cast.ToFloat64(record[8]),
		Vol:       cast.ToInt64(record[9]),
		Amount:    cast.ToFloat64(record[10]),
	}
}

// 批量插入到指定的表
func insertDataToTable(batch []*model.StockData, tableName string) error {
	// 执行插入到指定表
	if err := dao.InsertStockData(context.Background(), batch); err != nil {
		log.Printf("Failed to insert data into %s: %v", tableName, err)
		return err
	}
	return nil
}

func f1(filename string) {
	// 打开文件
	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	zap.S().Debugf("[Daily] [start] 1 ")

	// 创建 CSV Reader
	reader := csv.NewReader(file)

	zap.S().Debugf("[Daily] [start] 2 ")

	// 跳过标题行
	if _, err := reader.Read(); err != nil {
		log.Fatal(err)
	}

	zap.S().Debugf("[Daily] [start] 3 ")

	// 读取剩余的所有数据
	records, err := reader.ReadAll()
	if err != nil {
		log.Fatal(err)
	}

	zap.S().Debugf("[Daily] [start] 4 ")

	data := make([]*model.StockData, 0, 1001)
	tmp := ""
	flag := false
	// 输出数据
	for _, record := range records {
		if !flag {
			tmp = record[0]
			flag = true
		} else {
			if len(data) > 1000 {
				// 插入数据
				err := insertDataToTable(data, "")
				if err != nil {
					return
				}
				time.Sleep(3 * time.Second)
				data = make([]*model.StockData, 0, 1001)
			} else if tmp != record[0] {
				// 插入数据
				err := insertDataToTable(data, "")
				if err != nil {
					return
				}
				time.Sleep(3 * time.Second)
				data = make([]*model.StockData, 0, 1001)
				tmp = record[0]
			}
		}
		stockData := processRecord(record)
		data = append(data, stockData)
	}

	err = insertDataToTable(data, "")
	if err != nil {
		return
	}
}

func main() {
	router.HTTPRouter()
}
