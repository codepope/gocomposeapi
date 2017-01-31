// Copyright 2016 Compose, an IBM Company
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package composeapi

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/parnurzeal/gorequest"
	"log"
	"strconv"
)

var ()

const (
	apibase = "https://api.compose.io/2016-07/"
)

type Client struct {
	apiToken string
}

func NewClient(apiToken string) (*Client, error) {
	return &Client{
		apiToken: apiToken,
	}, nil
}

// Link structure for JSON+HAL links
type Link struct {
	HREF      string `json:"href"`
	Templated bool   `json:"templated"`
}

//Errors struct for parsing error returns
type Errors struct {
	Error string `json:"errors,omitempty"`
}

func printJSON(jsontext string) {
	var tempholder map[string]interface{}

	if err := json.Unmarshal([]byte(jsontext), &tempholder); err != nil {
		log.Fatal(err)
	}
	indentedjson, err := json.MarshalIndent(tempholder, "", "  ")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(indentedjson))
}

//SetAPIToken overrides the API token
func (c *Client) SetAPIToken(newtoken string) {
	c.apiToken = newtoken
}

//GetJSON Gets JSON string of content at an endpoint
func (c *Client) getJSON(endpoint string) (string, []error) {
	response, body, errs := gorequest.New().Get(apibase+endpoint).
		Set("Authorization", "Bearer "+c.apiToken).
		Set("Content-type", "json").
		End()
	if response.StatusCode != 200 {
		myerrors := Errors{}
		err := json.Unmarshal([]byte(body), &myerrors)
		if err != nil {
			errs = append(errs, errors.New("Unable to parse error - status code "+strconv.Itoa(response.StatusCode)))
		} else {
			errs = append(errs, errors.New(myerrors.Error))
		}
	}
	return body, errs
}

//GetAccountJSON gets JSON string from endpoint
func (c *Client) GetAccountJSON() (string, []error) { return c.getJSON("accounts") }

//GetAccount Gets first Account struct from account endpoint
func (c *Client) GetAccount() (*Account, []error) {
	body, errs := c.GetAccountJSON()

	if errs != nil {
		return nil, errs
	}

	accountsResponse := accountResponse{}
	json.Unmarshal([]byte(body), &accountsResponse)
	firstAccount := accountsResponse.Embedded.Accounts[0]

	return &firstAccount, nil
}

//GetDeploymentsJSON returns raw deployment
func (c *Client) GetDeploymentsJSON() (string, []error) { return c.getJSON("deployments") }

//GetDeployments returns deployment structure
func (c *Client) GetDeployments() (*[]Deployment, []error) {
	body, errs := c.GetDeploymentsJSON()

	if errs != nil {
		return nil, errs
	}

	deploymentResponse := deploymentsResponse{}
	json.Unmarshal([]byte(body), &deploymentResponse)
	deployments := deploymentResponse.Embedded.Deployments

	return &deployments, nil
}

//GetDeploymentJSON returns raw deployment
func (c *Client) GetDeploymentJSON(deploymentid string) (string, []error) {
	return c.getJSON("deployments/" + deploymentid)
}

//GetDeployment returns deployment structure
func (c *Client) GetDeployment(deploymentid string) (*Deployment, []error) {
	body, errs := c.GetDeploymentJSON(deploymentid)

	if errs != nil {
		return nil, errs
	}

	deployment := Deployment{}
	json.Unmarshal([]byte(body), &deployment)

	return &deployment, nil
}

//GetScalingsJSON returns raw scalings
func (c *Client) GetScalingsJSON(deploymentid string) (string, []error) {
	return c.getJSON("deployments/" + deploymentid + "/scalings")
}

//GetScalings returns deployment structure
func (c *Client) GetScalings(deploymentid string) (*Scalings, []error) {
	body, errs := c.GetScalingsJSON(deploymentid)

	if errs != nil {
		return nil, errs
	}

	scalings := Scalings{}
	json.Unmarshal([]byte(body), &scalings)

	return &scalings, nil
}

//SetScalingsJSON sets JSON scaling and returns string respones
func (c *Client) SetScalingsJSON(params ScalingsParams) (string, []error) {
	response, body, errs := gorequest.New().Post(apibase+"deployments/"+params.DeploymentID+"/scalings").
		Set("Authorization", "Bearer "+c.apiToken).
		Set("Content-type", "application/json; charset=utf-8").
		Send(params).
		End()

	if response.StatusCode != 200 { // Expect Accepted on success - assume error on anything else
		myerrors := Errors{}
		err := json.Unmarshal([]byte(body), &myerrors)
		if err != nil {
			errs = append(errs, errors.New("Unable to parse error - status code "+strconv.Itoa(response.StatusCode)))
		} else {
			errs = append(errs, errors.New(myerrors.Error))
		}
	}

	return body, errs
}

//SetScalings sets scale and returns recipe for scaling
func (c *Client) SetScalings(scalingsParams ScalingsParams) (*Recipe, []error) {
	body, errs := c.SetScalingsJSON(scalingsParams)
	if errs != nil {
		return nil, errs
	}

	recipe := Recipe{}
	json.Unmarshal([]byte(body), &recipe)

	return &recipe, nil
}

