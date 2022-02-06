// Copyright 2020 IBM Corp.
// SPDX-License-Identifier: Apache-2.0

package mockup

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/rs/zerolog/log"

	connectors "fybrik.io/fybrik/pkg/connectors/policymanager/clients"
	"fybrik.io/fybrik/pkg/model/policymanager"
	"fybrik.io/fybrik/pkg/model/taxonomy"
	"fybrik.io/fybrik/pkg/random"
)

// MockPolicyManager is a mock for PolicyManager interface used in tests
type MockPolicyManager struct {
	connectors.PolicyManager
}

func deserializeToTaxonomyAction(action map[string]interface{}, taxAction *taxonomy.Action) error {
	actionBytes, errJSON := json.MarshalIndent(action, "", "\t")
	if errJSON != nil {
		return fmt.Errorf("error Marshalling in deserializeToTaxonomyAction: %v", errJSON)
	}
	err := json.Unmarshal(actionBytes, taxAction)
	if err != nil {
		return fmt.Errorf("error in unmarshalling in deserializeToTaxonomyAction: %v", err)
	}
	return nil
}

// GetPoliciesDecisions implements the PolicyCompiler interface
func (m *MockPolicyManager) GetPoliciesDecisions(input *policymanager.GetPolicyDecisionsRequest, creds string) (*policymanager.GetPolicyDecisionsResponse, error) {
	log.Printf("Received OpenAPI request in mockup GetPoliciesDecisions: ")
	log.Printf("ProcessingGeography: %s", input.Action.ProcessingLocation)
	log.Printf("Destination: " + input.Action.Destination)

	datasetID := string(input.Resource.ID)
	log.Printf("   DataSetID: " + datasetID)
	respResult := []policymanager.ResultItem{}
	policyManagerResult := policymanager.ResultItem{}

	splittedID := strings.SplitN(datasetID, "/", 2)
	if len(splittedID) != 2 {
		panic(fmt.Sprintf("Invalid dataset ID for mock: %s", datasetID))
	}
	assetID := splittedID[1]
	switch assetID {
	case "allow-dataset":
		// empty result simulates allow
		// no need to construct any result item

	case "deny-dataset":
		actionOnDataset := taxonomy.Action{}
		action := make(map[string]interface{})
		action["name"] = "Deny"
		denyAction := map[string]interface{}{}
		action["Deny"] = denyAction

		err := deserializeToTaxonomyAction(action, &actionOnDataset)
		if err != nil {
			log.Print("error in deserializeToTaxonomyAction for scenario deny-dataset :", err)
			return nil, err
		}
		policyManagerResult.Action = actionOnDataset
		respResult = append(respResult, policyManagerResult)

	case "allow-theshire":
		if input.Action.Destination != "theshire" {
			actionOnDataset := taxonomy.Action{}
			action := make(map[string]interface{})
			action["name"] = "Deny"
			denyAction := map[string]interface{}{}
			action["Deny"] = denyAction

			err := deserializeToTaxonomyAction(action, &actionOnDataset)
			if err != nil {
				log.Print("error in deserializeToTaxonomyAction for scenario allow-theshire:", err)
				return nil, err
			}
			policyManagerResult.Action = actionOnDataset
			respResult = append(respResult, policyManagerResult)
		}

	case "deny-write":
		if input.Action.ActionType == policymanager.WRITE {
			actionOnDataset := taxonomy.Action{}
			action := make(map[string]interface{})
			action["name"] = "Deny"
			denyAction := map[string]interface{}{}
			action["Deny"] = denyAction

			err := deserializeToTaxonomyAction(action, &actionOnDataset)
			if err != nil {
				log.Print("error in deserializeToTaxonomyAction for scenario deny-write:", err)
				return nil, err
			}
			policyManagerResult.Action = actionOnDataset
			respResult = append(respResult, policyManagerResult)
		}

	default:
		actionOnCols := taxonomy.Action{}
		action := make(map[string]interface{})
		action["name"] = "RedactAction"
		redactAction := make(map[string]interface{})
		redactAction["columns"] = []string{"SSN"}
		action["RedactAction"] = redactAction

		err := deserializeToTaxonomyAction(action, &actionOnCols)
		if err != nil {
			log.Print("error in deserializeToTaxonomyAction for scenario default:", err)
			return nil, err
		}
		policyManagerResult.Action = actionOnCols
		respResult = append(respResult, policyManagerResult)
	}

	decisionID, _ := random.Hex(20)
	policyManagerResp := &policymanager.GetPolicyDecisionsResponse{DecisionID: decisionID, Result: respResult}

	res, err := json.MarshalIndent(policyManagerResp, "", "\t")
	if err != nil {
		log.Print("error in marshalling policy manager response :", err)
		return nil, err
	}
	log.Print("Marshalled policy manager response in mockup GetPoliciesDecisions:", string(res))

	return policyManagerResp, nil
}
