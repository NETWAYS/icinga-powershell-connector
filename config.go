package main

import (
	"crypto/tls"
	"errors"
	"fmt"
	flag "github.com/spf13/pflag"
	"net/http"
	"os"
)

const (
	DefaultAPI  = "https://localhost:5668"
	ProgramName = "icinga-powershell-connector"
)

type Config struct {
	API          string
	Command      string
	Arguments    map[string]interface{}
	CertName     string
	CAFile       string
	Insecure     bool
	Debug        bool
	PrintVersion bool
}

var (
	// ErrVersionRequested returned when --version flag is set
	ErrVersionRequested = errors.New("version was requested by a flag")

	// ErrNoCommand is returned when no PowerShell command could be parsed from flags.
	ErrNoCommand = errors.New("no command found for PowerShell execution")
)

func NewConfig() *Config {
	return &Config{
		API:      DefaultAPI,
		CertName: GetIcingaNodeName(),
		CAFile:   IcingaCAPath,
	}
}

// BuildFlags for a passed flag.FlagSet to store values inside Config.
func (c *Config) BuildFlags(fs *flag.FlagSet) {
	fs.StringVar(&c.Command, "command", c.Command, "Command to be executed")
	fs.StringVar(&c.API, "api", c.API, "API Endpoint")
	fs.StringVar(&c.CertName, "cert-name", c.CertName, "Certificate Name to be expected")
	fs.StringVar(&c.CAFile, "ca-file", c.CAFile, "Icinga CA file to be loaded")
	fs.BoolVar(&c.Insecure, "insecure", c.Insecure, "Ignore any certificate checks")
	fs.BoolVar(&c.Debug, "debug", c.Debug, "Enable debug logging")
	fs.BoolVar(&c.PrintVersion, "version", false, "Print program version")
}

// ParseConfigFromFlags to be called to parse CLI arguments and return the built Config struct.
func ParseConfigFromFlags(arguments []string) (config *Config, err error) {
	config = NewConfig()

	flags, powerShellArgs := SplitPowerShellArguments(arguments)

	fs := flag.NewFlagSet(os.Args[0], flag.ContinueOnError)
	config.BuildFlags(fs)
	fs.SortFlags = false

	err = fs.Parse(flags)
	if err != nil {
		return nil, err
	}

	if config.PrintVersion {
		_, _ = fmt.Fprintln(os.Stdout, ProgramName+" "+buildVersion())
		_, _ = fmt.Fprint(os.Stdout, License+"\n")

		return nil, ErrVersionRequested
	}

	// Parse Powershell arguments
	command, args := GetPowershellArgs(powerShellArgs)
	if command != "" {
		config.Command = command
	}

	config.Arguments = args

	if config.Command == "" {
		return config, ErrNoCommand
	}

	return
}

// SplitPowerShellArguments separate this commands flags from Powershell.exe arguments.
//
// Usually this starts shorthand flag.
func SplitPowerShellArguments(arguments []string) (flags, powerShell []string) {
	var isPowerShell bool

	for _, arg := range arguments {
		// Look for a shorthand argument
		if arg[0] == '-' && arg[1] != '-' {
			isPowerShell = true
		}

		if isPowerShell {
			powerShell = append(powerShell, arg)
		} else {
			flags = append(flags, arg)
		}
	}

	return
}

func (c Config) NewClient() *http.Client {
	tlsConfig := &tls.Config{
		RootCAs:            LoadIcingaCACert(c.CAFile),
		InsecureSkipVerify: c.Insecure, // nolint:gosec // intended configuration
		ServerName:         c.CertName,
	}

	transport := &http.Transport{TLSClientConfig: tlsConfig}

	return &http.Client{Transport: transport}
}
