package manifest

import (
	"os"
	"path/filepath"
	"testing"
)

func writeManifest(t *testing.T, contents string) string {
	t.Helper()

	file := filepath.Join(t.TempDir(), "concord.yml")

	err := os.WriteFile(file, []byte(contents), 0o600)
	if err != nil {
		t.Fatalf("failed to write manifest: %v", err)
	}

	return file
}

func TestReadManifest(t *testing.T) {
	file := writeManifest(t, `
organization:
  name: my-org
  permissions:
    base_permissions: read
`)

	m, err := ReadManifest(file)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if m.GetName() != "my-org" {
		t.Errorf("name: got %q, want %q", m.GetName(), "my-org")
	}

	if m.GetPermissions().GetBasePermissions() != "read" {
		t.Errorf("base permissions: got %q, want %q", m.GetPermissions().GetBasePermissions(), "read")
	}
}

func TestReadManifestValidation(t *testing.T) {
	tests := []struct {
		name     string
		contents string
	}{
		{
			name: "empty org name",
			contents: `
organization:
  name: ""
`,
		},
		{
			name: "invalid base permissions",
			contents: `
organization:
  name: my-org
  permissions:
    base_permissions: superuser
`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := ReadManifest(writeManifest(t, tt.contents))
			if err == nil {
				t.Fatal("expected validation error")
			}
		})
	}
}
