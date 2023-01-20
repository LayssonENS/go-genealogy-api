package repository

import (
	"database/sql"
	"github.com/LayssonENS/go-genealogy-api/domain"
	"github.com/gin-gonic/gin"
)

type postgresPersonRepo struct {
	DB *sql.DB
}

// NewPostgresRelationshipRepository will create an implementation of relationship.Repository
func NewPostgresRelationshipRepository(db *sql.DB) domain.RelationshipRepository {
	return &postgresPersonRepo{
		DB: db,
	}
}

func (p *postgresPersonRepo) GetRelationshipByID(c *gin.Context, id int64) (domain.Relationship, error) {
	var relationship domain.Relationship

	result, _ := p.getFamilyTree("laysson")
	println(result)

	err := p.DB.QueryRow(
		"SELECT id, person_id, related_person_id, relationship  FROM relationships WHERE id = $1", id).Scan(
		&relationship.ID, &relationship.PersonID, &relationship.Relationship, &relationship.RelatedPersonID)
	if err != nil {
		if err == sql.ErrNoRows {
			return relationship, domain.ErrNotFound
		}
		return relationship, err
	}
	return relationship, nil

}

func (p *postgresPersonRepo) CreateRelationship(c *gin.Context, relationship domain.Relationship) error {
	query := `INSERT INTO relationships (person_id, related_person_id, relationship) VALUES ($1, $2 , $3) `
	_, err := p.DB.Exec(query,
		relationship.PersonID,
		relationship.RelatedPersonID,
		relationship.Relationship,
	)
	if err != nil {
		return err
	}

	return nil
}

func (p *postgresPersonRepo) getFamilyTree(name string) ([]FamilyMember, error) {
	var family []FamilyMember
	var parents []Person
	var grandparents []Person
	var auntsUncles []Person
	var siblings []Person

	// Busca os pais
	rows, err := p.DB.Query("SELECT p.id, p.name FROM person p JOIN relationships r ON p.id = r.related_person_id WHERE r.person_id = (SELECT id FROM person WHERE name = $1) AND r.relationship = 'parent'", name)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var parent Person
		err := rows.Scan(&parent.ID, &parent.Name)
		if err != nil {
			return nil, err
		}
		parents = append(parents, parent)
	}

	// Adiciona os pais na lista de parentes
	for _, parent := range parents {
		family = append(family, FamilyMember{
			Person:       parent,
			Relationship: "parent",
		})
	}

	// Busca os avós
	for _, parent := range parents {
		rows, err := p.DB.Query("SELECT p.id, p.name FROM person p JOIN relationships r ON p.id = r.related_person_id WHERE r.person_id = $1 AND r.relationship = 'parent'", parent.ID)
		if err != nil {
			return nil, err
		}
		defer rows.Close()

		var grandparentsOfParent []Person
		for rows.Next() {
			var grandparent Person
			err := rows.Scan(&grandparent.ID, &grandparent.Name)
			if err != nil {
				return nil, err
			}
			grandparentsOfParent = append(grandparentsOfParent, grandparent)
		}
		grandparents = append(grandparents, grandparentsOfParent...)
	}

	// Adiciona os avós na lista de parentes
	for _, grandparent := range grandparents {
		family = append(family, FamilyMember{
			Person:       grandparent,
			Relationship: "grandparent",
		})
	}

	// Busca os tios e tias
	for _, grandparent := range grandparents {
		rows, err := p.DB.Query("SELECT p.id, p.name FROM person p JOIN relationships r ON p.id = r.related_person_id WHERE r.person_id = $1 AND r.relationship = 'children'", grandparent.ID)
		if err != nil {
			return nil, err
		}
		defer rows.Close()

		var auntsUnclesOfGrandparent []Person
		for rows.Next() {
			var auntsUncles Person
			err := rows.Scan(&auntsUncles.ID, &auntsUncles.Name)
			if err != nil {
				return nil, err
			}
			auntsUnclesOfGrandparent = append(auntsUnclesOfGrandparent, auntsUncles)
		}
		auntsUncles = append(auntsUncles, auntsUnclesOfGrandparent...)
	}

	// Adiciona os tios e tias na lista de parentes
	for _, auntUncle := range auntsUncles {
		family = append(family, FamilyMember{
			Person:       auntUncle,
			Relationship: "aunt/uncle",
		})
	}

	// Busca os irmãos
	for _, parent := range parents {
		rows, err := p.DB.Query("SELECT p.id, p.name FROM person p JOIN relationships r ON p.id = r.related_person_id WHERE r.person_id = $1 AND r.relationship = 'children'", parent.ID)
		if err != nil {
			return nil, err
		}
		defer rows.Close()

		var siblingParents []Person
		for rows.Next() {
			var sibling Person
			err := rows.Scan(&sibling.ID, &sibling.Name)
			if err != nil {
				return nil, err
			}
			siblingParents = append(siblingParents, sibling)
		}
		siblings = append(siblings, siblingParents...)
	}

	// Adiciona os irmãos na lista de parentes
	for _, sibling := range siblings {
		if sibling.Name != name { // para não incluir o individuo
			family = append(family, FamilyMember{
				Person:       sibling,
				Relationship: "sibling",
			})
		}
	}

	return family, nil
}

type FamilyMember struct {
	Person       Person
	Relationship string
}

type Person struct {
	ID   int
	Name string
}
