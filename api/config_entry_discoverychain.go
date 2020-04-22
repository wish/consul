package api

import (
	"encoding/json"
	"time"
)

type ServiceRouterConfigEntry struct {
	Kind      string
	Name      string
	Namespace string `json:",omitempty"`

	Routes []ServiceRoute `json:",omitempty"`

	CreateIndex uint64
	ModifyIndex uint64
}

func (e *ServiceRouterConfigEntry) GetKind() string        { return e.Kind }
func (e *ServiceRouterConfigEntry) GetName() string        { return e.Name }
func (e *ServiceRouterConfigEntry) GetCreateIndex() uint64 { return e.CreateIndex }
func (e *ServiceRouterConfigEntry) GetModifyIndex() uint64 { return e.ModifyIndex }

type ServiceRoute struct {
	Match       *ServiceRouteMatch       `json:",omitempty"`
	Destination *ServiceRouteDestination `json:",omitempty"`
}

type ServiceRouteMatch struct {
	HTTP *ServiceRouteHTTPMatch `json:",omitempty"`
}

type ServiceRouteHTTPMatch struct {
	PathExact  string `json:",omitempty"`
	PathPrefix string `json:",omitempty"`
	PathRegex  string `json:",omitempty"`

	Header     []ServiceRouteHTTPMatchHeader     `json:",omitempty"`
	QueryParam []ServiceRouteHTTPMatchQueryParam `json:",omitempty"`
	Methods    []string                          `json:",omitempty"`
}

func (m ServiceRouteHTTPMatch) DecodeKeyMapping() map[string]string {
	return map[string]string{
		"path_exact":  "pathexact",
		"path_prefix": "pathprefix",
		"path_regex":  "pathregex",
		"query_param": "queryparam",
	}
}

type ServiceRouteHTTPMatchHeader struct {
	Name    string
	Present bool   `json:",omitempty"`
	Exact   string `json:",omitempty"`
	Prefix  string `json:",omitempty"`
	Suffix  string `json:",omitempty"`
	Regex   string `json:",omitempty"`
	Invert  bool   `json:",omitempty"`
}

type ServiceRouteHTTPMatchQueryParam struct {
	Name    string
	Present bool   `json:",omitempty"`
	Exact   string `json:",omitempty"`
	Regex   string `json:",omitempty"`
}

type ServiceRouteDestination struct {
	Service               string        `json:",omitempty"`
	ServiceSubset         string        `json:",omitempty"`
	Namespace             string        `json:",omitempty"`
	PrefixRewrite         string        `json:",omitempty"`
	RequestTimeout        time.Duration `json:",omitempty"`
	NumRetries            uint32        `json:",omitempty"`
	RetryOnConnectFailure bool          `json:",omitempty"`
	RetryOnStatusCodes    []uint32      `json:",omitempty"`
}

func (e *ServiceRouteDestination) DecodeKeyMapping() map[string]string {
	return map[string]string{
		"service_subset":           "servicesubset",
		"prefix_rewrite":           "prefixrewrite",
		"request_timeout":          "requesttimeout",
		"num_retries":              "numretries",
		"retry_on_connect_failure": "retryonconnectfailure",
		"retry_on_status_codes":    "retryonstatuscodes",
	}
}

func (e *ServiceRouteDestination) MarshalJSON() ([]byte, error) {
	type Alias ServiceRouteDestination
	exported := &struct {
		RequestTimeout string `json:",omitempty"`
		*Alias
	}{
		RequestTimeout: e.RequestTimeout.String(),
		Alias:          (*Alias)(e),
	}
	if e.RequestTimeout == 0 {
		exported.RequestTimeout = ""
	}

	return json.Marshal(exported)
}

