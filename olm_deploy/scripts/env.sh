#!/bin/sh
set -eou pipefail

LOGGING_VERSION=${LOGGING_VERSION:-5.5}
LOGGING_FLUENTD_VERSION=${LOGGING_FLUENTD_VERSION:-1.14.5}
LOGGING_VECTOR_VERSION=${LOGGING_VECTOR_VERSION:-0.14.1}
LOGGING_LOG_FILE_METRIC_EXPORTER_VERSION=${LOGGING_LOG_FILE_METRIC_EXPORTER_VERSION:-1.0}
LOGGING_IS=${LOGGING_IS:-openshift-logging}
export IMAGE_CLUSTER_LOGGING_OPERATOR_REGISTRY=${IMAGE_CLUSTER_LOGGING_OPERATOR_REGISTRY:-quay.io/${LOGGING_IS}/cluster-logging-operator-registry:${LOGGING_VERSION}}
export IMAGE_CLUSTER_LOGGING_OPERATOR=${IMAGE_CLUSTER_LOGGING_OPERATOR:-quay.io/${LOGGING_IS}/cluster-logging-operator:${LOGGING_VERSION}}
export IMAGE_LOGGING_FLUENTD=${IMAGE_LOGGING_FLUENTD:-quay.io/${LOGGING_IS}/fluentd:${LOGGING_FLUENTD_VERSION}}
export IMAGE_LOGGING_VECTOR=timberio/vector:0.20.0-debian
export IMAGE_LOG_FILE_METRIC_EXPORTER=${IMAGE_LOG_FILE_METRIC_EXPORTER:-quay.io/${LOGGING_IS}/log-file-metric-exporter:${LOGGING_LOG_FILE_METRIC_EXPORTER_VERSION}}

export CLUSTER_LOGGING_OPERATOR_NAMESPACE=${CLUSTER_LOGGING_OPERATOR_NAMESPACE:-openshift-logging}
