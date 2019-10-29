package ItsOrgaApi

type Security struct {
	Wsse          string `xml:"wsse,attr"`
	Wsu           string `xml:"wsu,attr"`
	UsernameToken struct {
		Username string `xml:"Username"`
		Password struct {
			Text string `xml:",chardata"`
			Type string `xml:"Type,attr"`
		} `xml:"Password"`
	} `xml:"UsernameToken"`
}


type WSSESecurity struct {
	username string
	passwort string
}

func NewWSSESecurity (username, password string) *WSSESecurity{
	p := new(WSSESecurity)
	p.username = username
	p.passwort = password
	return p
}
