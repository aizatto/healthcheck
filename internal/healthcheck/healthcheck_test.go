package healthcheck

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockConfig struct {
	mock.Mock
}

func (m *MockConfig) offline(t TargetInterface, err error) {
	m.Called(t, err)
}

func (m *MockConfig) online(t TargetInterface) {
	m.Called(t)
}

type MockTarget struct {
	Error  error
	Online bool
	mock.Mock
}

func (m *MockTarget) name() string {
	return ""
}

func (m *MockTarget) healthcheck() error {
	return m.Error
}

func (m *MockTarget) isOnline() bool {
	return m.Online
}

func (m *MockTarget) setOnline(status bool) TargetInterface {
	m.Online = status
	return m
}

func TestHealthcheck(t *testing.T) {
	t.Run("default", func(t *testing.T) {
		assert := assert.New(t)
		err := HealthcheckTarget(&ConfigJSON{}, &MockTarget{})
		assert.NoError(err)
	})

	t.Run("target moves from online to online", func(t *testing.T) {
		assert := assert.New(t)
		config := &MockConfig{}
		target := &MockTarget{
			Online: true,
		}
		err := HealthcheckTarget(config, target)
		config.AssertExpectations(t)
		assert.NoError(err)
	})

	t.Run("target moves from online to offline", func(t *testing.T) {
		assert := assert.New(t)
		config := &MockConfig{}
		target := &MockTarget{
			Online: true,
			Error:  errors.New("error"),
		}
		config.On("offline", target, target.Error)
		err := HealthcheckTarget(config, target)
		config.AssertExpectations(t)
		assert.Error(err, "Error")
		assert.False(target.Online)
	})

	t.Run("target moves from offline to offline", func(t *testing.T) {
		assert := assert.New(t)
		config := &MockConfig{}
		target := &MockTarget{
			Online: false,
			Error:  errors.New("error"),
		}
		err := HealthcheckTarget(config, target)
		config.AssertExpectations(t)
		assert.Error(err, "Error")
		assert.False(target.Online)
	})

	t.Run("target moves from offline to online", func(t *testing.T) {
		assert := assert.New(t)
		config := &MockConfig{}
		target := &MockTarget{
			Online: false,
		}
		config.On("online", target)
		err := HealthcheckTarget(config, target)
		config.AssertExpectations(t)
		assert.NoError(err)
		assert.True(target.Online)
	})
}
