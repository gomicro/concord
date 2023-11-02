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

func handleError(c *cobra.Command, err error) error {
	c.SilenceUsage = true
	return err
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
	for _, gl := range o.Labels {
		for _, r := range o.Repositories {
			if !hasDefaultLabel(r.Labels, gl) {
				r.Labels = append(r.Labels, gl)
			}
		}
	}

	for _, gf := range o.Files {
		for _, r := range o.Repositories {
			if !hasDefaultFile(r.Files, gf) {
				r.Files = append(r.Files, gf)
			}
		}
	}

	if o.Defaults != nil {
		for _, r := range o.Repositories {
			if o.Defaults.DefaultBranch != nil {
				if r.DefaultBranch == nil {
					r.DefaultBranch = o.Defaults.DefaultBranch
				}
			}

			if o.Defaults.Private != nil {
				if r.Private == nil {
					r.Private = o.Defaults.Private
				}
			}

			for _, p := range o.Defaults.ProtectedBranches {
				if !hasDefaultProtectedBranch(r.ProtectedBranches, p) {
					r.ProtectedBranches = append(r.ProtectedBranches, p)
				} else {
					fillDefaultProtections(r.ProtectedBranches, p)
				}
			}
		}
	}
}

func hasDefaultLabel(labels []string, label string) bool {
	for _, l := range labels {
		if strings.EqualFold(l, label) {
			return true
		}
	}

	return false
}

func hasDefaultFile(files []*gh_pb.File, file *gh_pb.File) bool {
	for _, f := range files {
		if strings.EqualFold(f.Destination, file.Destination) {
			return true
		}
	}

	return false
}

func hasDefaultProtectedBranch(branches []*gh_pb.Branch, branch *gh_pb.Branch) bool {
	for _, b := range branches {
		if strings.EqualFold(b.Name, branch.Name) {
			return true
		}
	}

	return false
}

func fillDefaultProtections(branches []*gh_pb.Branch, branch *gh_pb.Branch) {
	for _, b := range branches {
		if strings.EqualFold(b.Name, branch.Name) {
			if b.Protection.RequirePr == nil {
				b.Protection.RequirePr = branch.Protection.RequirePr
			}

			if b.Protection.ChecksMustPass == nil {
				b.Protection.ChecksMustPass = branch.Protection.ChecksMustPass
			}

			if b.Protection.SignedCommits == nil {
				b.Protection.SignedCommits = branch.Protection.SignedCommits
			}

			if len(b.Protection.RequiredChecks) == 0 {
				b.Protection.RequiredChecks = branch.Protection.RequiredChecks
			} else {
				for _, rc := range branch.Protection.RequiredChecks {
					if !hasDefaultRequiredCheck(b.Protection.RequiredChecks, rc) {
						b.Protection.RequiredChecks = append(b.Protection.RequiredChecks, rc)
					}
				}
			}
		}
	}
}

func hasDefaultRequiredCheck(checks []string, check string) bool {
	for _, c := range checks {
		if strings.EqualFold(c, check) {
			return true
		}
	}

	return false
}
