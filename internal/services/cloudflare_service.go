package services

import (
	"context"
	"fmt"

	"github.com/cloudflare/cloudflare-go"
	"gwid.io/gwid-core/internal/config"
	"gwid.io/gwid-core/internal/types"
)

type CloudflareService struct {
	config *config.Config
}

func NewCloudflareService(config *config.Config) *CloudflareService {
	return &CloudflareService{
		config: config,
	}
}

func (s *CloudflareService) AddGatewayToCloudflare(ipAddress, gatewayName string) (*types.CreateDNSRecordResult, error) {
	cloudflareAPI, err := cloudflare.NewWithAPIToken(s.config.CloudflareAPIToken)
	if err != nil {
		return nil, err
	}

	ctx := context.Background()

	zoneID, err := cloudflareAPI.ZoneIDByName(s.config.CloudflareZoneName)
	if err != nil {
		return nil, err
	}

	subdomain := fmt.Sprintf("%s.%s", gatewayName, s.config.CloudflareZoneName)

	record := cloudflare.CreateDNSRecordParams{
		Type:      "A",
		Name:      subdomain,
		Content:   ipAddress,
		ID:        zoneID,
		TTL:       120,
		Proxied:   cloudflare.BoolPtr(false),
		Proxiable: false,
	}

	if _, err := cloudflareAPI.CreateDNSRecord(ctx, cloudflare.ZoneIdentifier(zoneID), record); err != nil {
		return nil, err
	}

	return &types.CreateDNSRecordResult{
		Subdomain: subdomain,
		ZoneID:    zoneID,
	}, nil
}

func (s *CloudflareService) RemoveGatewayFromCloudflare(zoneID, gatewayName string) error {
	cloudflareAPI, err := cloudflare.NewWithAPIToken(s.config.CloudflareAPIToken)
	if err != nil {
		return err
	}

	ctx := context.Background()

	subdomain := fmt.Sprintf("%s.%s", gatewayName, s.config.CloudflareZoneName)

	records, _, err := cloudflareAPI.ListDNSRecords(ctx, cloudflare.ZoneIdentifier(zoneID), cloudflare.ListDNSRecordsParams{
		Type: "A",
		Name: subdomain,
	})
	if err != nil {
		return err
	}

	if len(records) == 0 {
		return nil
	}

	recordID := records[0].ID

	if err := cloudflareAPI.DeleteDNSRecord(ctx, cloudflare.ZoneIdentifier(zoneID), recordID); err != nil {
		return err
	}

	return nil
}
