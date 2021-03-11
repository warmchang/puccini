# https://docs.ansible.com/ansible/latest/dev_guide/developing_inventory.html
# https://github.com/ansible/ansible/blob/devel/lib/ansible/plugins/inventory/yaml.py

from ansible.errors import AnsibleError, AnsibleParserError
from ansible.module_utils.common._collections_compat import MutableMapping
from ansible.module_utils._text import to_text
from ansible.plugins.inventory import BaseInventoryPlugin
import puccini.tosca

class InventoryModule(BaseInventoryPlugin):

    NAME = 'hosts'

    def __init__(self):
        super(InventoryModule, self).__init__()

    def verify_file(self, path):
        if super(InventoryModule, self).verify_file(path):
            if path.endswith(('tosca.yaml', 'tosca.yml')):
                return True
        return False

    def parse(self, inventory, loader, path, cache=True):
        super(InventoryModule, self).parse(inventory, loader, path, cache)

        try:
            data = self.loader.load_from_file(path, cache=False)
        except Exception as e:
            raise AnsibleParserError(e)

        if not data:
            raise AnsibleParserError('Parsed empty YAML file')
        elif not isinstance(data, MutableMapping):
            raise AnsibleParserError('YAML inventory has invalid structure, it should be a dictionary, got: %s' % type(data))

        try:
            tosca_group = self.inventory.add_group('tosca')
        except AnsibleError as e:
            raise AnsibleParserError("Unable to add group %s: %s" % ('tosca', to_text(e)))

        for service in data.get('services', []):
            group = service.get('group')
            if (not group) or (group == 'tosca'):
                service_group = tosca_group
            else:
                try:
                    service_group = self.inventory.add_group(group)
                    self.inventory.add_child(tosca_group, service_group)
                except AnsibleError as e:
                    raise AnsibleParserError("Unable to add group %s: %s" % (group, to_text(e)))

            template = service.get('template')
            try:
                clout = puccini.tosca.compile(template)
            except Exception as e:
                raise AnsibleError('TOSCA compilation error: %s' % to_text(e))

            for vertex in clout['vertexes'].values():
                try:
                    if vertex['metadata']['puccini']['kind'] == 'NodeTemplate':
                        node_template = vertex['properties']
                        self.inventory.add_host(node_template['name'], group=service_group)
                except:
                    pass
   
