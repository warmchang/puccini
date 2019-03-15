// This file was auto-generated from a YAML file

package v4_5

func init() {
	Profile["/cloudify/4.5/relationships.yaml"] = `
tosca_definitions_version: cloudify_dsl_1_3

# https://docs.cloudify.co/4.5.5/developer/blueprints/spec-relationships/
# https://docs.cloudify.co/4.5.5/developer/blueprints/built-in-types/

relationships:

  cloudify.relationships.depends_on: {}

  cloudify.relationships.contained_in:
    derived_from: cloudify.relationships.depends_on

  cloudify.relationships.connected_to:
    derived_from: cloudify.relationships.depends_on

  cloudify.relationships.file_system_depends_on_volume: {}

  cloudify.relationships.file_system_contained_in_compute: {}
`
}
