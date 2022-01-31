package tosca_v2_0

import (
	"fmt"

	"github.com/tliron/kutil/ard"
	"github.com/tliron/kutil/util"
	"github.com/tliron/puccini/tosca"
	"github.com/tliron/puccini/tosca/normal"
)

//
// RequirementAssignment
//
// [TOSCA-v2.0] @ ?
// [TOSCA-Simple-Profile-YAML-v1.3] @ 3.8.2
// [TOSCA-Simple-Profile-YAML-v1.2] @ 3.8.2
// [TOSCA-Simple-Profile-YAML-v1.1] @ 3.7.2
// [TOSCA-Simple-Profile-YAML-v1.0] @ 3.7.2
//

type RequirementAssignment struct {
	*Entity `name:"requirement"`
	Name    string

	TargetCapabilityNameOrTypeName   *string                 `read:"capability"`
	TargetNodeTemplateNameOrTypeName *string                 `read:"node"`
	TargetNodeFilter                 *NodeFilter             `read:"node_filter,NodeFilter"`
	Relationship                     *RelationshipAssignment `read:"relationship,RelationshipAssignment"`
	Occurrences                      *RangeEntity            `read:"occurrences,RangeEntity"` // introduced in TOSCA 1.3

	TargetCapabilityType *CapabilityType `lookup:"capability,?TargetCapabilityNameOrTypeName" json:"-" yaml:"-"`
	TargetNodeTemplate   *NodeTemplate   `lookup:"node,TargetNodeTemplateNameOrTypeName" json:"-" yaml:"-"`
	TargetNodeType       *NodeType       `lookup:"node,TargetNodeTemplateNameOrTypeName" json:"-" yaml:"-"`
}

func NewRequirementAssignment(context *tosca.Context) *RequirementAssignment {
	return &RequirementAssignment{
		Entity: NewEntity(context),
		Name:   context.Name,
	}
}

// tosca.Reader signature
func ReadRequirementAssignment(context *tosca.Context) tosca.EntityPtr {
	self := NewRequirementAssignment(context)

	if context.Is(ard.TypeMap) {
		// Long notation
		context.ValidateUnsupportedFields(context.ReadFields(self))
	} else if context.ValidateType(ard.TypeMap, ard.TypeString) {
		// Short notation
		self.TargetNodeTemplateNameOrTypeName = context.FieldChild("node", context.Data).ReadString()
	}

	return self
}

func NewDefaultRequirementAssignment(index int, definition *RequirementDefinition, context *tosca.Context) *RequirementAssignment {
	context = context.SequencedListChild(index, definition.Name, nil)
	context.Name = definition.Name
	self := NewRequirementAssignment(context)
	self.TargetNodeTemplateNameOrTypeName = definition.TargetNodeTypeName
	self.TargetNodeType = definition.TargetNodeType
	self.TargetCapabilityNameOrTypeName = definition.TargetCapabilityTypeName
	self.TargetCapabilityType = definition.TargetCapabilityType
	return self
}

func (self *RequirementAssignment) GetDefinition(nodeTemplate *NodeTemplate) (*RequirementDefinition, bool) {
	if nodeTemplate.NodeType == nil {
		return nil, false
	}
	definition, ok := nodeTemplate.NodeType.RequirementDefinitions[self.Name]
	return definition, ok
}

func (self *RequirementAssignment) Normalize(nodeTemplate *NodeTemplate, normalNodeTemplate *normal.NodeTemplate) *normal.Requirement {
	normalRequirement := normalNodeTemplate.NewRequirement(self.Name, normal.NewLocationForContext(self.Context))

	if self.TargetCapabilityType != nil {
		lock := self.TargetCapabilityType.GetEntityLock()
		lock.RLock()
		name := tosca.GetCanonicalName(self.TargetCapabilityType)
		lock.RUnlock()
		normalRequirement.CapabilityTypeName = &name
	} else if self.TargetCapabilityNameOrTypeName != nil {
		normalRequirement.CapabilityName = self.TargetCapabilityNameOrTypeName
	}

	if self.TargetNodeTemplate != nil {
		lock := self.TargetNodeTemplate.GetEntityLock()
		lock.RLock()
		normalRequirement.NodeTemplate, _ = normalNodeTemplate.ServiceTemplate.NodeTemplates[self.TargetNodeTemplate.Name]
		lock.RUnlock()
	}

	if self.TargetNodeType != nil {
		lock := self.TargetNodeType.GetEntityLock()
		lock.RLock()
		name := tosca.GetCanonicalName(self.TargetNodeType)
		lock.RUnlock()
		normalRequirement.NodeTypeName = &name
	}

	if nodeTemplate.RequirementTargetsNodeFilter != nil {
		lock := nodeTemplate.RequirementTargetsNodeFilter.GetEntityLock()
		lock.RLock()
		nodeTemplate.RequirementTargetsNodeFilter.Normalize(normalRequirement)
		lock.RUnlock()
	}

	if self.TargetNodeFilter != nil {
		lock := self.TargetNodeFilter.GetEntityLock()
		lock.RLock()
		self.TargetNodeFilter.Normalize(normalRequirement)
		lock.RUnlock()
	}

	if self.Relationship != nil {
		lock1 := self.Relationship.GetEntityLock()
		lock1.RLock()
		if definition, ok := self.GetDefinition(nodeTemplate); ok {
			lock2 := definition.GetEntityLock()
			lock2.RLock()
			var lock3 util.RWLocker
			if definition.RelationshipDefinition != nil {
				lock3 = definition.RelationshipDefinition.GetEntityLock()
				lock3.RLock()
			}
			self.Relationship.Normalize(definition.RelationshipDefinition, normalRequirement.NewRelationship())
			if lock3 != nil {
				lock3.RUnlock()
			}
			lock2.RUnlock()
		} else {
			self.Relationship.Normalize(nil, normalRequirement.NewRelationship())
		}
		lock1.RUnlock()
	}

	return normalRequirement
}

