package qualys

import (
	"fmt"
	"net/http"
	"reflect"
	"testing"
)

func TestListAssetGroups(t *testing.T) {

	cases := []struct {
		name     string
		response string
		expected []AssetGroup
		opts     *ListAssetGroupOptions
		isErr    bool
	}{
		{
			name:     "ListAssetGroups - single item, without list options",
			response: assetGroupsXMLSingleGroup,
			expected: []AssetGroup{
				{
					ID:    "1759735",
					Title: "AG - Elastic Cloud Dynamic Perimeter",
					IPs: AssetGroupIPs{
						IPs:      []string{"10.1.1.1", "10.10.10.11"},
						IPRanges: nil,
					},
				},
			},
			opts: nil,
		},
		{
			name:     "ListAssetGroups - single item, with list options",
			response: assetGroupsXMLSingleGroup,
			expected: []AssetGroup{
				{
					ID:    "1759735",
					Title: "AG - Elastic Cloud Dynamic Perimeter",
					IPs: AssetGroupIPs{
						IPs:      []string{"10.1.1.1", "10.10.10.11"},
						IPRanges: nil,
					},
				},
			},
			opts: &ListAssetGroupOptions{Ids: []string{}},
		},
		{
			name:     "ListAssetGroups - multi item",
			response: assetGroupsXMLMultiGroups,
			expected: []AssetGroup{
				{ID: "1759734", Title: "AG - New"},
				{ID: "1759735", Title: "AG - Elastic Cloud Dynamic Perimeter",
					IPs: AssetGroupIPs{
						IPs:      []string{"10.10.10.14"},
						IPRanges: []string{"10.10.10.3-10.10.10.6"},
					},
				},
			},
			opts: &ListAssetGroupOptions{Ids: []string{"1", "2"}},
		},
	}

	for _, c := range cases {
		setup()
		defer teardown()
		mux.HandleFunc("/asset/group/", func(w http.ResponseWriter, r *http.Request) {
			testMethod(t, r, "GET")
			fmt.Fprint(w, c.response)
		})

		assetGroups, _, err := client.Assets.ListAssetGroups(c.opts)
		if err != nil {
			t.Errorf("Assets.ListAssetGroups returned error: %v", err)
		}

		if !reflect.DeepEqual(assetGroups, c.expected) {
			t.Errorf("Assets.ListAssetGroups case: %s returned %+v, expected %+v", c.name, assetGroups, c.expected)
		}
	}
}

func TestGetAssetGroupByID(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/asset/group/", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, assetGroupsXMLSingleGroup)
	})

	groupID := "1759735"

	assetGroup, _, err := client.Assets.GetAssetGroupByID(groupID)
	if err != nil {
		t.Errorf("Assets.GetAssetGroupByID(%s) returned error: %v", groupID, err)
	}

	expected := &AssetGroup{
		ID:    "1759735",
		Title: "AG - Elastic Cloud Dynamic Perimeter",
		IPs: AssetGroupIPs{
			IPs:      []string{"10.1.1.1", "10.10.10.11"},
			IPRanges: nil,
		},
	}
	if !reflect.DeepEqual(assetGroup, expected) {
		t.Errorf("Assets.GetAssetGroupByID(%s) returned %+v, expected %+v", groupID, assetGroup, expected)
	}
}

func TestAddIPsToGroup(t *testing.T) {
	setup()
	defer teardown()

	groupID := "1759735"
	ip := "10.10.10.10"

	mux.HandleFunc("/asset/group/", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		if r.FormValue("add_ips") != ip {
			t.Errorf("Request form data did not include the correct IP")
		}
		if r.FormValue("id") != groupID {
			t.Errorf("Request form data did not include the correct asset group ID")
		}
		fmt.Fprint(w, assetGroupsAddIPsResponse)
	})
	opts := &AddIPsToGroupOptions{
		GroupID: groupID,
		IPs:     []string{ip},
	}

	_, err := client.Assets.AddIPsToGroup(opts)
	if err != nil {
		t.Errorf("Assets.AddIPsToGroup returned error: %v", err)
	}
}

