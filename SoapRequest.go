package ItsOrgaApi

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

type SoapRequest struct {
	XMLName xml.Name `xml:"Envelope"`
	X       string   `xml:"x,attr"`
	Tem     string   `xml:"tem,attr"`
	Its     string   `xml:"its,attr"`
	Header  Header   `xml:"Header"`
	Body    Body     `xml:"Body"`
}
type Message struct {
	AddMessage struct {
		DataMessage struct {
			Data struct {
				Message createLinkMessage `xml:"Message"`
			} `xml:"Data"`
			Type string `xml:"Type"`
		} `xml:"dataMessage"`
	}
}
type Header struct {
	Security Security `xml:"Security"`
}
type Body struct {
	Message Message `xml:"AddMessage"`
}

/*
NewSoapRequest creates a new SoapRequest
*/
func NewSoapRequest() *SoapRequest {
	out := new(SoapRequest)
	return out
}

/*
Add Security adds a WSSE Security to the code
*/
func (p *SoapRequest) AddSecurity(security Security) *SoapRequest {
	p.Header.Security = security
	return p
}

/*
It puts the Data for a New Instance Link into the SoapRequest
*/
func (p *SoapRequest) CreateNewLink(link Link) (err error, r *SoapRequest) {
	mes := new(createLinkMessage)
	mes.Xmlns = "urn:message-schema" //Todo
	mes.SyncKeys.SyncKey = link.Id
	mes.VendorId = link.basicData.vendorID
	if link.basicData.location.Library || (!link.basicData.location.Library && !link.basicData.location.Course) {
		mes.CreateExtensionInstance.Location = "Library"
	}
	if link.basicData.location.Course {
		mes.CreateExtensionInstance.Location = "Course"
	}
	mes.CreateExtensionInstance.ExtensionId = "5000"
	mes.CreateExtensionInstance.UserSyncKey = link.basicData.userSyncKey
	mes.CreateExtensionInstance.Title = link.Title
	mes.CreateExtensionInstance.Metadata.Description = link.Description
	mes.CreateExtensionInstance.Metadata.Language = link.Language
	if link.Format.Video {
		mes.CreateExtensionInstance.Metadata.Format = "Video"
	} else if link.Format.Interactive {
		mes.CreateExtensionInstance.Metadata.Format = "Interactive"
	} else if link.Format.Image {
		mes.CreateExtensionInstance.Metadata.Format = "Image"
	} else if link.Format.Audio {
		mes.CreateExtensionInstance.Metadata.Format = "Audio"
	} else if link.Format.Any {
		mes.CreateExtensionInstance.Metadata.Format = "Any"
	} else if link.Format.Text {
		mes.CreateExtensionInstance.Metadata.Format = "Text"
	} else {
		mes.CreateExtensionInstance.Metadata.Format = "Any"
	}
	mes.CreateExtensionInstance.Metadata.Keywords.Keyword = link.Keywords
	role := ""
	if link.IntendedEndUserRole.Mentor {
		role = role + " Mentor"
	}
	if link.IntendedEndUserRole.Learner {
		role = role + " Learner"
	}
	if link.IntendedEndUserRole.Instructor {
		role = role + " Instructor"
	}
	mes.CreateExtensionInstance.Metadata.IntendedEndUserRole = role
	mes.CreateExtensionInstance.Metadata.Grade = link.Grade
	mes.CreateExtensionInstance.Metadata.ThumbnailUrl = link.ThumbnailUrl
	edu := ""
	if link.EducationalIntent.ProfessionalDevelopment {
		edu = edu + " ProfessionalDevelopment"
	}
	if link.EducationalIntent.Practice {
		edu = edu + " Practice"
	}
	if link.EducationalIntent.Instructional {
		edu = edu + " Instructional"
	}
	if link.EducationalIntent.Assessment {
		edu = edu + " Assessment"
	}
	if link.EducationalIntent.Activity {
		edu = edu + " Activity"
	}
	mes.CreateExtensionInstance.Metadata.EducationalIntent = edu
	mes.CreateExtensionInstance.Metadata.Publisher = link.Publisher
	if link.basicData.scope.Custom {
		mes.CreateExtensionInstance.Sharing.Scope = "Custom"
	} else if link.basicData.scope.Community {
		mes.CreateExtensionInstance.Sharing.Scope = "Community"
	} else if link.basicData.scope.School {
		mes.CreateExtensionInstance.Sharing.Scope = "School"
	} else if link.basicData.scope.Site {
		mes.CreateExtensionInstance.Sharing.Scope = "Site"
	} else {
		mes.CreateExtensionInstance.Sharing.Scope = "Private"
	}
	mes.CreateExtensionInstance.Content.FileLinkContent.Description = link.Description
	mes.CreateExtensionInstance.Content.FileLinkContent.Link = link.Url
	mes.CreateExtensionInstance.Content.FileLinkContent.HideLink = true
	mes.CreateExtensionInstance.ElementProperties.Active = true
	/*
		byteAr, err := xml.Marshal(mes)
		if err != nil {
			r = p
			return
		}
		mess := fmt.Sprint(" <![CDATA[", string(byteAr), "]]>")
	*/
	p.Body.Message.AddMessage.DataMessage.Data.Message = *mes
	p.Body.Message.AddMessage.DataMessage.Type = "49"
	r = p
	return
}

