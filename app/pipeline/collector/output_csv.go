package collector

import (
	"encoding/csv"
	"fmt"
	"os"

	"github.com/ktktcom/pholcus/common/util"
	"github.com/ktktcom/pholcus/config"
	"github.com/ktktcom/pholcus/logs"
)

/************************ CSV 输出 ***************************/
func init() {
	Output["csv"] = func(self *Collector, dataIndex int) {
		defer func() {
			if err := recover(); err != nil {
				logs.Log.Error("%v", err)
			}
		}()
		var namespace = util.FileNameReplace(self.namespace())
		var sheets = make(map[string]*csv.Writer)
		for _, datacell := range self.DockerQueue.Dockers[dataIndex] {
			var subNamespace = util.FileNameReplace(self.subNamespace(datacell))
			if _, ok := sheets[subNamespace]; !ok {
				folder := config.COMM_PATH.TEXT + "/" + self.startTime.Format("2006年01月02日 15时04分05秒") + "/" + namespace + "__" + subNamespace
				filename := fmt.Sprintf("%v/%v-%v.csv", folder, self.sum[0], self.sum[1])

				// 创建/打开目录
				f, err := os.Stat(folder)
				if err != nil || !f.IsDir() {
					if err := os.MkdirAll(folder, 0777); err != nil {
						logs.Log.Error("Error: %v\n", err)
					}
				}

				// 按数据分类创建文件
				file, err := os.Create(filename)

				if err != nil {
					logs.Log.Error("%v", err)
					continue
				}

				file.WriteString("\xEF\xBB\xBF") // 写入UTF-8 BOM

				sheets[subNamespace] = csv.NewWriter(file)
				th := self.GetRule(datacell["RuleName"].(string)).GetOutFeild()
				th = append(th, "当前链接", "上级链接", "下载时间")
				sheets[subNamespace].Write(th)

				defer func(file *os.File) {
					// 发送缓存数据流
					sheets[subNamespace].Flush()
					// 关闭文件
					file.Close()
				}(file)
			}

			row := []string{}
			for _, title := range self.GetRule(datacell["RuleName"].(string)).GetOutFeild() {
				vd := datacell["Data"].(map[string]interface{})
				if v, ok := vd[title].(string); ok || vd[title] == nil {
					row = append(row, v)
				} else {
					row = append(row, util.JsonString(vd[title]))
				}
			}

			row = append(row, datacell["Url"].(string))
			row = append(row, datacell["ParentUrl"].(string))
			row = append(row, datacell["DownloadTime"].(string))
			sheets[subNamespace].Write(row)
		}
	}
}
