package core

import (
	"sync"

	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/service/ec2"
)

type MockEc2 struct {
	methodInvocations *sync.WaitGroup // whenever an expected method is invoked, call Done() on this waitgroup

	// Everytime a method is invoked on this MockEc2, a new message will be pushed
	// into this channel with the primary argument of the method invocation (eg,
	// it will be a *ec2.DescribeInstancesInput if DescribeInstances is invoked)
	methodInvocationsChan chan<- interface{}
}

func NewMockEc2(wg *sync.WaitGroup, mic chan<- interface{}) *MockEc2 {
	return &MockEc2{
		methodInvocations:     wg,
		methodInvocationsChan: mic,
	}
}

func (m *MockEc2) AcceptVpcPeeringConnectionRequest(*ec2.AcceptVpcPeeringConnectionInput) (*request.Request, *ec2.AcceptVpcPeeringConnectionOutput) {
	panic("Not implemented")
}

func (m *MockEc2) AcceptVpcPeeringConnection(*ec2.AcceptVpcPeeringConnectionInput) (*ec2.AcceptVpcPeeringConnectionOutput, error) {
	panic("Not implemented")
}

func (m *MockEc2) AllocateAddressRequest(*ec2.AllocateAddressInput) (*request.Request, *ec2.AllocateAddressOutput) {
	panic("Not implemented")
}

func (m *MockEc2) AllocateAddress(*ec2.AllocateAddressInput) (*ec2.AllocateAddressOutput, error) {
	panic("Not implemented")
}

func (m *MockEc2) AllocateHostsRequest(*ec2.AllocateHostsInput) (*request.Request, *ec2.AllocateHostsOutput) {
	panic("Not implemented")
}

func (m *MockEc2) AllocateHosts(*ec2.AllocateHostsInput) (*ec2.AllocateHostsOutput, error) {
	panic("Not implemented")
}

func (m *MockEc2) AssignPrivateIpAddressesRequest(*ec2.AssignPrivateIpAddressesInput) (*request.Request, *ec2.AssignPrivateIpAddressesOutput) {
	panic("Not implemented")
}

func (m *MockEc2) AssignPrivateIpAddresses(*ec2.AssignPrivateIpAddressesInput) (*ec2.AssignPrivateIpAddressesOutput, error) {
	panic("Not implemented")
}

func (m *MockEc2) AssociateAddressRequest(*ec2.AssociateAddressInput) (*request.Request, *ec2.AssociateAddressOutput) {
	panic("Not implemented")
}

func (m *MockEc2) AssociateAddress(*ec2.AssociateAddressInput) (*ec2.AssociateAddressOutput, error) {
	panic("Not implemented")
}

func (m *MockEc2) AssociateDhcpOptionsRequest(*ec2.AssociateDhcpOptionsInput) (*request.Request, *ec2.AssociateDhcpOptionsOutput) {
	panic("Not implemented")
}

func (m *MockEc2) AssociateDhcpOptions(*ec2.AssociateDhcpOptionsInput) (*ec2.AssociateDhcpOptionsOutput, error) {
	panic("Not implemented")
}

func (m *MockEc2) AssociateRouteTableRequest(*ec2.AssociateRouteTableInput) (*request.Request, *ec2.AssociateRouteTableOutput) {
	panic("Not implemented")
}

func (m *MockEc2) AssociateRouteTable(*ec2.AssociateRouteTableInput) (*ec2.AssociateRouteTableOutput, error) {
	panic("Not implemented")
}

func (m *MockEc2) AttachClassicLinkVpcRequest(*ec2.AttachClassicLinkVpcInput) (*request.Request, *ec2.AttachClassicLinkVpcOutput) {
	panic("Not implemented")
}

func (m *MockEc2) AttachClassicLinkVpc(*ec2.AttachClassicLinkVpcInput) (*ec2.AttachClassicLinkVpcOutput, error) {
	panic("Not implemented")
}

func (m *MockEc2) AttachInternetGatewayRequest(*ec2.AttachInternetGatewayInput) (*request.Request, *ec2.AttachInternetGatewayOutput) {
	panic("Not implemented")
}

func (m *MockEc2) AttachInternetGateway(*ec2.AttachInternetGatewayInput) (*ec2.AttachInternetGatewayOutput, error) {
	panic("Not implemented")
}

func (m *MockEc2) AttachNetworkInterfaceRequest(*ec2.AttachNetworkInterfaceInput) (*request.Request, *ec2.AttachNetworkInterfaceOutput) {
	panic("Not implemented")
}

func (m *MockEc2) AttachNetworkInterface(*ec2.AttachNetworkInterfaceInput) (*ec2.AttachNetworkInterfaceOutput, error) {
	panic("Not implemented")
}

func (m *MockEc2) AttachVolumeRequest(*ec2.AttachVolumeInput) (*request.Request, *ec2.VolumeAttachment) {
	panic("Not implemented")
}

func (m *MockEc2) AttachVolume(*ec2.AttachVolumeInput) (*ec2.VolumeAttachment, error) {
	panic("Not implemented")
}

func (m *MockEc2) AttachVpnGatewayRequest(*ec2.AttachVpnGatewayInput) (*request.Request, *ec2.AttachVpnGatewayOutput) {
	panic("Not implemented")
}

func (m *MockEc2) AttachVpnGateway(*ec2.AttachVpnGatewayInput) (*ec2.AttachVpnGatewayOutput, error) {
	panic("Not implemented")
}

func (m *MockEc2) AuthorizeSecurityGroupEgressRequest(*ec2.AuthorizeSecurityGroupEgressInput) (*request.Request, *ec2.AuthorizeSecurityGroupEgressOutput) {
	panic("Not implemented")
}

func (m *MockEc2) AuthorizeSecurityGroupEgress(*ec2.AuthorizeSecurityGroupEgressInput) (*ec2.AuthorizeSecurityGroupEgressOutput, error) {
	panic("Not implemented")
}

func (m *MockEc2) AuthorizeSecurityGroupIngressRequest(*ec2.AuthorizeSecurityGroupIngressInput) (*request.Request, *ec2.AuthorizeSecurityGroupIngressOutput) {
	panic("Not implemented")
}

func (m *MockEc2) AuthorizeSecurityGroupIngress(*ec2.AuthorizeSecurityGroupIngressInput) (*ec2.AuthorizeSecurityGroupIngressOutput, error) {
	panic("Not implemented")
}

func (m *MockEc2) BundleInstanceRequest(*ec2.BundleInstanceInput) (*request.Request, *ec2.BundleInstanceOutput) {
	panic("Not implemented")
}

func (m *MockEc2) BundleInstance(*ec2.BundleInstanceInput) (*ec2.BundleInstanceOutput, error) {
	panic("Not implemented")
}

func (m *MockEc2) CancelBundleTaskRequest(*ec2.CancelBundleTaskInput) (*request.Request, *ec2.CancelBundleTaskOutput) {
	panic("Not implemented")
}

func (m *MockEc2) CancelBundleTask(*ec2.CancelBundleTaskInput) (*ec2.CancelBundleTaskOutput, error) {
	panic("Not implemented")
}

func (m *MockEc2) CancelConversionTaskRequest(*ec2.CancelConversionTaskInput) (*request.Request, *ec2.CancelConversionTaskOutput) {
	panic("Not implemented")
}

func (m *MockEc2) CancelConversionTask(*ec2.CancelConversionTaskInput) (*ec2.CancelConversionTaskOutput, error) {
	panic("Not implemented")
}

func (m *MockEc2) CancelExportTaskRequest(*ec2.CancelExportTaskInput) (*request.Request, *ec2.CancelExportTaskOutput) {
	panic("Not implemented")
}

func (m *MockEc2) CancelExportTask(*ec2.CancelExportTaskInput) (*ec2.CancelExportTaskOutput, error) {
	panic("Not implemented")
}

func (m *MockEc2) CancelImportTaskRequest(*ec2.CancelImportTaskInput) (*request.Request, *ec2.CancelImportTaskOutput) {
	panic("Not implemented")
}

func (m *MockEc2) CancelImportTask(*ec2.CancelImportTaskInput) (*ec2.CancelImportTaskOutput, error) {
	panic("Not implemented")
}

func (m *MockEc2) CancelReservedInstancesListingRequest(*ec2.CancelReservedInstancesListingInput) (*request.Request, *ec2.CancelReservedInstancesListingOutput) {
	panic("Not implemented")
}

func (m *MockEc2) CancelReservedInstancesListing(*ec2.CancelReservedInstancesListingInput) (*ec2.CancelReservedInstancesListingOutput, error) {
	panic("Not implemented")
}

func (m *MockEc2) CancelSpotFleetRequestsRequest(*ec2.CancelSpotFleetRequestsInput) (*request.Request, *ec2.CancelSpotFleetRequestsOutput) {
	panic("Not implemented")
}

func (m *MockEc2) CancelSpotFleetRequests(*ec2.CancelSpotFleetRequestsInput) (*ec2.CancelSpotFleetRequestsOutput, error) {
	panic("Not implemented")
}

func (m *MockEc2) CancelSpotInstanceRequestsRequest(*ec2.CancelSpotInstanceRequestsInput) (*request.Request, *ec2.CancelSpotInstanceRequestsOutput) {
	panic("Not implemented")
}

func (m *MockEc2) CancelSpotInstanceRequests(*ec2.CancelSpotInstanceRequestsInput) (*ec2.CancelSpotInstanceRequestsOutput, error) {
	panic("Not implemented")
}

func (m *MockEc2) ConfirmProductInstanceRequest(*ec2.ConfirmProductInstanceInput) (*request.Request, *ec2.ConfirmProductInstanceOutput) {
	panic("Not implemented")
}

func (m *MockEc2) ConfirmProductInstance(*ec2.ConfirmProductInstanceInput) (*ec2.ConfirmProductInstanceOutput, error) {
	panic("Not implemented")
}

func (m *MockEc2) CopyImageRequest(*ec2.CopyImageInput) (*request.Request, *ec2.CopyImageOutput) {
	panic("Not implemented")
}

func (m *MockEc2) CopyImage(*ec2.CopyImageInput) (*ec2.CopyImageOutput, error) {
	panic("Not implemented")
}

func (m *MockEc2) CopySnapshotRequest(*ec2.CopySnapshotInput) (*request.Request, *ec2.CopySnapshotOutput) {
	panic("Not implemented")
}

func (m *MockEc2) CopySnapshot(*ec2.CopySnapshotInput) (*ec2.CopySnapshotOutput, error) {
	panic("Not implemented")
}

func (m *MockEc2) CreateCustomerGatewayRequest(*ec2.CreateCustomerGatewayInput) (*request.Request, *ec2.CreateCustomerGatewayOutput) {
	panic("Not implemented")
}

func (m *MockEc2) CreateCustomerGateway(*ec2.CreateCustomerGatewayInput) (*ec2.CreateCustomerGatewayOutput, error) {
	panic("Not implemented")
}

func (m *MockEc2) CreateDhcpOptionsRequest(*ec2.CreateDhcpOptionsInput) (*request.Request, *ec2.CreateDhcpOptionsOutput) {
	panic("Not implemented")
}

func (m *MockEc2) CreateDhcpOptions(*ec2.CreateDhcpOptionsInput) (*ec2.CreateDhcpOptionsOutput, error) {
	panic("Not implemented")
}

func (m *MockEc2) CreateFlowLogsRequest(*ec2.CreateFlowLogsInput) (*request.Request, *ec2.CreateFlowLogsOutput) {
	panic("Not implemented")
}