func TestAssetGroupContainsIP(t *testing.T) {
	cases := []struct {
		name     string
		ip       string
		group    *AssetGroup
		expected bool
	}{
		{
			name:     "AssetGroup.ContainsIP - nil",
			ip:       "10.1.1.1",
			group:    &AssetGroup{ID: "1759735", Title: "AG - Elastic Cloud Dynamic Perimeter"},
			expected: false,
		},
		{
			name: "AssetGroup.ContainsIP - empty",
			ip:   "10.1.1.1",
			group: &AssetGroup{
				ID:    "1759735",
				Title: "AG - Elastic Cloud Dynamic Perimeter",
				IPs:   AssetGroupIPs{}},
			expected: false,
		},
		{
			name: "AssetGroup.ContainsIP - single item list",
			ip:   "10.1.1.1",
			group: &AssetGroup{
				ID:    "1759735",
				Title: "AG - Elastic Cloud Dynamic Perimeter",
				IPs: AssetGroupIPs{
					IPs:      []string{"10.1.1.1"},
					IPRanges: []string{},
				},
			},
			expected: true,
		},
		{
			name: "AssetGroup.ContainsIP - multi item list",
			ip:   "10.1.1.1",
			group: &AssetGroup{
				ID:    "1759735",
				Title: "AG - Elastic Cloud Dynamic Perimeter",
				IPs: AssetGroupIPs{
					IPs:      []string{"10.1.1.1"},
					IPRanges: []string{"10.10.1.1-10.10.10.10"},
				},
			},
			expected: true,
		},
	}
	for _, c := range cases {
		contains := c.group.ContainsIP(c.ip)
		if contains != c.expected {
			t.Errorf("%s - AssetGroup.ContainsIP(%s) returned %v, expected %v", c.name, c.ip, contains, c.expected)
		}
	}
}

func TestAssetGroupIPsContainsIP(t *testing.T) {
	group := AssetGroupIPs{IPs: []string{"10.0.1.1"}, IPRanges: []string{"10.10.10.3-10.10.10.6"}}

	cases := []struct {
		name     string
		ip       string
		group    AssetGroupIPs
		expected bool
	}{
		{
			name:     "AssetGroupIPs.ContainsIP - IP value match",
			ip:       "10.0.1.1",
			group:    group,
			expected: true,
		},
		{
			name:     "AssetGroupIPs.ContainsIP - IP value no match",
			ip:       "192.0.1.1",
			group:    group,
			expected: false,
		},
		{
			name:     "AssetGroupIPs.ContainsIP - IP Range value match",
			ip:       "10.10.10.4",
			group:    group,
			expected: true,
		},
		{
			name:     "AssetGroupIPs.ContainsIP - IP Range value no match",
			ip:       "10.10.10.1",
			group:    group,
			expected: false,
		},
		{
			name:     "AssetGroupIPs.ContainsIP - IP Range value match",
			ip:       "10.10.0.4",
			group:    AssetGroupIPs{IPs: []string{"10.0.1.1"}, IPRanges: []string{"10.10.0.0-10.10.10.6"}},
			expected: true,
		},
		{
			name:     "AssetGroupIPs.ContainsIP - IP Range value no match",
			ip:       "10.10.0.4",
			group:    AssetGroupIPs{IPs: []string{"10.0.1.1"}, IPRanges: []string{"10.10.1.3-10.10.10.6"}},
			expected: false,
		},
	}

	for _, c := range cases {
		contains := c.group.ContainsIP(c.ip)
		if contains != c.expected {
			t.Errorf("%s - AssetGroupIPs.ContainsIP(%s) returned %v, expected %v", c.name, c.ip, contains, c.expected)
		}
	}
}
