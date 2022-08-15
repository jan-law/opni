package slo

import (
	"context"
	"fmt"
	"github.com/hashicorp/go-hclog"
	"github.com/prometheus/common/model"
	"github.com/rancher/opni/pkg/metrics/unmarshal"
	"github.com/rancher/opni/plugins/cortex/pkg/apis/cortexadmin"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gopkg.in/yaml.v3"
)

// Apply Cortex Rules to Cortex separately :
// - recording rules
// - metadata rules
// - alert rules
func applyCortexSLORules(
	p *Plugin,
	lg hclog.Logger,
	ctx context.Context,
	clusterId string,
	ruleSpec RuleGroupYAMLv2,
) error {
	out, err := yaml.Marshal(ruleSpec)
	if err != nil {
		return err
	}

	_, err = p.adminClient.Get().LoadRules(ctx, &cortexadmin.PostRuleRequest{
		YamlContent: string(out),
		ClusterId:   clusterId,
	})
	if err != nil {
		lg.Error(fmt.Sprintf(
			"Failed to load rules for cluster %s, rule : %s,",
			clusterId, string(out)))
	}
	return err
}

// }
func deleteCortexSLORules(
	p *Plugin,
	lg hclog.Logger,
	ctx context.Context,
	clusterId string,
	groupName string,
) error {
	_, err := p.adminClient.Get().DeleteRule(ctx, &cortexadmin.RuleRequest{
		ClusterId: clusterId,
		GroupName: groupName,
	})
	// we can ignore 404s here since if we can't find them,
	// then it will be impossible to delete them anyway
	if err != nil && status.Code(err) != codes.NotFound {
		lg.Error(fmt.Sprintf("Failed to delete rule group with clusterId, groupName (%v %v): %v", clusterId, groupName, err))
		return err
	}
	return nil
}

func QuerySLOComponentByRecordName(
	client cortexadmin.CortexAdminClient,
	ctx context.Context,
	recordName string,
	clusterId string,
) (*model.Vector, error) {
	resp, err := client.Query(ctx, &cortexadmin.QueryRequest{
		Tenants: []string{clusterId},
		Query:   recordName,
	})
	if err != nil {
		return nil, err
	}
	rawBytes := resp.Data
	qres, err := unmarshal.UnmarshalPrometheusResponse(rawBytes)
	if err != nil {
		return nil, err
	}
	dataVector, err := qres.GetVector()
	if err != nil {
		return nil, err
	}
	return dataVector, nil
}

func QuerySLOComponentByRawQuery(
	client cortexadmin.CortexAdminClient,
	ctx context.Context,
	rawQuery string,
	clusterId string,
) (*model.Vector, error) {
	resp, err := client.Query(ctx, &cortexadmin.QueryRequest{
		Tenants: []string{clusterId},
		Query:   rawQuery,
	})
	if err != nil {
		return nil, err
	}
	rawBytes := resp.Data
	qres, err := unmarshal.UnmarshalPrometheusResponse(rawBytes)
	if err != nil {
		return nil, err
	}
	dataVector, err := qres.GetVector()
	if err != nil {
		return nil, err
	}
	return dataVector, nil
}
