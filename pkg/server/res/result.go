package res

type VmDetail struct {
	ClusterId string `json:"clusterId,omitempty"`
	VmId   string `json:"vmId,omitempty"`
	VmName string `json:"vmName,omitempty"`
	Status string `json:"status,omitempty"`
	Uid string `json:"uid,omitempty"`
	CreationTime string `json:"creationTime,omitempty"`
}

type ImageDetail struct {
	ClusterId string `json:"clusterId,omitempty"`
	Uid string `json:"uid,omitempty"`
	ImageId   string `json:"imageId,omitempty"`
	ImageName string `json:"imageName,omitempty"`
	Status string `json:"status,omitempty"`
}

type InstanceDetail struct {
	ClusterId string `json:"clusterId,omitempty"`
	Uid string `json:"uid,omitempty"`
	InstanceId   string `json:"instanceId,omitempty"`
	InstanceName string `json:"instanceName,omitempty"`
	Status string `json:"status,omitempty"`
}

type NetworkDetail struct {
	ClusterId string `json:"clusterId,omitempty"`
	Uid string `json:"uid,omitempty"`
	NetworkId   string `json:"networkId,omitempty"`
	NetworkName string `json:"networkName,omitempty"`
	Status string `json:"status,omitempty"`
}

type TemplateDetail struct {
	ClusterId string `json:"clusterId,omitempty"`
	Uid string `json:"uid,omitempty"`
	TemplateId   string `json:"templateId,omitempty"`
	TemplateName string `json:"templateName,omitempty"`
	Status string `json:"status,omitempty"`
}

type VolumeDetail struct {
	ClusterId string `json:"clusterId,omitempty"`
	Uid string `json:"uid,omitempty"`
	VolumeId   string `json:"volumeId,omitempty"`
	VolumeName string `json:"volumeName,omitempty"`
	Status string `json:"status,omitempty"`
}
