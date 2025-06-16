package service

import (
	"context"
	"github.com/cebas/go-util/google/auth"
	"github.com/cebas/go-util/util"
	"google.golang.org/api/option"
	"google.golang.org/api/people/v1"
)

// People is a service for managing Google People API

type People struct {
	auth    auth.Gauth
	service *people.Service
}

func NewPeople(ctx context.Context, credentialsFile string) *People {
	newGauth := auth.NewGauth(credentialsFile, people.ContactsReadonlyScope)

	httpClient, err := newGauth.HttpClient()
	if err != nil {
		util.FatalErrorCheck(err)
	}

	newService, err := people.NewService(ctx, option.WithHTTPClient(httpClient))
	if err != nil {
		util.FatalErrorCheck(err)
	}

	return &People{
		auth:    newGauth,
		service: newService,
	}
}

func (people *People) Connections() []*people.Person {
	personFields := "addresses,ageRanges,biographies,birthdays,calendarUrls,clientData,coverPhotos,emailAddresses,events,externalIds,genders,imClients,interests,locales,locations,memberships,metadata,miscKeywords,names,nicknames,occupations,organizations,phoneNumbers,photos,relations,sipAddresses,skills,urls,userDefined"

	connections, err := people.service.People.Connections.
		List("people/me").
		PageSize(10).
		PersonFields(personFields).
		Do()
	util.FatalErrorCheck(err)

	return connections.Connections
}
