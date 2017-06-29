package main

import "encoding/xml"

type XMLUploadFile struct {
	XMLName xml.Name		`xml:"ФайлВыгрузки"`
	UploadDate string		`xml:"ДатаВыгрузки,attr"`
	UploadBegin string		`xml:"НачалоПериодаВыгрузки,attr"`
	UploadEnd string		`xml:"ОкончаниеПериодаВыгрузки,attr"`
	Docs []*XMLDoc			`xml:"Документ"` 
}

type XMLDoc struct {
	XMLName xml.Name		`xml:"Документ"`
	DocName string			`xml:"НаименованиеДокумента,attr"`
	Elements []*XMLElement 	`xml:"Элемент"`
}

type XMLElement struct {
	XMLName xml.Name 		`xml:"Элемент"`
	Props []*XMLProp 		`xml:"Реквизит"`
	TabParts []*XMLTabPart 	`xml:"ТабличнаяЧасть"` 
}

type XMLProp struct {
	XMLName xml.Name 		`xml:"Реквизит"`
	Name string				`xml:"Имя,attr"`
	PropType string			`xml:"Тип,attr"`
	LinkPropName string		`xml:"НаименованиеТипаРеквизитаПоСсылке,omitempty"`
	Code string				`xml:"Код,omitempty"`
	Value string			`xml:"Значение"`
}

type XMLTabPart struct {
	XMLName xml.Name 		`xml:"ТабличнаяЧасть"`
	Name string				`xml:"Имя,attr"`
	Rows []*XMLTabPartRow	`xml:"СтрокаТабличнойЧасти"`
}

type XMLTabPartRow struct {
	XMLName xml.Name 		`xml:"СтрокаТабличнойЧасти"`
	TabPartProps []*XMLTabPartProp `xml:"РеквизитТабличнойЧасти"`
}

type XMLTabPartProp struct {
	XMLName xml.Name 		`xml:"РеквизитТабличнойЧасти"`
	Name string				`xml:"Имя,attr"`
	PropType string			`xml:"Тип,attr"`
	LinkPropName string		`xml:"НаименованиеТипаРеквизитаПоСсылке,omitempty"`
	Code string				`xml:"Код,omitempty"`
	Value string			`xml:"Значение"`
}
