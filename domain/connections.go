package domain

var FamilyConnection = map[string]int64{
	"self":         0,
	"parent":       1,
	"children":     1,
	"sibling":      2,
	"grandparent":  2,
	"grandchild":   2,
	"aunt_uncle":   3,
	"niece_nephew": 3,
	"cousin":       4,
}
