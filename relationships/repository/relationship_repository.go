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
				LEFT JOIN relationships r ON p.id = r.related_person_id OR p.id = r.person_id
				WHERE r.person_id = $1 OR r.related_person_id = $1
				), grandparents AS (
					SELECT p.id, p.name, 'grandparent' AS relationship FROM person p
					JOIN relationships r ON p.id = r.related_person_id
					WHERE r.person_id IN (SELECT id FROM parents where relationship != 'children') AND r.relationship = 'parent' AND p.id NOT in (SELECT id FROM parents)
				), aunts_uncles AS (
					SELECT p.id, p.name, 'aunt_uncle' AS relationship FROM person p
					JOIN relationships r ON p.id = r.related_person_id
					WHERE r.person_id IN (SELECT id FROM grandparents) AND r.relationship = 'children' AND p.id NOT in (SELECT id FROM parents)
				), cousin AS (
					SELECT p.id, p.name, 'cousin' AS relationship FROM person p
					JOIN relationships r ON p.id = r.related_person_id
					WHERE r.person_id IN (SELECT id FROM aunts_uncles) AND r.relationship = 'children' AND p.id NOT in (SELECT id FROM parents)
				),siblings AS (
					SELECT p.id, p.name, 'sibling' AS relationship FROM person p
					JOIN relationships r ON p.id = r.related_person_id
					WHERE r.person_id IN (SELECT id FROM parents where relationship != 'children') AND r.relationship = 'children' AND p.id NOT in (SELECT id FROM parents)
				), nieces_nephews AS (
					SELECT p.id, p.name, 'nieces_nephews' AS relationship FROM person p
					JOIN relationships r ON p.id = r.related_person_id
					WHERE r.person_id IN (SELECT id FROM siblings) AND r.relationship = 'children' AND p.id NOT in (SELECT id FROM parents)
				), grandchild AS (
					SELECT p.id, p.name, 'grandchild' AS relationship FROM person p
					JOIN relationships r ON p.id = r.related_person_id
					WHERE r.person_id IN (SELECT id FROM parents where relationship = 'children') AND r.relationship = 'children' AND p.id NOT in (SELECT id FROM parents)
				 
				)
				SELECT distinct * FROM parents 
				UNION
				SELECT * FROM grandparents
				UNION
				SELECT * FROM aunts_uncles
				UNION
				SELECT * FROM cousin
				UNION
				SELECT * FROM siblings
				UNION
				SELECT * FROM nieces_nephews
				UNION
				SELECT * FROM grandchild`
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

		parent.FamilyConnection = domain.FamilyConnection[parent.Relationship]
		if parent.ID == personId {
			familyMembers.ID = parent.ID
			familyMembers.Name = parent.Name
			continue
		}

		relationships = append(relationships, parent)
	}

	if len(relationships) == 0 {
		return nil, domain.ErrRelationNotFound
	}

	familyMembers.Members = relationships

	return familyMembers, nil

}

func (p *postgresPersonRepo) CreateRelationship(personId int64, relationship domain.Relationship) error {

	if personId == relationship.ChildrenId || personId == relationship.ParentId {
		return domain.ErrInvalidSelfRelation
	}

	prepareQuery, err := p.DB.Prepare("INSERT INTO relationships (person_id, related_person_id, relationship) VALUES($1, $2, $3)")
	if err != nil {
		return err
	}
	defer prepareQuery.Close()

	if relationship.ChildrenId != 0 {
		family, err := p.GetRelationshipByID(relationship.ChildrenId)
		if err != nil && err != domain.ErrRelationNotFound {
			return err
		}

		if family != nil {
			for _, familyMember := range family.Members {
				if familyMember.ID == personId && familyMember.Relationship == "parent" {
					return domain.ErrDuplicateRelation
				}
			}
		}

		_, err = prepareQuery.Exec(personId, relationship.ChildrenId, "children")
		if err != nil {
			return err
		}
	}

	if relationship.ParentId != 0 {
		family, err := p.GetRelationshipByID(relationship.ParentId)
		if err != nil && err != domain.ErrRelationNotFound {
			return err
		}

		if len(family.Members) > 0 {
			for _, familyMember := range family.Members {
				if familyMember.ID == personId && familyMember.Relationship == "children" {
					return domain.ErrDuplicateRelation
				}
			}
		}

		_, err = prepareQuery.Exec(personId, relationship.ParentId, "parent")
		if err != nil {
			return err
		}
	}
	return nil
}
