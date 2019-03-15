// This file was auto-generated from a YAML file

package v4_5

func init() {
	Profile["/cloudify/4.5/nodes.yaml"] = `
tosca_definitions_version: cloudify_dsl_1_3

# https://docs.cloudify.co/4.5.5/developer/blueprints/built-in-types/

imports:
- data.yaml

node_types:

  cloudify.nodes.Root:
    interfaces:
      cloudify.interfaces.lifecycle: {}
      cloudify.interfaces.validation: {}
      cloudify.interfaces.monitoring_agent: {}
      cloudify.interfaces.monitoring: {}
  cloudify.nodes.Tier: {}
  cloudify.nodes.Compute: {}
  cloudify.nodes.Container: {}
  cloudify.nodes.Network: {}
  cloudify.nodes.Subnet: {}
  cloudify.nodes.Router: {}
  cloudify.nodes.Port: {}
  cloudify.nodes.VirtualIP: {}
  cloudify.nodes.SecurityGroup: {}
  cloudify.nodes.LoadBalancer: {}
  cloudify.nodes.Volume: {}
  cloudify.nodes.FileSystem:
    properties:
      use_external_resource:
        type: boolean
        default: false
      partition_type:
        type: integer
        default: 83
      fs_type:
        type: string
      fs_mount_path:
        type: string
  cloudify.nodes.ObjectStorage: {}
  cloudify.nodes.SoftwareComponent: {}
  cloudify.nodes.WebServer:
    properties:
      port:
        type: integer
  cloudify.nodes.ApplicationServer: {}
  cloudify.nodes.DBMS: {}
  cloudify.nodes.MessageBugServer: {}
  cloudify.nodes.ApplicationModule: {}
`
}
