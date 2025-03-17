package utils

import (
	"fmt"
	"github.com/samber/lo"
	"github.com/xuri/excelize/v2"
)

const xlsxSheetName = "Sheet1"

type Page struct {
	Name  string
	Table [][]string
}

type XLSXView interface {
	ToXLSX() XLSXView
}

func ExportXLSX(table [][]string) (*excelize.File, error) {
	file := excelize.NewFile()
	defer func() {
		if err := file.Close(); err != nil {
			fmt.Println(err)
		}
	}()

	sheet, err := addSheet(file, xlsxSheetName, table)
	if err != nil {
		return nil, err
	}

	file.SetActiveSheet(sheet)

	return file, nil
}

func ExportMultiPageXLSX(pages []Page) (*excelize.File, error) {
	pages = lo.Map(pages, func(item Page, index int) Page {
		item.Name = truncateSheetName(item.Name)

		return item
	})

	file := excelize.NewFile()
	defer func() {
		if err := file.Close(); err != nil {
			fmt.Println(err)
		}
	}()

	for index, page := range pages {
		sheet, err := addSheet(file, page.Name, page.Table)
		if err != nil {
			return nil, err
		}

		if index == 0 {
			file.SetActiveSheet(sheet)
		}
	}

	// delete default sheet if not used
	sheetNames := lo.Map(pages, func(item Page, index int) string {
		return item.Name
	})

	if lo.Contains(sheetNames, xlsxSheetName) {
		return file, nil
	}

	if err := file.DeleteSheet(xlsxSheetName); err != nil {
		return nil, err
	}

	return file, nil
}

func addSheet(file *excelize.File, sheetName string, table [][]string) (int, error) {
	sheet, err := file.NewSheet(sheetName)
	if err != nil {
		return sheet, err
	}

	style, err := file.NewStyle(&excelize.Style{})
	if err != nil {
		return sheet, err
	}

	for i, row := range table {
		for j, cell := range row {
			if i == 0 {
				if err = file.SetCellStyle(sheetName, getCellName(i, j), getCellName(i, j), style); err != nil {
					return sheet, err
				}
			} else {
				if v, _ := file.GetCellValue(sheetName, getCellName(i-1, j)); v == "" {
					if err = file.SetCellStyle(sheetName, getCellName(i, j), getCellName(i, j), style); err != nil {
						return sheet, err
					}
				}
			}

			if err = file.SetCellValue(sheetName, getCellName(i, j), cell); err != nil {
				return sheet, err
			}
		}
	}

	return sheet, nil
}

func getCellName(i, j int) string {
	return fmt.Sprintf("%s%d", ColumnNames[j], i+1)
}

func truncateSheetName(sheetName string) string {
	if len(sheetName) < excelize.MaxSheetNameLength {
		return sheetName
	}

	return sheetName[:excelize.MaxSheetNameLength]
}
