package config

// Configurations ...
type Configurations struct {
	Server   ServerConfigurations
	Database DatabaseConfigurations
	Logger   LoggerConfigurations
}

// ServerConfigurations exported
type ServerConfigurations struct {
	Host        string
	Port        int
	Environment string
}

// DatabaseConfigurations exported
type DatabaseConfigurations struct {
	URL            string
	DBName         string
	CollectionName string
}

// LoggerConfigurations exported
type LoggerConfigurations struct {
	OutputPath string
}
