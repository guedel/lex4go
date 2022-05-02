package generator

import (
	"encoding/xml"
)

type Compare struct {
	Compare CompareInterface
}

func (r *Compare) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	for {
		t, err := d.Token()
		if err != nil {
			return err
		}
		var i CompareInterface
		switch tt := t.(type) {
		case xml.StartElement:
			// fmt.Println("Test.UnmarshalXML.StartElement: ", tt)
			switch tt.Name.Local {
			case "charset":
				i = new(CompareCharset)
			case "like":
				i = new(CompareLike)
			case "eos":
				i = new(CompareEos)
			}
			if i != nil {
				err = d.DecodeElement(i, &tt)
				if err != nil {
					return err
				}
				r.Compare = i
			}
		case xml.EndElement:
			if tt == start.End() {
				return nil
			}
		}

	}
}
