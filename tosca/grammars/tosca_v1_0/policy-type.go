package tosca_v1_0

import (
	"github.com/tliron/go-puccini/tosca/grammars/tosca_v2_0"
	"github.com/tliron/go-puccini/tosca/parsing"
)

//
// PolicyType
//
// [TOSCA-Simple-Profile-YAML-v1.0] @ 3.6.11
//

// ([parsing.Reader] signature)
func ReadPolicyType(context *parsing.Context) parsing.EntityPtr {
	context.SetReadTag("TriggerDefinitions", "")

	return tosca_v2_0.ReadPolicyType(context)
}
