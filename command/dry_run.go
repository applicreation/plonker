package command

import (
	"github.com/applicreation/plonker/config"
	"github.com/applicreation/plonker/connection"
	_ "github.com/go-sql-driver/mysql"
	"log"
)

type DryRunCommand struct {
	Config     *config.Config
	Connection *connection.Connection
}

type Table struct {
	Config config.Table
	Count  int
	Keys   []string
}

func (c *DryRunCommand) Run(args []string) int {
	var tables = make(map[string]Table)

	for _, table := range c.Config.Tables {
		tables[table.Name] = Table{
			Config: table,
		}
	}

	log.Println("--------------------")

	for key, table := range tables {
		log.Printf("Table: %s", key)

		if table.Config.Range != (config.Range{}) {
			count, _ := c.Connection.Count(table.Config)
			table.Count += count
		}

		if len(table.Config.Relations) > 0 {
			for _, relationship := range table.Config.Relations {
				log.Printf("Relationship: %s", relationship.Table)

				relationshipTable := tables[relationship.Table]

				count, _ := c.Connection.RelationshipCount(table.Config, relationship, table.Count)

				relationshipTable.Count += count

				tables[relationship.Table] = relationshipTable
			}
		}

		tables[key] = table

		log.Println("--------------------")
	}

	log.Println("Estimated record counts")
	for _, table := range tables {
		log.Printf("%s: %d", table.Config.Name, table.Count)
	}

	return 0
}

func (c *DryRunCommand) Help() string {
	return ""
}

func (c *DryRunCommand) Synopsis() string {
	return ""
}
