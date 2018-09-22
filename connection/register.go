package connection

import (
	"fmt"
	"log"
)

const (
	requestTypeRegister    = "register"
	pairingTypePrompt      = "PROMPT"
	responseTypeRegistered = "registered"
)

type registerRequestPayload struct {
	PairingType string           `json:"pairingType"`
	Manifest    registerManifest `json:"manifest"`
	ClientKey   string           `json:"client-key"`
}

type registerManifest struct {
	Permissions []string `json:"permissions"`
}

type registeredPayload struct {
	Key string `json:"client-key"`
}

// Register sends registration request to the LG TV
func (c *Connection) Register(id int, key string, permissions []string) error {
	c.addResponseChannel(id)
	defer c.removeResponseChannel(id)

	payload := registerRequestPayload{
		PairingType: pairingTypePrompt,
		Manifest: registerManifest{
			Permissions: permissions,
		},
		ClientKey: key,
	}

	if err := c.sendRequest(id, requestTypeRegister, "", payload); err != nil {
		return err
	}

	resp, err := c.waitForResponse(id, nil, c.timeouts.Response)
	if err != nil {
		return fmt.Errorf("Registering on LG TV at %s using key %s returned an error: %s", c.url, key, err)
	}
	if resp.Type == responseTypeRegistered {
		return nil
	}

	result := registeredPayload{}
	resp, err = c.waitForResponse(id, &result, c.timeouts.Register)
	if err != nil {
		return fmt.Errorf("Registering on LG TV at %s using key %s returned an error while waiting for response: %s", c.url, key, err)
	}

	if key == "" {
		log.Printf("New key received from LG TV at %s, please use it on the next connection: %s", c.url, result.Key)
	}

	return nil
}
