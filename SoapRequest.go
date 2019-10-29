package ItsOrgaApi

import (
	"encoding/xml"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

type SoapRequest struct {
	username   string
	password   string
	method uint
	/*
	createLink = 1
	updateLink = 2
	deleteLink = 4
	 */
	mesage     string
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
func (p *SoapRequest) AddSecurity(sec WSSESecurity) *SoapRequest {
	p.username = sec.username
	p.password = sec.passwort
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
	byteAr, err := xml.Marshal(mes)
	if err != nil {
		r = p
		return
	}
	mess := string(byteAr)
	p.mesage = mess
	p.method = 1
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
Sends the SoapRequest and printed out the response.
The response is a number for the Queue when success is true
*/
func (p *SoapRequest) Send() (response string, success bool, err error) {
	const cont = `text/xml; charset="utf-8"`
	switch p.method {
	case 1:
		action := "http://tempuri.org/IDataService/AddMessage"
		begin := CreateSoapBegin(*p, p.username, p.password)
		tail := CreateSoapTailInput(p)
		soap := fmt.Sprint(begin, p.mesage, tail)
		httpMethod := "POST"
		req, err := http.NewRequest(httpMethod, "https://migra.itsltest.com/DataService.svc", strings.NewReader(soap))
		if err != nil {
			return "",false, err
		}
		req.Header.Set("Content-type", cont)
		req.Header.Set("SOAPAction", action)

		client := &http.Client{}

		res, err := client.Do(req)
		if err != nil {
			return "",false, err
		}

		bodyBytes, err := ioutil.ReadAll(res.Body)
		response = string(bodyBytes)


		// check if it works
		if !strings.Contains(res.Status, "200") {
			if strings.Contains(response, "faultstring") {
				faultstringStruct := faultstringStruct{}
				tmpData := []byte(response)
				err = xml.Unmarshal(tmpData, &faultstringStruct)
				if err != nil {
					return response, false, err
				}
				return faultstringStruct.Body.Fault.Faultstring.Text, false, nil
			} else {
				success = false
				return response, false, err
			}
		}
		//
		responseMessage := responseMessage{}
		data := []byte(response)
		err = xml.Unmarshal(data, &responseMessage)
		if err != nil {
			return response,false, err
		}
		return responseMessage.Body.AddMessageResponse.AddMessageResult.MessageId,true, nil


	case 2:
		action := "http://tempuri.org/IDataService/AddMessage"
		begin := CreateSoapBegin(*p, p.username, p.password)
		tail := CreateSoapTailInput(p)
		soap := fmt.Sprint(begin, p.mesage, tail)
		httpMethod := "POST"
		req, err := http.NewRequest(httpMethod, "https://migra.itsltest.com/DataService.svc", strings.NewReader(soap))
		fmt.Println(err)

		req.Header.Set("Content-type", cont)
		req.Header.Set("SOAPAction", action)

		client := &http.Client{}

		res, err := client.Do(req)

		bodyBytes, err := ioutil.ReadAll(res.Body)
		response = string(bodyBytes)
		// check if it works
		if !strings.Contains(res.Status, "200") {
			if strings.Contains(response, "faultstring") {
				faultstringStruct := faultstringStruct{}
				tmpData := []byte(response)
				err = xml.Unmarshal(tmpData, &faultstringStruct)
				if err != nil {
					return response, false, err
				}
				return faultstringStruct.Body.Fault.Faultstring.Text, false, nil
			} else {
				success = false
				return response, false, err
			}
		}
		//
		responseMessage := responseMessage{}
		data := []byte(response)
		err = xml.Unmarshal(data, &responseMessage)
		if err != nil {
			return response,false, err
		}
		return responseMessage.Body.AddMessageResponse.AddMessageResult.MessageId,true, nil



		return response,true, nil
	case 4:
		action := "http://tempuri.org/IDataService/AddMessage"
		begin := CreateSoapBegin(*p, p.username, p.password)
		tail := CreateSoapTailInput(p)
		soap := fmt.Sprint(begin, p.mesage, tail)
		fmt.Println("soap")
		fmt.Println(soap)
		httpMethod := "POST"
		req, err := http.NewRequest(httpMethod, "https://migra.itsltest.com/DataService.svc", strings.NewReader(soap))
		fmt.Println(err)

		req.Header.Set("Content-type", cont)
		req.Header.Set("SOAPAction", action)

		client := &http.Client{}

		res, err := client.Do(req)

		bodyBytes, err := ioutil.ReadAll(res.Body)
		response = string(bodyBytes)
		// check if it works
		if !strings.Contains(res.Status, "200") {
			if strings.Contains(response, "faultstring") {
				faultstringStruct := faultstringStruct{}
				tmpData := []byte(response)
				err = xml.Unmarshal(tmpData, &faultstringStruct)
				if err != nil {
					return response, false, err
				}
				return faultstringStruct.Body.Fault.Faultstring.Text, false, nil
			} else {
				success = false
				return response, false, err
			}
		}
		//
		responseMessage := responseMessage{}
		data := []byte(response)
		err = xml.Unmarshal(data, &responseMessage)
		if err != nil {
			return response,false, err
		}
		return responseMessage.Body.AddMessageResponse.AddMessageResult.MessageId,true, nil

	default:
		return "",false, nil

	}

	return
}


func CreateSoapBegin(soapRequest SoapRequest, username string, password string) string {
	switch soapRequest.method {
	case 1:
		begin := `<x:Envelope xmlns:x="http://schemas.xmlsoap.org/soap/envelope/" xmlns:tem="http://tempuri.org/" xmlns:its="http://schemas.datacontract.org/2004/07/Itslearning.Integration.ContentImport.Services.Entities">
    <x:Header>
        <wsse:Security xmlns:wsse="http://docs.oasis-open.org/wss/2004/01/oasis-200401-wss-wssecurity-secext-1.0.xsd" xmlns:wsu="http://docs.oasis-open.org/wss/2004/01/oasis-200401-wss-wssecurity-utility-1.0.xsd">
            <wsse:UsernameToken>
                <wsse:Username>` + username + `</wsse:Username>
                <wsse:Password Type="http://docs.oasis-open.org/wss/2004/01/oasis-200401-wss-username-token-profile-1.0#PasswordText">` + password + `</wsse:Password>
            </wsse:UsernameToken>
        </wsse:Security>
    </x:Header>
    <x:Body>
        <tem:AddMessage>
            <tem:dataMessage>
                <its:Data>
                    <![CDATA[`
		return begin
	case 2:
		begin := `<x:Envelope xmlns:x="http://schemas.xmlsoap.org/soap/envelope/" xmlns:tem="http://tempuri.org/" xmlns:its="http://schemas.datacontract.org/2004/07/Itslearning.Integration.ContentImport.Services.Entities">
<x:Header>
	<wsse:Security xmlns:wsse="http://docs.oasis-open.org/wss/2004/01/oasis-200401-wss-wssecurity-secext-1.0.xsd" xmlns:wsu="http://docs.oasis-open.org/wss/2004/01/oasis-200401-wss-wssecurity-utility-1.0.xsd">
<wsse:UsernameToken>
                <wsse:Username>` + username + `</wsse:Username>
                <wsse:Password Type="http://docs.oasis-open.org/wss/2004/01/oasis-200401-wss-username-token-profile-1.0#PasswordText">` + password + `</wsse:Password>
	</wsse:UsernameToken>
	</wsse:Security>
	</x:Header>
	<x:Body>
	<tem:AddMessage>
	<tem:dataMessage>
	<its:Data>
	<![CDATA[`
		return begin
	case 4:
		begin := `
<x:Envelope xmlns:x="http://schemas.xmlsoap.org/soap/envelope/" xmlns:tem="http://tempuri.org/" xmlns:its="http://schemas.datacontract.org/2004/07/Itslearning.Integration.ContentImport.Services.Entities">
    <x:Header>
        <wsse:Security xmlns:wsse="http://docs.oasis-open.org/wss/2004/01/oasis-200401-wss-wssecurity-secext-1.0.xsd" xmlns:wsu="http://docs.oasis-open.org/wss/2004/01/oasis-200401-wss-wssecurity-utility-1.0.xsd">
            <wsse:UsernameToken>
                <wsse:Username>` + username + `</wsse:Username>
                <wsse:Password Type="http://docs.oasis-open.org/wss/2004/01/oasis-200401-wss-username-token-profile-1.0#PasswordText">` + password + `</wsse:Password>
            </wsse:UsernameToken>
        </wsse:Security>
    </x:Header>
    <x:Body>
        <tem:AddMessage>
            <tem:dataMessage>
                <its:Data>
                           <![CDATA[
`
			return begin
		default:
		return ""
	}
	return ""
}
func CreateSoapTailInput(p *SoapRequest) (tail string) {
	switch p.method {
	case 1:
		tail =  `
				]]>
                </its:Data>
                <its:Type>49</its:Type>
            </tem:dataMessage>
        </tem:AddMessage>
    </x:Body>
</x:Envelope>`
		return
	case 2:
		tail = `
				]]>
			</its:Data>
                <its:Type>54</its:Type>
            </tem:dataMessage>
        </tem:AddMessage>
    </x:Body>
</x:Envelope>	
`
		return
	case 4:
		tail = `
				]]>
                    
                </its:Data>
                <its:Type>55</its:Type>
            </tem:dataMessage>
        </tem:AddMessage>
    </x:Body>
</x:Envelope>
       `
		return
	default:
		return
	}
	return
}


/*
Update an existing Link in itslearning.
*/
func (p *SoapRequest) UpdateLink(link Link) (err error, r *SoapRequest) {
	mes := new(updateLinkMessage)
	mes.Xmlns = "urn:message-schema"
	mes.VendorId = link.basicData.vendorID

	mes.UpdateExtensionInstance.ContentSyncKey = link.Id
	mes.UpdateExtensionInstance.UserSyncKey = link.basicData.userSyncKey
	mes.UpdateExtensionInstance.Title = link.Title
	mes.UpdateExtensionInstance.Metadata.Description = link.Description
	mes.UpdateExtensionInstance.Metadata.Language = link.Language
	if link.Format.Video {
		mes.UpdateExtensionInstance.Metadata.Format = "Video"
	} else if link.Format.Interactive {
		mes.UpdateExtensionInstance.Metadata.Format = "Interactive"
	} else if link.Format.Image {
		mes.UpdateExtensionInstance.Metadata.Format = "Image"
	} else if link.Format.Audio {
		mes.UpdateExtensionInstance.Metadata.Format = "Audio"
	} else if link.Format.Any {
		mes.UpdateExtensionInstance.Metadata.Format = "Any"
	} else if link.Format.Text {
		mes.UpdateExtensionInstance.Metadata.Format = "Text"
	} else {
		mes.UpdateExtensionInstance.Metadata.Format = "Any"
	}
	mes.UpdateExtensionInstance.Metadata.Keywords.Keyword = link.Keywords
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
	mes.UpdateExtensionInstance.Metadata.IntendedEndUserRole = role
	mes.UpdateExtensionInstance.Metadata.Grade = link.Grade
	mes.UpdateExtensionInstance.Metadata.ThumbnailUrl = link.ThumbnailUrl
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
	mes.UpdateExtensionInstance.Metadata.EducationalIntent = edu
	mes.UpdateExtensionInstance.Metadata.Publisher = link.Publisher
	if link.basicData.scope.Custom {
		mes.UpdateExtensionInstance.Sharing.Scope = "Custom"
	} else if link.basicData.scope.Community {
		mes.UpdateExtensionInstance.Sharing.Scope = "Community"
	} else if link.basicData.scope.School {
		mes.UpdateExtensionInstance.Sharing.Scope = "School"
	} else if link.basicData.scope.Site {
		mes.UpdateExtensionInstance.Sharing.Scope = "Site"
	} else {
		mes.UpdateExtensionInstance.Sharing.Scope = "Private"
	}
	mes.UpdateExtensionInstance.Content.FileLinkContent.Description = link.Description
	mes.UpdateExtensionInstance.Content.FileLinkContent.Link = link.Url
	mes.UpdateExtensionInstance.Content.FileLinkContent.HideLink = true
	mes.UpdateExtensionInstance.ElementProperties.Active = true

	byteAr, err := xml.Marshal(mes)
	if err != nil {
		r = p
		return
	}
	mess := string(byteAr)
	p.mesage = mess
	p.method = 2
	r = p
	return
}
type updateLinkMessage struct {
	XMLName                 xml.Name `xml:"Message"`
	Xmlns                   string   `xml:"xmlns,attr"`
	VendorId                string   `xml:"VendorId"`
	UpdateExtensionInstance struct {
		ContentSyncKey string `xml:"ContentSyncKey"`
		UserSyncKey    string `xml:"UserSyncKey"`
		Title          string `xml:"Title"`
		Metadata       struct {
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
	} `xml:"UpdateExtensionInstance"`
}


/*
Delete a Link in itslearning
 */
func (p *SoapRequest) DeleteLink (link Link) (err error, r *SoapRequest) {
	mes := new(deleteLinkMessage)
	mes.Xmlns = "urn:message-schema"
	mes.VendorId = link.basicData.vendorID
	mes.DeleteExtensionInstance.ContentSyncKey = link.Id
	mes.DeleteExtensionInstance.UserSyncKey =  link.basicData.userSyncKey

	byteAr, err := xml.Marshal(mes)
	if err != nil {
		r = p
		return
	}
	mess := string(byteAr)

	p.mesage = mess
	p.method = 4
	r = p
	return
}

type deleteLinkMessage struct {
	XMLName                 xml.Name `xml:"Message"`
	Xmlns                   string   `xml:"xmlns,attr"`
	VendorId                string   `xml:"VendorId"`
	DeleteExtensionInstance struct {
	ContentSyncKey string `xml:"ContentSyncKey"`
	UserSyncKey    string `xml:"UserSyncKey"`
}


}

type responseMessage struct {
	XMLName xml.Name `xml:"Envelope"`
	Text    string   `xml:",chardata"`
	S       string   `xml:"s,attr"`
	U       string   `xml:"u,attr"`
	Header  struct {
	Text     string `xml:",chardata"`
	Security struct {
	Text           string `xml:",chardata"`
	MustUnderstand string `xml:"mustUnderstand,attr"`
	O              string `xml:"o,attr"`
	Timestamp      struct {
	Text    string `xml:",chardata"`
	ID      string `xml:"Id,attr"`
	Created string `xml:"Created"`
	Expires string `xml:"Expires"`
	} `xml:"Timestamp"`
	} `xml:"Security"`
	} `xml:"Header"`
	Body struct {
	Text               string `xml:",chardata"`
	AddMessageResponse struct {
	Text             string `xml:",chardata"`
	Xmlns            string `xml:"xmlns,attr"`
	AddMessageResult struct {
	Text          string `xml:",chardata"`
	A             string `xml:"a,attr"`
	I             string `xml:"i,attr"`
	MessageId     string `xml:"MessageId"`
	Status        string `xml:"Status"`
	StatusDetails struct {
	Text string `xml:",chardata"`
	Nil  string `xml:"nil,attr"`
	} `xml:"StatusDetails"`
	} `xml:"AddMessageResult"`
	} `xml:"AddMessageResponse"`
	} `xml:"Body"`
	}



func (p *SoapRequest) GetMessageResult (messageID string)(result string, err error){
	action := "http://tempuri.org/IDataService/GetMessageResult"
	soap := `
<x:Envelope xmlns:x="http://schemas.xmlsoap.org/soap/envelope/" xmlns:tem="http://tempuri.org/">
    <x:Header>
        <wsse:Security xmlns:wsse="http://docs.oasis-open.org/wss/2004/01/oasis-200401-wss-wssecurity-secext-1.0.xsd" xmlns:wsu="http://docs.oasis-open.org/wss/2004/01/oasis-200401-wss-wssecurity-utility-1.0.xsd">
            <wsse:UsernameToken>
                <wsse:Username>`+p.username+`</wsse:Username>
                <wsse:Password Type="http://docs.oasis-open.org/wss/2004/01/oasis-200401-wss-username-token-profile-1.0#PasswordText">`+p.password+`</wsse:Password>
            </wsse:UsernameToken>
        </wsse:Security>
    </x:Header>
    <x:Body>
        <tem:GetMessageResult>
            <tem:messageId>`+messageID+`</tem:messageId>
        </tem:GetMessageResult>
    </x:Body>
</x:Envelope>`

	httpMethod := "POST"
	req, err := http.NewRequest(httpMethod, "https://migra.itsltest.com/DataService.svc", strings.NewReader(soap))
	if err != nil {
		return
	}
	req.Header.Set("Content-type", "text/xml; charset=utf-8")
	req.Header.Set("SOAPAction", action)

	client := &http.Client{}

	res, err := client.Do(req)

	bodyBytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return
	}
	tmpStruct := GetMessageResultStruct{}
	xml.Unmarshal(bodyBytes, &tmpStruct)
	return tmpStruct.Body.GetMessageResultResponse.GetMessageResultResult.StatusDetails.DataMessageStatusDetail.Message, nil
	
}

type GetMessageResultStruct struct {
	XMLName xml.Name `xml:"Envelope"`
	Text    string   `xml:",chardata"`
	S       string   `xml:"s,attr"`
	U       string   `xml:"u,attr"`
	Header  struct {
		Text     string `xml:",chardata"`
		Security struct {
			Text           string `xml:",chardata"`
			MustUnderstand string `xml:"mustUnderstand,attr"`
			O              string `xml:"o,attr"`
			Timestamp      struct {
				Text    string `xml:",chardata"`
				ID      string `xml:"Id,attr"`
				Created string `xml:"Created"`
				Expires string `xml:"Expires"`
			} `xml:"Timestamp"`
		} `xml:"Security"`
	} `xml:"Header"`
	Body struct {
		Text                     string `xml:",chardata"`
		GetMessageResultResponse struct {
			Text                   string `xml:",chardata"`
			Xmlns                  string `xml:"xmlns,attr"`
			GetMessageResultResult struct {
				Text          string `xml:",chardata"`
				A             string `xml:"a,attr"`
				I             string `xml:"i,attr"`
				MessageId     string `xml:"MessageId"`
				Status        string `xml:"Status"`
				StatusDetails struct {
					Text                    string `xml:",chardata"`
					DataMessageStatusDetail struct {
						Text    string `xml:",chardata"`
						Entity  string `xml:"Entity"`
						Message string `xml:"Message"`
						SyncKey string `xml:"SyncKey"`
						Type    string `xml:"Type"`
					} `xml:"DataMessageStatusDetail"`
				} `xml:"StatusDetails"`
			} `xml:"GetMessageResultResult"`
		} `xml:"GetMessageResultResponse"`
	} `xml:"Body"`
}


type faultstringStruct struct {
	XMLName xml.Name `xml:"Envelope"`
	Text    string   `xml:",chardata"`
	S       string   `xml:"s,attr"`
	U       string   `xml:"u,attr"`
	Header  struct {
		Text     string `xml:",chardata"`
		Security struct {
			Text           string `xml:",chardata"`
			MustUnderstand string `xml:"mustUnderstand,attr"`
			O              string `xml:"o,attr"`
			Timestamp      struct {
				Text    string `xml:",chardata"`
				ID      string `xml:"Id,attr"`
				Created string `xml:"Created"`
				Expires string `xml:"Expires"`
			} `xml:"Timestamp"`
		} `xml:"Security"`
	} `xml:"Header"`
	Body struct {
		Text  string `xml:",chardata"`
		Fault struct {
			Text        string `xml:",chardata"`
			Faultcode   string `xml:"faultcode"`
			Faultstring struct {
				Text string `xml:",chardata"`
				Lang string `xml:"lang,attr"`
			} `xml:"faultstring"`
		} `xml:"Fault"`
	} `xml:"Body"`
}
