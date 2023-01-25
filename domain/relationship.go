package domain

import "encoding/xml"

type Relationship struct {
	ParentId   int64 `json:"parent"`
	ChildrenId int64 `json:"children"`
}

type Family struct {
	ID               int64      `json:"id" xml:"id"`
	Name             string     `json:"name" xml:"name"`
	Relationship     string     `json:"relationship" xml:"relationship"`
	FamilyConnection int64      `json:"family_connection" xml:"family_connection"`
	Relationships    []Relation `json:"relationships" xml:"relationships>relation"`
}

type Relation struct {
	ID           int64  `json:"id" xml:"id"`
	Name         string `json:"name" xml:"name"`
	Relationship string `json:"relationship" xml:"relationship"`
}

type Member struct {
	XMLName xml.Name `json:"-" xml:"family"`
	ID      int64    `json:"id" xml:"id"`
	Name    string   `json:"name" xml:"name"`
	Members []Family `json:"members" xml:"members>member"`
}

type RelationshipUseCase interface {
	GetRelationshipByID(personId int64) (*Member, error)
	CreateRelationship(personId int64, person Relationship) error
}

type RelationshipRepository interface {
	GetRelationshipByID(personId int64) (*Member, error)
	CreateRelationship(personId int64, relationship Relationship) error
}
