package main

import (
	"crypto/x509"
	"encoding/json"
	"os"

	log "github.com/sirupsen/logrus"
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
		log.WithError(err).WithField("path", path).Debug("could not read Icinga CA certificate")

		return nil
	}

	// Build pool
	pool := x509.NewCertPool()
	if !pool.AppendCertsFromPEM(data) {
		log.WithField("path", path).Warning("could not append any CA certificates to pool")
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
		log.WithError(err).WithField("path", path).Debug("could not read vars file")
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