func (m *MockEc2) CreateFlowLogs(*ec2.CreateFlowLogsInput) (*ec2.CreateFlowLogsOutput, error) {
	panic("Not implemented")
}

func (m *MockEc2) CreateImageRequest(*ec2.CreateImageInput) (*request.Request, *ec2.CreateImageOutput) {
	panic("Not implemented")
}

func (m *MockEc2) CreateImage(*ec2.CreateImageInput) (*ec2.CreateImageOutput, error) {
	panic("Not implemented")
}

func (m *MockEc2) CreateInstanceExportTaskRequest(*ec2.CreateInstanceExportTaskInput) (*request.Request, *ec2.CreateInstanceExportTaskOutput) {
	panic("Not implemented")
}

func (m *MockEc2) CreateInstanceExportTask(*ec2.CreateInstanceExportTaskInput) (*ec2.CreateInstanceExportTaskOutput, error) {
	panic("Not implemented")
}

func (m *MockEc2) CreateInternetGatewayRequest(*ec2.CreateInternetGatewayInput) (*request.Request, *ec2.CreateInternetGatewayOutput) {
	panic("Not implemented")
}

func (m *MockEc2) CreateInternetGateway(*ec2.CreateInternetGatewayInput) (*ec2.CreateInternetGatewayOutput, error) {
	panic("Not implemented")
}

func (m *MockEc2) CreateKeyPairRequest(*ec2.CreateKeyPairInput) (*request.Request, *ec2.CreateKeyPairOutput) {
	panic("Not implemented")
}

func (m *MockEc2) CreateKeyPair(*ec2.CreateKeyPairInput) (*ec2.CreateKeyPairOutput, error) {
	panic("Not implemented")
}

func (m *MockEc2) CreateNatGatewayRequest(*ec2.CreateNatGatewayInput) (*request.Request, *ec2.CreateNatGatewayOutput) {
	panic("Not implemented")
}

func (m *MockEc2) CreateNatGateway(*ec2.CreateNatGatewayInput) (*ec2.CreateNatGatewayOutput, error) {
	panic("Not implemented")
}

func (m *MockEc2) CreateNetworkAclRequest(*ec2.CreateNetworkAclInput) (*request.Request, *ec2.CreateNetworkAclOutput) {
	panic("Not implemented")
}

func (m *MockEc2) CreateNetworkAcl(*ec2.CreateNetworkAclInput) (*ec2.CreateNetworkAclOutput, error) {
	panic("Not implemented")
}

func (m *MockEc2) CreateNetworkAclEntryRequest(*ec2.CreateNetworkAclEntryInput) (*request.Request, *ec2.CreateNetworkAclEntryOutput) {
	panic("Not implemented")
}

func (m *MockEc2) CreateNetworkAclEntry(*ec2.CreateNetworkAclEntryInput) (*ec2.CreateNetworkAclEntryOutput, error) {
	panic("Not implemented")
}

func (m *MockEc2) CreateNetworkInterfaceRequest(*ec2.CreateNetworkInterfaceInput) (*request.Request, *ec2.CreateNetworkInterfaceOutput) {
	panic("Not implemented")
}

func (m *MockEc2) CreateNetworkInterface(*ec2.CreateNetworkInterfaceInput) (*ec2.CreateNetworkInterfaceOutput, error) {
	panic("Not implemented")
}

func (m *MockEc2) CreatePlacementGroupRequest(*ec2.CreatePlacementGroupInput) (*request.Request, *ec2.CreatePlacementGroupOutput) {
	panic("Not implemented")
}

func (m *MockEc2) CreatePlacementGroup(*ec2.CreatePlacementGroupInput) (*ec2.CreatePlacementGroupOutput, error) {
	panic("Not implemented")
}

func (m *MockEc2) CreateReservedInstancesListingRequest(*ec2.CreateReservedInstancesListingInput) (*request.Request, *ec2.CreateReservedInstancesListingOutput) {
	panic("Not implemented")
}

func (m *MockEc2) CreateReservedInstancesListing(*ec2.CreateReservedInstancesListingInput) (*ec2.CreateReservedInstancesListingOutput, error) {
	panic("Not implemented")
}

func (m *MockEc2) CreateRouteRequest(*ec2.CreateRouteInput) (*request.Request, *ec2.CreateRouteOutput) {
	panic("Not implemented")
}

func (m *MockEc2) CreateRoute(*ec2.CreateRouteInput) (*ec2.CreateRouteOutput, error) {
	panic("Not implemented")
}

func (m *MockEc2) CreateRouteTableRequest(*ec2.CreateRouteTableInput) (*request.Request, *ec2.CreateRouteTableOutput) {
	panic("Not implemented")
}

func (m *MockEc2) CreateRouteTable(*ec2.CreateRouteTableInput) (*ec2.CreateRouteTableOutput, error) {
	panic("Not implemented")
}

func (m *MockEc2) CreateSecurityGroupRequest(*ec2.CreateSecurityGroupInput) (*request.Request, *ec2.CreateSecurityGroupOutput) {
	panic("Not implemented")
}

func (m *MockEc2) CreateSecurityGroup(*ec2.CreateSecurityGroupInput) (*ec2.CreateSecurityGroupOutput, error) {
	panic("Not implemented")
}

func (m *MockEc2) CreateSnapshotRequest(*ec2.CreateSnapshotInput) (*request.Request, *ec2.Snapshot) {
	panic("Not implemented")
}

func (m *MockEc2) CreateSnapshot(*ec2.CreateSnapshotInput) (*ec2.Snapshot, error) {
	panic("Not implemented")
}

func (m *MockEc2) CreateSpotDatafeedSubscriptionRequest(*ec2.CreateSpotDatafeedSubscriptionInput) (*request.Request, *ec2.CreateSpotDatafeedSubscriptionOutput) {
	panic("Not implemented")
}

func (m *MockEc2) CreateSpotDatafeedSubscription(*ec2.CreateSpotDatafeedSubscriptionInput) (*ec2.CreateSpotDatafeedSubscriptionOutput, error) {
	panic("Not implemented")
}

func (m *MockEc2) CreateSubnetRequest(*ec2.CreateSubnetInput) (*request.Request, *ec2.CreateSubnetOutput) {
	panic("Not implemented")
}

func (m *MockEc2) CreateSubnet(*ec2.CreateSubnetInput) (*ec2.CreateSubnetOutput, error) {
	panic("Not implemented")
}

func (m *MockEc2) CreateTagsRequest(*ec2.CreateTagsInput) (*request.Request, *ec2.CreateTagsOutput) {
	panic("Not implemented")
}

func (m *MockEc2) CreateTags(*ec2.CreateTagsInput) (*ec2.CreateTagsOutput, error) {
	panic("Not implemented")
}

func (m *MockEc2) CreateVolumeRequest(*ec2.CreateVolumeInput) (*request.Request, *ec2.Volume) {
	panic("Not implemented")
}

func (m *MockEc2) CreateVolume(*ec2.CreateVolumeInput) (*ec2.Volume, error) {
	panic("Not implemented")
}

func (m *MockEc2) CreateVpcRequest(*ec2.CreateVpcInput) (*request.Request, *ec2.CreateVpcOutput) {
	panic("Not implemented")
}

func (m *MockEc2) CreateVpc(*ec2.CreateVpcInput) (*ec2.CreateVpcOutput, error) {
	panic("Not implemented")
}

func (m *MockEc2) CreateVpcEndpointRequest(*ec2.CreateVpcEndpointInput) (*request.Request, *ec2.CreateVpcEndpointOutput) {
	panic("Not implemented")
}

func (m *MockEc2) CreateVpcEndpoint(*ec2.CreateVpcEndpointInput) (*ec2.CreateVpcEndpointOutput, error) {
	panic("Not implemented")
}

func (m *MockEc2) CreateVpcPeeringConnectionRequest(*ec2.CreateVpcPeeringConnectionInput) (*request.Request, *ec2.CreateVpcPeeringConnectionOutput) {
	panic("Not implemented")
}

func (m *MockEc2) CreateVpcPeeringConnection(*ec2.CreateVpcPeeringConnectionInput) (*ec2.CreateVpcPeeringConnectionOutput, error) {
	panic("Not implemented")
}

func (m *MockEc2) CreateVpnConnectionRequest(*ec2.CreateVpnConnectionInput) (*request.Request, *ec2.CreateVpnConnectionOutput) {
	panic("Not implemented")
}

func (m *MockEc2) CreateVpnConnection(*ec2.CreateVpnConnectionInput) (*ec2.CreateVpnConnectionOutput, error) {
	panic("Not implemented")
}

func (m *MockEc2) CreateVpnConnectionRouteRequest(*ec2.CreateVpnConnectionRouteInput) (*request.Request, *ec2.CreateVpnConnectionRouteOutput) {
	panic("Not implemented")
}

func (m *MockEc2) CreateVpnConnectionRoute(*ec2.CreateVpnConnectionRouteInput) (*ec2.CreateVpnConnectionRouteOutput, error) {
	panic("Not implemented")
}

func (m *MockEc2) CreateVpnGatewayRequest(*ec2.CreateVpnGatewayInput) (*request.Request, *ec2.CreateVpnGatewayOutput) {
	panic("Not implemented")
}

func (m *MockEc2) CreateVpnGateway(*ec2.CreateVpnGatewayInput) (*ec2.CreateVpnGatewayOutput, error) {
	panic("Not implemented")
}

func (m *MockEc2) DeleteCustomerGatewayRequest(*ec2.DeleteCustomerGatewayInput) (*request.Request, *ec2.DeleteCustomerGatewayOutput) {
	panic("Not implemented")
}

func (m *MockEc2) DeleteCustomerGateway(*ec2.DeleteCustomerGatewayInput) (*ec2.DeleteCustomerGatewayOutput, error) {
	panic("Not implemented")
}

func (m *MockEc2) DeleteDhcpOptionsRequest(*ec2.DeleteDhcpOptionsInput) (*request.Request, *ec2.DeleteDhcpOptionsOutput) {
	panic("Not implemented")
}

func (m *MockEc2) DeleteDhcpOptions(*ec2.DeleteDhcpOptionsInput) (*ec2.DeleteDhcpOptionsOutput, error) {
	panic("Not implemented")
}

func (m *MockEc2) DeleteFlowLogsRequest(*ec2.DeleteFlowLogsInput) (*request.Request, *ec2.DeleteFlowLogsOutput) {
	panic("Not implemented")
}

func (m *MockEc2) DeleteFlowLogs(*ec2.DeleteFlowLogsInput) (*ec2.DeleteFlowLogsOutput, error) {
	panic("Not implemented")
}

func (m *MockEc2) DeleteInternetGatewayRequest(*ec2.DeleteInternetGatewayInput) (*request.Request, *ec2.DeleteInternetGatewayOutput) {
	panic("Not implemented")
}

func (m *MockEc2) DeleteInternetGateway(*ec2.DeleteInternetGatewayInput) (*ec2.DeleteInternetGatewayOutput, error) {
	panic("Not implemented")
}

func (m *MockEc2) DeleteKeyPairRequest(*ec2.DeleteKeyPairInput) (*request.Request, *ec2.DeleteKeyPairOutput) {
	panic("Not implemented")
}

func (m *MockEc2) DeleteKeyPair(*ec2.DeleteKeyPairInput) (*ec2.DeleteKeyPairOutput, error) {
	panic("Not implemented")
}

func (m *MockEc2) DeleteNatGatewayRequest(*ec2.DeleteNatGatewayInput) (*request.Request, *ec2.DeleteNatGatewayOutput) {
	panic("Not implemented")
}

func (m *MockEc2) DeleteNatGateway(*ec2.DeleteNatGatewayInput) (*ec2.DeleteNatGatewayOutput, error) {
	panic("Not implemented")
}

