package manifest

import (
	"context"
	"encoding/json"
	"errors"
	"os"
	"path"
	"strings"

	"github.com/bufbuild/protovalidate-go"
	gh_pb "github.com/gomicro/concord/github/v1"
	"google.golang.org/protobuf/encoding/protojson"
	"gopkg.in/yaml.v3"
)

type ctxKey string

const (
	manifestKey ctxKey = "manifest"
)

var (
	ErrManifestnotFound    = errors.New("manifest not found")
	ErrManifestOrgRequried = errors.New("organization is required")
)

func WithManifest(ctx context.Context, file string) context.Context {
	ctx, cancel := context.WithCancelCause(ctx)

	p, err := os.Getwd()
	if err != nil {
		cancel(err)
	}

	b, err := os.ReadFile(path.Join(p, file))
	if err != nil {
		cancel(err)
	}

	var v map[string]interface{}
	err = yaml.Unmarshal(b, &v)
	if err != nil {
		cancel(err)
	}

	if v["organization"] == nil {
		cancel(ErrManifestOrgRequried)
	}

	j, err := json.Marshal(v["organization"])
	if err != nil {
		cancel(err)
	}

	var m gh_pb.Organization
	err = protojson.Unmarshal(j, &m)
	if err != nil {
		cancel(err)
	}

	validator, err := protovalidate.New()
	if err != nil {
		cancel(err)
	}

	err = validator.Validate(&m)
	if err != nil {
		cancel(err)
	}

	fillDefaults(&m)

	return context.WithValue(ctx, manifestKey, &m)
}

func OrgFromContext(ctx context.Context) (*gh_pb.Organization, error) {
	m, ok := ctx.Value(manifestKey).(*gh_pb.Organization)
	if !ok {
		return nil, ErrManifestnotFound
	}

	return m, nil
}

func fillDefaults(o *gh_pb.Organization) {
	for _, gl := range o.Labels {
		for _, r := range o.Repositories {
			if !hasDefaultLabel(r.Labels, gl) {
				r.Labels = append(r.Labels, gl)
			}
		}
	}

	if o.Defaults != nil {
		for _, r := range o.Repositories {
			if o.Defaults.Private != nil {
				if r.Private == nil {
					r.Private = o.Defaults.Private
				}
			}

			if o.Defaults.DefaultBranch != nil {
				if r.DefaultBranch == nil {
					r.DefaultBranch = o.Defaults.DefaultBranch
				}
			}

			if o.Defaults.AllowAutoMerge != nil {
				if r.AllowAutoMerge == nil {
					r.AllowAutoMerge = o.Defaults.AllowAutoMerge
				}
			}

			if o.Defaults.AutoDeleteHeadBranches != nil {
				if r.AutoDeleteHeadBranches == nil {
					r.AutoDeleteHeadBranches = o.Defaults.AutoDeleteHeadBranches
				}
			}

			for _, p := range o.Defaults.ProtectedBranches {
				if !hasDefaultProtectedBranch(r.ProtectedBranches, p) {
					r.ProtectedBranches = append(r.ProtectedBranches, p)
				} else {
					fillDefaultProtections(r.ProtectedBranches, p)
				}
			}

			for _, gf := range o.Defaults.Files {
				for _, r := range o.Repositories {
					if !hasDefaultFile(r.Files, gf) {
						r.Files = append(r.Files, gf)
					}
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
