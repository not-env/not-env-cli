package commands

import (
	"fmt"

	"not-env-cli/internal/client"
	"not-env-cli/internal/config"
)

// VarList lists all variables
func VarList() error {
	cfg, err := config.Load()
	if err != nil {
		return err
	}

	cl := client.NewClient(cfg.URL, cfg.APIKey)

	resp, err := cl.Get("/variables")
	if err != nil {
		return err
	}

	if resp.StatusCode != 200 {
		return client.ParseResponse(resp, nil)
	}

	var result struct {
		Variables []struct {
			Key       string `json:"key"`
			Value     string `json:"value"`
			CreatedAt string `json:"created_at"`
			UpdatedAt string `json:"updated_at"`
		} `json:"variables"`
	}

	if err := client.ParseResponse(resp, &result); err != nil {
		return err
	}

	if len(result.Variables) == 0 {
		fmt.Println("No variables found.")
		return nil
	}

	fmt.Println("Variables:")
	for _, v := range result.Variables {
		fmt.Printf("  %s=%s\n", v.Key, v.Value)
	}

	return nil
}

// VarGet gets a single variable
func VarGet(key string) error {
	cfg, err := config.Load()
	if err != nil {
		return err
	}

	cl := client.NewClient(cfg.URL, cfg.APIKey)

	resp, err := cl.Get(fmt.Sprintf("/variables/%s", key))
	if err != nil {
		return err
	}

	if resp.StatusCode != 200 {
		return client.ParseResponse(resp, nil)
	}

	var v struct {
		Key       string `json:"key"`
		Value     string `json:"value"`
		CreatedAt string `json:"created_at"`
		UpdatedAt string `json:"updated_at"`
	}

	if err := client.ParseResponse(resp, &v); err != nil {
		return err
	}

	fmt.Printf("%s=%s\n", v.Key, v.Value)
	return nil
}

// VarSet sets a variable
func VarSet(key, value string) error {
	cfg, err := config.Load()
	if err != nil {
		return err
	}

	cl := client.NewClient(cfg.URL, cfg.APIKey)

	resp, err := cl.Put(fmt.Sprintf("/variables/%s", key), map[string]interface{}{
		"value": value,
	})
	if err != nil {
		return err
	}

	if resp.StatusCode != 204 {
		return client.ParseResponse(resp, nil)
	}

	fmt.Printf("Variable %s set successfully!\n", key)
	return nil
}

// VarDelete deletes a variable
func VarDelete(key string) error {
	cfg, err := config.Load()
	if err != nil {
		return err
	}

	cl := client.NewClient(cfg.URL, cfg.APIKey)

	resp, err := cl.Delete(fmt.Sprintf("/variables/%s", key))
	if err != nil {
		return err
	}

	if resp.StatusCode != 204 {
		return client.ParseResponse(resp, nil)
	}

	fmt.Printf("Variable %s deleted successfully!\n", key)
	return nil
}
