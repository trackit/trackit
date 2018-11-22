package reports

type style interface { apply(*cell) }

func (c cell) addStyle(options ...style) (cell){
	for _, option := range options {
		option.apply(&c)
	}
	return c
}

type textBoldStyle struct { style }
var textBold = textBoldStyle{}

func (textBoldStyle) apply(item *cell) {
	item.style.Font.Bold = true
}

type textItalicStyle struct { style }
var textItalic = textItalicStyle{}

func (textItalicStyle) apply(item *cell) {
	item.style.Font.Italic = true
}

type textCenterStyle struct { style }
var textCenter = textCenterStyle{}

func (textCenterStyle) apply(item *cell) {
	item.style.Alignment.Horizontal = "center"
}

type backgroundGreenStyle struct { style }
var backgroundGreen = backgroundGreenStyle{}

func (backgroundGreenStyle) apply(item *cell) {
	item.style.Fill.BgColor = "B9F6CA00"
}

type backgroundRedStyle struct { style }
var backgroundRed = backgroundRedStyle{}

func (backgroundRedStyle) apply(item *cell) {
	item.style.Fill.BgColor = "FF8A80FF"
}
