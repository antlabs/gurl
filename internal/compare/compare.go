package compare

import (
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/antlabs/gurl/internal/asserts"
	"github.com/antlabs/gurl/internal/config"
	"github.com/antlabs/gurl/internal/parser"
	"github.com/tidwall/gjson"
)

// AssertionResult represents the result of a single compare assertion.
type AssertionResult struct {
	// PairLabel 标识此断言所属的请求对（例如 "base_cache_on vs cache_off"）
	PairLabel   string
	Line        string
	OK          bool
	BaseValue   string
	TargetValue string
	Message     string
}

// RunScenario finds a scenario by name in cfg and executes it.
// Currently only mode "one_to_one" (or empty) is supported.
func RunScenario(cfg *config.CompareConfig, scenarioName string) (results []AssertionResult, passed int, failed int, err error) {
	var scenario *config.CompareScenario
	for i := range cfg.Scenarios {
		if cfg.Scenarios[i].Name == scenarioName {
			scenario = &cfg.Scenarios[i]
			break
		}
	}
	if scenario == nil {
		return nil, 0, 0, fmt.Errorf("compare scenario '%s' not found", scenarioName)
	}

	mode := scenario.Mode
	if mode == "" {
		mode = "one_to_one"
	}

	switch mode {
	case "one_to_one":
		if scenario.Base == "" || scenario.Target == "" {
			return nil, 0, 0, fmt.Errorf("one_to_one mode requires 'base' and 'target'")
		}

		baseReqDef, err := findRequestByName(cfg.Requests, scenario.Base)
		if err != nil {
			return nil, 0, 0, err
		}
		targetReqDef, err := findRequestByName(cfg.Requests, scenario.Target)
		if err != nil {
			return nil, 0, 0, err
		}

		pairLabel := fmt.Sprintf("%s vs %s", scenario.Base, scenario.Target)
		baseResp, err := doSingleRequest(baseReqDef.Curl)
		if err != nil {
			return nil, 0, 0, fmt.Errorf("base request '%s' failed: %w", scenario.Base, err)
		}
		targetResp, err := doSingleRequest(targetReqDef.Curl)
		if err != nil {
			return nil, 0, 0, fmt.Errorf("target request '%s' failed: %w", scenario.Target, err)
		}

		pairResults := evaluateCompareAssertions(scenario.ResponseCompare, baseResp, targetResp)
		for _, r := range pairResults {
			r.PairLabel = pairLabel
			results = append(results, r)
			if r.OK {
				passed++
			} else {
				failed++
			}
		}

		return results, passed, failed, nil

	case "one_to_many":
		if scenario.Base == "" || len(scenario.Targets) == 0 {
			return nil, 0, 0, fmt.Errorf("one_to_many mode requires 'base' and non-empty 'targets'")
		}

		baseReqDef, err := findRequestByName(cfg.Requests, scenario.Base)
		if err != nil {
			return nil, 0, 0, err
		}

		// 按设计：base 请求只发送一次，其响应在多个 target 间复用
		baseResp, err := doSingleRequest(baseReqDef.Curl)
		if err != nil {
			return nil, 0, 0, fmt.Errorf("base request '%s' failed: %w", scenario.Base, err)
		}

		for _, targetName := range scenario.Targets {
			targetReqDef, err := findRequestByName(cfg.Requests, targetName)
			if err != nil {
				return nil, 0, 0, err
			}

			pairLabel := fmt.Sprintf("%s vs %s", scenario.Base, targetName)
			targetResp, err := doSingleRequest(targetReqDef.Curl)
			if err != nil {
				return nil, 0, 0, fmt.Errorf("target request '%s' failed: %w", targetName, err)
			}

			pairResults := evaluateCompareAssertions(scenario.ResponseCompare, baseResp, targetResp)
			for _, r := range pairResults {
				r.PairLabel = pairLabel
				results = append(results, r)
				if r.OK {
					passed++
				} else {
					failed++
				}
			}
		}

		return results, passed, failed, nil

	case "pair_by_index":
		leftList, okL := cfg.RequestSets[scenario.LeftSet]
		rightList, okR := cfg.RequestSets[scenario.RightSet]
		if !okL || !okR {
			return nil, 0, 0, fmt.Errorf("request_sets '%s' or '%s' not found", scenario.LeftSet, scenario.RightSet)
		}
		if len(leftList) != len(rightList) {
			return nil, 0, 0, fmt.Errorf("request_sets '%s' and '%s' must have the same length", scenario.LeftSet, scenario.RightSet)
		}

		for i := range leftList {
			leftReq := leftList[i]
			rightReq := rightList[i]
			pairLabel := fmt.Sprintf("%s vs %s", leftReq.Name, rightReq.Name)

			baseResp, err := doSingleRequest(leftReq.Curl)
			if err != nil {
				return nil, 0, 0, fmt.Errorf("left request '%s' failed: %w", leftReq.Name, err)
			}
			targetResp, err := doSingleRequest(rightReq.Curl)
			if err != nil {
				return nil, 0, 0, fmt.Errorf("right request '%s' failed: %w", rightReq.Name, err)
			}

			pairResults := evaluateCompareAssertions(scenario.ResponseCompare, baseResp, targetResp)
			for _, r := range pairResults {
				r.PairLabel = pairLabel
				results = append(results, r)
				if r.OK {
					passed++
				} else {
					failed++
				}
			}
		}

		return results, passed, failed, nil

	case "group_by_field":
		if scenario.GroupField == "" {
			return nil, 0, 0, fmt.Errorf("group_by_field mode requires 'group_field'")
		}
		if scenario.BaseRole == "" || scenario.TargetRole == "" {
			return nil, 0, 0, fmt.Errorf("group_by_field mode requires 'base_role' and 'target_role'")
		}

		// 当前实现仅支持使用 CompareRequest.Group 作为分组字段
		if scenario.GroupField != "group" {
			return nil, 0, 0, fmt.Errorf("group_field '%s' is not supported yet (only 'group')", scenario.GroupField)
		}

		// group_value -> {base, target}
		type pairInGroup struct {
			baseReq   *config.CompareRequest
			targetReq *config.CompareRequest
		}
		groups := make(map[string]*pairInGroup)

		for i := range cfg.Requests {
			req := &cfg.Requests[i]
			grp := req.Group
			if grp == "" {
				continue
			}
			g := groups[grp]
			if g == nil {
				g = &pairInGroup{}
				groups[grp] = g
			}
			if req.Role == scenario.BaseRole {
				g.baseReq = req
			} else if req.Role == scenario.TargetRole {
				g.targetReq = req
			}
		}

		for grp, pair := range groups {
			if pair.baseReq == nil || pair.targetReq == nil {
				return nil, 0, 0, fmt.Errorf("group '%s' does not have both base_role '%s' and target_role '%s'", grp, scenario.BaseRole, scenario.TargetRole)
			}

			pairLabel := fmt.Sprintf("%s vs %s (group=%s)", pair.baseReq.Name, pair.targetReq.Name, grp)
			baseResp, err := doSingleRequest(pair.baseReq.Curl)
			if err != nil {
				return nil, 0, 0, fmt.Errorf("base request '%s' in group '%s' failed: %w", pair.baseReq.Name, grp, err)
			}
			targetResp, err := doSingleRequest(pair.targetReq.Curl)
			if err != nil {
				return nil, 0, 0, fmt.Errorf("target request '%s' in group '%s' failed: %w", pair.targetReq.Name, grp, err)
			}

			pairResults := evaluateCompareAssertions(scenario.ResponseCompare, baseResp, targetResp)
			for _, r := range pairResults {
				r.PairLabel = pairLabel
				results = append(results, r)
				if r.OK {
					passed++
				} else {
					failed++
				}
			}
		}

		return results, passed, failed, nil

	default:
		return nil, 0, 0, fmt.Errorf("compare mode '%s' is not supported", mode)
	}
}

