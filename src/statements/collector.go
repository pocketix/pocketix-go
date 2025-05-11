package statements

import "github.com/pocketix/pocketix-go/src/models"

type Collector interface {
	NewCollectorBasedOnType(collectorType Collector, target *[]Statement) Collector
	Collect(statement Statement)
	Type() Collector
	GetTarget() *[]Statement
	CollectDeviceCommand(deviceCommand models.SDInformationFromBackend)
	IsDeviceCommandCollected(deviceCommand models.DeviceCommand) (models.SDInformationFromBackend, bool)
}

type NoOpCollector struct {
	Target []Statement
}

func (c *NoOpCollector) Collect(statement Statement) {
	if statement.GetId() == "case" {
		c.Target = append(c.Target, statement)
	}
}

func (c *NoOpCollector) Type() Collector {
	return c
}

func (c *NoOpCollector) NewCollectorBasedOnType(collectorType Collector, _ *[]Statement) Collector {
	return &NoOpCollector{}
}

func (c *NoOpCollector) GetTarget() *[]Statement {
	return &c.Target
}

func (c *NoOpCollector) CollectDeviceCommand(deviceCommand models.SDInformationFromBackend) {
	// No operation for NoOpCollector
}

func (c *NoOpCollector) IsDeviceCommandCollected(deviceCommand models.DeviceCommand) (models.SDInformationFromBackend, bool) {
	// No operation for NoOpCollector
	return models.SDInformationFromBackend{}, false
}

type ASTCollector struct {
	Target         *[]Statement
	DeviceCommands []models.SDInformationFromBackend
}

func (c *ASTCollector) Collect(statement Statement) {
	*c.Target = append(*c.Target, statement)
}

func (c *ASTCollector) Type() Collector {
	return c
}

func (c *ASTCollector) NewCollectorBasedOnType(collectorType Collector, target *[]Statement) Collector {
	return &ASTCollector{
		Target:         target,
		DeviceCommands: []models.SDInformationFromBackend{},
	}
}

func (c *ASTCollector) GetTarget() *[]Statement {
	return c.Target
}

func (c *ASTCollector) CollectDeviceCommand(deviceCommand models.SDInformationFromBackend) {
	c.DeviceCommands = append(c.DeviceCommands, deviceCommand)
}

func (c *ASTCollector) IsDeviceCommandCollected(deviceCommand models.DeviceCommand) (models.SDInformationFromBackend, bool) {
	for _, cmd := range c.DeviceCommands {
		if cmd.DeviceUID == deviceCommand.DeviceUID && cmd.Command.CommandDenotation == deviceCommand.CommandDenotation {
			return cmd, true
		}
	}
	return models.SDInformationFromBackend{}, false
}