func (m *MockEc2) DeleteNetworkAclRequest(*ec2.DeleteNetworkAclInput) (*request.Request, *ec2.DeleteNetworkAclOutput) {
	panic("Not implemented")
}

func (m *MockEc2) DeleteNetworkAcl(*ec2.DeleteNetworkAclInput) (*ec2.DeleteNetworkAclOutput, error) {
	panic("Not implemented")
}

func (m *MockEc2) DeleteNetworkAclEntryRequest(*ec2.DeleteNetworkAclEntryInput) (*request.Request, *ec2.DeleteNetworkAclEntryOutput) {
	panic("Not implemented")
}

func (m *MockEc2) DeleteNetworkAclEntry(*ec2.DeleteNetworkAclEntryInput) (*ec2.DeleteNetworkAclEntryOutput, error) {
	panic("Not implemented")
}

func (m *MockEc2) DeleteNetworkInterfaceRequest(*ec2.DeleteNetworkInterfaceInput) (*request.Request, *ec2.DeleteNetworkInterfaceOutput) {
	panic("Not implemented")
}

func (m *MockEc2) DeleteNetworkInterface(*ec2.DeleteNetworkInterfaceInput) (*ec2.DeleteNetworkInterfaceOutput, error) {
	panic("Not implemented")
}

func (m *MockEc2) DeletePlacementGroupRequest(*ec2.DeletePlacementGroupInput) (*request.Request, *ec2.DeletePlacementGroupOutput) {
	panic("Not implemented")
}

func (m *MockEc2) DeletePlacementGroup(*ec2.DeletePlacementGroupInput) (*ec2.DeletePlacementGroupOutput, error) {
	panic("Not implemented")
}

func (m *MockEc2) DeleteRouteRequest(*ec2.DeleteRouteInput) (*request.Request, *ec2.DeleteRouteOutput) {
	panic("Not implemented")
}

func (m *MockEc2) DeleteRoute(*ec2.DeleteRouteInput) (*ec2.DeleteRouteOutput, error) {
	panic("Not implemented")
}

func (m *MockEc2) DeleteRouteTableRequest(*ec2.DeleteRouteTableInput) (*request.Request, *ec2.DeleteRouteTableOutput) {
	panic("Not implemented")
}

func (m *MockEc2) DeleteRouteTable(*ec2.DeleteRouteTableInput) (*ec2.DeleteRouteTableOutput, error) {
	panic("Not implemented")
}

func (m *MockEc2) DeleteSecurityGroupRequest(*ec2.DeleteSecurityGroupInput) (*request.Request, *ec2.DeleteSecurityGroupOutput) {
	panic("Not implemented")
}

func (m *MockEc2) DeleteSecurityGroup(*ec2.DeleteSecurityGroupInput) (*ec2.DeleteSecurityGroupOutput, error) {
	panic("Not implemented")
}

func (m *MockEc2) DeleteSnapshotRequest(*ec2.DeleteSnapshotInput) (*request.Request, *ec2.DeleteSnapshotOutput) {
	panic("Not implemented")
}

func (m *MockEc2) DeleteSnapshot(*ec2.DeleteSnapshotInput) (*ec2.DeleteSnapshotOutput, error) {
	panic("Not implemented")
}

func (m *MockEc2) DeleteSpotDatafeedSubscriptionRequest(*ec2.DeleteSpotDatafeedSubscriptionInput) (*request.Request, *ec2.DeleteSpotDatafeedSubscriptionOutput) {
	panic("Not implemented")
}

func (m *MockEc2) DeleteSpotDatafeedSubscription(*ec2.DeleteSpotDatafeedSubscriptionInput) (*ec2.DeleteSpotDatafeedSubscriptionOutput, error) {
	panic("Not implemented")
}

func (m *MockEc2) DeleteSubnetRequest(*ec2.DeleteSubnetInput) (*request.Request, *ec2.DeleteSubnetOutput) {
	panic("Not implemented")
}

func (m *MockEc2) DeleteSubnet(*ec2.DeleteSubnetInput) (*ec2.DeleteSubnetOutput, error) {
	panic("Not implemented")
}

func (m *MockEc2) DeleteTagsRequest(*ec2.DeleteTagsInput) (*request.Request, *ec2.DeleteTagsOutput) {
	panic("Not implemented")
}

func (m *MockEc2) DeleteTags(*ec2.DeleteTagsInput) (*ec2.DeleteTagsOutput, error) {
	panic("Not implemented")
}

func (m *MockEc2) DeleteVolumeRequest(*ec2.DeleteVolumeInput) (*request.Request, *ec2.DeleteVolumeOutput) {
	panic("Not implemented")
}

func (m *MockEc2) DeleteVolume(*ec2.DeleteVolumeInput) (*ec2.DeleteVolumeOutput, error) {
	panic("Not implemented")
}

func (m *MockEc2) DeleteVpcRequest(*ec2.DeleteVpcInput) (*request.Request, *ec2.DeleteVpcOutput) {
	panic("Not implemented")
}

func (m *MockEc2) DeleteVpc(*ec2.DeleteVpcInput) (*ec2.DeleteVpcOutput, error) {
	panic("Not implemented")
}

func (m *MockEc2) DeleteVpcEndpointsRequest(*ec2.DeleteVpcEndpointsInput) (*request.Request, *ec2.DeleteVpcEndpointsOutput) {
	panic("Not implemented")
}

func (m *MockEc2) DeleteVpcEndpoints(*ec2.DeleteVpcEndpointsInput) (*ec2.DeleteVpcEndpointsOutput, error) {
	panic("Not implemented")
}

func (m *MockEc2) DeleteVpcPeeringConnectionRequest(*ec2.DeleteVpcPeeringConnectionInput) (*request.Request, *ec2.DeleteVpcPeeringConnectionOutput) {
	panic("Not implemented")
}

func (m *MockEc2) DeleteVpcPeeringConnection(*ec2.DeleteVpcPeeringConnectionInput) (*ec2.DeleteVpcPeeringConnectionOutput, error) {
	panic("Not implemented")
}

func (m *MockEc2) DeleteVpnConnectionRequest(*ec2.DeleteVpnConnectionInput) (*request.Request, *ec2.DeleteVpnConnectionOutput) {
	panic("Not implemented")
}

func (m *MockEc2) DeleteVpnConnection(*ec2.DeleteVpnConnectionInput) (*ec2.DeleteVpnConnectionOutput, error) {
	panic("Not implemented")
}

func (m *MockEc2) DeleteVpnConnectionRouteRequest(*ec2.DeleteVpnConnectionRouteInput) (*request.Request, *ec2.DeleteVpnConnectionRouteOutput) {
	panic("Not implemented")
}

func (m *MockEc2) DeleteVpnConnectionRoute(*ec2.DeleteVpnConnectionRouteInput) (*ec2.DeleteVpnConnectionRouteOutput, error) {
	panic("Not implemented")
}

func (m *MockEc2) DeleteVpnGatewayRequest(*ec2.DeleteVpnGatewayInput) (*request.Request, *ec2.DeleteVpnGatewayOutput) {
	panic("Not implemented")
}

func (m *MockEc2) DeleteVpnGateway(*ec2.DeleteVpnGatewayInput) (*ec2.DeleteVpnGatewayOutput, error) {
	panic("Not implemented")
}

func (m *MockEc2) DeregisterImageRequest(*ec2.DeregisterImageInput) (*request.Request, *ec2.DeregisterImageOutput) {
	panic("Not implemented")
}

func (m *MockEc2) DeregisterImage(*ec2.DeregisterImageInput) (*ec2.DeregisterImageOutput, error) {
	panic("Not implemented")
}

func (m *MockEc2) DescribeAccountAttributesRequest(*ec2.DescribeAccountAttributesInput) (*request.Request, *ec2.DescribeAccountAttributesOutput) {
	panic("Not implemented")
}

func (m *MockEc2) DescribeAccountAttributes(*ec2.DescribeAccountAttributesInput) (*ec2.DescribeAccountAttributesOutput, error) {
	panic("Not implemented")
}

func (m *MockEc2) DescribeAddressesRequest(*ec2.DescribeAddressesInput) (*request.Request, *ec2.DescribeAddressesOutput) {
	panic("Not implemented")
}

func (m *MockEc2) DescribeAddresses(*ec2.DescribeAddressesInput) (*ec2.DescribeAddressesOutput, error) {
	panic("Not implemented")
}

func (m *MockEc2) DescribeAvailabilityZonesRequest(*ec2.DescribeAvailabilityZonesInput) (*request.Request, *ec2.DescribeAvailabilityZonesOutput) {
	panic("Not implemented")
}

func (m *MockEc2) DescribeAvailabilityZones(*ec2.DescribeAvailabilityZonesInput) (*ec2.DescribeAvailabilityZonesOutput, error) {
	panic("Not implemented")
}

func (m *MockEc2) DescribeBundleTasksRequest(*ec2.DescribeBundleTasksInput) (*request.Request, *ec2.DescribeBundleTasksOutput) {
	panic("Not implemented")
}

func (m *MockEc2) DescribeBundleTasks(*ec2.DescribeBundleTasksInput) (*ec2.DescribeBundleTasksOutput, error) {
	panic("Not implemented")
}

func (m *MockEc2) DescribeClassicLinkInstancesRequest(*ec2.DescribeClassicLinkInstancesInput) (*request.Request, *ec2.DescribeClassicLinkInstancesOutput) {
	panic("Not implemented")
}

func (m *MockEc2) DescribeClassicLinkInstances(*ec2.DescribeClassicLinkInstancesInput) (*ec2.DescribeClassicLinkInstancesOutput, error) {
	panic("Not implemented")
}

func (m *MockEc2) DescribeConversionTasksRequest(*ec2.DescribeConversionTasksInput) (*request.Request, *ec2.DescribeConversionTasksOutput) {
	panic("Not implemented")
}

func (m *MockEc2) DescribeConversionTasks(*ec2.DescribeConversionTasksInput) (*ec2.DescribeConversionTasksOutput, error) {
	panic("Not implemented")
}

func (m *MockEc2) DescribeCustomerGatewaysRequest(*ec2.DescribeCustomerGatewaysInput) (*request.Request, *ec2.DescribeCustomerGatewaysOutput) {
	panic("Not implemented")
}

func (m *MockEc2) DescribeCustomerGateways(*ec2.DescribeCustomerGatewaysInput) (*ec2.DescribeCustomerGatewaysOutput, error) {
	panic("Not implemented")
}

func (m *MockEc2) DescribeDhcpOptionsRequest(*ec2.DescribeDhcpOptionsInput) (*request.Request, *ec2.DescribeDhcpOptionsOutput) {
	panic("Not implemented")
}

func (m *MockEc2) DescribeDhcpOptions(*ec2.DescribeDhcpOptionsInput) (*ec2.DescribeDhcpOptionsOutput, error) {
	panic("Not implemented")
}

func (m *MockEc2) DescribeExportTasksRequest(*ec2.DescribeExportTasksInput) (*request.Request, *ec2.DescribeExportTasksOutput) {
	panic("Not implemented")
}

func (m *MockEc2) DescribeExportTasks(*ec2.DescribeExportTasksInput) (*ec2.DescribeExportTasksOutput, error) {
	panic("Not implemented")
}

func (m *MockEc2) DescribeFlowLogsRequest(*ec2.DescribeFlowLogsInput) (*request.Request, *ec2.DescribeFlowLogsOutput) {
	panic("Not implemented")
}

