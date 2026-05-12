package file

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/xuri/excelize/v2"
	"io"
	"os"
	"strconv"
)

// Excel 处理excel文件
type Excel struct {
}

// ImportData excel文件导入
func (e *Excel) ImportData(c *gin.Context) (headers []string, list [][]string, err error) {
	formFile, _, _ := c.Request.FormFile("file")
	defer formFile.Close()
	tempFile, _ := os.CreateTemp("/tmp/", "importTemp-*.xlsx")
	defer os.RemoveAll(tempFile.Name())
	_, err = io.Copy(tempFile, formFile)
	if err != nil {
		return nil, nil, errors.New("创建临时文件失败")
	}
	f, err := excelize.OpenFile(tempFile.Name())
	if err != nil {
		return nil, nil, errors.New("打开临时文件失败")
	}
	sheetList := f.GetSheetList()
	if len(sheetList) == 0 {
		return nil, nil, errors.New("文件中没有任何工作表")
	}
	rows, err := f.GetRows(sheetList[0])
	if err != nil {
		return nil, nil, errors.New("读取数据失败")
	}
	if len(rows) == 0 {
		err = errors.New("请确认文件中是否有Sheet1，且此sheet中是否有数据")
		return
	}
	// 遍历数据
	for line, val := range rows {
		if line == 0 {
			headers = val
			continue
		}
		list = append(list, val)
	}
	return
}

var excelLine = []string{"A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M", "N", "O", "P", "Q", "R", "S", "T", "U", "V", "W", "X", "Y", "Z", "AA", "AB", "AC", "AD", "AE", "AF", "AG", "AH", "AI", "AJ", "AK", "AL", "AM", "AN", "AO", "AP", "AQ", "AR", "AS", "AT", "AU", "AV", "AW", "AX", "AY", "AZ", "BA", "BB", "BC", "BD", "BE", "BF", "BG", "BH", "BI", "BJ", "BK", "BL", "BM", "BN", "BO", "BP", "BQ", "BR", "BS", "BT", "BU", "BV", "BW", "BX", "BY", "BZ"}

func (e *Excel) Export(outputFileName string, header []string, list [][]string) error {
	ex := excelize.NewFile()
	defer ex.Close()
	// 写表头
	for k, item := range header {
		// A1 B1 C1 D1 ...
		err := ex.SetCellValue("Sheet1", excelLine[k]+"1", item)
		if err != nil {
			return err
		}
	}
	// 写内容
	for pk, item := range list {
		// A2 B2 C2 D2 ...
		// A3 B3 C3 D3 ...
		// A4 B4 C4 D4 ...
		// ...
		for k, point := range item {
			En := excelLine[k]
			err := ex.SetCellValue("Sheet1", En+strconv.Itoa(pk+2), point)
			if err != nil {
				return err
			}
		}
	}
	if err := ex.SaveAs(outputFileName); err != nil {
		return err
	}
	return nil
}

func (e *Excel) ExportEX(outputFileName string, info map[string]string, header []string, list [][]string) error {
	ex := excelize.NewFile()
	defer ex.Close()
	var (
		err error
		pos = 1
	)
	// 写额外信息
	for k, v := range info {
		err = ex.SetCellValue("Sheet1", fmt.Sprintf("%s%d", excelLine[0], pos), k)
		if err != nil {
			return err
		}
		err = ex.SetCellValue("Sheet1", fmt.Sprintf("%s%d", excelLine[1], pos), v)
		pos++
	}
	// 写表头
	for k, item := range header {
		err = ex.SetCellValue("Sheet1", fmt.Sprintf("%s%d", excelLine[k], pos), item)
		if err != nil {
			return err
		}
	}
	// 写内容
	for pk, item := range list {
		for k, point := range item {
			err = ex.SetCellValue("Sheet1", fmt.Sprintf("%s%d", excelLine[k], pos+pk+1), point)
			if err != nil {
				return err
			}
		}
	}
	// save
	if err = ex.SaveAs(outputFileName); err != nil {
		return err
	}
	return nil
}
