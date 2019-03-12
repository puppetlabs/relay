package state

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/puppetlabs/nebula/pkg/errors"
)

type Manager interface {
	Load(name string) (*Resource, errors.Error)
	Save(resource *Resource) errors.Error
}

// State is a temporary state management model for the prototype-1 demo
type State struct {
	Resources map[string]*Resource `json:"resources"`
}

type Resource struct {
	Name  string          `json:"name"`
	Value json.RawMessage `json:"value"`
}

type FilesystemStateManager struct {
	path string
}

func (fsm FilesystemStateManager) loadState() (*State, errors.Error) {
	var state State

	f, err := os.Open(fsm.path)
	if err != nil {
		if os.IsNotExist(err) {
			state.Resources = make(map[string]*Resource)

			return &state, nil
		}

		return nil, errors.NewStateLoadError().WithCause(err)
	}

	defer f.Close()

	if err := json.NewDecoder(f).Decode(&state); err != nil {
		return nil, errors.NewStateLoadError().WithCause(err)
	}

	return &state, nil
}

func (fsm FilesystemStateManager) Load(name string) (*Resource, errors.Error) {
	state, err := fsm.loadState()
	if err != nil {
		return nil, err
	}

	resource, ok := state.Resources[name]
	if !ok {
		return nil, errors.NewStateResourceNotExists(name)
	}

	return resource, nil
}

func (fsm FilesystemStateManager) Save(resource *Resource) errors.Error {
	state, err := fsm.loadState()
	if err != nil {
		if !errors.IsStateResourceNotExists(err) {
			return err
		}

		state = &State{}
	}

	state.Resources[resource.Name] = resource

	f, ferr := os.OpenFile(fsm.path, os.O_RDWR|os.O_CREATE, 0755)
	if ferr != nil {
		return errors.NewStateSaveError().WithCause(ferr).Bug()
	}

	defer f.Close()

	if err := json.NewEncoder(f).Encode(state); err != nil {
		return errors.NewStateSaveError().WithCause(err)
	}

	return nil
}

func NewFilesystemStateManager(path string) (*FilesystemStateManager, errors.Error) {
	if path == "" {
		path = filepath.Join(os.Getenv("HOME"), ".local", "share", "nebula", "state.json")
	}

	dataDir, _ := filepath.Split(path)

	// TODO is a quick hack to make sure the datadir exists
	if err := os.MkdirAll(dataDir, 0755); err != nil {
		return nil, errors.NewStateUnknownError().WithCause(err)
	}

	return &FilesystemStateManager{
		path: path,
	}, nil
}