func (m *MockEc2) DescribeFlowLogs(*ec2.DescribeFlowLogsInput) (*ec2.DescribeFlowLogsOutput, error) {
	panic("Not implemented")
}

func (m *MockEc2) DescribeHostReservationOfferingsRequest(*ec2.DescribeHostReservationOfferingsInput) (*request.Request, *ec2.DescribeHostReservationOfferingsOutput) {
	panic("Not implemented")
}

func (m *MockEc2) DescribeHostReservationOfferings(*ec2.DescribeHostReservationOfferingsInput) (*ec2.DescribeHostReservationOfferingsOutput, error) {
	panic("Not implemented")
}

func (m *MockEc2) DescribeHostReservationsRequest(*ec2.DescribeHostReservationsInput) (*request.Request, *ec2.DescribeHostReservationsOutput) {
	panic("Not implemented")
}

func (m *MockEc2) DescribeHostReservations(*ec2.DescribeHostReservationsInput) (*ec2.DescribeHostReservationsOutput, error) {
	panic("Not implemented")
}

func (m *MockEc2) DescribeHostsRequest(*ec2.DescribeHostsInput) (*request.Request, *ec2.DescribeHostsOutput) {
	panic("Not implemented")
}

func (m *MockEc2) DescribeHosts(*ec2.DescribeHostsInput) (*ec2.DescribeHostsOutput, error) {
	panic("Not implemented")
}

func (m *MockEc2) DescribeIdFormatRequest(*ec2.DescribeIdFormatInput) (*request.Request, *ec2.DescribeIdFormatOutput) {
	panic("Not implemented")
}

func (m *MockEc2) DescribeIdFormat(*ec2.DescribeIdFormatInput) (*ec2.DescribeIdFormatOutput, error) {
	panic("Not implemented")
}

func (m *MockEc2) DescribeIdentityIdFormatRequest(*ec2.DescribeIdentityIdFormatInput) (*request.Request, *ec2.DescribeIdentityIdFormatOutput) {
	panic("Not implemented")
}

func (m *MockEc2) DescribeIdentityIdFormat(*ec2.DescribeIdentityIdFormatInput) (*ec2.DescribeIdentityIdFormatOutput, error) {
	panic("Not implemented")
}

func (m *MockEc2) DescribeImageAttributeRequest(*ec2.DescribeImageAttributeInput) (*request.Request, *ec2.DescribeImageAttributeOutput) {
	panic("Not implemented")
}

func (m *MockEc2) DescribeImageAttribute(*ec2.DescribeImageAttributeInput) (*ec2.DescribeImageAttributeOutput, error) {
	panic("Not implemented")
}

func (m *MockEc2) DescribeImagesRequest(*ec2.DescribeImagesInput) (*request.Request, *ec2.DescribeImagesOutput) {
	panic("Not implemented")
}

func (m *MockEc2) DescribeImages(*ec2.DescribeImagesInput) (*ec2.DescribeImagesOutput, error) {
	panic("Not implemented")
}

func (m *MockEc2) DescribeImportImageTasksRequest(*ec2.DescribeImportImageTasksInput) (*request.Request, *ec2.DescribeImportImageTasksOutput) {
	panic("Not implemented")
}

func (m *MockEc2) DescribeImportImageTasks(*ec2.DescribeImportImageTasksInput) (*ec2.DescribeImportImageTasksOutput, error) {
	panic("Not implemented")
}

func (m *MockEc2) DescribeImportSnapshotTasksRequest(*ec2.DescribeImportSnapshotTasksInput) (*request.Request, *ec2.DescribeImportSnapshotTasksOutput) {
	panic("Not implemented")
}

func (m *MockEc2) DescribeImportSnapshotTasks(*ec2.DescribeImportSnapshotTasksInput) (*ec2.DescribeImportSnapshotTasksOutput, error) {
	panic("Not implemented")
}

func (m *MockEc2) DescribeInstanceAttributeRequest(*ec2.DescribeInstanceAttributeInput) (*request.Request, *ec2.DescribeInstanceAttributeOutput) {
	panic("Not implemented")
}

func (m *MockEc2) DescribeInstanceAttribute(*ec2.DescribeInstanceAttributeInput) (*ec2.DescribeInstanceAttributeOutput, error) {
	panic("Not implemented")
}

func (m *MockEc2) DescribeInstanceStatusRequest(*ec2.DescribeInstanceStatusInput) (*request.Request, *ec2.DescribeInstanceStatusOutput) {
	panic("Not implemented")
}

func (m *MockEc2) DescribeInstanceStatus(*ec2.DescribeInstanceStatusInput) (*ec2.DescribeInstanceStatusOutput, error) {
	panic("Not implemented")
}

func (m *MockEc2) DescribeInstanceStatusPages(*ec2.DescribeInstanceStatusInput, func(*ec2.DescribeInstanceStatusOutput, bool) bool) error {
	panic("Not implemented")
}

func (m *MockEc2) DescribeInstancesRequest(*ec2.DescribeInstancesInput) (*request.Request, *ec2.DescribeInstancesOutput) {
	panic("Not implemented")
}

func (m *MockEc2) DescribeInstances(input *ec2.DescribeInstancesInput) (*ec2.DescribeInstancesOutput, error) {

	logger.Info("MockEc2 DescribeInstances", "DescribeInstancesInput", input)
	defer func() {
		m.methodInvocationsChan <- input
		m.methodInvocations.Done()
	}()

	az := "us-east-1a"
	instanceState := ec2.InstanceStateNamePending

	instance := ec2.Instance{
		InstanceId: input.InstanceIds[0],
		Placement: &ec2.Placement{
			AvailabilityZone: &az,
		},
		State: &ec2.InstanceState{
			Name: &instanceState,
		},
	}
	reservation := ec2.Reservation{
		Instances: []*ec2.Instance{
			&instance,
		},
	}
	output := ec2.DescribeInstancesOutput{
		Reservations: []*ec2.Reservation{
			&reservation,
		},
	}

	return &output, nil
}

func (m *MockEc2) DescribeInstancesPages(*ec2.DescribeInstancesInput, func(*ec2.DescribeInstancesOutput, bool) bool) error {
	panic("Not implemented")
}

func (m *MockEc2) DescribeInternetGatewaysRequest(*ec2.DescribeInternetGatewaysInput) (*request.Request, *ec2.DescribeInternetGatewaysOutput) {
	panic("Not implemented")
}

func (m *MockEc2) DescribeInternetGateways(*ec2.DescribeInternetGatewaysInput) (*ec2.DescribeInternetGatewaysOutput, error) {
	panic("Not implemented")
}

func (m *MockEc2) DescribeKeyPairsRequest(*ec2.DescribeKeyPairsInput) (*request.Request, *ec2.DescribeKeyPairsOutput) {
	panic("Not implemented")
}

func (m *MockEc2) DescribeKeyPairs(*ec2.DescribeKeyPairsInput) (*ec2.DescribeKeyPairsOutput, error) {
	panic("Not implemented")
}

func (m *MockEc2) DescribeMovingAddressesRequest(*ec2.DescribeMovingAddressesInput) (*request.Request, *ec2.DescribeMovingAddressesOutput) {
	panic("Not implemented")
}

func (m *MockEc2) DescribeMovingAddresses(*ec2.DescribeMovingAddressesInput) (*ec2.DescribeMovingAddressesOutput, error) {
	panic("Not implemented")
}

func (m *MockEc2) DescribeNatGatewaysRequest(*ec2.DescribeNatGatewaysInput) (*request.Request, *ec2.DescribeNatGatewaysOutput) {
	panic("Not implemented")
}

func (m *MockEc2) DescribeNatGateways(*ec2.DescribeNatGatewaysInput) (*ec2.DescribeNatGatewaysOutput, error) {
	panic("Not implemented")
}

func (m *MockEc2) DescribeNetworkAclsRequest(*ec2.DescribeNetworkAclsInput) (*request.Request, *ec2.DescribeNetworkAclsOutput) {
	panic("Not implemented")
}

func (m *MockEc2) DescribeNetworkAcls(*ec2.DescribeNetworkAclsInput) (*ec2.DescribeNetworkAclsOutput, error) {
	panic("Not implemented")
}

func (m *MockEc2) DescribeNetworkInterfaceAttributeRequest(*ec2.DescribeNetworkInterfaceAttributeInput) (*request.Request, *ec2.DescribeNetworkInterfaceAttributeOutput) {
	panic("Not implemented")
}

func (m *MockEc2) DescribeNetworkInterfaceAttribute(*ec2.DescribeNetworkInterfaceAttributeInput) (*ec2.DescribeNetworkInterfaceAttributeOutput, error) {
	panic("Not implemented")
}

func (m *MockEc2) DescribeNetworkInterfacesRequest(*ec2.DescribeNetworkInterfacesInput) (*request.Request, *ec2.DescribeNetworkInterfacesOutput) {
	panic("Not implemented")
}

func (m *MockEc2) DescribeNetworkInterfaces(*ec2.DescribeNetworkInterfacesInput) (*ec2.DescribeNetworkInterfacesOutput, error) {
	panic("Not implemented")
}

func (m *MockEc2) DescribePlacementGroupsRequest(*ec2.DescribePlacementGroupsInput) (*request.Request, *ec2.DescribePlacementGroupsOutput) {
	panic("Not implemented")
}

func (m *MockEc2) DescribePlacementGroups(*ec2.DescribePlacementGroupsInput) (*ec2.DescribePlacementGroupsOutput, error) {
	panic("Not implemented")
}

func (m *MockEc2) DescribePrefixListsRequest(*ec2.DescribePrefixListsInput) (*request.Request, *ec2.DescribePrefixListsOutput) {
	panic("Not implemented")
}

func (m *MockEc2) DescribePrefixLists(*ec2.DescribePrefixListsInput) (*ec2.DescribePrefixListsOutput, error) {
	panic("Not implemented")
}

func (m *MockEc2) DescribeRegionsRequest(*ec2.DescribeRegionsInput) (*request.Request, *ec2.DescribeRegionsOutput) {
	panic("Not implemented")
}

func (m *MockEc2) DescribeRegions(*ec2.DescribeRegionsInput) (*ec2.DescribeRegionsOutput, error) {
	panic("Not implemented")
}

func (m *MockEc2) DescribeReservedInstancesRequest(*ec2.DescribeReservedInstancesInput) (*request.Request, *ec2.DescribeReservedInstancesOutput) {
	panic("Not implemented")
}

func (m *MockEc2) DescribeReservedInstances(*ec2.DescribeReservedInstancesInput) (*ec2.DescribeReservedInstancesOutput, error) {
	panic("Not implemented")
}

func (m *MockEc2) DescribeReservedInstancesListingsRequest(*ec2.DescribeReservedInstancesListingsInput) (*request.Request, *ec2.DescribeReservedInstancesListingsOutput) {
	panic("Not implemented")
}

func (m *MockEc2) DescribeReservedInstancesListings(*ec2.DescribeReservedInstancesListingsInput) (*ec2.DescribeReservedInstancesListingsOutput, error) {
	panic("Not implemented")
}

func (m *MockEc2) DescribeReservedInstancesModificationsRequest(*ec2.DescribeReservedInstancesModificationsInput) (*request.Request, *ec2.DescribeReservedInstancesModificationsOutput) {
	panic("Not implemented")
}

func (m *MockEc2) DescribeReservedInstancesModifications(*ec2.DescribeReservedInstancesModificationsInput) (*ec2.DescribeReservedInstancesModificationsOutput, error) {
	panic("Not implemented")
}

func (m *MockEc2) DescribeReservedInstancesModificationsPages(*ec2.DescribeReservedInstancesModificationsInput, func(*ec2.DescribeReservedInstancesModificationsOutput, bool) bool) error {
	panic("Not implemented")
}

