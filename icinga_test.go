package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadIcingaVariables(t *testing.T) {
	vars := LoadIcingaVariables("testdata/icinga2.vars")

	assert.Contains(t, vars, "PluginDir")
	assert.Contains(t, vars, "TicketSalt")
	assert.Contains(t, vars, "NodeName")

	assert.Equal(t, "icinga.example.com", vars["NodeName"])
	assert.Equal(t, "secret", vars["TicketSalt"])
}
