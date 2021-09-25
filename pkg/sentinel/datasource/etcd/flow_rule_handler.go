package etcd

import (
	"github.com/alibaba/sentinel-golang/core/flow"
	"github.com/alibaba/sentinel-golang/ext/datasource"
	"github.com/pkg/errors"
)

type FlowRuleHandler struct {
	*fixedFlowRule
	*datasource.DefaultPropertyHandler
}

type fixedFlowRule struct {
	Rules []*flow.Rule
}

func NewFlowRuleHandler(fixedRule []*flow.Rule) *FlowRuleHandler {
	RuleList := &fixedFlowRule{
		Rules: fixedRule,
	}
	res := &FlowRuleHandler{
		fixedFlowRule:          RuleList,
		DefaultPropertyHandler: datasource.NewDefaultPropertyHandler(datasource.FlowRuleJsonArrayParser, RuleList.FlowRulesUpdater),
	}
	return res
}

// FlowRulesUpdater load the newest []flow.Rule to downstream flow component.
func (f *fixedFlowRule) FlowRulesUpdater(data interface{}) error {
	if data == nil {
		return flow.ClearRules()
	}

	rules := make([]*flow.Rule, 0, 8)
	if val, ok := data.([]flow.Rule); ok {
		for _, v := range val {
			rules = append(rules, &v)
		}
	} else if val, ok := data.([]*flow.Rule); ok {
		rules = val
	} else {
		return errors.Errorf("Fail to type assert data to []flow.FlowRule or []*flow.FlowRule, in fact, data: %+v", data)
	}
	rules = append(rules, f.Rules...)
	_, err := flow.LoadRules(rules)
	if err == nil {
		return err
	}
	return nil
}
