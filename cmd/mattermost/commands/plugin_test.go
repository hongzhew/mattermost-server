package commands

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/mattermost/mattermost-server/config"
	"github.com/mattermost/mattermost-server/utils/fileutils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPlugin(t *testing.T) {
	th := Setup().InitBasic()
	defer th.TearDown()

	cfg := th.Config()
	*cfg.PluginSettings.EnableUploads = true
	*cfg.PluginSettings.Directory = "./test-plugins"
	*cfg.PluginSettings.ClientDirectory = "./test-client-plugins"
	th.SetConfig(cfg)

	os.MkdirAll("./test-plugins", os.ModePerm)
	os.MkdirAll("./test-client-plugins", os.ModePerm)

	path, _ := fileutils.FindDir("tests")

	th.CheckCommand(t, "plugin", "add", filepath.Join(path, "testplugin.tar.gz"))

	th.CheckCommand(t, "plugin", "enable", "testplugin")
	fs, err := config.NewFileStore(th.ConfigPath(), false)
	require.Nil(t, err)
	assert.True(t, fs.Get().PluginSettings.PluginStates["testplugin"].Enable)
	fs.Close()

	th.CheckCommand(t, "plugin", "disable", "testplugin")
	fs, err = config.NewFileStore(th.ConfigPath(), false)
	require.Nil(t, err)
	assert.False(t, fs.Get().PluginSettings.PluginStates["testplugin"].Enable)
	fs.Close()

	th.CheckCommand(t, "plugin", "list")

	th.CheckCommand(t, "plugin", "delete", "testplugin")
}

func TestPluginPublicKeys(t *testing.T) {
	th := Setup().InitBasic()
	defer th.TearDown()

	cfg := th.Config()
	cfg.PluginSettings.SignaturePublicKeyFiles = []string{"public-key"}
	th.SetConfig(cfg)

	output := th.CheckCommand(t, "plugin", "keys")
	assert.Contains(t, output, "public-key")
	assert.NotContains(t, output, "Plugin name:")

}

func TestPluginPublicKeyDetails(t *testing.T) {
	th := Setup().InitBasic()
	defer th.TearDown()

	cfg := th.Config()
	cfg.PluginSettings.SignaturePublicKeyFiles = []string{"public-key"}

	th.SetConfig(cfg)

	output := th.CheckCommand(t, "plugin", "keys", "--verbose", "true")
	assert.Contains(t, output, "Plugin name: public-key")
	output = th.CheckCommand(t, "plugin", "keys", "--verbose")
	assert.Contains(t, output, "Plugin name: public-key")
}

func TestAddPluginPublicKeys(t *testing.T) {
	th := Setup().InitBasic()
	defer th.TearDown()

	cfg := th.Config()
	cfg.PluginSettings.SignaturePublicKeyFiles = []string{"public-key"}
	th.SetConfig(cfg)

	output := th.CheckCommand(t, "plugin", "keys", "add", "pk1.gpg")
	assert.Contains(t, output, "pk1.gpg.plugin.gpg: no such file or directory")
}

func TestDeletePluginPublicKeys(t *testing.T) {
	th := Setup().InitBasic()
	defer th.TearDown()

	cfg := th.Config()
	cfg.PluginSettings.SignaturePublicKeyFiles = []string{"pk1.gpg"}
	th.SetConfig(cfg)

	output := th.CheckCommand(t, "plugin", "keys", "delete", "pk1.gpg")
	assert.Contains(t, output, "Deleted public key: pk1.gpg")
}
