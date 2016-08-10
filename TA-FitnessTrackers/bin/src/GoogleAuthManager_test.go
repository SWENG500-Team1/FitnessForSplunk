package main

import (
	"testing"

	splunk "github.com/AndyNortrup/GoSplunk"
)

//Test if the app is able to exchange a stale token for a valid client.
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
