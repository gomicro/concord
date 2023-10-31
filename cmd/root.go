package cmd

import (
	"encoding/json"
	"errors"
	"os"
	"path"
	"strings"

	protovalidate "github.com/bufbuild/protovalidate-go"
	"github.com/gomicro/concord/client"
	gh_pb "github.com/gomicro/concord/github/v1"
	"github.com/spf13/cobra"
	"google.golang.org/protobuf/encoding/protojson"
	"gopkg.in/yaml.v3"
)

var (
	clt *client.Client
)

func init() {
	cobra.OnInitialize(initEnvs)
}

func initEnvs() {
}

var rootCmd = &cobra.Command{
	Use:   "concord",
	Short: "concord is a tool to manage your Github repositories",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func setupClient(cmd *cobra.Command, args []string) error {
	tkn := os.Getenv("GITHUB_TOKEN")

	var err error
	clt, err = client.New(tkn)
	if err != nil {
		return err
	}

	return nil
}

func readManifest(file string) (*gh_pb.Organization, error) {
	p, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	b, err := os.ReadFile(path.Join(p, file))
	if err != nil {
		return nil, err
	}

	var v map[string]interface{}
	err = yaml.Unmarshal(b, &v)
	if err != nil {
		return nil, err
	}

	if v["organization"] == nil {
		return nil, errors.New("organization is required")
	}

	j, err := json.Marshal(v["organization"])
	if err != nil {
		return nil, err
	}

	var o gh_pb.Organization
	err = protojson.Unmarshal(j, &o)
	if err != nil {
		return nil, err
	}

	validator, err := protovalidate.New()
	if err != nil {
		return nil, err
	}

	err = validator.Validate(&o)
	if err != nil {
		return nil, err
	}

	fillDefaults(&o)

	return &o, nil
}

func fillDefaults(o *gh_pb.Organization) {
	// labels
	for _, gl := range o.Labels {
		for _, r := range o.Repositories {
			if !hasLabel(r.Labels, gl) {
				r.Labels = append(r.Labels, gl)
			}
		}
	}

	return &o, nil
}

func hasLabel(labels []string, label string) bool {
	for _, l := range labels {
		if strings.EqualFold(l, label) {
			return true
		}
	}

	return false
}
