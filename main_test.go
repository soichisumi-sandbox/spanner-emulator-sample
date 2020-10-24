package main

import (
	"context"
	"fmt"
	"log"
	"testing"

	"go.uber.org/zap"

	"github.com/GincoInc/go-util/logger"
	"google.golang.org/api/iterator"

	database "cloud.google.com/go/spanner/admin/database/apiv1"
	instance "cloud.google.com/go/spanner/admin/instance/apiv1"
	instancepb "google.golang.org/genproto/googleapis/spanner/admin/instance/v1"
)

var (
	instanceAdmin *instance.InstanceAdminClient
	databaseAdmin *database.DatabaseAdminClient

	testProjectID  = "test-projectID"
	testInstanceID = "test-instance"
)

func cleanupInstances() {
	logger.Info("cleanup")
	if instanceAdmin == nil {
		// Integration tests skipped.
		return
	}

	ctx := context.Background()
	parent := fmt.Sprintf("projects/%v", testProjectID)
	iter := instanceAdmin.ListInstances(ctx, &instancepb.ListInstancesRequest{
		Parent: parent,
		Filter: "name:*",
	})

	for {
		inst, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			panic(err)
		}

		if err := instanceAdmin.DeleteInstance(ctx, &instancepb.DeleteInstanceRequest{Name: inst.Name}); err != nil {
			log.Printf("failed to delete instance %s (error %v), might need a manual removal",
				inst.Name, err)
		}
		log.Printf("instance %s deleted", inst.Name)
	}
}

func setup() {
	ctx := context.Background()
	instanceClient, err := instance.NewInstanceAdminClient(ctx)
	if err != nil {
		logger.Fatal(err.Error(), zap.Error(err))
	}
	instanceAdmin = instanceClient

	cleanupInstances()

	configIterator := instanceAdmin.ListInstanceConfigs(ctx, &instancepb.ListInstanceConfigsRequest{
		Parent: fmt.Sprintf("projects/%s", testProjectID),
	})

	config, err := configIterator.Next()
	if err != nil {
		logger.Fatal(err.Error())
	}

	op, err := instanceAdmin.CreateInstance(ctx, &instancepb.CreateInstanceRequest{
		Parent:     fmt.Sprintf("projects/%s", testProjectID),
		InstanceId: testInstanceID,
		Instance: &instancepb.Instance{
			Config:      config.Name,
			DisplayName: "testdb",
			NodeCount:   1,
		},
	})
	if err != nil {
		logger.Fatal(err.Error(), zap.Error(err))
	}
	i, err := op.Wait(ctx)
	if err != nil {
		logger.Fatal(err.Error(), zap.Error(err))
	}
	if i.State != instancepb.Instance_READY {
		log.Fatalf("could not create instance with id %s: %v", fmt.Sprintf("projects/%s/instanceIterator/%s", testProjectID, testInstanceID), err)
	}
	fmt.Printf("instance %s created\n", i.Name)
}

func teardown() {
	cleanupInstances()
}

func TestMain(m *testing.M) {
	cleanupInstances()
	setup()
	defer teardown()

	m.Run()
}

func Test_SpannerEmulator(t *testing.T) {
	tests := []struct {
		name string
	}{
		// TODO: test cases
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

		})
	}
}