func (m *MockEc2) DescribeReservedInstancesOfferingsRequest(*ec2.DescribeReservedInstancesOfferingsInput) (*request.Request, *ec2.DescribeReservedInstancesOfferingsOutput) {
	panic("Not implemented")
}

func (m *MockEc2) DescribeReservedInstancesOfferings(*ec2.DescribeReservedInstancesOfferingsInput) (*ec2.DescribeReservedInstancesOfferingsOutput, error) {
	panic("Not implemented")
}

func (m *MockEc2) DescribeReservedInstancesOfferingsPages(*ec2.DescribeReservedInstancesOfferingsInput, func(*ec2.DescribeReservedInstancesOfferingsOutput, bool) bool) error {
	panic("Not implemented")
}

func (m *MockEc2) DescribeRouteTablesRequest(*ec2.DescribeRouteTablesInput) (*request.Request, *ec2.DescribeRouteTablesOutput) {
	panic("Not implemented")
}

func (m *MockEc2) DescribeRouteTables(*ec2.DescribeRouteTablesInput) (*ec2.DescribeRouteTablesOutput, error) {
	panic("Not implemented")
}

func (m *MockEc2) DescribeScheduledInstanceAvailabilityRequest(*ec2.DescribeScheduledInstanceAvailabilityInput) (*request.Request, *ec2.DescribeScheduledInstanceAvailabilityOutput) {
	panic("Not implemented")
}

func (m *MockEc2) DescribeScheduledInstanceAvailability(*ec2.DescribeScheduledInstanceAvailabilityInput) (*ec2.DescribeScheduledInstanceAvailabilityOutput, error) {
	panic("Not implemented")
}

func (m *MockEc2) DescribeScheduledInstancesRequest(*ec2.DescribeScheduledInstancesInput) (*request.Request, *ec2.DescribeScheduledInstancesOutput) {
	panic("Not implemented")
}

func (m *MockEc2) DescribeScheduledInstances(*ec2.DescribeScheduledInstancesInput) (*ec2.DescribeScheduledInstancesOutput, error) {
	panic("Not implemented")
}

func (m *MockEc2) DescribeSecurityGroupReferencesRequest(*ec2.DescribeSecurityGroupReferencesInput) (*request.Request, *ec2.DescribeSecurityGroupReferencesOutput) {
	panic("Not implemented")
}

func (m *MockEc2) DescribeSecurityGroupReferences(*ec2.DescribeSecurityGroupReferencesInput) (*ec2.DescribeSecurityGroupReferencesOutput, error) {
	panic("Not implemented")
}

func (m *MockEc2) DescribeSecurityGroupsRequest(*ec2.DescribeSecurityGroupsInput) (*request.Request, *ec2.DescribeSecurityGroupsOutput) {
	panic("Not implemented")
}

func (m *MockEc2) DescribeSecurityGroups(*ec2.DescribeSecurityGroupsInput) (*ec2.DescribeSecurityGroupsOutput, error) {
	panic("Not implemented")
}

func (m *MockEc2) DescribeSnapshotAttributeRequest(*ec2.DescribeSnapshotAttributeInput) (*request.Request, *ec2.DescribeSnapshotAttributeOutput) {
	panic("Not implemented")
}

func (m *MockEc2) DescribeSnapshotAttribute(*ec2.DescribeSnapshotAttributeInput) (*ec2.DescribeSnapshotAttributeOutput, error) {
	panic("Not implemented")
}

func (m *MockEc2) DescribeSnapshotsRequest(*ec2.DescribeSnapshotsInput) (*request.Request, *ec2.DescribeSnapshotsOutput) {
	panic("Not implemented")
}

func (m *MockEc2) DescribeSnapshots(*ec2.DescribeSnapshotsInput) (*ec2.DescribeSnapshotsOutput, error) {
	panic("Not implemented")
}

func (m *MockEc2) DescribeSnapshotsPages(*ec2.DescribeSnapshotsInput, func(*ec2.DescribeSnapshotsOutput, bool) bool) error {
	panic("Not implemented")
}

func (m *MockEc2) DescribeSpotDatafeedSubscriptionRequest(*ec2.DescribeSpotDatafeedSubscriptionInput) (*request.Request, *ec2.DescribeSpotDatafeedSubscriptionOutput) {
	panic("Not implemented")
}

func (m *MockEc2) DescribeSpotDatafeedSubscription(*ec2.DescribeSpotDatafeedSubscriptionInput) (*ec2.DescribeSpotDatafeedSubscriptionOutput, error) {
	panic("Not implemented")
}

func (m *MockEc2) DescribeSpotFleetInstancesRequest(*ec2.DescribeSpotFleetInstancesInput) (*request.Request, *ec2.DescribeSpotFleetInstancesOutput) {
	panic("Not implemented")
}

func (m *MockEc2) DescribeSpotFleetInstances(*ec2.DescribeSpotFleetInstancesInput) (*ec2.DescribeSpotFleetInstancesOutput, error) {
	panic("Not implemented")
}

func (m *MockEc2) DescribeSpotFleetRequestHistoryRequest(*ec2.DescribeSpotFleetRequestHistoryInput) (*request.Request, *ec2.DescribeSpotFleetRequestHistoryOutput) {
	panic("Not implemented")
}

func (m *MockEc2) DescribeSpotFleetRequestHistory(*ec2.DescribeSpotFleetRequestHistoryInput) (*ec2.DescribeSpotFleetRequestHistoryOutput, error) {
	panic("Not implemented")
}

func (m *MockEc2) DescribeSpotFleetRequestsRequest(*ec2.DescribeSpotFleetRequestsInput) (*request.Request, *ec2.DescribeSpotFleetRequestsOutput) {
	panic("Not implemented")
}

func (m *MockEc2) DescribeSpotFleetRequests(*ec2.DescribeSpotFleetRequestsInput) (*ec2.DescribeSpotFleetRequestsOutput, error) {
	panic("Not implemented")
}

func (m *MockEc2) DescribeSpotFleetRequestsPages(*ec2.DescribeSpotFleetRequestsInput, func(*ec2.DescribeSpotFleetRequestsOutput, bool) bool) error {
	panic("Not implemented")
}

func (m *MockEc2) DescribeSpotInstanceRequestsRequest(*ec2.DescribeSpotInstanceRequestsInput) (*request.Request, *ec2.DescribeSpotInstanceRequestsOutput) {
	panic("Not implemented")
}

func (m *MockEc2) DescribeSpotInstanceRequests(*ec2.DescribeSpotInstanceRequestsInput) (*ec2.DescribeSpotInstanceRequestsOutput, error) {
	panic("Not implemented")
}

func (m *MockEc2) DescribeSpotPriceHistoryRequest(*ec2.DescribeSpotPriceHistoryInput) (*request.Request, *ec2.DescribeSpotPriceHistoryOutput) {
	panic("Not implemented")
}

func (m *MockEc2) DescribeSpotPriceHistory(*ec2.DescribeSpotPriceHistoryInput) (*ec2.DescribeSpotPriceHistoryOutput, error) {
	panic("Not implemented")
}

func (m *MockEc2) DescribeSpotPriceHistoryPages(*ec2.DescribeSpotPriceHistoryInput, func(*ec2.DescribeSpotPriceHistoryOutput, bool) bool) error {
	panic("Not implemented")
}

func (m *MockEc2) DescribeStaleSecurityGroupsRequest(*ec2.DescribeStaleSecurityGroupsInput) (*request.Request, *ec2.DescribeStaleSecurityGroupsOutput) {
	panic("Not implemented")
}

func (m *MockEc2) DescribeStaleSecurityGroups(*ec2.DescribeStaleSecurityGroupsInput) (*ec2.DescribeStaleSecurityGroupsOutput, error) {
	panic("Not implemented")
}

func (m *MockEc2) DescribeSubnetsRequest(*ec2.DescribeSubnetsInput) (*request.Request, *ec2.DescribeSubnetsOutput) {
	panic("Not implemented")
}

func (m *MockEc2) DescribeSubnets(*ec2.DescribeSubnetsInput) (*ec2.DescribeSubnetsOutput, error) {
	panic("Not implemented")
}

func (m *MockEc2) DescribeTagsRequest(*ec2.DescribeTagsInput) (*request.Request, *ec2.DescribeTagsOutput) {
	panic("Not implemented")
}

func (m *MockEc2) DescribeTags(*ec2.DescribeTagsInput) (*ec2.DescribeTagsOutput, error) {
	panic("Not implemented")
}

func (m *MockEc2) DescribeTagsPages(*ec2.DescribeTagsInput, func(*ec2.DescribeTagsOutput, bool) bool) error {
	panic("Not implemented")
}

func (m *MockEc2) DescribeVolumeAttributeRequest(*ec2.DescribeVolumeAttributeInput) (*request.Request, *ec2.DescribeVolumeAttributeOutput) {
	panic("Not implemented")
}

func (m *MockEc2) DescribeVolumeAttribute(*ec2.DescribeVolumeAttributeInput) (*ec2.DescribeVolumeAttributeOutput, error) {
	panic("Not implemented")
}

func (m *MockEc2) DescribeVolumeStatusRequest(*ec2.DescribeVolumeStatusInput) (*request.Request, *ec2.DescribeVolumeStatusOutput) {
	panic("Not implemented")
}

func (m *MockEc2) DescribeVolumeStatus(*ec2.DescribeVolumeStatusInput) (*ec2.DescribeVolumeStatusOutput, error) {
	panic("Not implemented")
}

func (m *MockEc2) DescribeVolumeStatusPages(*ec2.DescribeVolumeStatusInput, func(*ec2.DescribeVolumeStatusOutput, bool) bool) error {
	panic("Not implemented")
}

func (m *MockEc2) DescribeVolumesRequest(*ec2.DescribeVolumesInput) (*request.Request, *ec2.DescribeVolumesOutput) {
	panic("Not implemented")
}

func (m *MockEc2) DescribeVolumes(*ec2.DescribeVolumesInput) (*ec2.DescribeVolumesOutput, error) {
	panic("Not implemented")
}

func (m *MockEc2) DescribeVolumesPages(*ec2.DescribeVolumesInput, func(*ec2.DescribeVolumesOutput, bool) bool) error {
	panic("Not implemented")
}

func (m *MockEc2) DescribeVpcAttributeRequest(*ec2.DescribeVpcAttributeInput) (*request.Request, *ec2.DescribeVpcAttributeOutput) {
	panic("Not implemented")
}

func (m *MockEc2) DescribeVpcAttribute(*ec2.DescribeVpcAttributeInput) (*ec2.DescribeVpcAttributeOutput, error) {
	panic("Not implemented")
}

func (m *MockEc2) DescribeVpcClassicLinkRequest(*ec2.DescribeVpcClassicLinkInput) (*request.Request, *ec2.DescribeVpcClassicLinkOutput) {
	panic("Not implemented")
}

func (m *MockEc2) DescribeVpcClassicLink(*ec2.DescribeVpcClassicLinkInput) (*ec2.DescribeVpcClassicLinkOutput, error) {
	panic("Not implemented")
}

func (m *MockEc2) DescribeVpcClassicLinkDnsSupportRequest(*ec2.DescribeVpcClassicLinkDnsSupportInput) (*request.Request, *ec2.DescribeVpcClassicLinkDnsSupportOutput) {
	panic("Not implemented")
}

func (m *MockEc2) DescribeVpcClassicLinkDnsSupport(*ec2.DescribeVpcClassicLinkDnsSupportInput) (*ec2.DescribeVpcClassicLinkDnsSupportOutput, error) {
	panic("Not implemented")
}

