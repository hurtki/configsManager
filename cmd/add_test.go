package cmd_test

import (
	"testing"

	"github.com/hurtki/configsManager/cmd"
	"github.com/hurtki/configsManager/mocks"
	"github.com/hurtki/configsManager/services"
	gomock "go.uber.org/mock/gomock"
)

// test for valid config adding one arg
func TestAddCmd_ValidConfigAddOneParam(t *testing.T) {

	ctrl := gomock.NewController(t)
	// creating mock dependencies
	mockAppConfigService := mocks.NewMockAppConfigService(ctrl)
	mockConfigsListService := mocks.NewMockConfigsListService(ctrl)
	mockInputService := mocks.NewMockInputService(ctrl)
	mockOsService := mocks.NewMockOsService(ctrl)

	args := []string{"config.json"}
	mockInputService.EXPECT().GetPipedInput().Return("", false)
	returnAppConfig := services.NewDefaultAppConfig()
	returnConfigsList := services.GetDefaultConfigsList("")

	mockAppConfigService.EXPECT().Load().Return(returnAppConfig, nil)
	mockConfigsListService.EXPECT().Load().Return(returnConfigsList, nil)

	mockOsService.EXPECT().FileExists("config.json").Return(true, nil)

	mockOsService.EXPECT().GetAbsolutePath("config.json").Return("config.json", nil)

	mockConfigsListService.EXPECT().Save(returnConfigsList).Return(nil)
	addCmd := cmd.NewAddCmd(mockAppConfigService, mockInputService, mockConfigsListService, mockOsService)
	err := addCmd.Command.RunE(addCmd.Command, args)
	if err != nil {
		t.Errorf("excpected no error while adding, got %d", err)
	}
}

// test for valid config adding two args
func TestAddCmd_ValidConfigAddTwoParam(t *testing.T) {

	ctrl := gomock.NewController(t)
	// creating mock dependencies
	mockAppConfigService := mocks.NewMockAppConfigService(ctrl)
	mockConfigsListService := mocks.NewMockConfigsListService(ctrl)
	mockInputService := mocks.NewMockInputService(ctrl)
	mockOsService := mocks.NewMockOsService(ctrl)

	args := []string{"config", "config.json"}
	returnAppConfig := services.NewDefaultAppConfig()
	returnConfigsList := services.GetDefaultConfigsList("")
	mockInputService.EXPECT().GetPipedInput().Return("", false)
	mockAppConfigService.EXPECT().Load().Return(returnAppConfig, nil)
	mockConfigsListService.EXPECT().Load().Return(returnConfigsList, nil)

	mockOsService.EXPECT().FileExists("config.json").Return(true, nil)

	mockOsService.EXPECT().GetAbsolutePath("config.json").Return("config.json", nil)

	mockConfigsListService.EXPECT().Save(returnConfigsList).Return(nil)
	addCmd := cmd.NewAddCmd(mockAppConfigService, mockInputService, mockConfigsListService, mockOsService)
	err := addCmd.Command.RunE(addCmd.Command, args)
	if err != nil {
		t.Errorf("excpected no error while adding, got %d", err)
	}
}

// test for valid config adding no args with pipe WIP

func TestAddCmd_ValidConfigAddNoParamWithPipe(t *testing.T) {

	ctrl := gomock.NewController(t)
	// creating mock dependencies
	mockAppConfigService := mocks.NewMockAppConfigService(ctrl)
	mockConfigsListService := mocks.NewMockConfigsListService(ctrl)
	mockInputService := mocks.NewMockInputService(ctrl)
	mockOsService := mocks.NewMockOsService(ctrl)

	args := []string{}
	pipe := "config.json\n"

	mockInputService.EXPECT().GetPipedInput().Return(pipe, true)
	returnAppConfig := services.NewDefaultAppConfig()
	returnConfigsList := services.GetDefaultConfigsList("")

	mockAppConfigService.EXPECT().Load().Return(returnAppConfig, nil)
	mockConfigsListService.EXPECT().Load().Return(returnConfigsList, nil)

	mockOsService.EXPECT().FileExists("config.json\n").Return(true, nil)

	mockOsService.EXPECT().GetAbsolutePath("config.json\n").Return("home/usr/config.json", nil)

	mockConfigsListService.EXPECT().Save(returnConfigsList).Return(nil)
	addCmd := cmd.NewAddCmd(mockAppConfigService, mockInputService, mockConfigsListService, mockOsService)
	err := addCmd.Command.RunE(addCmd.Command, args)
	if err != nil {
		t.Errorf("excpected no error while adding, got %d", err)
	}
}
