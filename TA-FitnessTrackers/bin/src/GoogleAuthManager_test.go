package main

import (
	"testing"

	splunk "github.com/AndyNortrup/GoSplunk"
)

// TestGetUsersFromKVStore tests if users can be retrieved from the Splunk KV
// store.
func TestGetUsersFromKVStore(t *testing.T) {

	sessionKey, err := splunk.NewSessionKey(accountName,
		password, splunk.LocalSplunkMgmntURL)

	if err != nil {
		t.Fatalf("Unable to get session key: %v", err)
	}

	users, err := getUsers(splunk.LocalSplunkMgmntURL,
		sessionKey.SessionKey,
		strategyGoogle)

	if len(users) == 0 {
		t.Fail()
		t.Error("No users returned from KV Store")
	}

}