//
// RequirementAssignments
//

type RequirementAssignments []*RequirementAssignment

func (self *RequirementAssignments) Render(definitions RequirementDefinitions, context *tosca.Context) {
	// TODO: currently have no idea what to do with "occurrences" keyword in the requirement
	// assignment, because we interpret "occurrences" in the definition to mean how many times
	// it would be assigned

	for _, definition := range definitions {
		lock := definition.GetEntityLock()
		lock.RLock()

		// The TOSCA spec says that definition occurrences has an "implied default of [1,1]"
		occurrences := definition.Occurrences
		if occurrences == nil {
			occurrencesContext := definition.Context.FieldChild("occurrences", ard.List{1, 1})
			occurrences = ReadRangeEntity(occurrencesContext).(*RangeEntity)
		}

		count := self.Count(definition.Name)

		// Automatically add missing assignments
		for index := count; index < occurrences.Range.Lower; index++ {
			*self = append(*self, NewDefaultRequirementAssignment(len(*self), definition, context))
			count++
		}

		if !occurrences.Range.InRange(count) {
			context.ReportNotInRange(fmt.Sprintf("number of requirement %q assignments", definition.Name), count, occurrences.Range.Lower, occurrences.Range.Upper)
		}

		lock.RUnlock()
	}

	for _, assignment := range *self {
		lock1 := assignment.GetEntityLock()
		lock1.Lock()

		if definition, ok := definitions[assignment.Name]; ok {
			lock2 := definition.GetEntityLock()
			lock2.RLock()

			if assignment.TargetCapabilityNameOrTypeName == nil {
				assignment.TargetCapabilityNameOrTypeName = definition.TargetCapabilityTypeName
			}

			if assignment.TargetCapabilityType == nil {
				assignment.TargetCapabilityType = definition.TargetCapabilityType
			}

			if assignment.TargetNodeTemplateNameOrTypeName == nil {
				assignment.TargetNodeTemplateNameOrTypeName = definition.TargetNodeTypeName
			}

			if assignment.TargetNodeType == nil {
				assignment.TargetNodeType = definition.TargetNodeType
			}

			if definition.RelationshipDefinition != nil {
				lock3 := definition.RelationshipDefinition.GetEntityLock()
				lock3.RLock()

				if assignment.Relationship == nil {
					assignment.Relationship = definition.RelationshipDefinition.NewDefaultAssignment(assignment.Context.FieldChild("relationship", nil))
				}

				lock4 := assignment.Relationship.GetEntityLock()
				lock4.Lock()

				if assignment.Relationship.RelationshipTemplateNameOrTypeName == nil {
					// Note: the definition can only specify a relationship type, not a relationship template
					assignment.Relationship.RelationshipTemplateNameOrTypeName = definition.RelationshipDefinition.RelationshipTypeName
				}

				if (assignment.Relationship.RelationshipType == nil) && (assignment.Relationship.RelationshipTemplate == nil) {
					// Note: we are careful not set the relationship type if the assignment uses a relationship template
					assignment.Relationship.RelationshipType = definition.RelationshipDefinition.RelationshipType
				}

				assignment.Relationship.Render(definition.RelationshipDefinition)

				lock4.Unlock()
				lock3.RUnlock()
			}

			lock2.RUnlock()
		} else {
			assignment.Context.ReportUndeclared("requirement")
		}

		lock1.Unlock()
	}
}

func (self RequirementAssignments) Normalize(nodeTemplate *NodeTemplate, normalNodeTemplate *normal.NodeTemplate) {
	for _, requirement := range self {
		lock := requirement.GetEntityLock()
		lock.RLock()
		requirement.Normalize(nodeTemplate, normalNodeTemplate)
		lock.RUnlock()
	}
}

func (self *RequirementAssignments) Count(name string) uint64 {
	var count uint64
	for _, assignment := range *self {
		if assignment.Name == name {
			count++
		}
	}
	return count
}
