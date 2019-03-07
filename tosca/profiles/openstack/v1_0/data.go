// This file was auto-generated from a YAML file

package v1_0

func init() {
	Profile["/tosca/openstack/1.0/data.yaml"] = `
tosca_definitions_version: tosca_simple_yaml_1_2

data_types:

  openstack.IpAddress:
    derived_from:
      string

  openstack.MacAddress:
    derived_from:
      string

  openstack.NetCidr:
    derived_from:
      string

  openstack.Neutron.Network:
    derived_from:
      string

  openstack.Neutron.Port:
    derived_from:
      string

  openstack.Neutron.QosPolicy:
    derived_from:
      string

  openstack.Glance.Image:
    derived_from:
      string

  openstack.Cinder.Snapshot:
    derived_from:
      string

  openstack.Cinder.Volume:
    derived_from:
      string

  openstack.nova.Flavor:
    derived_from:
      string

  openstack.nova.Keypair:
    derived_from:
      string

  openstack.nova.Server.Network:
    derived_from: Root
    properties:
      allocate_network:
        description: >-
          The special string values of network, auto: means either a network that is already
          available to the project will be used, or if one does not exist, will be automatically
          created for the project; none: means no networking will be allocated for the created
          server. Supported by Nova API since version "2.37". This property can not be used with
          other network keys.
        type: string
        constraints:
        - valid_values: [ none, auto ]
      fixed_ip:
        description: >-
          Fixed IP address to specify for the port created on the requested network.
        type: openstack.IpAddress
      network:
        description: >-
          Name or ID of network to create a port on.
        type: openstack.Neutron.Network
      floating_ip:
        description: >-
          ID of the floating IP to associate.
        type: string
      port:
        description: >-
          ID of an existing port to associate with this server.
        type: openstack.Neutron.Port
      port_extra_properties:
        description: >-
          Dict, which has expand properties for port. Used only if port property is not specified
          for creating port.
        type: map
        entry_schema: openstack.nova.Server.Port
      subnet:
        description: >-
          Subnet in which to allocate the IP address for port. Used for creating port, based on
          derived properties. If subnet is specified, network property becomes optional.
        type: string
      tag:
        description: >-
          Port tag. Heat ignores any update on this property as nova does not support it.
        type: string

  openstack.nova.Server.Port:
    derived_from: Root
    properties:
      admin_state_up:
        description: >-
          The administrative state of this port.
        type: boolean
        default: true
      allowed_address_pairs:
        description: >-
          Additional MAC/IP address pairs allowed to pass through the port.
        type: list
        entry_schema: openstack.nova.Server.AddressPair
      binding.vnic_type:
        description: >-
          The vnic type to be bound on the neutron port. To support SR-IOV PCI passthrough
          networking, you can request that the neutron port to be realized as normal (virtual nic),
          direct (pci passthrough), or macvtap (virtual interface with a tap-like software
          interface). Note that this only works for Neutron deployments that support the bindings
          extension.
        type: string
        default: normal
        constraints:
        - valid_values: [ normal, direct, macvtap, direct-physical, baremetal ]
      mac_address:
        description: >-
          MAC address to give to this port. The default update policy of this property in neutron is
          that allow admin role only.
        type: openstack.MacAddress
      port_security_enabled:
        description: >-
          Flag to enable/disable port security on the port. When disable this feature (set it to
          False), there will be no packages filtering, like security-group and address-pairs.
        type: boolean
      qos_policy:
        description: >-
          The name or ID of QoS policy to attach to this port.
        type: openstack.Neutron.QosPolicy
      value_specs:
        description: >-
          Extra parameters to include in the request.
        type: map
        entry_schema: string # TODO

  openstack.nova.Server.AddressPair:
    derived_from: Root
    properties:
      ip_address:
        description: >-
          IP address to allow through this port.
        type: openstack.NetCidr
      mac_address:
        description: >-
          MAC address to allow through this port.
        type: openstack.MacAddress

  openstack.nova.Server.SwiftData:
    derived_from: Root
    properties:
      container:
        description: >-
          Name of the container.
        type: string
        constraints:
        - min_length: 1
      object:
        description: >-
          Name of the object.
        type: string
        constraints:
        - min_length: 1

  openstack.nova.Server.BlockDevice:
    derived_from: Root
    properties:
      delete_on_termination:
        description: >-
          Indicate whether the volume should be deleted when the server is terminated.
        type: boolean
      device_name:
        description: >-
          A device name where the volume will be attached in the system at /dev/device_name. This
          value is typically vda.
        type: string
      snapshot_id:
        description: >-
          The ID of the snapshot to create a volume from.
        type: openstack.Cinder.Snapshot
      volume_id:
        description: >-
          The ID of the volume to boot from. Only one of volume_id or snapshot_id should be provided.
        type: openstack.Cinder.Volume
      volume_size:
        description: >-
          The size of the volume, in GB. It is safe to leave this blank and have the Compute service
          infer the size.
        type: scalar-unit.size

  openstack.nova.Server.BlockDevice2:
    derived_from: Root
    properties:
      boot_index:
        description: >-
          Integer used for ordering the boot disks. If it is not specified, value "0" will be set
          for bootable sources (volume, snapshot, image); value "-1" will be set for non-bootable
          sources.
        type: integer
      delete_on_termination:
        description: >-
          Indicate whether the volume should be deleted when the server is terminated.
        type: boolean
      device_name:
        description: >-
          A device name where the volume will be attached in the system at /dev/device_name. This
          value is typically vda.
        type: string
      device_type:
        description: >-
          Device type: at the moment we can make distinction only between disk and cdrom.
        type: string
        constraints:
        - valid_values: [ cdrom, disk ]
      disk_bus:
        description: >-
          Bus of the device: hypervisor driver chooses a suitable default if omitted.
        type: string
        constraints:
        - valid_values: [ ide, lame_bus, scsi, usb, virtio ]
      ephemeral_format:
        description: >-
          The format of the local ephemeral block device. If no format is specified, uses default
          value, defined in nova configuration file.
        type: string
        constraints:
        - valid_values: [ ext2, ext3, ext4, xfs, ntfs ]
      ephemeral_size:
        description: >-
          The size of the local ephemeral block device, in GB.
        type: scalar-unit.size
        constraints:
        - greater_or_equal: 1
      image:
        description: >-
          The ID or name of the image to create a volume from.
        type: openstack.Glance.Image
      snapshot_id:
        description: >-
          The ID of the snapshot to create a volume from.
        type: openstack.Cinder.Snapshot
      swap_size:
        description: >-
          The size of the swap, in MB.
        type: scalar-unit.size
      volume_id:
        description: >-
          The volume_id can be boot or non-boot device to the server.
        type: openstack.Cinder.Volume
      volume_size:
        description: >-
          Size of the block device in GB. If it is omitted, hypervisor driver calculates size.
        type: scalar-unit.size
`
}
