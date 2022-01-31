package tosca_v2_0

import (
	"github.com/tliron/kutil/ard"
	"github.com/tliron/puccini/tosca"
	"github.com/tliron/puccini/tosca/normal"
)

//
// NotificationAssignment
//
// [TOSCA-v2.0] @ ?
// [TOSCA-Simple-Profile-YAML-v1.3] @ 3.6.19
//

type NotificationAssignment struct {
	*Entity `name:"notification"`
	Name    string

	Description    *string                  `read:"description"`
	Implementation *InterfaceImplementation `read:"implementation,InterfaceImplementation"`
	Outputs        OutputMappings           `read:"outputs,OutputMapping"`
}

func NewNotificationAssignment(context *tosca.Context) *NotificationAssignment {
	return &NotificationAssignment{
		Entity:  NewEntity(context),
		Name:    context.Name,
		Outputs: make(OutputMappings),
	}
}

// tosca.Reader signature
func ReadNotificationAssignment(context *tosca.Context) tosca.EntityPtr {
	self := NewNotificationAssignment(context)

	if context.Is(ard.TypeMap) {
		// Long notation
		context.ValidateUnsupportedFields(context.ReadFields(self))
	} else if context.ValidateType(ard.TypeMap, ard.TypeString) {
		// Short notation
		self.Implementation = ReadInterfaceImplementation(context.FieldChild("implementation", context.Data)).(*InterfaceImplementation)
	}

	return self
}

// tosca.Mappable interface
func (self *NotificationAssignment) GetKey() string {
	return self.Name
}

func (self *NotificationAssignment) Normalize(normalInterface *normal.Interface) *normal.Notification {
	logNormalize.Debugf("notification: %s", self.Name)

	normalNotification := normalInterface.NewNotification(self.Name)

	if self.Description != nil {
		normalNotification.Description = *self.Description
	}

	if self.Implementation != nil {
		self.Implementation.NormalizeNotification(normalNotification)
	}

	if normalInterface.NodeTemplate != nil {
		self.Outputs.NormalizeForNodeTemplate(normalInterface.NodeTemplate.ServiceTemplate, normalNotification.Outputs)
	} else if normalInterface.Relationship != nil {
		self.Outputs.NormalizeForRelationship(normalInterface.Relationship, normalNotification.Outputs)
	} else if normalInterface.Group != nil {
		self.Outputs.NormalizeForGroup(normalInterface.Group.ServiceTemplate, normalNotification.Outputs)
	}

	return normalNotification
}

//
// NotificationAssignments
//

type NotificationAssignments map[string]*NotificationAssignment

func (self NotificationAssignments) CopyUnassigned(assignments NotificationAssignments) {
	for key, assignment := range assignments {
		lock1 := assignment.GetEntityLock()
		lock1.Lock()
		if selfAssignment, ok := self[key]; ok {
			lock2 := selfAssignment.GetEntityLock()
			lock2.Lock()
			selfAssignment.Outputs.CopyUnassigned(assignment.Outputs)
			lock2.Unlock()
		} else {
			self[key] = assignment
		}
		lock1.Unlock()
	}
}

func (self NotificationAssignments) RenderForNodeTemplate(nodeTemplate *NodeTemplate, definitions NotificationDefinitions, context *tosca.Context) {
	self.render(definitions, context)
	for _, assignment := range self {
		lock := assignment.GetEntityLock()
		lock.Lock()
		assignment.Outputs.RenderForNodeTemplate(nodeTemplate)
		lock.Unlock()
	}
}

func (self NotificationAssignments) RenderForRelationship(relationship *RelationshipAssignment, definitions NotificationDefinitions, context *tosca.Context) {
	self.render(definitions, context)
	for _, assignment := range self {
		lock := assignment.GetEntityLock()
		lock.Lock()
		assignment.Outputs.RenderForRelationship(relationship)
		lock.Unlock()
	}
}

func (self NotificationAssignments) RenderForGroup(definitions NotificationDefinitions, context *tosca.Context) {
	self.render(definitions, context)
	for _, assignment := range self {
		lock := assignment.GetEntityLock()
		lock.Lock()
		assignment.Outputs.RenderForGroup()
		lock.Unlock()
	}
}

func (self NotificationAssignments) render(definitions NotificationDefinitions, context *tosca.Context) {
	for key, definition := range definitions {
		lock1 := definition.GetEntityLock()
		lock1.RLock()

		assignment, ok := self[key]

		if !ok {
			assignment = NewNotificationAssignment(context.FieldChild(key, nil))
			self[key] = assignment
		}

		lock2 := assignment.GetEntityLock()
		lock2.Lock()

		if assignment.Description == nil {
			assignment.Description = definition.Description
		}

		if (assignment.Implementation == nil) && (definition.Implementation != nil) {
			// If the definition has an implementation then we must have one, too
			assignment.Implementation = NewInterfaceImplementation(assignment.Context.FieldChild("implementation", nil))
		}

		if assignment.Implementation != nil {
			lock3 := assignment.Implementation.GetEntityLock()
			lock3.Lock()
			assignment.Implementation.Render(definition.Implementation)
			lock3.Unlock()
		}

		assignment.Outputs.Inherit(definition.Outputs)

		lock2.Unlock()
		lock1.RUnlock()
	}

	for key, assignment := range self {
		if _, ok := definitions[key]; !ok {
			assignment.Context.ReportUndeclared("notification")
			delete(self, key)
		}
	}
}

func (self NotificationAssignments) Normalize(normalInterface *normal.Interface) {
	for key, notification := range self {
		lock := notification.GetEntityLock()
		lock.RLock()
		normalInterface.Notifications[key] = notification.Normalize(normalInterface)
		lock.RUnlock()
	}
}
