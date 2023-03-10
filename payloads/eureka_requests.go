package payloads

type InstanceRegistrationRequest struct {
	Instance *EurekaInstance `json:"instance" validate:"required"`
}

type EurekaApplication struct {
	Name      string           `xml:"name" json:"name" validate:"required"`
	Instances []EurekaInstance `xml:"instance" json:"instance" validate:"required"`
}

type EurekaInstance struct {
	HostName                      string            `xml:"hostName" json:"hostName"`
	HomePageUrl                   string            `xml:"homePageUrl,omitempty" json:"homePageUrl,omitempty"`
	StatusPageUrl                 string            `xml:"statusPageUrl" json:"statusPageUrl"`
	HealthCheckUrl                string            `xml:"healthCheckUrl,omitempty" json:"healthCheckUrl,omitempty"`
	App                           string            `xml:"app" json:"app"`
	IpAddress                     string            `xml:"ipAddr" json:"ipAddr"`
	VipAddress                    string            `xml:"vipAddress" json:"vipAddress"`
	SecureVipAddress              string            `xml:"secureVipAddress,omitempty" json:"secureVipAddress,omitempty"`
	Status                        string            `xml:"status" json:"status"`
	Port                          *EurekaPort       `xml:"port,omitempty" json:"port,omitempty"`
	SecurePort                    *EurekaPort       `xml:"securePort,omitempty" json:"securePort,omitempty"`
	DataCenterInfo                *EurekaDatacenter `xml:"dataCenterInfo" json:"dataCenterInfo"`
	LeaseInfo                     *EurekaLease      `xml:"leaseInfo,omitempty" json:"leaseInfo,omitempty"`
	Metadata                      map[string]string `xml:"metadata,omitempty" json:"metadata,omitempty"`
	IsCoordinatingDiscoveryServer bool              `xml:"isCoordinatingDiscoveryServer,omitempty" json:"isCoordinatingDiscoveryServer,omitempty"`
	LastUpdatedTimestamp          int               `xml:"lastUpdatedTimestamp,omitempty" json:"lastUpdatedTimestamp,omitempty"`
	LastDirtyTimestamp            int               `xml:"lastDirtyTimestamp,omitempty" json:"lastDirtyTimestamp,omitempty"`
	ActionType                    string            `xml:"actionType,omitempty" json:"actionType,omitempty"`
	OverriddenStatus              string            `xml:"overriddenstatus,omitempty" json:"overriddenstatus,omitempty"`
	CountryId                     int               `xml:"countryId,omitempty" json:"countryId,omitempty"`
	InstanceId                    string            `xml:"instanceId,omitempty" json:"instanceId,omitempty"`
}

type EurekaPort struct {
	Port int `xml:",chardata" json:"$"`
}

type EurekaLease struct {
	EvictionDurationInSecs uint `xml:"evictionDurationInSecs,omitempty" json:"evictionDurationInSecs,omitempty"`
	RenewalIntervalInSecs  int  `xml:"renewalIntervalInSecs,omitempty" json:"renewalIntervalInSecs,omitempty"`
	DurationInSecs         int  `xml:"durationInSecs,omitempty" json:"durationInSecs,omitempty"`
	RegistrationTimestamp  int  `xml:"registrationTimestamp,omitempty" json:"registrationTimestamp,omitempty"`
	LastRenewalTimestamp   int  `xml:"lastRenewalTimestamp,omitempty" json:"lastRenewalTimestamp,omitempty"`
	EvictionTimestamp      int  `xml:"evictionTimestamp,omitempty" json:"evictionTimestamp,omitempty"`
	ServiceUpTimestamp     int  `xml:"serviceUpTimestamp,omitempty" json:"serviceUpTimestamp,omitempty"`
}

type EurekaDatacenter struct {
	Name     string                    `xml:"name" json:"name"`
	Class    string                    `xml:"class,attr" json:"@class"`
	Metadata *EurekaDatacenterMetadata `xml:"metadata,omitempty" json:"metadata,omitempty"`
}

type EurekaDatacenterMetadata struct {
	AmiLaunchIndex   string `xml:"ami-launch-index,omitempty" json:"ami-launch-index,omitempty"`
	LocalHostname    string `xml:"local-hostname,omitempty" json:"local-hostname,omitempty"`
	AvailabilityZone string `xml:"availability-zone,omitempty" json:"availability-zone,omitempty"`
	InstanceId       string `xml:"instance-id,omitempty" json:"instance-id,omitempty"`
	PublicIpv4       string `xml:"public-ipv4,omitempty" json:"public-ipv4,omitempty"`
	PublicHostname   string `xml:"public-hostname,omitempty" json:"public-hostname,omitempty"`
	AmiManifestPath  string `xml:"ami-manifest-path,omitempty" json:"ami-manifest-path,omitempty"`
	LocalIpv4        string `xml:"local-ipv4,omitempty" json:"local-ipv4,omitempty"`
	Hostname         string `xml:"hostname,omitempty" json:"hostname,omitempty"`
	AmiId            string `xml:"ami-id,omitempty" json:"ami-id,omitempty"`
	InstanceType     string `xml:"instance-type,omitempty" json:"instance-type,omitempty"`
}
