package action

import (
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/mudler/luet/pkg/helpers/docker"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestMissingParamsUnpack(t *testing.T) {
	data := `{"data": {"Image": "testimage"}}`
	eventDispatcher := NewEventDispatcherAction("image.post.unpack", data)
	out, err := eventDispatcher.Run()
	commonKeysAssertUnpack(t, out)

	assert.NotEmpty(t, err)
	assert.Error(t, err)
	assert.Equal(t, err.Error(), "fields missing from payload")

}

func TestMalformedPayloadUnpack(t *testing.T) {
	data := `{`
	eventDispatcher := NewEventDispatcherAction("image.post.unpack", data)
	out, err := eventDispatcher.Run()
	commonKeysAssertUnpack(t, out)

	assert.NotEmpty(t, err)
	assert.Error(t, err)
}

func TestBlacklistUnpack(t *testing.T) {
	for _, image := range ImageBlacklist {
		data := fmt.Sprintf("{\"data\": {\"Image\": \"%s\", \"Dest\": \"/tmp/upgrade\"}}", image)
		eventDispatcher := NewEventDispatcherAction("image.post.unpack", data)
		out, err := eventDispatcher.Run()
		commonKeysAssertUnpack(t, out)
		assert.Equal(t, out["state"], "All checks succeeded")
		assert.NoError(t, err)
		assert.Empty(t, err)
	}

	moreData := []string{
		`{"data": {"Image": "raccos/releases-opensuse/repository.yaml.gz", "Dest": "/tmp/upgrade"}}`,
		`{"data": {"Image": "raccos/releases-opensuse/tree.tar.yaml.gz", "Dest": "/tmp/upgrade"}}`,
		`{"data": {"Image": "raccos/releases-opensuse/repository.meta.yaml.tar.zstd", "Dest": "/tmp/upgrade"}}`,
		`{"data": {"Image": "raccos/releases-opensuse/compilertree.tar", "Dest": "/tmp/upgrade"}}`,
		`{"data": {"Image": "raccos/releases-opensuse/compilertree.tar.gz", "Dest": "/tmp/upgrade"}}`,
	}

	for _, data := range moreData {
		eventDispatcher := NewEventDispatcherAction("image.post.unpack", data)
		out, err := eventDispatcher.Run()
		commonKeysAssertUnpack(t, out)
		assert.Equal(t, out["state"], "All checks succeeded")
		assert.NoError(t, err)
		assert.Empty(t, err)
	}
}

func TestUnpackFailureToValidate(t *testing.T) {
	image := "quay.io/costoolkit/releases-opensuse:syslinux-live-6.03"
	dest, _ := os.MkdirTemp("", "luet-mtree-test-unpack")
	defer os.RemoveAll(dest)
	data := fmt.Sprintf("{\"data\": {\"Image\": \"%s\", \"Dest\": \"%s/\"}}", image, dest)
	eventDispatcher := NewEventDispatcherAction("image.post.unpack", data)
	out, err := eventDispatcher.Run()
	commonKeysAssertUnpack(t, out)

	assert.Equal(t, out["state"], "Checks failed")
	assert.Error(t, err)
	assert.Contains(t, out["error"], "Validation failed")
	assert.NotEmpty(t, err)
}

func TestUnpackSuccess(t *testing.T) {
	// Image to download and check the mtree values
	image := "quay.io/costoolkit/releases-opensuse:systemd-boot-live-26"
	// Dir to temporally extract the full image
	extractTempDir, _ := os.MkdirTemp("", "luet-mtree-docker-extract")
	defer os.RemoveAll(extractTempDir)
	// Destination dir for the image contents
	imageTmp, _ := os.MkdirTemp("", "luet-mtree-image-extract")
	defer os.RemoveAll(imageTmp)
	// Empty auth
	auth := &types.AuthConfig{
		Username:      "",
		Password:      "",
		ServerAddress: "",
		Auth:          "",
		IdentityToken: "",
		RegistryToken: "",
	}
	_, err := docker.DownloadAndExtractDockerImage(extractTempDir, image, imageTmp, auth, false)
	if err != nil {
		t.Fatal(err)
	}

	// Now that we have download and extracted the image, call the unpack event on it like luet does
	data := fmt.Sprintf("{\"data\": {\"Image\": \"%s\", \"Dest\": \"%s/\"}}", image, imageTmp)
	eventDispatcher := NewEventDispatcherAction("image.post.unpack", data)
	out, err := eventDispatcher.Run()
	commonKeysAssertUnpack(t, out)
	assert.Equal(t, out["state"], "All checks succeeded")
	assert.Empty(t, out["error"])
	assert.Empty(t, out["data"])
}

func commonKeysAssertUnpack(t *testing.T, out map[string]string) {
	assert.NotEmpty(t, out)
	assert.Contains(t, out, "data")
	assert.Contains(t, out, "error")
	assert.Contains(t, out, "state")
}
