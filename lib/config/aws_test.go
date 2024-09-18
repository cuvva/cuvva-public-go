package config

import (
	"context"
	"os"
	"strconv"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/sts"
	"github.com/stretchr/testify/require"
)

func TestAWS_SessionV2(t *testing.T) {
	if ok, _ := strconv.ParseBool(os.Getenv("TEST_REMOTE_APIS")); !ok {
		t.Skip("Skipping remote API tests without TEST_REMOTE_APIS=true")
	}

	a := AWS{
		Region: "eu-west-1",
	}

	cfg, err := a.SessionV2(context.Background())
	require.NoError(t, err)

	client2 := sts.NewFromConfig(cfg)
	_, err = client2.GetCallerIdentity(context.Background(), &sts.GetCallerIdentityInput{})
	require.NoError(t, err)
}

func TestAWS_SessionV2_Static(t *testing.T) {
	if ok, _ := strconv.ParseBool(os.Getenv("TEST_REMOTE_APIS")); !ok {
		t.Skip("Skipping remote API tests without TEST_REMOTE_APIS=true")
	}

	a := AWS{
		Region:          "eu-west-1",
		AccessKeyID:     os.Getenv("TEST_AWS_ACCESS_KEY_ID"),
		AccessKeySecret: os.Getenv("TEST_AWS_SECRET_ACCESS_KEY"),
	}

	if a.AccessKeyID == "" {
		t.Fatal("TEST_AWS_ACCESS_KEY_ID not provided")
	}

	cfg, err := a.SessionV2(context.Background())
	require.NoError(t, err)

	client2 := sts.NewFromConfig(cfg)
	_, err = client2.GetCallerIdentity(context.Background(), &sts.GetCallerIdentityInput{})
	require.NoError(t, err)
}
