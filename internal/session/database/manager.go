package database

import (
	"fmt"
	"sync"

	"demodesk/neko/internal/types"
)

func New() *MembersDatabaseCtx {
	return &MembersDatabaseCtx{
		profiles: make(map[string]types.MemberProfile),
		mu:       sync.Mutex{},
	}
}

type MembersDatabaseCtx struct {
	profiles map[string]types.MemberProfile
	mu       sync.Mutex
}

func (manager *MembersDatabaseCtx) Insert(id string, profile types.MemberProfile) error {
	manager.mu.Lock()
	defer manager.mu.Unlock()

	_, ok := manager.profiles[id]
	if ok {
		return fmt.Errorf("Member ID already exists.")
	}

	manager.profiles[id] = profile
	return nil
}

func (manager *MembersDatabaseCtx) Update(id string, profile types.MemberProfile) error {
	manager.mu.Lock()
	defer manager.mu.Unlock()

	_, ok := manager.profiles[id]
	if !ok {
		return fmt.Errorf("Member ID does not exist.")
	}

	manager.profiles[id] = profile
	return nil
}

func (manager *MembersDatabaseCtx) Delete(id string) error {
	manager.mu.Lock()
	defer manager.mu.Unlock()

	_, ok := manager.profiles[id]
	if !ok {
		return fmt.Errorf("Member ID does not exist.")
	}

	delete(manager.profiles, id)
	return nil
}

func (manager *MembersDatabaseCtx) Select(id string) (types.MemberProfile, bool) {
	manager.mu.Lock()
	defer manager.mu.Unlock()

	profile, ok := manager.profiles[id]
	return profile, ok
}
