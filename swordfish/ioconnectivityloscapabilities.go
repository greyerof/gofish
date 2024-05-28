//
// SPDX-License-Identifier: BSD-3-Clause
//

package swordfish

import (
	"encoding/json"
	"reflect"

	"github.com/stmcginnis/gofish/common"
)

// IOConnectivityLoSCapabilities describes capabilities of the system to
// support various IO Connectivity service options.
type IOConnectivityLoSCapabilities struct {
	common.Entity

	// ODataContext is the odata context.
	ODataContext string `json:"@odata.context"`
	// ODataType is the odata type.
	ODataType string `json:"@odata.type"`
	// Description provides a description of this resource.
	Description string
	// Identifier shall be unique within the managed ecosystem.
	Identifier common.Identifier
	// MaxSupportedBytesPerSecond shall be the maximum bytes per second that a
	// connection can support.
	MaxSupportedBytesPerSecond int64
	// MaxSupportedIOPS shall be the maximum IOPS that a connection can support.
	MaxSupportedIOPS int
	// Oem shall contain the OEM extensions. All values for properties that this object contains shall conform to the
	// Redfish Specification-described requirements.
	OEM json.RawMessage `json:"Oem"`
	// SupportedAccessProtocols is Access protocols supported by this service
	// option. NOTE: SMB+NFS* requires that SMB and at least one of NFSv3 or
	// NFXv4 are also selected, (i.e. {'SMB', 'NFSv4', 'SMB+NFS*'}).
	SupportedAccessProtocols []common.Protocol
	// SupportedLinesOfService shall contain known and
	// supported IOConnectivityLinesOfService.
	SupportedLinesOfService []IOConnectivityLineOfService
	// SupportedLinesOfServiceCount is the number of IOConnectivityLineOfServices.
	SupportedLinesOfServiceCount int `json:"SupportedLinesOfService@odata.count"`
	// rawData holds the original serialized JSON so we can compare updates.
	rawData []byte
}

// UnmarshalJSON unmarshals a IOConnectivityLoSCapabilities object from the raw JSON.
func (ioconnectivityloscapabilities *IOConnectivityLoSCapabilities) UnmarshalJSON(b []byte) error {
	type temp IOConnectivityLoSCapabilities
	var t struct {
		temp
	}

	err := json.Unmarshal(b, &t)
	if err != nil {
		return err
	}

	*ioconnectivityloscapabilities = IOConnectivityLoSCapabilities(t.temp)

	// Extract the links to other entities for later

	// This is a read/write object, so we need to save the raw object data for later
	ioconnectivityloscapabilities.rawData = b

	return nil
}

// Update commits updates to this object's properties to the running system.
func (ioconnectivityloscapabilities *IOConnectivityLoSCapabilities) Update() error {
	// Get a representation of the object's original state so we can find what
	// to update.
	original := new(IOConnectivityLoSCapabilities)
	err := original.UnmarshalJSON(ioconnectivityloscapabilities.rawData)
	if err != nil {
		return err
	}

	readWriteFields := []string{
		"MaxSupportedBytesPerSecond",
		"MaxSupportedIOPS",
		"SupportedAccessProtocols",
		"SupportedLinesOfService",
	}

	originalElement := reflect.ValueOf(original).Elem()
	currentElement := reflect.ValueOf(ioconnectivityloscapabilities).Elem()

	return ioconnectivityloscapabilities.Entity.Update(originalElement, currentElement, readWriteFields)
}

// GetIOConnectivityLoSCapabilities will get a IOConnectivityLoSCapabilities
// instance from the service.
func GetIOConnectivityLoSCapabilities(c common.Client, uri string) (*IOConnectivityLoSCapabilities, error) {
	var ioConnectivityLoSCapabilities IOConnectivityLoSCapabilities
	return &ioConnectivityLoSCapabilities, ioConnectivityLoSCapabilities.Get(c, uri, &ioConnectivityLoSCapabilities)
}

// ListReferencedIOConnectivityLoSCapabilitiess gets the collection of
// IOConnectivityLoSCapabilities from a provided reference.
func ListReferencedIOConnectivityLoSCapabilitiess(c common.Client, link string) ([]*IOConnectivityLoSCapabilities, error) { //nolint:dupl
	if link == "" {
		return nil, nil
	}

	type GetResult struct {
		Item  *IOConnectivityLoSCapabilities
		Link  string
		Error error
	}

	ch := make(chan GetResult)
	collectionError := common.NewCollectionError()
	get := func(link string) {
		ioconnectivityloscapabilities, err := GetIOConnectivityLoSCapabilities(c, link)
		ch <- GetResult{Item: ioconnectivityloscapabilities, Link: link, Error: err}
	}

	var links []string
	var err error
	go func() {
		links, err = common.CollectList(get, c, link)
		if err != nil {
			collectionError.Failures[link] = err
		}
		close(ch)
	}()

	// Save unordered results into link-to-IOConnectivityLoSCapabilities helper map.
	unorderedResults := map[string]*IOConnectivityLoSCapabilities{}
	for r := range ch {
		if r.Error != nil {
			collectionError.Failures[r.Link] = r.Error
		} else {
			unorderedResults[r.Link] = r.Item
		}
	}

	if !collectionError.Empty() {
		return nil, collectionError
	}
	// Build the final ordered slice based on the original order from the links list.
	results := make([]*IOConnectivityLoSCapabilities, len(links))
	for i, link := range links {
		results[i] = unorderedResults[link]
	}

	return results, nil
}
