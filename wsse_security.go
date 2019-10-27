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

/*

 */
func NewSecurity() *Security {
	out := new(Security)
	out.Wsse = "http://docs.oasis-open.org/wss/2004/01/oasis-200401-wss-wssecurity-secext-1.0.xsd"
	out.Wsu = "http://docs.oasis-open.org/wss/2004/01/oasis-200401-wss-wssecurity-utility-1.0.xsd"
	out.UsernameToken.Password.Type = "http://docs.oasis-open.org/wss/2004/01/oasis-200401-wss-username-token-profile-1.0#PasswordText"
	return out
}

/*

 */
func NewSecurityWithCredentials(username string, password string) *Security {
	out := NewSecurity()
	out.UsernameToken.Username = username
	out.UsernameToken.Password.Text = password
	return out
}

/*

 */
func (p *Security) SetCredentials(username string, password string) *Security {
	p.UsernameToken.Username = username
	p.UsernameToken.Password.Text = password
	return p
}
