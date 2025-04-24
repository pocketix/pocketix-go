package commands

type Collector interface {
	NewCollectorBasedOnType(collectorType Collector, target *[]Command) Collector
	Collect(statement Command)
	Type() Collector
	GetTarget() *[]Command
}

type NoOpCollector struct {
	Target []Command
}

func (c *NoOpCollector) Collect(statement Command) {
	if statement.GetId() == "case" {
		c.Target = append(c.Target, statement)
	}
}

func (c *NoOpCollector) Type() Collector {
	return c
}

func (c *NoOpCollector) NewCollectorBasedOnType(collectorType Collector, _ *[]Command) Collector {
	return &NoOpCollector{}
}

func (c *NoOpCollector) GetTarget() *[]Command {
	return &c.Target
}

type ASTCollector struct {
	Target *[]Command
}

func (c *ASTCollector) Collect(statement Command) {
	*c.Target = append(*c.Target, statement)
}

func (c *ASTCollector) Type() Collector {
	return c
}

func (c *ASTCollector) NewCollectorBasedOnType(collectorType Collector, target *[]Command) Collector {
	return &ASTCollector{
		Target: target,
	}
}

func (c *ASTCollector) GetTarget() *[]Command {
	return c.Target
}
