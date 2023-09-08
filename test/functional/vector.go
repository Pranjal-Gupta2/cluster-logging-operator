// Default to vector if fluentd is not requested
//go:build !fluentd

package functional

import (
	logging "github.com/openshift/cluster-logging-operator/apis/logging/v1"
)

const LogCollectionType = logging.LogCollectionTypeVector
