package ItsOrgaApi

import (
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
)

type Link struct {
	Id          string
	Description string
	Language    string
	FormatAny         bool
	FormatAudio       bool
	FormatImage       bool
	FormatInteractive bool
	FormatText        bool
	FormatVideo       bool
	Keywords            string // List with commas
	IntendedEndUserRoleLearner    bool
	IntendedEndUserRoleInstructor bool
	IntendedEndUserRoleMentor     bool
	Grade             string
	ThumbnailUrl      string
	EducationalIntentInstructional           bool
	EducationalIntentPractice                bool
	EducationalIntentProfessionalDevelopment bool
	EducationalIntentAssessment              bool
	EducationalIntentActivity                bool
	Publisher string
	Title     string
	Url       string
	basicData ItslearningBasicData
}

type ItslearningBasicData struct {
	userSyncKey string
	vendorID    string
	location    struct {
		Course  bool
		Library bool
	}
	scope struct {
		Private   bool
		School    bool
		Site      bool
		Community bool
		Custom    bool
	}
}

func NewLink() *Link {
	return new(Link)
}

func NewItslearningBasicData() *ItslearningBasicData {
	return new(ItslearningBasicData)
}

func (p *ItslearningBasicData) SetItslearningBasicData(vendorID, location, userSyncKey, scope string, locationCourse, locationLibrary, scopePrivate, scopeSchool, scopeSite, scopeCommunity, scopeCustom bool) (err error, r *ItslearningBasicData) {
	p.vendorID = vendorID
	p.userSyncKey = userSyncKey
	p.location.Course = locationCourse
	p.location.Library = locationLibrary
	count := 0
	if scopeSite {
		count++
	}
	if scopeSchool {
		count++
	}
	if scopeCommunity {
		count++
	}
	if scopeCustom {
		count++
	}
	if scopePrivate {
		count++
	}
	if count > 1 {
		err = errors.New("Scope can only one Item!")
	}

	p.scope.Site = scopeSite
	p.scope.School = scopeSchool
	p.scope.Community = scopeCommunity
	p.scope.Custom = scopeCustom
	p.scope.Private = scopePrivate
	r = p
	return
}

func (p *Link) SetLinkData(title, description, language, format, intendedEndUserRole, grade, thumbnailUrl, educationalIntent, publisher, url, id string, keywords []string, EducationalIntentInstructional, EducationalIntentPractice, EducationalIntentProfessionalDevelopment, EducationalIntentAssessment, EducationalIntentActivity, IntendedEndUserRoleLearner, IntendedEndUserRoleInstructor, IntendedEndUserRoleMentor bool, FormatAny, FormatAudio, FortmatImage, FormatInteractive, FormatText, FormatVideo bool) (err error, r *Link) {
	p.Id = id
	p.Title = title
	p.Description = description
	p.Language = language

	count := 0
	if FormatAny {
		count++
	}
	if FormatAudio {
		count++
	}
	if FormatInteractive {
		count++
	}
	if FormatText {
		count++
	}
	if FormatVideo {
		count++
	}
	if FortmatImage {
		count++
	}
	if count > 1 {
		err = errors.New("Only one Format is allowed!")
	}
	p.FormatText = FormatText
	p.FormatAny = FormatAny
	p.FormatAudio = FormatAudio
	p.FormatImage = FortmatImage
	p.FormatInteractive = FormatInteractive
	p.FormatVideo = FormatVideo
	tmp := ""
	for i := 0; i < len (keywords); i++ {
		if i == 0 {
			tmp = keywords[i]
		}else {
			tmp = tmp + "," + keywords[i]
		}
	}
	p.Keywords = tmp
	p.IntendedEndUserRoleInstructor = IntendedEndUserRoleInstructor
	p.IntendedEndUserRoleLearner = IntendedEndUserRoleLearner
	p.IntendedEndUserRoleMentor = IntendedEndUserRoleMentor
	p.Grade = grade
	p.ThumbnailUrl = thumbnailUrl
	p.EducationalIntentActivity = EducationalIntentActivity
	p.EducationalIntentAssessment = EducationalIntentAssessment
	p.EducationalIntentInstructional = EducationalIntentInstructional
	p.EducationalIntentPractice = EducationalIntentPractice
	p.EducationalIntentProfessionalDevelopment = EducationalIntentProfessionalDevelopment
	p.Publisher = publisher
	p.Url = url
	r = p
	return
}

func (p *Link) SetItslearningBasicData(data ItslearningBasicData) *Link {
	p.basicData = data
	return p
}

//TODO:
func (p *Link) StoreInDataBase(db *gorm.DB) *Link {

	return p
}

