package repository

import (
	"database/sql"

	"github.com/LayssonENS/go-genealogy-api/domain"
)

type postgresRelationshipRepo struct {
	DB *sql.DB
}

// NewPostgresRelationshipRepository will create an implementation of relationship.Repository
func NewPostgresRelationshipRepository(db *sql.DB) domain.RelationshipRepository {
	return &postgresRelationshipRepo{
		DB: db,
	}
}

// GetRelationshipByID : Retrieves all relationships of a person from the Postgres repository
func (p *postgresRelationshipRepo) GetRelationshipByID(personId int64) (*domain.Member, error) {
	var relationships []domain.Family
	familyMembers := &domain.Member{}
	query := `WITH parents AS (
					SELECT p.id, p.name, 'parent' AS relationship
					FROM person p
				    JOIN relationships r ON p.id = r.person_id  
					WHERE r.related_person_id = $1
			   ), self_search AS (
					SELECT p.id, p.name, 'self' AS relationship
					FROM person p
					WHERE p.id = $1
			   ), childrens AS (
					SELECT p.id, p.name, 'children' AS relationship
					FROM person p
				    JOIN relationships r ON p.id = r.related_person_id 
					WHERE r.person_id = $1
				), grandparents AS (
					SELECT p.id, p.name, 'grandparent' AS relationship
					FROM person p
				    JOIN relationships r ON p.id = r.person_id 
					WHERE r.related_person_id IN (SELECT id from parents)
				), grandchild AS (
				    SELECT p.id, p.name, 'grandchild' AS relationship
					FROM person p
				    JOIN relationships r ON p.id = r.related_person_id 
					WHERE r.person_id IN (SELECT id from childrens)
				), aunts_uncles AS (
				    SELECT p.id, p.name, 'aunt_uncle' AS relationship
					FROM person p
				    JOIN relationships r ON p.id = r.related_person_id 
					WHERE r.person_id IN (SELECT id from grandparents)
					AND p.id NOT in (SELECT id FROM parents)
				), cousins AS (
				    SELECT p.id, p.name, 'cousin' AS relationship
					FROM person p
				    JOIN relationships r ON p.id = r.related_person_id 
					WHERE r.person_id IN (SELECT id from aunts_uncles)
				), siblings AS (
				    SELECT p.id, p.name, 'sibling' AS relationship
					FROM person p
				    JOIN relationships r ON p.id = r.related_person_id 
					WHERE r.person_id IN (SELECT id from parents) 
					AND p.id <> $1
				), nieces_nephews AS (
				    SELECT p.id, p.name, 'niece_nephew' AS relationship
					FROM person p
				    JOIN relationships r ON p.id = r.related_person_id 
					WHERE r.person_id IN (SELECT id from siblings)
					AND p.id NOT in (SELECT id FROM childrens) 
				)
				SELECT * FROM parents 
				union
				SELECT * FROM self_search
				union
				SELECT * FROM childrens 
				UNION
				SELECT * FROM grandparents
				UNION
				SELECT * FROM grandchild
				UNION
				SELECT * FROM aunts_uncles
				UNION
				SELECT * FROM cousins
				UNION
				SELECT * FROM siblings
				UNION
				SELECT * FROM nieces_nephews`
	rows, err := p.DB.Query(query, personId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var parent domain.Family
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

		parents, err := p.getRelationships(parent.ID)
		if err != nil {
			return nil, err
		}
		parent.Relationships = parents

		relationships = append(relationships, parent)
	}

	if len(relationships) == 0 {
		return nil, domain.ErrRelationNotFound
	}

	familyMembers.Members = relationships

	return familyMembers, nil

}

// CreateRelationship : creates a new parent-child relationship between two persons in the Postgres repository
func (p *postgresRelationshipRepo) CreateRelationship(personId int64, relationship domain.Relationship) error {

	if personId == relationship.ChildrenId || personId == relationship.ParentId {
		return domain.ErrInvalidSelfRelation
	}

	query := `INSERT INTO relationships (person_id, related_person_id, relationship) VALUES ($1, $2, $3) `

	if relationship.ChildrenId != 0 {
		family, err := p.GetRelationshipByID(relationship.ChildrenId)
		if err != nil && err != domain.ErrRelationNotFound {
			return err
		}

		if family != nil {
			for _, familyMember := range family.Members {
				if familyMember.ID == personId && familyMember.Relationship == domain.ParentName {
					return domain.ErrDuplicateRelation
				} else if familyMember.ID == personId {
					return domain.ErrIncestuousRelation
				}
			}
		}

		_, err = p.DB.Exec(query, personId, relationship.ChildrenId, domain.ChildrenName)
		if err != nil {
			return err
		}
	}

	if relationship.ParentId != 0 {
		family, err := p.GetRelationshipByID(relationship.ParentId)
		if err != nil && err != domain.ErrRelationNotFound {
			return err
		}

		if family != nil {
			for _, familyMember := range family.Members {
				if familyMember.ID == personId && familyMember.Relationship == domain.ChildrenName {
					return domain.ErrDuplicateRelation
				} else if familyMember.ID == personId {
					return domain.ErrIncestuousRelation
				}
			}
		}

		_, err = p.DB.Exec(query, relationship.ParentId, personId, "children")
		if err != nil {
			return err
		}
	}
	return nil
}

// getRelationships : Retrieves relationships (parents and children) of a person from the Postgres repository
func (p *postgresRelationshipRepo) getRelationships(personId int64) ([]domain.Relation, error) {
	var parents []domain.Relation
	query := `WITH parents AS (
					SELECT p.id, p.name, 'parent' AS relationship
					FROM person p
				    JOIN relationships r ON p.id = r.person_id  
					WHERE r.related_person_id = $1
			   ), childrens AS (
					SELECT p.id, p.name, 'children' AS relationship
					FROM person p
				    JOIN relationships r ON p.id = r.related_person_id 
					WHERE r.person_id = $1
				)
				SELECT * FROM parents 
				UNION
				SELECT * FROM childrens`
	rows, err := p.DB.Query(query, personId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var parent domain.Relation
		err := rows.Scan(&parent.ID, &parent.Name, &parent.Relationship)
		if err != nil {
			return nil, err
		}
		parents = append(parents, parent)
	}

	return parents, nil

}

// DeleteRelationship : Deletes a relationship between two people in the Postgres repository
func (p *postgresRelationshipRepo) DeleteRelationship(personId int64) error {
	query := `DELETE FROM relationships
			  WHERE person_id = $1 OR related_person_id = $1`

	_, err := p.DB.Exec(query, personId)
	if err != nil {
		return err
	}

	return nil
}