func (m *MockEc2) DescribeVpcEndpointServicesRequest(*ec2.DescribeVpcEndpointServicesInput) (*request.Request, *ec2.DescribeVpcEndpointServicesOutput) {
	panic("Not implemented")
}

func (m *MockEc2) DescribeVpcEndpointServices(*ec2.DescribeVpcEndpointServicesInput) (*ec2.DescribeVpcEndpointServicesOutput, error) {
	panic("Not implemented")
}

func (m *MockEc2) DescribeVpcEndpointsRequest(*ec2.DescribeVpcEndpointsInput) (*request.Request, *ec2.DescribeVpcEndpointsOutput) {
	panic("Not implemented")
}

func (m *MockEc2) DescribeVpcEndpoints(*ec2.DescribeVpcEndpointsInput) (*ec2.DescribeVpcEndpointsOutput, error) {
	panic("Not implemented")
}

func (m *MockEc2) DescribeVpcPeeringConnectionsRequest(*ec2.DescribeVpcPeeringConnectionsInput) (*request.Request, *ec2.DescribeVpcPeeringConnectionsOutput) {
	panic("Not implemented")
}

func (m *MockEc2) DescribeVpcPeeringConnections(*ec2.DescribeVpcPeeringConnectionsInput) (*ec2.DescribeVpcPeeringConnectionsOutput, error) {
	panic("Not implemented")
}

func (m *MockEc2) DescribeVpcsRequest(*ec2.DescribeVpcsInput) (*request.Request, *ec2.DescribeVpcsOutput) {
	panic("Not implemented")
}

func (m *MockEc2) DescribeVpcs(*ec2.DescribeVpcsInput) (*ec2.DescribeVpcsOutput, error) {
	panic("Not implemented")
}

func (m *MockEc2) DescribeVpnConnectionsRequest(*ec2.DescribeVpnConnectionsInput) (*request.Request, *ec2.DescribeVpnConnectionsOutput) {
	panic("Not implemented")
}

func (m *MockEc2) DescribeVpnConnections(*ec2.DescribeVpnConnectionsInput) (*ec2.DescribeVpnConnectionsOutput, error) {
	panic("Not implemented")
}

func (m *MockEc2) DescribeVpnGatewaysRequest(*ec2.DescribeVpnGatewaysInput) (*request.Request, *ec2.DescribeVpnGatewaysOutput) {
	panic("Not implemented")
}

func (m *MockEc2) DescribeVpnGateways(*ec2.DescribeVpnGatewaysInput) (*ec2.DescribeVpnGatewaysOutput, error) {
	panic("Not implemented")
}

func (m *MockEc2) DetachClassicLinkVpcRequest(*ec2.DetachClassicLinkVpcInput) (*request.Request, *ec2.DetachClassicLinkVpcOutput) {
	panic("Not implemented")
}

func (m *MockEc2) DetachClassicLinkVpc(*ec2.DetachClassicLinkVpcInput) (*ec2.DetachClassicLinkVpcOutput, error) {
	panic("Not implemented")
}

func (m *MockEc2) DetachInternetGatewayRequest(*ec2.DetachInternetGatewayInput) (*request.Request, *ec2.DetachInternetGatewayOutput) {
	panic("Not implemented")
}

func (m *MockEc2) DetachInternetGateway(*ec2.DetachInternetGatewayInput) (*ec2.DetachInternetGatewayOutput, error) {
	panic("Not implemented")
}

func (m *MockEc2) DetachNetworkInterfaceRequest(*ec2.DetachNetworkInterfaceInput) (*request.Request, *ec2.DetachNetworkInterfaceOutput) {
	panic("Not implemented")
}

func (m *MockEc2) DetachNetworkInterface(*ec2.DetachNetworkInterfaceInput) (*ec2.DetachNetworkInterfaceOutput, error) {
	panic("Not implemented")
}

func (m *MockEc2) DetachVolumeRequest(*ec2.DetachVolumeInput) (*request.Request, *ec2.VolumeAttachment) {
	panic("Not implemented")
}

func (m *MockEc2) DetachVolume(*ec2.DetachVolumeInput) (*ec2.VolumeAttachment, error) {
	panic("Not implemented")
}

func (m *MockEc2) DetachVpnGatewayRequest(*ec2.DetachVpnGatewayInput) (*request.Request, *ec2.DetachVpnGatewayOutput) {
	panic("Not implemented")
}

func (m *MockEc2) DetachVpnGateway(*ec2.DetachVpnGatewayInput) (*ec2.DetachVpnGatewayOutput, error) {
	panic("Not implemented")
}

func (m *MockEc2) DisableVgwRoutePropagationRequest(*ec2.DisableVgwRoutePropagationInput) (*request.Request, *ec2.DisableVgwRoutePropagationOutput) {
	panic("Not implemented")
}

func (m *MockEc2) DisableVgwRoutePropagation(*ec2.DisableVgwRoutePropagationInput) (*ec2.DisableVgwRoutePropagationOutput, error) {
	panic("Not implemented")
}

func (m *MockEc2) DisableVpcClassicLinkRequest(*ec2.DisableVpcClassicLinkInput) (*request.Request, *ec2.DisableVpcClassicLinkOutput) {
	panic("Not implemented")
}

func (m *MockEc2) DisableVpcClassicLink(*ec2.DisableVpcClassicLinkInput) (*ec2.DisableVpcClassicLinkOutput, error) {
	panic("Not implemented")
}

func (m *MockEc2) DisableVpcClassicLinkDnsSupportRequest(*ec2.DisableVpcClassicLinkDnsSupportInput) (*request.Request, *ec2.DisableVpcClassicLinkDnsSupportOutput) {
	panic("Not implemented")
}

func (m *MockEc2) DisableVpcClassicLinkDnsSupport(*ec2.DisableVpcClassicLinkDnsSupportInput) (*ec2.DisableVpcClassicLinkDnsSupportOutput, error) {
	panic("Not implemented")
}

func (m *MockEc2) DisassociateAddressRequest(*ec2.DisassociateAddressInput) (*request.Request, *ec2.DisassociateAddressOutput) {
	panic("Not implemented")
}

func (m *MockEc2) DisassociateAddress(*ec2.DisassociateAddressInput) (*ec2.DisassociateAddressOutput, error) {
	panic("Not implemented")
}

func (m *MockEc2) DisassociateRouteTableRequest(*ec2.DisassociateRouteTableInput) (*request.Request, *ec2.DisassociateRouteTableOutput) {
	panic("Not implemented")
}

func (m *MockEc2) DisassociateRouteTable(*ec2.DisassociateRouteTableInput) (*ec2.DisassociateRouteTableOutput, error) {
	panic("Not implemented")
}

func (m *MockEc2) EnableVgwRoutePropagationRequest(*ec2.EnableVgwRoutePropagationInput) (*request.Request, *ec2.EnableVgwRoutePropagationOutput) {
	panic("Not implemented")
}

func (m *MockEc2) EnableVgwRoutePropagation(*ec2.EnableVgwRoutePropagationInput) (*ec2.EnableVgwRoutePropagationOutput, error) {
	panic("Not implemented")
}

func (m *MockEc2) EnableVolumeIORequest(*ec2.EnableVolumeIOInput) (*request.Request, *ec2.EnableVolumeIOOutput) {
	panic("Not implemented")
}

func (m *MockEc2) EnableVolumeIO(*ec2.EnableVolumeIOInput) (*ec2.EnableVolumeIOOutput, error) {
	panic("Not implemented")
}

func (m *MockEc2) EnableVpcClassicLinkRequest(*ec2.EnableVpcClassicLinkInput) (*request.Request, *ec2.EnableVpcClassicLinkOutput) {
	panic("Not implemented")
}

func (m *MockEc2) EnableVpcClassicLink(*ec2.EnableVpcClassicLinkInput) (*ec2.EnableVpcClassicLinkOutput, error) {
	panic("Not implemented")
}

func (m *MockEc2) EnableVpcClassicLinkDnsSupportRequest(*ec2.EnableVpcClassicLinkDnsSupportInput) (*request.Request, *ec2.EnableVpcClassicLinkDnsSupportOutput) {
	panic("Not implemented")
}

func (m *MockEc2) EnableVpcClassicLinkDnsSupport(*ec2.EnableVpcClassicLinkDnsSupportInput) (*ec2.EnableVpcClassicLinkDnsSupportOutput, error) {
	panic("Not implemented")
}

func (m *MockEc2) GetConsoleOutputRequest(*ec2.GetConsoleOutputInput) (*request.Request, *ec2.GetConsoleOutputOutput) {
	panic("Not implemented")
}

func (m *MockEc2) GetConsoleOutput(*ec2.GetConsoleOutputInput) (*ec2.GetConsoleOutputOutput, error) {
	panic("Not implemented")
}

func (m *MockEc2) GetConsoleScreenshotRequest(*ec2.GetConsoleScreenshotInput) (*request.Request, *ec2.GetConsoleScreenshotOutput) {
	panic("Not implemented")
}

func (m *MockEc2) GetConsoleScreenshot(*ec2.GetConsoleScreenshotInput) (*ec2.GetConsoleScreenshotOutput, error) {
	panic("Not implemented")
}

func (m *MockEc2) GetHostReservationPurchasePreviewRequest(*ec2.GetHostReservationPurchasePreviewInput) (*request.Request, *ec2.GetHostReservationPurchasePreviewOutput) {
	panic("Not implemented")
}

func (m *MockEc2) GetHostReservationPurchasePreview(*ec2.GetHostReservationPurchasePreviewInput) (*ec2.GetHostReservationPurchasePreviewOutput, error) {
	panic("Not implemented")
}

func (m *MockEc2) GetPasswordDataRequest(*ec2.GetPasswordDataInput) (*request.Request, *ec2.GetPasswordDataOutput) {
	panic("Not implemented")
}

func (m *MockEc2) GetPasswordData(*ec2.GetPasswordDataInput) (*ec2.GetPasswordDataOutput, error) {
	panic("Not implemented")
}

func (m *MockEc2) ImportImageRequest(*ec2.ImportImageInput) (*request.Request, *ec2.ImportImageOutput) {
	panic("Not implemented")
}

func (m *MockEc2) ImportImage(*ec2.ImportImageInput) (*ec2.ImportImageOutput, error) {
	panic("Not implemented")
}

func (m *MockEc2) ImportInstanceRequest(*ec2.ImportInstanceInput) (*request.Request, *ec2.ImportInstanceOutput) {
	panic("Not implemented")
}

func (m *MockEc2) ImportInstance(*ec2.ImportInstanceInput) (*ec2.ImportInstanceOutput, error) {
	panic("Not implemented")
}

func (m *MockEc2) ImportKeyPairRequest(*ec2.ImportKeyPairInput) (*request.Request, *ec2.ImportKeyPairOutput) {
	panic("Not implemented")
}

func (m *MockEc2) ImportKeyPair(*ec2.ImportKeyPairInput) (*ec2.ImportKeyPairOutput, error) {
	panic("Not implemented")
}

func (m *MockEc2) ImportSnapshotRequest(*ec2.ImportSnapshotInput) (*request.Request, *ec2.ImportSnapshotOutput) {
	panic("Not implemented")
}

func (m *MockEc2) ImportSnapshot(*ec2.ImportSnapshotInput) (*ec2.ImportSnapshotOutput, error) {
	panic("Not implemented")
}

func (m *MockEc2) ImportVolumeRequest(*ec2.ImportVolumeInput) (*request.Request, *ec2.ImportVolumeOutput) {
	panic("Not implemented")
}

