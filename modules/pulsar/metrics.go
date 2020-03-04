package pulsar

// sad truth: https://github.com/apache/pulsar/issues/3289#issuecomment-578475801

// root@c0f367a0af21:/pulsar/conf# grep stats broker.conf
// statsUpdateFrequencyInSecs=60
// statsUpdateInitialDelayInSecs=60
// exposeTopicLevelMetricsInPrometheus=true
// exposeConsumerLevelMetricsInPrometheus=false
// exposePublisherStats=true

// Monitoring
// http://pulsar.apache.org/docs/en/deploy-monitoring/

// REST API
// http://pulsar.apache.org/admin-rest-api/?version=2.5.0

// Grafana Dashboards
// https://github.com/apache/pulsar/tree/master/docker/grafana/dashboards
