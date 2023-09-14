package main

import (
	"crypto/x509"
	"encoding/json"
	"log/slog"
	"os"
)

const (
	IcingaDataPath = IcingaStatePrefix + "/lib/icinga2"
	IcingaCAPath   = IcingaDataPath + "/certs/ca.crt"
	IcingaVarsFile = IcingaStatePrefix + "/cache/icinga2/icinga2.vars"
)

type IcingaVar struct {
	Name  string
	Value string
}

func LoadIcingaCACert(path string) *x509.CertPool {
	if path == "" {
		path = IcingaCAPath
	}

	// Load contents
	data, err := os.ReadFile(path)
	if err != nil {
		slog.Error("could not read Icinga CA certificate", "path", path)

		return nil
	}

	// Build pool
	pool := x509.NewCertPool()
	if !pool.AppendCertsFromPEM(data) {
		slog.Error("could not append any CA certificates to pool", "path", path)
	}

	return pool
}

func GetIcingaNodeName() string {
	vars := LoadIcingaVariables("")
	return vars["NodeName"]
}

func LoadIcingaVariables(path string) (vars map[string]string) {
	if path == "" {
		path = IcingaVarsFile
	}

	vars = map[string]string{}

	fh, err := os.Open(path)
	if err != nil {
		slog.Error("could not read vars file", "path", path)

		return nil
	}

	var (
		entry []byte
		v     IcingaVar
	)

	for {
		entry, err = ParseNetstring(fh)
		if err != nil || entry == nil {
			// TODO: handle error?
			break
		}

		err = json.Unmarshal(entry, &v)
		if err != nil {
			// TODO: handle error? - non string can not be parsed currently
			continue
		}

		vars[v.Name] = v.Value
	}

	return
}