//GetRecipeJSON Gets raw JSON for recipeid
func (c *Client) GetRecipeJSON(recipeid string) (string, []error) {
	return c.getJSON("recipes/" + recipeid)
}

//GetRecipe gets status of Recipe
func (c *Client) GetRecipe(recipeid string) (*Recipe, []error) {
	body, errs := c.GetRecipeJSON(recipeid)

	if errs != nil {
		return nil, errs
	}

	recipe := Recipe{}
	json.Unmarshal([]byte(body), &recipe)

	return &recipe, nil
}

//GetRecipesForDeploymentJSON returns raw JSON for getRecipesforDeployment
func (c *Client) GetRecipesForDeploymentJSON(deploymentid string) (string, []error) {
	return c.getJSON("deployments/" + deploymentid + "/recipes")
}

//GetRecipesForDeployment gets deployment recipe life
func (c *Client) GetRecipesForDeployment(deploymentid string) (*[]Recipe, []error) {
	body, errs := c.GetRecipesForDeploymentJSON(deploymentid)

	if errs != nil {
		return nil, errs
	}

	recipeResponse := Recipe{}
	json.Unmarshal([]byte(body), &recipeResponse)
	recipes := recipeResponse.Embedded.Recipes

	return &recipes, nil
}

//GetVersionsForDeploymentJSON returns raw JSON for getVersionsforDeployment
func (c *Client) GetVersionsForDeploymentJSON(deploymentid string) (string, []error) {
	return c.getJSON("deployments/" + deploymentid + "/versions")
}

//GetVersionsForDeployment gets deployment recipe life
func (c *Client) GetVersionsForDeployment(deploymentid string) (*[]VersionTransition, []error) {
	body, errs := c.GetVersionsForDeploymentJSON(deploymentid)

	if errs != nil {
		return nil, errs
	}

	versionsResponse := versionsResponse{}
	json.Unmarshal([]byte(body), &versionsResponse)
	versionTransitions := versionsResponse.Embedded.VersionTransitions

	return &versionTransitions, nil
}

//GetClustersJSON gets clusters available
func (c *Client) GetClustersJSON() (string, []error) {
	return c.getJSON("clusters")
}

//GetClusters gets clusters available
func (c *Client) GetClusters() (*[]Cluster, []error) {
	body, errs := c.GetClustersJSON()

	if errs != nil {
		return nil, errs
	}

	clustersResponse := clusterResponse{}
	json.Unmarshal([]byte(body), &clustersResponse)
	clusters := clustersResponse.Embedded.Clusters

	return &clusters, nil
}

//GetDatacentersJSON gets datacenters available as a string
func (c *Client) GetDatacentersJSON() (string, []error) {
	return c.getJSON("datacenters")
}

//GetDatacenters gets datacenters available as a Go struct
func (c *Client) GetDatacenters() (*[]Datacenter, []error) {
	body, errs := c.GetDatacentersJSON()

	if errs != nil {
		return nil, errs
	}

	datacenterResponse := datacentersResponse{}
	json.Unmarshal([]byte(body), &datacenterResponse)
	datacenters := datacenterResponse.Embedded.Datacenters

	return &datacenters, nil
}

//GetDatabasesJSON gets databases available as a string
func (c *Client) GetDatabasesJSON() (string, []error) {
	return c.getJSON("databases")
}

//GetDatabases gets databases available as a Go struct
func (c *Client) GetDatabases() (*[]Database, []error) {
	body, errs := c.GetDatabasesJSON()

	if errs != nil {
		return nil, errs
	}

	databaseResponse := databasesResponse{}
	json.Unmarshal([]byte(body), &databaseResponse)
	databases := databaseResponse.Embedded.Databases

	return &databases, nil
}

//GetUserJSON returns user JSON string
func (c *Client) GetUserJSON() (string, []error) {
	return c.getJSON("user")
}

//GetUser Gets information about user
func (c *Client) GetUser() (*User, []error) {
	body, errs := c.GetUserJSON()

	if errs != nil {
		return nil, errs
	}

	user := User{}
	json.Unmarshal([]byte(body), &user)
	return &user, nil
}

//CreateDeploymentJSON performs the call
func (c *Client) CreateDeploymentJSON(params CreateDeploymentParams) (string, []error) {
	response, body, errs := gorequest.New().Post(apibase+"deployments").
		Set("Authorization", "Bearer "+c.apiToken).
		Set("Content-type", "application/json; charset=utf-8").
		Send(params).
		End()

	if response.StatusCode != 202 { // Expect Accepted on success - assume error on anything else
		myerrors := Errors{}
		err := json.Unmarshal([]byte(body), &myerrors)
		if err != nil {
			errs = append(errs, errors.New("Unable to parse error - status code "+strconv.Itoa(response.StatusCode)))
		} else {
			errs = append(errs, errors.New(myerrors.Error))
		}
	}

	return body, errs
}

//CreateDeployment creates a deployment
func (c *Client) CreateDeployment(params CreateDeploymentParams) (*Deployment, []error) {

	// This is a POST not a GET, so it builds its own request

	body, errs := c.CreateDeploymentJSON(params)

	if errs != nil {
		return nil, errs
	}

	deployed := Deployment{}
	json.Unmarshal([]byte(body), &deployed)

	return &deployed, nil
}
