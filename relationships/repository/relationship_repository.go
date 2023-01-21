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

func (p *postgresPersonRepo) GetRelationshipByID(c *gin.Context, id int64) (*domain.FamilyMembers, error) {
	var relationships []domain.FamilyMember
	familyMembers := &domain.FamilyMembers{}

	rows, err := p.DB.Query("WITH parents AS ( "+
		"SELECT p.id, p.name, 'parent' AS relationship FROM person p "+
		"JOIN relationships r ON p.id = r.related_person_id "+
		"WHERE r.person_id = (SELECT id FROM person WHERE id = $1) AND r.relationship = 'parent' "+
		"), grandparents AS ( "+
		"SELECT p.id, p.name, 'grandparent' AS relationship FROM person p "+
		"JOIN relationships r ON p.id = r.related_person_id "+
		"WHERE r.person_id IN (SELECT id FROM parents) AND r.relationship = 'parent' "+
		"), aunts_uncles AS ( "+
		"SELECT p.id, p.name, 'aunt/uncle' AS relationship FROM person p "+
		"JOIN relationships r ON p.id = r.related_person_id "+
		"WHERE r.person_id IN (SELECT id FROM grandparents) AND r.relationship = 'children' AND p.id NOT in (SELECT id FROM parents) "+
		"), siblings AS ( "+
		"SELECT p.id, p.name, 'sibling' AS relationship FROM person p "+
		"JOIN relationships r ON p.id = r.related_person_id "+
		"WHERE r.person_id IN (SELECT id FROM parents) AND r.relationship = 'children' "+
		"), nieces_nephews AS ( "+
		"SELECT p.id, p.name, 'cousin' AS relationship FROM person p "+
		"JOIN relationships r ON p.id = r.related_person_id "+
		"WHERE r.person_id IN (SELECT id FROM siblings) AND r.relationship = 'children' "+
		") "+
		"SELECT * FROM parents "+
		"UNION "+
		"SELECT * FROM grandparents "+
		"UNION "+
		"SELECT * FROM aunts_uncles "+
		"UNION "+
		"SELECT * FROM siblings "+
		"UNION "+
		"SELECT * FROM nieces_nephews ", id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var parent domain.FamilyMember
		err := rows.Scan(&parent.ID, &parent.Name, &parent.Relationship)
		if err != nil {
			return nil, err
		}

		if parent.ID == id {
			familyMembers.ID = parent.ID
			familyMembers.Name = parent.Name
			continue
		}

		relationships = append(relationships, parent)
	}

	familyMembers.Members = relationships

	return familyMembers, nil

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

func (p *postgresPersonRepo) getFamilyTree(name string) (*domain.FamilyMembers, error) {
	var relationships []domain.FamilyMember
	familyMembers := &domain.FamilyMembers{}

	rows, err := p.DB.Query("WITH parents AS ( " +
		"SELECT p.id, p.name, 'parent' AS relationship FROM person p " +
		"JOIN relationships r ON p.id = r.related_person_id " +
		"WHERE r.person_id = (SELECT id FROM person WHERE name = 'laysson') AND r.relationship = 'parent' " +
		"), grandparents AS ( " +
		"SELECT p.id, p.name, 'grandparent' AS relationship FROM person p " +
		"JOIN relationships r ON p.id = r.related_person_id " +
		"WHERE r.person_id IN (SELECT id FROM parents) AND r.relationship = 'parent' " +
		"), aunts_uncles AS ( " +
		"SELECT p.id, p.name, 'aunt/uncle' AS relationship FROM person p " +
		"JOIN relationships r ON p.id = r.related_person_id " +
		"WHERE r.person_id IN (SELECT id FROM grandparents) AND r.relationship = 'children' AND p.id NOT in (SELECT id FROM parents) " +
		"), siblings AS ( " +
		"SELECT p.id, p.name, 'sibling' AS relationship FROM person p " +
		"JOIN relationships r ON p.id = r.related_person_id " +
		"WHERE r.person_id IN (SELECT id FROM parents) AND r.relationship = 'children' " +
		"), nieces_nephews AS ( " +
		"SELECT p.id, p.name, 'cousin' AS relationship FROM person p " +
		"JOIN relationships r ON p.id = r.related_person_id " +
		"WHERE r.person_id IN (SELECT id FROM siblings) AND r.relationship = 'children' " +
		") " +
		"SELECT * FROM parents " +
		"UNION " +
		"SELECT * FROM grandparents " +
		"UNION " +
		"SELECT * FROM aunts_uncles " +
		"UNION " +
		"SELECT * FROM siblings " +
		"UNION " +
		"SELECT * FROM nieces_nephews ")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var parent domain.FamilyMember
		err := rows.Scan(&parent.ID, &parent.Name, &parent.Relationship)
		if err != nil {
			return nil, err
		}

		if parent.Name == name {
			familyMembers.ID = parent.ID
			familyMembers.Name = parent.Name
			continue
		}

		relationships = append(relationships, parent)
	}

	familyMembers.Members = relationships

	return familyMembers, nil
}
