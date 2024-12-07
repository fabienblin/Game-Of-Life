package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

type dynamicSelect struct {
	widget.Select
	options func() []string
}

func newDynamicSelect(options func() []string, changed func(string)) *dynamicSelect {
	ds := &dynamicSelect{options: options}
	ds.Select.OnChanged = changed

	ds.ExtendBaseWidget(ds)

	return ds
}

func (ds *dynamicSelect) Tapped(point *fyne.PointEvent) {
	ds.Options = ds.options()
	ds.Select.Tapped(point)
}
