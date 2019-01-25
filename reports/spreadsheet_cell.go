//   Copyright 2018 MSolution.IO
//
//   Licensed under the Apache License, Version 2.0 (the "License");
//   you may not use this file except in compliance with the License.
//   You may obtain a copy of the License at
//
//       http://www.apache.org/licenses/LICENSE-2.0
//
//   Unless required by applicable law or agreed to in writing, software
//   distributed under the License is distributed on an "AS IS" BASIS,
//   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//   See the License for the specific language governing permissions and
//   limitations under the License.

package reports

import (
	"fmt"
	"strings"

	"github.com/360EntSecGroup-Skylar/excelize"
)

type cell struct {
	value              interface{}
	formula            string
	location           string
	merge              string
	styles             []string
	conditionalFormats conditionalFormats
}

type cells []cell

type columnWidth struct {
	from  string
	to    string
	width float64
}

type columnsWidth []columnWidth

func newCell(value interface {}, location string) cell {
	return cell{value, "", location, "", []string{}, conditionalFormats{}}
}

func newFormula(formula string, location string) cell {
	return cell{"", formula, location, "", []string{}, conditionalFormats{}}
}

func (c cell)mergeTo(merge string) cell {
	c.merge = merge
	return c
}

func (c cell)addStyles(styles ...string) cell {
	for _, style := range styles {
		c.styles = append(c.styles, style)
	}
	return c
}

func (cs cells)addStyles(styles ...string) cells {
	for index := range cs {
		cs[index] = cs[index].addStyles(styles...)
	}
	return cs
}

func (c cell)addConditionalFormat(name string, styles ...string) cell {
	c.conditionalFormats = append(c.conditionalFormats, conditionalFormat{name, false, styles})
	return c
}

func (c cell)addCustomConditionalFormat(condition string, styles ...string) cell {
	c.conditionalFormats = append(c.conditionalFormats, conditionalFormat{condition, true, styles})
	return c
}

func (cs cells)addConditionalFormat(name string, styles ...string) cells {
	for index := range cs {
		cs[index] = cs[index].addConditionalFormat(name, styles...)
	}
	return cs
}

func (cs cells)addCustomConditionalFormat(condition string, styles ...string) cells {
	for index := range cs {
		cs[index] = cs[index].addCustomConditionalFormat(condition, styles...)
	}
	return cs
}

func (cs cells)setValues(file *excelize.File, sheet string) {
	for _, cell := range cs {
		cell.setValue(file, sheet)
	}
}

/* TODO: Update error handing (Errors should not interrupt spreadsheet generation since it is only styling issue) */
func (c cell)setValue(file *excelize.File, sheet string) {
	if len(c.formula) > 0 {
		file.SetCellFormula(sheet, c.location, c.formula)
	} else {
		file.SetCellValue(sheet, c.location, c.value)
	}
	endCellLocation := c.location
	if len(c.merge) > 0 {
		file.MergeCell(sheet, c.location, c.merge)
		endCellLocation = c.merge
	}
	if len(c.styles) > 0 {
		styleId, err := getStyleId(file, c.styles)
		if err == nil {
			file.SetCellStyle(sheet, c.location, endCellLocation, styleId)
		} else {
			fmt.Println(err, c.styles)
		}
	}
	if len(c.conditionalFormats) > 0 {
		formattedConditions, err := c.conditionalFormats.getConditionalFormatting(file)
		if err == nil {
			err = file.SetConditionalFormat(sheet, strings.Join([]string{c.location, endCellLocation}, ":"), formattedConditions)
			if err != nil {
				fmt.Println(err, formattedConditions)
			}
		} else {
			fmt.Println(err, c.conditionalFormats)
		}
	}
}

func newColumnWidth(column string, width float64) columnWidth {
	return columnWidth{column, column, width}
}

func (c columnWidth)toColumn(column string) columnWidth {
	c.to = column
	return c
}

func (c columnWidth)setValue(file *excelize.File, sheet string) {
	file.SetColWidth(sheet, c.from, c.to, c.width)
}

func (cs columnsWidth)setValues(file *excelize.File, sheet string) {
	for _, col := range cs {
		col.setValue(file, sheet)
	}
}