func (e *ServiceRouteDestination) UnmarshalJSON(data []byte) error {
	type Alias ServiceRouteDestination
	aux := &struct {
		RequestTimeout string
		*Alias
	}{
		Alias: (*Alias)(e),
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	var err error
	if aux.RequestTimeout != "" {
		if e.RequestTimeout, err = time.ParseDuration(aux.RequestTimeout); err != nil {
			return err
		}
	}
	return nil
}

type ServiceSplitterConfigEntry struct {
	Kind      string
	Name      string
	Namespace string `json:",omitempty"`

	Splits []ServiceSplit `json:",omitempty"`

	CreateIndex uint64
	ModifyIndex uint64
}

func (e *ServiceSplitterConfigEntry) GetKind() string        { return e.Kind }
func (e *ServiceSplitterConfigEntry) GetName() string        { return e.Name }
func (e *ServiceSplitterConfigEntry) GetCreateIndex() uint64 { return e.CreateIndex }
func (e *ServiceSplitterConfigEntry) GetModifyIndex() uint64 { return e.ModifyIndex }

type ServiceSplit struct {
	Weight        float32
	Service       string `json:",omitempty"`
	ServiceSubset string `json:",omitempty"`
	Namespace     string `json:",omitempty"`
}

func (s ServiceSplit) DecodeKeyMapping() map[string]string {
	return map[string]string{
		"service_subset": "servicesubset",
	}
}

type ServiceResolverConfigEntry struct {
	Kind      string
	Name      string
	Namespace string `json:",omitempty"`

	DefaultSubset  string                             `json:",omitempty"`
	Subsets        map[string]ServiceResolverSubset   `json:",omitempty"`
	Redirect       *ServiceResolverRedirect           `json:",omitempty"`
	Failover       map[string]ServiceResolverFailover `json:",omitempty"`
	ConnectTimeout time.Duration                      `json:",omitempty"`

	CreateIndex uint64
	ModifyIndex uint64
}

func (e *ServiceResolverConfigEntry) DecodeKeyMapping() map[string]string {
	return map[string]string{
		"connect_timeout": "connecttimeout",
		"default_subset":  "defaultsubset",
	}
}

func (e *ServiceResolverConfigEntry) MarshalJSON() ([]byte, error) {
	type Alias ServiceResolverConfigEntry
	exported := &struct {
		ConnectTimeout string `json:",omitempty"`
		*Alias
	}{
		ConnectTimeout: e.ConnectTimeout.String(),
		Alias:          (*Alias)(e),
	}
	if e.ConnectTimeout == 0 {
		exported.ConnectTimeout = ""
	}

	return json.Marshal(exported)
}

func (e *ServiceResolverConfigEntry) UnmarshalJSON(data []byte) error {
	type Alias ServiceResolverConfigEntry
	aux := &struct {
		ConnectTimeout string
		*Alias
	}{
		Alias: (*Alias)(e),
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	var err error
	if aux.ConnectTimeout != "" {
		if e.ConnectTimeout, err = time.ParseDuration(aux.ConnectTimeout); err != nil {
			return err
		}
	}
	return nil
}

func (e *ServiceResolverConfigEntry) GetKind() string        { return e.Kind }
func (e *ServiceResolverConfigEntry) GetName() string        { return e.Name }
func (e *ServiceResolverConfigEntry) GetCreateIndex() uint64 { return e.CreateIndex }
func (e *ServiceResolverConfigEntry) GetModifyIndex() uint64 { return e.ModifyIndex }

type ServiceResolverSubset struct {
	Filter      string `json:",omitempty"`
	OnlyPassing bool   `json:",omitempty"`
}

func (s ServiceResolverSubset) DecodeKeyMapping() map[string]string {
	return map[string]string{
		"only_passing": "onlypassing",
	}
}

type ServiceResolverRedirect struct {
	Service       string `json:",omitempty"`
	ServiceSubset string `json:",omitempty"`
	Namespace     string `json:",omitempty"`
	Datacenter    string `json:",omitempty"`
}

func (s ServiceResolverRedirect) DecodeKeyMapping() map[string]string {
	return map[string]string{
		"service_subset": "servicesubset",
	}
}

type ServiceResolverFailover struct {
	Service       string   `json:",omitempty"`
	ServiceSubset string   `json:",omitempty"`
	Namespace     string   `json:",omitempty"`
	Datacenters   []string `json:",omitempty"`
}

func (s ServiceResolverFailover) DecodeKeyMapping() map[string]string {
	return map[string]string{
		"service_subset": "servicesubset",
	}
}
