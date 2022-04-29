package main

type Rule struct {
	Id      string           `xml:"id,attr"`
	From    string           `xml:"from,attr"`
	To      string           `xml:"to,attr"`
	Compare CompareInterface `xml:"test"`
	Repeat  int              `xml:"repeat,attr"`
	Final   bool             `xml:"final"`
	Concat  bool             `xml:"concat"`
	Reset   bool             `xml:"reset"`
	Action  string           `xml:"action"`
}

/*
func (r *Rule) UnmarshalText(text []byte) error {
	fmt.Println("Rule.UnmarshalText: ", string(text))
	return nil
}
*/
