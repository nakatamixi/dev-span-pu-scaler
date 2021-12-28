package spanner

import (
	"context"
	"errors"
	"fmt"

	database "cloud.google.com/go/spanner/admin/database/apiv1"
	instance "cloud.google.com/go/spanner/admin/instance/apiv1"
	"google.golang.org/api/iterator"
	databasepb "google.golang.org/genproto/googleapis/spanner/admin/database/v1"
	instancepb "google.golang.org/genproto/googleapis/spanner/admin/instance/v1"
	fieldmaskpb "google.golang.org/protobuf/types/known/fieldmaskpb"
)

type Client struct {
	projectID, instanceID string
	dbAdmin               *database.DatabaseAdminClient
	instanceAdmin         *instance.InstanceAdminClient
}

func NewClient(ctx context.Context, projectID, instanceID string) (*Client, error) {
	iAdmin, err := instance.NewInstanceAdminClient(ctx)
	if err != nil {
		return nil, err
	}
	dAdmin, err := database.NewDatabaseAdminClient(ctx)
	if err != nil {
		return nil, err
	}
	return &Client{
		projectID:     projectID,
		instanceID:    instanceID,
		instanceAdmin: iAdmin,
		dbAdmin:       dAdmin,
	}, nil
}

func (s *Client) Close() error {
	err := s.dbAdmin.Close()
	if err != nil {
		return err
	}
	return s.instanceAdmin.Close()
}

func (s *Client) FQDN() string {
	return fmt.Sprintf("projects/%s/instances/%s", s.projectID, s.instanceID)
}

func (s *Client) Instance(ctx context.Context) (*instancepb.Instance, error) {
	return s.instanceAdmin.GetInstance(ctx, &instancepb.GetInstanceRequest{
		Name: s.FQDN(),
	})
}

func (s *Client) DBCount(ctx context.Context) (int, error) {
	iter := s.dbAdmin.ListDatabases(ctx, &databasepb.ListDatabasesRequest{
		Parent: s.FQDN(),
	})
	dbCount := 0
	for {
		_, err := iter.Next()
		if err == iterator.Done {
			return dbCount, nil
		}
		if err != nil {
			return 0, err
		}
		dbCount++
	}
	return 0, errors.New("cant get dbcount")
}

func (s *Client) UpdatePU(ctx context.Context, ins *instancepb.Instance) error {
	op, err := s.instanceAdmin.UpdateInstance(ctx, &instancepb.UpdateInstanceRequest{
		Instance: ins,
		FieldMask: &fieldmaskpb.FieldMask{
			Paths: []string{"processing_units"},
		},
	})
	if err != nil {
		return err
	}
	if _, err := op.Wait(ctx); err != nil {
		return err
	}
	return nil
}