func (m *MockEc2) ImportVolume(*ec2.ImportVolumeInput) (*ec2.ImportVolumeOutput, error) {
	panic("Not implemented")
}

func (m *MockEc2) ModifyHostsRequest(*ec2.ModifyHostsInput) (*request.Request, *ec2.ModifyHostsOutput) {
	panic("Not implemented")
}

func (m *MockEc2) ModifyHosts(*ec2.ModifyHostsInput) (*ec2.ModifyHostsOutput, error) {
	panic("Not implemented")
}

func (m *MockEc2) ModifyIdFormatRequest(*ec2.ModifyIdFormatInput) (*request.Request, *ec2.ModifyIdFormatOutput) {
	panic("Not implemented")
}

func (m *MockEc2) ModifyIdFormat(*ec2.ModifyIdFormatInput) (*ec2.ModifyIdFormatOutput, error) {
	panic("Not implemented")
}

func (m *MockEc2) ModifyIdentityIdFormatRequest(*ec2.ModifyIdentityIdFormatInput) (*request.Request, *ec2.ModifyIdentityIdFormatOutput) {
	panic("Not implemented")
}

func (m *MockEc2) ModifyIdentityIdFormat(*ec2.ModifyIdentityIdFormatInput) (*ec2.ModifyIdentityIdFormatOutput, error) {
	panic("Not implemented")
}

func (m *MockEc2) ModifyImageAttributeRequest(*ec2.ModifyImageAttributeInput) (*request.Request, *ec2.ModifyImageAttributeOutput) {
	panic("Not implemented")
}

func (m *MockEc2) ModifyImageAttribute(*ec2.ModifyImageAttributeInput) (*ec2.ModifyImageAttributeOutput, error) {
	panic("Not implemented")
}

func (m *MockEc2) ModifyInstanceAttributeRequest(*ec2.ModifyInstanceAttributeInput) (*request.Request, *ec2.ModifyInstanceAttributeOutput) {
	panic("Not implemented")
}

func (m *MockEc2) ModifyInstanceAttribute(*ec2.ModifyInstanceAttributeInput) (*ec2.ModifyInstanceAttributeOutput, error) {
	panic("Not implemented")
}

func (m *MockEc2) ModifyInstancePlacementRequest(*ec2.ModifyInstancePlacementInput) (*request.Request, *ec2.ModifyInstancePlacementOutput) {
	panic("Not implemented")
}

func (m *MockEc2) ModifyInstancePlacement(*ec2.ModifyInstancePlacementInput) (*ec2.ModifyInstancePlacementOutput, error) {
	panic("Not implemented")
}

func (m *MockEc2) ModifyNetworkInterfaceAttributeRequest(*ec2.ModifyNetworkInterfaceAttributeInput) (*request.Request, *ec2.ModifyNetworkInterfaceAttributeOutput) {
	panic("Not implemented")
}

func (m *MockEc2) ModifyNetworkInterfaceAttribute(*ec2.ModifyNetworkInterfaceAttributeInput) (*ec2.ModifyNetworkInterfaceAttributeOutput, error) {
	panic("Not implemented")
}

func (m *MockEc2) ModifyReservedInstancesRequest(*ec2.ModifyReservedInstancesInput) (*request.Request, *ec2.ModifyReservedInstancesOutput) {
	panic("Not implemented")
}

func (m *MockEc2) ModifyReservedInstances(*ec2.ModifyReservedInstancesInput) (*ec2.ModifyReservedInstancesOutput, error) {
	panic("Not implemented")
}

func (m *MockEc2) ModifySnapshotAttributeRequest(*ec2.ModifySnapshotAttributeInput) (*request.Request, *ec2.ModifySnapshotAttributeOutput) {
	panic("Not implemented")
}

func (m *MockEc2) ModifySnapshotAttribute(*ec2.ModifySnapshotAttributeInput) (*ec2.ModifySnapshotAttributeOutput, error) {
	panic("Not implemented")
}

func (m *MockEc2) ModifySpotFleetRequestRequest(*ec2.ModifySpotFleetRequestInput) (*request.Request, *ec2.ModifySpotFleetRequestOutput) {
	panic("Not implemented")
}

func (m *MockEc2) ModifySpotFleetRequest(*ec2.ModifySpotFleetRequestInput) (*ec2.ModifySpotFleetRequestOutput, error) {
	panic("Not implemented")
}

func (m *MockEc2) ModifySubnetAttributeRequest(*ec2.ModifySubnetAttributeInput) (*request.Request, *ec2.ModifySubnetAttributeOutput) {
	panic("Not implemented")
}

func (m *MockEc2) ModifySubnetAttribute(*ec2.ModifySubnetAttributeInput) (*ec2.ModifySubnetAttributeOutput, error) {
	panic("Not implemented")
}

func (m *MockEc2) ModifyVolumeAttributeRequest(*ec2.ModifyVolumeAttributeInput) (*request.Request, *ec2.ModifyVolumeAttributeOutput) {
	panic("Not implemented")
}

func (m *MockEc2) ModifyVolumeAttribute(*ec2.ModifyVolumeAttributeInput) (*ec2.ModifyVolumeAttributeOutput, error) {
	panic("Not implemented")
}

func (m *MockEc2) ModifyVpcAttributeRequest(*ec2.ModifyVpcAttributeInput) (*request.Request, *ec2.ModifyVpcAttributeOutput) {
	panic("Not implemented")
}

func (m *MockEc2) ModifyVpcAttribute(*ec2.ModifyVpcAttributeInput) (*ec2.ModifyVpcAttributeOutput, error) {
	panic("Not implemented")
}

func (m *MockEc2) ModifyVpcEndpointRequest(*ec2.ModifyVpcEndpointInput) (*request.Request, *ec2.ModifyVpcEndpointOutput) {
	panic("Not implemented")
}

func (m *MockEc2) ModifyVpcEndpoint(*ec2.ModifyVpcEndpointInput) (*ec2.ModifyVpcEndpointOutput, error) {
	panic("Not implemented")
}

func (m *MockEc2) ModifyVpcPeeringConnectionOptionsRequest(*ec2.ModifyVpcPeeringConnectionOptionsInput) (*request.Request, *ec2.ModifyVpcPeeringConnectionOptionsOutput) {
	panic("Not implemented")
}

func (m *MockEc2) ModifyVpcPeeringConnectionOptions(*ec2.ModifyVpcPeeringConnectionOptionsInput) (*ec2.ModifyVpcPeeringConnectionOptionsOutput, error) {
	panic("Not implemented")
}

func (m *MockEc2) MonitorInstancesRequest(*ec2.MonitorInstancesInput) (*request.Request, *ec2.MonitorInstancesOutput) {
	panic("Not implemented")
}

func (m *MockEc2) MonitorInstances(*ec2.MonitorInstancesInput) (*ec2.MonitorInstancesOutput, error) {
	panic("Not implemented")
}

func (m *MockEc2) MoveAddressToVpcRequest(*ec2.MoveAddressToVpcInput) (*request.Request, *ec2.MoveAddressToVpcOutput) {
	panic("Not implemented")
}

func (m *MockEc2) MoveAddressToVpc(*ec2.MoveAddressToVpcInput) (*ec2.MoveAddressToVpcOutput, error) {
	panic("Not implemented")
}

func (m *MockEc2) PurchaseHostReservationRequest(*ec2.PurchaseHostReservationInput) (*request.Request, *ec2.PurchaseHostReservationOutput) {
	panic("Not implemented")
}

func (m *MockEc2) PurchaseHostReservation(*ec2.PurchaseHostReservationInput) (*ec2.PurchaseHostReservationOutput, error) {
	panic("Not implemented")
}

func (m *MockEc2) PurchaseReservedInstancesOfferingRequest(*ec2.PurchaseReservedInstancesOfferingInput) (*request.Request, *ec2.PurchaseReservedInstancesOfferingOutput) {
	panic("Not implemented")
}

func (m *MockEc2) PurchaseReservedInstancesOffering(*ec2.PurchaseReservedInstancesOfferingInput) (*ec2.PurchaseReservedInstancesOfferingOutput, error) {
	panic("Not implemented")
}

func (m *MockEc2) PurchaseScheduledInstancesRequest(*ec2.PurchaseScheduledInstancesInput) (*request.Request, *ec2.PurchaseScheduledInstancesOutput) {
	panic("Not implemented")
}

func (m *MockEc2) PurchaseScheduledInstances(*ec2.PurchaseScheduledInstancesInput) (*ec2.PurchaseScheduledInstancesOutput, error) {
	panic("Not implemented")
}

func (m *MockEc2) RebootInstancesRequest(*ec2.RebootInstancesInput) (*request.Request, *ec2.RebootInstancesOutput) {
	panic("Not implemented")
}

func (m *MockEc2) RebootInstances(*ec2.RebootInstancesInput) (*ec2.RebootInstancesOutput, error) {
	panic("Not implemented")
}

func (m *MockEc2) RegisterImageRequest(*ec2.RegisterImageInput) (*request.Request, *ec2.RegisterImageOutput) {
	panic("Not implemented")
}

func (m *MockEc2) RegisterImage(*ec2.RegisterImageInput) (*ec2.RegisterImageOutput, error) {
	panic("Not implemented")
}

func (m *MockEc2) RejectVpcPeeringConnectionRequest(*ec2.RejectVpcPeeringConnectionInput) (*request.Request, *ec2.RejectVpcPeeringConnectionOutput) {
	panic("Not implemented")
}

func (m *MockEc2) RejectVpcPeeringConnection(*ec2.RejectVpcPeeringConnectionInput) (*ec2.RejectVpcPeeringConnectionOutput, error) {
	panic("Not implemented")
}

func (m *MockEc2) ReleaseAddressRequest(*ec2.ReleaseAddressInput) (*request.Request, *ec2.ReleaseAddressOutput) {
	panic("Not implemented")
}

func (m *MockEc2) ReleaseAddress(*ec2.ReleaseAddressInput) (*ec2.ReleaseAddressOutput, error) {
	panic("Not implemented")
}

func (m *MockEc2) ReleaseHostsRequest(*ec2.ReleaseHostsInput) (*request.Request, *ec2.ReleaseHostsOutput) {
	panic("Not implemented")
}

func (m *MockEc2) ReleaseHosts(*ec2.ReleaseHostsInput) (*ec2.ReleaseHostsOutput, error) {
	panic("Not implemented")
}

func (m *MockEc2) ReplaceNetworkAclAssociationRequest(*ec2.ReplaceNetworkAclAssociationInput) (*request.Request, *ec2.ReplaceNetworkAclAssociationOutput) {
	panic("Not implemented")
}

func (m *MockEc2) ReplaceNetworkAclAssociation(*ec2.ReplaceNetworkAclAssociationInput) (*ec2.ReplaceNetworkAclAssociationOutput, error) {
	panic("Not implemented")
}

func (m *MockEc2) ReplaceNetworkAclEntryRequest(*ec2.ReplaceNetworkAclEntryInput) (*request.Request, *ec2.ReplaceNetworkAclEntryOutput) {
	panic("Not implemented")
}

func (m *MockEc2) ReplaceNetworkAclEntry(*ec2.ReplaceNetworkAclEntryInput) (*ec2.ReplaceNetworkAclEntryOutput, error) {
	panic("Not implemented")
}

func (m *MockEc2) ReplaceRouteRequest(*ec2.ReplaceRouteInput) (*request.Request, *ec2.ReplaceRouteOutput) {
	panic("Not implemented")
}

func (m *MockEc2) ReplaceRoute(*ec2.ReplaceRouteInput) (*ec2.ReplaceRouteOutput, error) {
	panic("Not implemented")
}