type createLinkMessage struct {
	XMLName  xml.Name `xml:"Message"`
	Xmlns    string   `xml:"xmlns,attr"`
	SyncKeys struct {
		SyncKey string `xml:"SyncKey"`
	} `xml:"SyncKeys"`
	VendorId                string `xml:"VendorId"`
	CreateExtensionInstance struct {
		Location    string `xml:"Location"`
		ExtensionId string `xml:"ExtensionId"`
		UserSyncKey string `xml:"UserSyncKey"`
		Title       string `xml:"Title"`
		Metadata    struct {
			Description string `xml:"Description"`
			Language    string `xml:"Language"`
			Format      string `xml:"Format"`
			Keywords    struct {
				Keyword []string `xml:"Keyword"`
			} `xml:"Keywords"`
			IntendedEndUserRole string `xml:"IntendedEndUserRole"`
			Grade               string `xml:"Grade"`
			ThumbnailUrl        string `xml:"ThumbnailUrl"`
			EducationalIntent   string `xml:"EducationalIntent"`
			Publisher           string `xml:"Publisher"`
		} `xml:"Metadata"`
		Sharing struct {
			Scope string `xml:"Scope"`
		} `xml:"Sharing"`
		Content struct {
			FileLinkContent struct {
				Description string `xml:"Description"`
				HideLink    bool   `xml:"HideLink"`
				Link        string `xml:"Link"`
			} `xml:"FileLinkContent"`
		} `xml:"Content"`
		ElementProperties struct {
			Active bool `xml:"Active"`
		} `xml:"ElementProperties"`
	} `xml:"CreateExtensionInstance"`
}

/*
Update an existing Link in itslearning.
*/
func (p *SoapRequest) UpdateLink(link Link) (err error, r *SoapRequest) {
	return
}

type updateLinkMessage struct {
}

/*
Sends the SoapRequest and printed out the response
*/
func (p *SoapRequest) Send(url, soapAction string) (err error, response string) {
	soapbyte, err := xml.Marshal(p)
	if err != nil {
		return
	}
	soap := string(soapbyte)
	fmt.Println("======================================")
	fmt.Println(soap)
	fmt.Println("======================================")
	httpMethod := "POST"
	req, err := http.NewRequest(httpMethod, url, strings.NewReader(soap))
	if err != nil {
		return
	}
	req.Header.Set("Content-type", `text/xml; charset="utf-8"`)
	req.Header.Set("SOAPAction", soapAction)

	client := &http.Client{}

	res, err := client.Do(req)
	if err != nil {
		return
	}

	bodyBytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return
	}

	response = string(bodyBytes)
	return
}
