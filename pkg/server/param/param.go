package param

type CreateVmParam struct {
	ClusterId     string        `json:"clusterId"`
	Namespace     string        `json:"namespace"`
	Name          string        `json:"name"`
	Image         string        `json:"image"`
	Cpu           string        `json:"cpu"`
	Memory        string        `json:"memory"`
	SystemDisk    int           `json:"systemDisk"`
	DataDisks     []int         `json:"dataDisks"`
	LoginConfig   LoginConfig   `json:"loginConfig"`
	NetworkConfig NetworkConfig `json:"networkConfig"`
}

type LoginConfig struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type NetworkConfig struct {
}

type ImageParam struct {
	ClusterId string `json:"clusterId,omitempty"`
	ImageId   string `json:"imageId,omitempty"`
	ImageName string `json:"imageName,omitempty"`
}

type InstanceParam struct {
	ClusterId    string `json:"clusterId,omitempty"`
	InstanceId   string `json:"instanceId,omitempty"`
	InstanceName string `json:"instanceName,omitempty"`
}

type NetworkParam struct {
	ClusterId   string `json:"clusterId,omitempty"`
	NetworkId   string `json:"networkId,omitempty"`
	NetworkName string `json:"networkName,omitempty"`
}

type TemplateParam struct {
	ClusterId    string `json:"clusterId,omitempty"`
	TemplateId   string `json:"templateId,omitempty"`
	TemplateName string `json:"templateName,omitempty"`
}

type VolumeParam struct {
	ClusterId  string `json:"clusterId,omitempty"`
	VolumeId   string `json:"volumeId,omitempty"`
	VolumeName string `json:"volumeName,omitempty"`
}
