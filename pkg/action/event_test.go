package action

import (
	"encoding/json"
	"github.com/docker/docker/api/types"
	"github.com/mudler/luet/pkg/helpers/docker"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

// createPayload creates a valid payload structure and marshals it to a string for use in testing
func createPayload(image string, dest string, t *testing.T) string {
	event := EventData{
		Image: image,
		Dest:  dest,
	}
	eventBytes, err := json.Marshal(event)
	if err != nil {
		t.Fatal("Error while marshalling event")
	}
	unpackEvent := UnpackEvent{
		Data: string(eventBytes),
	}

	unpackBytes, err := json.Marshal(unpackEvent)

	if err != nil {
		t.Fatal("Error while marshalling event payload")
	}

	return string(unpackBytes)
}

func TestMissingParamsUnpack(t *testing.T) {
	data := `{"data": "{\"Image\": \"testimage\"}"}`
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
	dest := "/tmp/upgrade"
	moreData := []string{
		createPayload("raccos/releases-opensuse/repository.yaml.gz", dest, t),
		createPayload("raccos/releases-opensuse/tree.tar.yaml.gz", dest, t),
		createPayload("raccos/releases-opensuse/repository.meta.yaml.tar.zstd", dest, t),
		createPayload("raccos/releases-opensuse/compilertree.tar", dest, t),
		createPayload("raccos/releases-opensuse/compilertree.tar.gz", dest, t),
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
	dest, _ := os.MkdirTemp("", "luet-mtree-test-unpack")
	defer os.RemoveAll(dest)
	payload := createPayload("quay.io/costoolkit/releases-opensuse:syslinux-live-6.03", dest, t)

	eventDispatcher := NewEventDispatcherAction("image.post.unpack", payload)
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
	payload := createPayload(image, imageTmp, t)
	eventDispatcher := NewEventDispatcherAction("image.post.unpack", payload)
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