func (m *MockEc2) ReplaceRouteTableAssociationRequest(*ec2.ReplaceRouteTableAssociationInput) (*request.Request, *ec2.ReplaceRouteTableAssociationOutput) {
	panic("Not implemented")
}

func (m *MockEc2) ReplaceRouteTableAssociation(*ec2.ReplaceRouteTableAssociationInput) (*ec2.ReplaceRouteTableAssociationOutput, error) {
	panic("Not implemented")
}

func (m *MockEc2) ReportInstanceStatusRequest(*ec2.ReportInstanceStatusInput) (*request.Request, *ec2.ReportInstanceStatusOutput) {
	panic("Not implemented")
}

func (m *MockEc2) ReportInstanceStatus(*ec2.ReportInstanceStatusInput) (*ec2.ReportInstanceStatusOutput, error) {
	panic("Not implemented")
}

func (m *MockEc2) RequestSpotFleetRequest(*ec2.RequestSpotFleetInput) (*request.Request, *ec2.RequestSpotFleetOutput) {
	panic("Not implemented")
}

func (m *MockEc2) RequestSpotFleet(*ec2.RequestSpotFleetInput) (*ec2.RequestSpotFleetOutput, error) {
	panic("Not implemented")
}

func (m *MockEc2) RequestSpotInstancesRequest(*ec2.RequestSpotInstancesInput) (*request.Request, *ec2.RequestSpotInstancesOutput) {
	panic("Not implemented")
}

func (m *MockEc2) RequestSpotInstances(*ec2.RequestSpotInstancesInput) (*ec2.RequestSpotInstancesOutput, error) {
	panic("Not implemented")
}

func (m *MockEc2) ResetImageAttributeRequest(*ec2.ResetImageAttributeInput) (*request.Request, *ec2.ResetImageAttributeOutput) {
	panic("Not implemented")
}

func (m *MockEc2) ResetImageAttribute(*ec2.ResetImageAttributeInput) (*ec2.ResetImageAttributeOutput, error) {
	panic("Not implemented")
}

func (m *MockEc2) ResetInstanceAttributeRequest(*ec2.ResetInstanceAttributeInput) (*request.Request, *ec2.ResetInstanceAttributeOutput) {
	panic("Not implemented")
}

func (m *MockEc2) ResetInstanceAttribute(*ec2.ResetInstanceAttributeInput) (*ec2.ResetInstanceAttributeOutput, error) {
	panic("Not implemented")
}

func (m *MockEc2) ResetNetworkInterfaceAttributeRequest(*ec2.ResetNetworkInterfaceAttributeInput) (*request.Request, *ec2.ResetNetworkInterfaceAttributeOutput) {
	panic("Not implemented")
}

func (m *MockEc2) ResetNetworkInterfaceAttribute(*ec2.ResetNetworkInterfaceAttributeInput) (*ec2.ResetNetworkInterfaceAttributeOutput, error) {
	panic("Not implemented")
}

func (m *MockEc2) ResetSnapshotAttributeRequest(*ec2.ResetSnapshotAttributeInput) (*request.Request, *ec2.ResetSnapshotAttributeOutput) {
	panic("Not implemented")
}

func (m *MockEc2) ResetSnapshotAttribute(*ec2.ResetSnapshotAttributeInput) (*ec2.ResetSnapshotAttributeOutput, error) {
	panic("Not implemented")
}

func (m *MockEc2) RestoreAddressToClassicRequest(*ec2.RestoreAddressToClassicInput) (*request.Request, *ec2.RestoreAddressToClassicOutput) {
	panic("Not implemented")
}

func (m *MockEc2) RestoreAddressToClassic(*ec2.RestoreAddressToClassicInput) (*ec2.RestoreAddressToClassicOutput, error) {
	panic("Not implemented")
}

func (m *MockEc2) RevokeSecurityGroupEgressRequest(*ec2.RevokeSecurityGroupEgressInput) (*request.Request, *ec2.RevokeSecurityGroupEgressOutput) {
	panic("Not implemented")
}

func (m *MockEc2) RevokeSecurityGroupEgress(*ec2.RevokeSecurityGroupEgressInput) (*ec2.RevokeSecurityGroupEgressOutput, error) {
	panic("Not implemented")
}

func (m *MockEc2) RevokeSecurityGroupIngressRequest(*ec2.RevokeSecurityGroupIngressInput) (*request.Request, *ec2.RevokeSecurityGroupIngressOutput) {
	panic("Not implemented")
}

func (m *MockEc2) RevokeSecurityGroupIngress(*ec2.RevokeSecurityGroupIngressInput) (*ec2.RevokeSecurityGroupIngressOutput, error) {
	panic("Not implemented")
}

func (m *MockEc2) RunInstancesRequest(*ec2.RunInstancesInput) (*request.Request, *ec2.Reservation) {
	panic("Not implemented")
}

func (m *MockEc2) RunInstances(*ec2.RunInstancesInput) (*ec2.Reservation, error) {
	panic("Not implemented")
}

func (m *MockEc2) RunScheduledInstancesRequest(*ec2.RunScheduledInstancesInput) (*request.Request, *ec2.RunScheduledInstancesOutput) {
	panic("Not implemented")
}

func (m *MockEc2) RunScheduledInstances(*ec2.RunScheduledInstancesInput) (*ec2.RunScheduledInstancesOutput, error) {
	panic("Not implemented")
}

func (m *MockEc2) StartInstancesRequest(*ec2.StartInstancesInput) (*request.Request, *ec2.StartInstancesOutput) {
	panic("Not implemented")
}

func (m *MockEc2) StartInstances(*ec2.StartInstancesInput) (*ec2.StartInstancesOutput, error) {
	panic("Not implemented")
}

func (m *MockEc2) StopInstancesRequest(*ec2.StopInstancesInput) (*request.Request, *ec2.StopInstancesOutput) {
	panic("Not implemented")
}

func (m *MockEc2) StopInstances(*ec2.StopInstancesInput) (*ec2.StopInstancesOutput, error) {
	panic("Not implemented")
}

func (m *MockEc2) TerminateInstancesRequest(*ec2.TerminateInstancesInput) (*request.Request, *ec2.TerminateInstancesOutput) {
	panic("Not implemented")
}

func (m *MockEc2) TerminateInstances(*ec2.TerminateInstancesInput) (*ec2.TerminateInstancesOutput, error) {
	panic("Not implemented")
}

func (m *MockEc2) UnassignPrivateIpAddressesRequest(*ec2.UnassignPrivateIpAddressesInput) (*request.Request, *ec2.UnassignPrivateIpAddressesOutput) {
	panic("Not implemented")
}

func (m *MockEc2) UnassignPrivateIpAddresses(*ec2.UnassignPrivateIpAddressesInput) (*ec2.UnassignPrivateIpAddressesOutput, error) {
	panic("Not implemented")
}

func (m *MockEc2) UnmonitorInstancesRequest(*ec2.UnmonitorInstancesInput) (*request.Request, *ec2.UnmonitorInstancesOutput) {
	panic("Not implemented")
}

func (m *MockEc2) UnmonitorInstances(*ec2.UnmonitorInstancesInput) (*ec2.UnmonitorInstancesOutput, error) {
	panic("Not implemented")
}

func (m *MockEc2) WaitUntilBundleTaskComplete(*ec2.DescribeBundleTasksInput) error {
	panic("Not implemented")
}

func (m *MockEc2) WaitUntilConversionTaskCancelled(*ec2.DescribeConversionTasksInput) error {
	panic("Not implemented")
}

func (m *MockEc2) WaitUntilConversionTaskCompleted(*ec2.DescribeConversionTasksInput) error {
	panic("Not implemented")
}

func (m *MockEc2) WaitUntilConversionTaskDeleted(*ec2.DescribeConversionTasksInput) error {
	panic("Not implemented")
}

func (m *MockEc2) WaitUntilCustomerGatewayAvailable(*ec2.DescribeCustomerGatewaysInput) error {
	panic("Not implemented")
}

func (m *MockEc2) WaitUntilExportTaskCancelled(*ec2.DescribeExportTasksInput) error {
	panic("Not implemented")
}

func (m *MockEc2) WaitUntilExportTaskCompleted(*ec2.DescribeExportTasksInput) error {
	panic("Not implemented")
}

func (m *MockEc2) WaitUntilImageAvailable(*ec2.DescribeImagesInput) error {
	panic("Not implemented")
}

func (m *MockEc2) WaitUntilImageExists(*ec2.DescribeImagesInput) error {
	panic("Not implemented")
}

func (m *MockEc2) WaitUntilInstanceExists(*ec2.DescribeInstancesInput) error {
	panic("Not implemented")
}

func (m *MockEc2) WaitUntilInstanceRunning(*ec2.DescribeInstancesInput) error {
	panic("Not implemented")
}

func (m *MockEc2) WaitUntilInstanceStatusOk(*ec2.DescribeInstanceStatusInput) error {
	panic("Not implemented")
}

func (m *MockEc2) WaitUntilInstanceStopped(*ec2.DescribeInstancesInput) error {
	panic("Not implemented")
}

func (m *MockEc2) WaitUntilInstanceTerminated(*ec2.DescribeInstancesInput) error {
	panic("Not implemented")
}

func (m *MockEc2) WaitUntilKeyPairExists(*ec2.DescribeKeyPairsInput) error {
	panic("Not implemented")
}

func (m *MockEc2) WaitUntilNatGatewayAvailable(*ec2.DescribeNatGatewaysInput) error {
	panic("Not implemented")
}

func (m *MockEc2) WaitUntilNetworkInterfaceAvailable(*ec2.DescribeNetworkInterfacesInput) error {
	panic("Not implemented")
}

func (m *MockEc2) WaitUntilPasswordDataAvailable(*ec2.GetPasswordDataInput) error {
	panic("Not implemented")
}

func (m *MockEc2) WaitUntilSnapshotCompleted(*ec2.DescribeSnapshotsInput) error {
	panic("Not implemented")
}

func (m *MockEc2) WaitUntilSpotInstanceRequestFulfilled(*ec2.DescribeSpotInstanceRequestsInput) error {
	panic("Not implemented")
}

func (m *MockEc2) WaitUntilSubnetAvailable(*ec2.DescribeSubnetsInput) error {
	panic("Not implemented")
}

func (m *MockEc2) WaitUntilSystemStatusOk(*ec2.DescribeInstanceStatusInput) error {
	panic("Not implemented")
}

func (m *MockEc2) WaitUntilVolumeAvailable(*ec2.DescribeVolumesInput) error {
	panic("Not implemented")
}

func (m *MockEc2) WaitUntilVolumeDeleted(*ec2.DescribeVolumesInput) error {
	panic("Not implemented")
}

func (m *MockEc2) WaitUntilVolumeInUse(*ec2.DescribeVolumesInput) error {
	panic("Not implemented")
}

func (m *MockEc2) WaitUntilVpcAvailable(*ec2.DescribeVpcsInput) error {
	panic("Not implemented")
}

func (m *MockEc2) WaitUntilVpcExists(*ec2.DescribeVpcsInput) error {
	panic("Not implemented")
}

func (m *MockEc2) WaitUntilVpcPeeringConnectionExists(*ec2.DescribeVpcPeeringConnectionsInput) error {
	panic("Not implemented")
}

func (m *MockEc2) WaitUntilVpnConnectionAvailable(*ec2.DescribeVpnConnectionsInput) error {
	panic("Not implemented")
}

func (m *MockEc2) WaitUntilVpnConnectionDeleted(*ec2.DescribeVpnConnectionsInput) error {
	panic("Not implemented")
}
