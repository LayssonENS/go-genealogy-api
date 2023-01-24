package domain

type Relationship struct {
	ParentId   int64 `json:"parent"`
	ChildrenId int64 `json:"children"`
}

type FamilyMember struct {
	ID               int64  `json:"id"`
	Name             string `json:"name"`
	Relationship     string `json:"relationship"`
	FamilyConnection int64  `json:"family_connection"`
}

type FamilyMembers struct {
	ID      int64          `json:"id"`
	Name    string         `json:"name"`
	Members []FamilyMember `json:"members"`
}

type RelationshipUseCase interface {
	GetRelationshipByID(personId int64) (*FamilyMembers, error)
	CreateRelationship(personId int64, person Relationship) error
}

type RelationshipRepository interface {
	GetRelationshipByID(personId int64) (*FamilyMembers, error)
	CreateRelationship(personId int64, relationship Relationship) error
}
