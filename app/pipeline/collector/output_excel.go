package collector

import (
	"fmt"
	"os"

	"github.com/ktktcom/pholcus/common/util"
	"github.com/ktktcom/pholcus/config"
	"github.com/ktktcom/pholcus/logs"
	"github.com/tealeg/xlsx"
)

/************************ excel 输出 ***************************/
func init() {
	Output["excel"] = func(self *Collector, dataIndex int) {
		defer func() {
			if err := recover(); err != nil {
				logs.Log.Error("%v", err)
			}
		}()

		var file *xlsx.File
		var sheets = make(map[string]*xlsx.Sheet)
		var row *xlsx.Row
		var cell *xlsx.Cell
		var err error

		// 创建文件
		file = xlsx.NewFile()

		// 添加分类数据工作表
		for _, datacell := range self.DockerQueue.Dockers[dataIndex] {
			var subNamespace = util.FileNameReplace(self.subNamespace(datacell))
			if _, ok := sheets[subNamespace]; !ok {
				// 添加工作表
				sheet, err := file.AddSheet(subNamespace)
				if err != nil {
					logs.Log.Error("%v", err)
					continue
				}
				sheets[subNamespace] = sheet
				// 写入表头
				row = sheets[subNamespace].AddRow()
				for _, title := range self.GetRule(datacell["RuleName"].(string)).GetOutFeild() {
					cell = row.AddCell()
					cell.Value = title
				}
				cell = row.AddCell()
				cell.Value = "当前链接"
				cell = row.AddCell()
				cell.Value = "上级链接"
				cell = row.AddCell()
				cell.Value = "下载时间"
			}

			row = sheets[subNamespace].AddRow()
			for _, title := range self.GetRule(datacell["RuleName"].(string)).GetOutFeild() {
				cell = row.AddCell()
				vd := datacell["Data"].(map[string]interface{})
				if v, ok := vd[title].(string); ok || vd[title] == nil {
					cell.Value = v
				} else {
					cell.Value = util.JsonString(vd[title])
				}
			}
			cell = row.AddCell()
			cell.Value = datacell["Url"].(string)
			cell = row.AddCell()
			cell.Value = datacell["ParentUrl"].(string)
			cell = row.AddCell()
			cell.Value = datacell["DownloadTime"].(string)
		}

		folder := config.COMM_PATH.TEXT + "/" + self.startTime.Format("2006年01月02日 15时04分05秒")
		filename := fmt.Sprintf("%v/%v__%v-%v.xlsx", folder, util.FileNameReplace(self.namespace()), self.sum[0], self.sum[1])

		// 创建/打开目录
		f2, err := os.Stat(folder)
		if err != nil || !f2.IsDir() {
			if err := os.MkdirAll(folder, 0777); err != nil {
				logs.Log.Error("Error: %v\n", err)
			}
		}

		// 保存文件
		err = file.Save(filename)
		if err != nil {
			logs.Log.Error("%v", err)
		}
	}
}