func findRequestByName(list []config.CompareRequest, name string) (*config.CompareRequest, error) {
	for i := range list {
		if list[i].Name == name {
			return &list[i], nil
		}
	}
	return nil, fmt.Errorf("request '%s' not found", name)
}

func doSingleRequest(curl string) (*asserts.HTTPResponse, error) {
	if strings.TrimSpace(curl) == "" {
		return nil, fmt.Errorf("curl command is empty")
	}

	req, err := parser.ParseCurl(curl)
	if err != nil {
		return nil, fmt.Errorf("failed to parse curl: %w", err)
	}

	client := &http.Client{}
	start := time.Now()
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	dur := time.Since(start)

	return &asserts.HTTPResponse{
		Status:   resp.StatusCode,
		Headers:  resp.Header,
		Body:     body,
		Duration: dur,
	}, nil
}

var (
	reStatusEq     = regexp.MustCompile(`^status\s*==\s*status$`)
	reHeaderEq     = regexp.MustCompile(`^header\[(.+)\]\s*==\s*header\[(.+)\]$`)
	reHeaderIgnore = regexp.MustCompile(`^header\[(.+)\]\s+ignore$`)
	reGJSONEq      = regexp.MustCompile(`^gjson\s+"([^"]+)"\s*==\s*gjson\s+"([^"]+)"$`)
)

