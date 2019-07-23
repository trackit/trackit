//   Copyright 2019 MSolution.IO
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

import "github.com/tealeg/xlsx"

type style interface{ apply(*cell) }

func (c cell) addStyle(options ...style) cell {
	for _, option := range options {
		option.apply(&c)
	}
	return c
}

type defaultStyle struct{ style }

func (defaultStyle) apply(item *cell) {
	item.style.Border = *xlsx.NewBorder("thin", "thin", "thin", "thin")
	item.style.ApplyBorder = true
}

type textBoldStyle struct{ style }

var textBold = textBoldStyle{}

func (textBoldStyle) apply(item *cell) {
	item.style.Font.Bold = true
	item.style.ApplyFont = true
}

type textItalicStyle struct{ style }

var textItalic = textItalicStyle{}

func (textItalicStyle) apply(item *cell) {
	item.style.Font.Italic = true
	item.style.ApplyFont = true
}

type textCenterStyle struct{ style }

var textCenter = textCenterStyle{}

func (textCenterStyle) apply(item *cell) {
	item.style.Alignment.Horizontal = "center"
	item.style.ApplyAlignment = true
}

type backgroundGreenStyle struct{ style }

var backgroundGreen = backgroundGreenStyle{}

func (backgroundGreenStyle) apply(item *cell) {
	item.style.Fill.PatternType = "solid"
	item.style.Fill.FgColor = "FFB9F6CA"
	item.style.Font.Color = "FF005005"
}

type backgroundRedStyle struct{ style }

var backgroundRed = backgroundRedStyle{}

func (backgroundRedStyle) apply(item *cell) {
	item.style.Fill.PatternType = "solid"
	item.style.Fill.FgColor = "FFFF8A80"
	item.style.Font.Color = "FF8E0000"
}

type backgroundGreyStyle struct{ style }

var backgroundGrey = backgroundGreyStyle{}

func (backgroundGreyStyle) apply(item *cell) {
	item.style.Fill.PatternType = "solid"
	item.style.Fill.FgColor = "FFCCCCCC"
	item.style.Font.Color = "FF000000"
}

type backgroundLightGreyStyle struct{ style }

var backgroundLightGrey = backgroundLightGreyStyle{}

func (backgroundLightGreyStyle) apply(item *cell) {
	item.style.Fill.PatternType = "solid"
	item.style.Fill.FgColor = "FFEEEEEE"
	item.style.Font.Color = "FF000000"
}
