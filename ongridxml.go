package main

import (
	"fmt"
	_ "github.com/nakagami/firebirdsql"
    "github.com/jmoiron/sqlx"
    "database/sql"
    "encoding/xml"
    "os"
    "log"
    "time"
    "strconv"
)

type InwPart struct {
	Id 			int 			`db:"ID"`
	DocNum 		sql.NullString	`db:"DOCNUM"`
	Date 		sql.NullString	`db:"DOCDATA"`
	TraiderId 	int 			`db:"TRAIDER"`
	TraiderName string 			`db:"ORGMNE"`
	Wrhouse 	int 			`db:"WRHOUSE"`
	Parts 		[]Part
}

type Part struct {
	Id 			int 	`db:"ID"`
	CatalogN 	int 	`db:"CATALOGN"`
	PartName	string 	`db:"NAME"`
	Qty 		float64 `db:"QUANTITY"`
	Price 		float64 `db:"PRICE"`
}

func main() {
	

	db, err := sqlx.Connect("firebirdsql", "sysdba:masterkey@localhost:3050/c:/database/infiniti/inter_akos.gdb")
	if err != nil {
		log.Fatalln(err)
        return
	}
	defer db.Close()
    
    inwparts := []InwPart{}

    rows, err := db.Queryx("select first 2 ip.id, ip.docnum, ip.docdata, ip.traider, ip.wrhouse, o.orgmne from inwpart ip join orgbase o on (o.id = ip.traider) where ip.isfolder <> -1")
    defer rows.Close()
    for rows.Next() {
    		
    	inwpart := InwPart{}
    	err := rows.StructScan(&inwpart)
    	if err != nil {
    		log.Println(err)
    		continue
    	}

   		parts, err := db.Queryx("select it.id, it.catalogn, it.quantity, it.price, w.catalogn as name from inwpitem it join wrhpart w on (w.id = it.catalogn) where it.parent = ? and it.deleted <> ?", inwpart.Id, "D")
   		defer parts.Close()
    	if err != nil {
    		log.Println(err)
    	}
    	
    	for parts.Next() {
    		inwpitem := Part{}
    		if err = parts.StructScan(&inwpitem); err != nil {
    			log.Println(err)
    		}
    		inwpart.Parts = append(inwpart.Parts, inwpitem)
    	}

    	inwparts = append(inwparts, inwpart)

    	//fmt.Printf("%#v\n", inwpart)
    }

    v := &XMLUploadFile{UploadDate: time.Now().String(), UploadBegin: "2017-05-01T10:00:00", UploadEnd: "2017-05-31T10:00:00"}

	d := &XMLDoc{DocName: "Поступление товаров"}

    for _, inwpart := range inwparts {
    	e := &XMLElement{}

    	p1 := &XMLProp{Name: "Номер", PropType: "Строка", Value: inwpart.DocNum.String}
		p2 := &XMLProp{Name: "Дата", PropType: "Дата", Value: inwpart.Date.String}
		p3 := &XMLProp{Name: "Контрагент", PropType: "Строка", LinkPropName: "Контрагенты", Code: strconv.Itoa(inwpart.TraiderId), Value: inwpart.TraiderName}

		d.Elements = append(d.Elements, e)

		e.Props = append(e.Props, p1)
		e.Props = append(e.Props, p2)
		e.Props = append(e.Props, p3)

		tp := &XMLTabPart{Name: "Товары"}

		for _, part := range inwpart.Parts {
			tprow := &XMLTabPartRow{}
			
			tpprop1 := &XMLTabPartProp{Name: "Номенклатура", PropType: "Ссылка", LinkPropName: "Номенклатура", Code: strconv.Itoa(part.CatalogN), Value: part.PartName}
			tpprop2 := &XMLTabPartProp{Name: "Количество", PropType: "Число", Value: strconv.FormatFloat(part.Qty, 'f', 2, 32)}
			tpprop3 := &XMLTabPartProp{Name: "Цена", PropType: "Число", Value: strconv.FormatFloat(part.Price, 'f', 2, 32)}
			tprow.TabPartProps = append(tprow.TabPartProps, tpprop1)
			tprow.TabPartProps = append(tprow.TabPartProps, tpprop2)
			tprow.TabPartProps = append(tprow.TabPartProps, tpprop3)
			tp.Rows = append(tp.Rows, tprow)

		}
    	e.TabParts = append(e.TabParts, tp)

    	//fmt.Printf("%#v\n", inwpart)

    }

	v.Docs = append(v.Docs, d)



	enc := xml.NewEncoder(os.Stdout)

	enc.Indent("  ", "    ")
	if err := enc.Encode(v); err != nil {
		fmt.Printf("error: %v\n", err)
	}
	
}