// This file was auto-generated from a YAML file

package v1_2

func init() {
	Profile["/tosca/simple/1.2/profile.yaml"] = `
# Modified from a file that was distributed with this NOTICE:
#
#   Apache AriaTosca
#   Copyright 2016-2017 The Apache Software Foundation
#
#   This product includes software developed at
#   The Apache Software Foundation (http://www.apache.org/).

tosca_definitions_version: tosca_simple_yaml_1_2

metadata:
  puccini-js.import.tosca.resolve: js/resolve.js
  puccini-js.import.tosca.coerce: js/coerce.js
  puccini-js.import.tosca.visualize: js/visualize.js
  puccini-js.import.tosca.utils: js/utils.js
  puccini-js.import.tosca.helpers: js/helpers.js

imports:
- artifacts.yaml
- groups.yaml
- nodes.yaml
- policies.yaml
`
}
