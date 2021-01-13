// DO NOT EDIT: This file is generated by var_gen.go

package beluga

const (
	varApplication      = "BELUGA_APPLICATION"
	varComposeFile      = "COMPOSE_FILE"
	varComposeTemplate  = "BELUGA_COMPOSE_TEMPLATE"
	varContext          = "BELUGA_CONTEXT"
	varDefaultBranch    = "BELUGA_DEFAULT_BRANCH"
	varDeployDSN        = "BELUGA_DEPLOY_DSN"
	varDockerHost       = "DOCKER_HOST"
	varDockerfile       = "BELUGA_DOCKERFILE"
	varDomain           = "BELUGA_DOMAIN"
	varEnvironment      = "BELUGA_ENVIRONMENT"
	varImage            = "BELUGA_IMAGE"
	varImages           = "BELUGA_IMAGES"
	varImagesTemplate   = "BELUGA_IMAGES_TEMPLATE"
	varOverrides        = "BELUGA_OVERRIDES"
	varRegistry         = "BELUGA_REGISTRY"
	varRegistryPassword = "BELUGA_REGISTRY_PASSWORD"
	varRegistryUsername = "BELUGA_REGISTRY_USERNAME"
	varStackName        = "BELUGA_STACK_NAME"
	varVersion          = "BELUGA_VERSION"
)

var knownVarNames = []string{
	varApplication,
	varComposeFile,
	varComposeTemplate,
	varContext,
	varDefaultBranch,
	varDeployDSN,
	varDockerHost,
	varDockerfile,
	varDomain,
	varEnvironment,
	varImage,
	varImages,
	varImagesTemplate,
	varOverrides,
	varRegistry,
	varRegistryPassword,
	varRegistryUsername,
	varStackName,
	varVersion,
}

// If provided, name of the (sub)application to compile
func (e Environment) Application() string {
	return e[varApplication]
}

// Compose file(s) for deploying
func (e Environment) ComposeFile() string {
	return e[varComposeFile]
}

// A template docker-compose file that may contain mutations for the compose
// file
func (e Environment) ComposeTemplate() string {
	return e[varComposeTemplate]
}

// Docker build context
func (e Environment) Context() string {
	return e[varContext]
}

// Docker instance for deploying
func (e Environment) DeployDSN() string {
	return e[varDeployDSN]
}

// Dockerfile to build
func (e Environment) Dockerfile() string {
	return e[varDockerfile]
}

// Domain name of the stack
func (e Environment) Domain() string {
	return e[varDomain]
}

// Environment name
func (e Environment) Environment() string {
	return e[varEnvironment]
}

// First image listed in BELUGA_IMAGES; doesn't affect pushng
func (e Environment) Image() string {
	return e[varImage]
}

// Docker images to push after build
func (e Environment) Images() string {
	return e[varImages]
}

// Go template for a space-separated list of Docker images to push after build
func (e Environment) ImagesTemplate() string {
	return e[varImagesTemplate]
}

// YAML document with environment names or patterns as keys and variables to
// override as values
func (e Environment) Overrides() string {
	return e[varOverrides]
}

// Docker registry for pushing
func (e Environment) Registry() string {
	return e[varRegistry]
}

// Password for Docker registry
func (e Environment) RegistryPassword() string {
	return e[varRegistryPassword]
}

// Username for Docker registry
func (e Environment) RegistryUsername() string {
	return e[varRegistryUsername]
}

// Name of the compose/swamrm/etc. stack
func (e Environment) StackName() string {
	return e[varStackName]
}

// Version of the application being built/deployed
func (e Environment) Version() string {
	return e[varVersion]
}
