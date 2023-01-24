package repository

import (
	"database/sql"
	"github.com/LayssonENS/go-genealogy-api/domain"
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

func (p *postgresPersonRepo) GetRelationshipByID(personId int64) (*domain.FamilyMembers, error) {
	var relationships []domain.FamilyMember
	familyMembers := &domain.FamilyMembers{}
	query := `WITH parents AS (
			SELECT DISTINCT ON (p.id) p.id, p.name, 
				(CASE 
					WHEN r.related_person_id = $1 and relationship='children' THEN 
						CASE 
							WHEN p.id = $1 THEN 'self'
							ELSE 'parent'
						END
					WHEN r.person_id = $1 and relationship='children' THEN 
						CASE 
							WHEN p.id = $1 THEN 'self'
							ELSE relationship
						END
					ELSE relationship
				END) as relationship 
			FROM person p
			INNER JOIN relationships r ON p.id = r.related_person_id OR p.id = r.person_id
			WHERE r.person_id = $1 OR r.related_person_id = $1
			), grandparents AS (
				SELECT p.id, p.name, 'grandparent' AS relationship FROM person p
				INNER JOIN relationships r ON p.id = r.related_person_id
				WHERE r.person_id IN (SELECT id FROM parents where relationship != 'children') AND r.relationship = 'parent' AND p.id NOT in (SELECT id FROM parents)
			), aunts_uncles AS (
				SELECT p.id, p.name, 'aunt/uncle' AS relationship FROM person p
				INNER JOIN relationships r ON p.id = r.related_person_id
				WHERE r.person_id IN (SELECT id FROM grandparents) AND r.relationship = 'children' AND p.id NOT in (SELECT id FROM parents)
			), siblings AS (
				SELECT p.id, p.name, 'sibling' AS relationship FROM person p
				INNER JOIN relationships r ON p.id = r.related_person_id
				WHERE r.person_id IN (SELECT id FROM parents where relationship != 'children') AND r.relationship = 'children' AND p.id NOT in (SELECT id FROM parents)
			), nieces_nephews AS (
				SELECT p.id, p.name, 'cousin' AS relationship FROM person p
				INNER JOIN relationships r ON p.id = r.related_person_id
				WHERE r.person_id IN (SELECT id FROM siblings) AND r.relationship = 'children' AND p.id NOT in (SELECT id FROM parents)
			)
			SELECT distinct * FROM parents 
			UNION
			SELECT * FROM grandparents
			UNION
			SELECT * FROM aunts_uncles
			UNION
			SELECT * FROM siblings
			UNION
			SELECT * FROM nieces_nephews
			`
	rows, err := p.DB.Query(query, personId)
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

		if parent.ID == personId {
			familyMembers.ID = parent.ID
			familyMembers.Name = parent.Name
			continue
		}

		relationships = append(relationships, parent)
	}

	familyMembers.Members = relationships

	return familyMembers, nil

}

func (p *postgresPersonRepo) CreateRelationship(personId int64, relationship domain.Relationship) error {

	prepareQuery, err := p.DB.Prepare("INSERT INTO relationships (person_id, related_person_id, relationship) VALUES($1, $2, $3)")
	if err != nil {
		return err
	}
	defer prepareQuery.Close()

	if relationship.ChildrenId != 0 {
		_, err = prepareQuery.Exec(personId, relationship.ChildrenId, "children")
		if err != nil {
			return err
		}
	}
	if relationship.ParentId != 0 {
		_, err = prepareQuery.Exec(personId, relationship.ParentId, "parent")
		if err != nil {
			return err
		}
	}
	return nil
}