func evaluateCompareAssertions(text string, baseResp, targetResp *asserts.HTTPResponse) []AssertionResult {
	lines := strings.Split(text, "\n")
	results := make([]AssertionResult, 0, len(lines))

	for _, raw := range lines {
		line := strings.TrimSpace(raw)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// status == status
		if reStatusEq.MatchString(line) {
			res := AssertionResult{Line: line}
			res.BaseValue = fmt.Sprintf("%d", baseResp.Status)
			res.TargetValue = fmt.Sprintf("%d", targetResp.Status)
			if baseResp.Status == targetResp.Status {
				res.OK = true
			} else {
				res.Message = "status codes are different"
			}
			results = append(results, res)
			continue
		}

		// header[Name] == header[Name]
		if m := reHeaderEq.FindStringSubmatch(line); m != nil {
			nameLeft := strings.TrimSpace(m[1])
			nameRight := strings.TrimSpace(m[2])
			res := AssertionResult{Line: line}
			baseVal := baseResp.Headers.Get(nameLeft)
			targetVal := targetResp.Headers.Get(nameRight)
			res.BaseValue = baseVal
			res.TargetValue = targetVal
			if baseVal == targetVal {
				res.OK = true
			} else {
				res.Message = "header values are different"
			}
			results = append(results, res)
			continue
		}

		// header[Name] ignore
		if m := reHeaderIgnore.FindStringSubmatch(line); m != nil {
			res := AssertionResult{Line: line, OK: true}
			results = append(results, res)
			continue
		}

		// gjson "path" == gjson "path"
		if m := reGJSONEq.FindStringSubmatch(line); m != nil {
			pathLeft := m[1]
			pathRight := m[2]
			res := AssertionResult{Line: line}
			baseVal := gjson.GetBytes(baseResp.Body, pathLeft)
			targetVal := gjson.GetBytes(targetResp.Body, pathRight)
			res.BaseValue = baseVal.Raw
			res.TargetValue = targetVal.Raw
			if baseVal.Raw == targetVal.Raw {
				res.OK = true
			} else {
				res.Message = "gjson values are different"
			}
			results = append(results, res)
			continue
		}

		// default: treat as single-response assert applied to both base and target
		res := AssertionResult{Line: line}
		if err := asserts.Evaluate(line, baseResp); err != nil {
			res.Message = fmt.Sprintf("base failed: %v", err)
			results = append(results, res)
			continue
		}
		if err := asserts.Evaluate(line, targetResp); err != nil {
			res.Message = fmt.Sprintf("target failed: %v", err)
			results = append(results, res)
			continue
		}
		res.OK = true
		results = append(results, res)
	}

	return results
}
