package sentry

import (
	"fmt"
	"time"
)

const (
	// OrgEndpointName is set to roganizations
	OrgEndpointName = "organizations"
)

// Quota is your quote for a project limit and max rate
type Quota struct {
	ProjectLimit int `json:"projectLimit"`
	MaxRate      int `json:"maxRate"`
}

// Organization is your sentry organization.
type Organization struct {
	PendingAccessRequest *int       `json:"pendingAccessRequests,omitempty"`
	Slug                 *string    `json:"slug,omitempty"`
	Name                 string     `json:"name"`
	Quota                *Quota     `json:"quota,omitempty"`
	DateCreated          *time.Time `json:"dateCreated,omitempty"`
	Teams                *[]Team    `json:"teams,omitempty"`
	ID                   *string    `json:"id,omitempty"`
	IsEarlyAdopter       *bool      `json:"isEarlyAdopter,omitempty"`
	Features             *[]string  `json:"features,omitempty"`
}

// SCIMResource represents the structure for SCIM API response.
type SCIMResource struct {
	Schemas      []string   `json:"schemas"`
	TotalResults int        `json:"totalResults"`
	StartIndex   int        `json:"startIndex"`
	ItemsPerPage int        `json:"itemsPerPage"`
	Resources    []Resource `json:"Resources"`
}

// Resource represents the individual resources within the SCIM API response.
type Resource struct {
	Schemas     []string     `json:"schemas"`
	ID          string       `json:"id"`
	DisplayName string       `json:"displayName"`
	Meta        Meta         `json:"meta"`
	Members     []SCIMMember `json:"members,omitempty"`
}

// Meta provides metadata about SCIM resources.
type Meta struct {
	ResourceType string `json:"resourceType"`
}

// Member represents the members of a SCIM resource group.
type SCIMMember struct {
	Value   string `json:"value"`
	Display string `json:"display"`
}

// GetOrganization takes a org slug and returns back the org
func (c *Client) GetOrganization(orgslug string) (Organization, error) {
	var org Organization

	err := c.do("GET", fmt.Sprintf("%s/%s", OrgEndpointName, orgslug), &org, nil)
	return org, err
}

// GetOrganizations will return back every organization in the sentry instance
func (c *Client) GetOrganizations() ([]Organization, *Link, error) {
	orgs := make([]Organization, 0)
	link, err := c.doWithPagination("GET", OrgEndpointName, &orgs, nil)
	return orgs, link, err
}

// CreateOrganization creates a organization with a name
func (c *Client) CreateOrganization(orgname string) (Organization, error) {
	var org Organization
	orgreq := &Organization{
		Name: orgname,
	}
	err := c.do("POST", OrgEndpointName, &org, orgreq)
	return org, err
}

// UpdateOrganization takes a organization and updates it on the server side
func (c *Client) UpdateOrganization(org Organization) error {
	return c.do("PUT", fmt.Sprintf("%s/%s", OrgEndpointName, *org.Slug), &org, &org)
}

// DeleteOrganization will delete the Org. There is not way to revert this if you do.
func (c *Client) DeleteOrganization(org Organization) error {
	return c.do("DELETE", fmt.Sprintf("%s/%s", OrgEndpointName, *org.Slug), nil, nil)
}

// GetOrganizationTeams will fetch all teams for this org
func (c *Client) GetOrganizationTeams(o Organization) ([]Team, error) {
	teams := make([]Team, 0)
	err := c.do("GET", fmt.Sprintf("organizations/%s/teams", *o.Slug), &teams, nil)
	return teams, err
}

// GetOrganizationTeams will fetch all teams for this org
func (c *Client) GetOrganizationMembers(o Organization) ([]Member, error) {
	members := make([]Member, 0)
	err := c.do("GET", fmt.Sprintf("organizations/%s/members", *o.Slug), &members, nil)
	return members, err
}

// GetTeamProjects fetchs all projects for a Team
func (c *Client) GetOrganizationTeamsSCIM(o Organization) (SCIMResource, error) {
	// members := make([]Member, 0)
	scimResponse := SCIMResource{}
	err := c.do("GET", fmt.Sprintf("organizations/%s/scim/v2/Groups", *o.Slug), &scimResponse, nil)
	return scimResponse, err
}
