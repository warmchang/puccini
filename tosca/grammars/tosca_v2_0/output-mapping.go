package tosca_v2_0

import (
	"reflect"

	"github.com/tliron/puccini/tosca"
	"github.com/tliron/puccini/tosca/normal"
)

//
// OutputMapping
//
// Attaches to notifications and operations
//

type OutputMapping struct {
	*Entity `name:"output mapping"`
	Name    string

	EntityName    *string `require:"0"`
	AttributeName *string `require:"1"`

	NodeTemplate *NodeTemplate           `traverse:"ignore" json:"-" yaml:"-"`
	Relationship *RelationshipAssignment `traverse:"ignore" json:"-" yaml:"-"`
}

func NewOutputMapping(context *tosca.Context) *OutputMapping {
	return &OutputMapping{
		Entity: NewEntity(context),
		Name:   context.Name,
	}
}

// tosca.Reader signature
func ReadOutputMapping(context *tosca.Context) tosca.EntityPtr {
	self := NewOutputMapping(context)

	if strings := context.ReadStringListFixed(2); strings != nil {
		self.EntityName = &(*strings)[0]
		self.AttributeName = &(*strings)[1]
	}

	return self
}

// tosca.Mappable interface
func (self *OutputMapping) GetKey() string {
	return self.Name
}

func (self *OutputMapping) RenderForNodeTemplate(nodeTemplate *NodeTemplate) {
	logRender.Debugf("output mapping: %s", self.Name)

	if (self.EntityName == nil) || (self.AttributeName == nil) {
		return
	}

	entityName := *self.EntityName
	if entityName == "SELF" {
		self.NodeTemplate = nodeTemplate
	} else {
		var nodeTemplateType *NodeTemplate
		if nodeTemplate, ok := self.Context.Namespace.LookupForType(entityName, reflect.TypeOf(nodeTemplateType)); ok {
			self.NodeTemplate = nodeTemplate.(*NodeTemplate)
		} else {
			self.Context.ListChild(0, entityName).ReportUnknown("node template")
		}
	}

	// Attributes should already have been rendered
	attributeName := *self.AttributeName
	if _, ok := self.NodeTemplate.Attributes[attributeName]; !ok {
		self.Context.ListChild(1, attributeName).ReportReferenceNotFound("attribute", self.NodeTemplate)
	}
}

func (self *OutputMapping) RenderForRelationship(relationship *RelationshipAssignment) {
	logRender.Debugf("output mapping: %s", self.Name)

	if (self.EntityName == nil) || (self.AttributeName == nil) {
		return
	}

	entityName := *self.EntityName
	if entityName == "SELF" {
		self.Relationship = relationship

		// Attributes should already have been rendered
		attributeName := *self.AttributeName
		if _, ok := self.Relationship.Attributes[attributeName]; !ok {
			self.Context.ListChild(1, attributeName).ReportReferenceNotFound("attribute", self.NodeTemplate)
		}
	} else {
		var nodeTemplateType *NodeTemplate
		if nodeTemplate, ok := self.Context.Namespace.LookupForType(entityName, reflect.TypeOf(nodeTemplateType)); ok {
			self.NodeTemplate = nodeTemplate.(*NodeTemplate)

			// Attributes should already have been rendered
			attributeName := *self.AttributeName
			if _, ok := self.NodeTemplate.Attributes[attributeName]; !ok {
				self.Context.ListChild(1, attributeName).ReportReferenceNotFound("attribute", self.NodeTemplate)
			}
		} else {
			self.Context.ListChild(0, entityName).ReportUnknown("node template")
		}
	}
}

func (self *OutputMapping) RenderForGroup() {
	logRender.Debugf("output mapping: %s", self.Name)

	if (self.EntityName == nil) || (self.AttributeName == nil) {
		return
	}

	entityName := *self.EntityName
	if entityName == "SELF" {
		self.Context.ListChild(0, entityName).ReportValueInvalid("modelable entity name", "cannot be used in groups")
	} else {
		var nodeTemplateType *NodeTemplate
		if nodeTemplate, ok := self.Context.Namespace.LookupForType(entityName, reflect.TypeOf(nodeTemplateType)); ok {
			self.NodeTemplate = nodeTemplate.(*NodeTemplate)

			// Attributes should already have been rendered
			attributeName := *self.AttributeName
			if _, ok := self.NodeTemplate.Attributes[attributeName]; !ok {
				self.Context.ListChild(1, attributeName).ReportReferenceNotFound("attribute", self.NodeTemplate)
			}
		} else {
			self.Context.ListChild(0, entityName).ReportUnknown("node template")
		}
	}
}

func (self *OutputMapping) NormalizeForNodeTemplate(normalServiceTemplate *normal.ServiceTemplate, normalOutputs normal.Outputs) {
	if (self.NodeTemplate == nil) || (self.AttributeName == nil) {
		return
	}

	if normalTargetNodeTemplate, ok := normalServiceTemplate.NodeTemplates[self.NodeTemplate.Name]; ok {
		normalOutputs[self.Name] = normalTargetNodeTemplate.NewMapping(*self.AttributeName)
	}
}

func (self *OutputMapping) NormalizeForRelationship(normalRelationship *normal.Relationship, normalOutputs normal.Outputs) {
	if (self.EntityName != nil) && (*self.EntityName == "SELF") {
		// TODO
		//normalOutputs[self.Name] = normalRelationship.NewMapping(*self.AttributeName)
	} else {
		self.NormalizeForNodeTemplate(normalRelationship.Requirement.NodeTemplate.ServiceTemplate, normalOutputs)
	}
}

//
// OutputMappings
//

type OutputMappings map[string]*OutputMapping

func (self OutputMappings) CopyUnassigned(outputMappings OutputMappings) {
	for key, outputMapping := range outputMappings {
		if _, ok := self[key]; !ok {
			self[key] = outputMapping
		}
	}
}

func (self OutputMappings) Inherit(parent OutputMappings) {
	for name, outputMapping := range parent {
		if _, ok := self[name]; !ok {
			self[name] = outputMapping
		}
	}
}

func (self OutputMappings) RenderForNodeTemplate(nodeTemplate *NodeTemplate) {
	for _, outputMapping := range self {
		outputMapping.RenderForNodeTemplate(nodeTemplate)
	}
}

func (self OutputMappings) RenderForRelationship(relationship *RelationshipAssignment) {
	for _, outputMapping := range self {
		outputMapping.RenderForRelationship(relationship)
	}
}

func (self OutputMappings) RenderForGroup() {
	for _, outputMapping := range self {
		outputMapping.RenderForGroup()
	}
}

func (self OutputMappings) NormalizeForNodeTemplate(normalServiceTemplate *normal.ServiceTemplate, normalOutputs normal.Outputs) {
	for _, outputMapping := range self {
		outputMapping.NormalizeForNodeTemplate(normalServiceTemplate, normalOutputs)
	}
}

func (self OutputMappings) NormalizeForRelationship(normalRelationship *normal.Relationship, normalOutputs normal.Outputs) {
	for _, outputMapping := range self {
		outputMapping.NormalizeForRelationship(normalRelationship, normalOutputs)
	}
}
