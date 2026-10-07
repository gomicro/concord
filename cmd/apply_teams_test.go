package cmd

import (
	"slices"
	"testing"

	gh_pb "github.com/gomicro/concord/github/v1"
	"github.com/google/go-github/v92/github"
)

func TestGetTeamMembersBreakdown(t *testing.T) {
	people := []*gh_pb.People{
		{Username: "alice", Teams: []string{"Backend"}},
		{Username: "bob", Teams: []string{"Frontend"}},
		{Username: "carol", Teams: []string{"backend", "Frontend"}},
		{Username: "dave"},
	}

	var members []*github.User
	for _, login := range []string{"Alice", "bob", "dave", "eve"} {
		members = append(members, &github.User{Login: &login})
	}

	missing, managed, unmanaged := getTeamMembersBreakdown("Backend", people, members)

	if !slices.Equal(missing, []string{"carol"}) {
		t.Errorf("missing: got %v, want [carol]", missing)
	}

	if !slices.Equal(managed, []string{"Alice"}) {
		t.Errorf("managed: got %v, want [Alice]", managed)
	}

	// bob and dave are in the manifest but not assigned to this team, eve isn't in the manifest at all
	if !slices.Equal(unmanaged, []string{"bob", "dave", "eve"}) {
		t.Errorf("unmanaged: got %v, want [bob dave eve]", unmanaged)
	}
}
