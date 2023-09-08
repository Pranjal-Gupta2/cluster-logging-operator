package v1

import auditv1 "k8s.io/apiserver/pkg/apis/audit/v1"

// APIAudit filter Kube API server audit logs, as described in [Kubernetes Auditing].
//
// # Policy Filtering
//
// Policy event rules are the same format as the [Kube Audit Policy] with some minor extensions.
// The extensions are described here, see the [Kube Audit Policy] for the standard rule behavior.
// Rules are checked in order, checking stops at the first matching rule.
//
// The `level` of the rule determines how the event is treated:
//   - None: the event is dropped.
//   - Metadata: Audit metadata only, request and response objects are removed.
//   - Request: Audit metadata and request are retained, response is removed.
//   - RequestResponse: All data is retained, including the full response. This can lead to very large events.
//     For example the command `oc list -A pods` results in a response containing the YAML description of every pod in the cluster.
//
// # Extensions
//
// The following features are extensions to the standard [Kube Audit Policy]
//
// ## Wildcards
//
// Names of users, groups, namespaces, and resources can have a leading or trailing '*' character.
// For example namespace 'openshift-*' matches 'openshift-apiserver' or 'openshift-authentication.
// Resource '*/status' matches 'Pod/status' or 'Deployment/status'
//
// ## Default Rules
//
// Events that do not match any rule in the policy are filtered as follows:
// - User events (ie. non-system and non-serviceaccount) are forwarded
// - Read-only system events (get/list/watch etc) are dropped
// - Service account write events that occur within the same namespace as the service account are dropped
// - All other events are forwarded, subject to any configured [rate limits][#rate-lmiting]
//
// If you want to disable these defaults, end your rules list with rule that has only a `level` field.
// An empty rule matches any event, and prevents the defaults from taking effect.
//
// ## Omit Response Codes
//
// You can drop events based on the HTTP status code in the response. See the OmitResponseCodes field.
//
// [Kube Audit Policy]: https://kubernetes.io/docs/reference/config-api/apiserver-audit.v1/#audit-k8s-io-v1-Policy
// [Kubernetes Auditing]: https://kubernetes.io/docs/tasks/debug/debug-cluster/audit/
type APIAudit struct {

	// Rules specify the audit Level a request should be recorded at.
	// A request may match multiple rules, in which case the FIRST matching rule is used.
	// PolicyRules are strictly ordered.
	//
	// If Rules is empty or missing default rules apply, see [APIAudit]
	Rules []auditv1.PolicyRule `json:"rules,omitempty"`

	// OmitStages is a list of stages for which no events are created.
	// Note that this can also be specified per rule in which case the union of both are omitted.
	// +optional
	OmitStages []auditv1.Stage `json:"omitStages,omitempty"`

	// OmitResponseCodes is a list of response codes for which no events are created.
	// If this field is missing or null, the default value used is [404, 409, 422, 429]
	// (NotFound, Conflict, UnprocessableEntity, TooManyRequests)
	// If it is the empty list [], then no status codes are omitted.
	// Otherwise this field should be a list of integer status codes to omit.
	//
	// +optional
	OmitResponseCodes *[]int `json:"omitResponseCodes,omitempty"`
}
