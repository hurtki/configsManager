package cmd_test

import (
	"testing"

	"github.com/hurtki/configsManager/cmd"
	"github.com/hurtki/configsManager/mocks"
	"github.com/hurtki/configsManager/services"
	gomock "go.uber.org/mock/gomock"
)

// Test for valid path checking
func TestPathCmd_ValidPathAdding(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAppConfigService := mocks.NewMockAppConfigService(ctrl)
	mockConfigsListService := mocks.NewMockConfigsListService(ctrl)
	key := "cm_config"
	path := "some_path"
	returnConfigsList := services.GetDefaultConfigsList("")
	returnConfigsList.SetConfig(key, path)
	mockConfigsListService.EXPECT().Load().Return(returnConfigsList, nil)

	pathCmd := cmd.NewPathCmd(mockAppConfigService, mockConfigsListService)
	args := []string{key}
	err := pathCmd.Command.RunE(pathCmd.Command, args)
	if err != nil {
		t.Errorf("excpected no error on getting path, got %v", err)
	}
}

// Test for not existing key
func TestPathCmd_KeyNotFound_ReturnsError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockConfigsList := mocks.NewMockConfigsListService(ctrl)
	cl := services.GetDefaultConfigsList("")
	cl.SetConfig("mykey", "/path/to/file")
	mockConfigsList.EXPECT().Load().Return(cl, nil)

	cmd := cmd.NewPathCmd(nil, mockConfigsList)

	err := cmd.Command.RunE(cmd.Command, []string{"not_exist_key"})
	if err == nil || err.Error() != "key not found" {
		t.Errorf("expected 'key not found' error, got: %v", err)
	}
}

// Test for not enough args
func TestPathCmd_NotEnoughArgs_ReturnsError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockConfigsList := mocks.NewMockConfigsListService(ctrl)
	cl := services.GetDefaultConfigsList("")
	cl.SetConfig("mykey", "/path/to/file")
	mockConfigsList.EXPECT().Load().Return(cl, nil)

	cmd := cmd.NewPathCmd(nil, mockConfigsList)

	err := cmd.Command.RunE(cmd.Command, []string{"not_exist_key"})
	if err == nil || err.Error() != "key not found" {
		t.Errorf("expected 'key not found' error, got: %v", err)
	}
}
