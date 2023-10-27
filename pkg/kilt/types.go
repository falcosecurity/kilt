package kilt

const (
	DefaultUserID      = 0
	DefaultGroupID     = 0
	DefaultPermissions = 0755
)

type TargetInfo struct {
	Image                string            `json:"image"`
	ContainerName        string            `json:"container_name"`
	ContainerGroupName   string            `json:"container_group_name"`
	EntryPoint           []string          `json:"entry_point"`
	Command              []string          `json:"command"`
	EnvironmentVariables map[string]string `json:"environment_variables"`
	Metadata             map[string]string `json:"metadata"`
}

type BuildResource struct {
	Name                 string
	Image                string
	Volumes              []string
	EntryPoint           []string
	EnvironmentVariables []map[string]interface{}
}

type Build struct {
	Image                string
	EntryPoint           []string
	Command              []string
	EnvironmentVariables map[string]string
	Capabilities         []string

	Resources []BuildResource
}

type RuntimeUpload struct {
	Payload     *Payload
	Destination string
	Uid         uint16
	Gid         uint16
	Permissions uint32
}

type RuntimeExecutable struct {
	Run []string
}

type Runtime struct {
	Uploads     []RuntimeUpload
	Executables []RuntimeExecutable
}

type PayloadType string

const (
	URL       PayloadType = "url"
	LocalPath PayloadType = "local-path"
	Base64    PayloadType = "base64"
	Text      PayloadType = "text"
	Unknown   PayloadType = "unknown"
)

type Payload struct {
	Contents string
	Type     PayloadType
	Gzipped  bool
}

type LanguageInterface interface {
	Build(info *TargetInfo) (*Build, error)
	Runtime(info *TargetInfo) (*Runtime, error)
}
