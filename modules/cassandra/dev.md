`Cassandra` setup:

- Install [cassandra](https://cassandra.apache.org/_/index.html) on your environment.

- Get latest [jmx_exporter](https://repo1.maven.org/maven2/io/prometheus/jmx/jmx_prometheus_javaagent/) jar file and
  install it in directory where `cassandra` can access it.

- Add [jmx_exporter.yaml](jmx_exporter.yaml) inside `/etc/cassandra`.

- Change `/etc/cassandra/cassandra-env.sh` adding a `javaagent`:

  ```
  JVM_OPTS="$JVM_OPTS $JVM_EXTRA_OPTS -javaagent:/opt/jmx_exporter/jmx_exporter.jar=7072:/etc/cassandra/jmx_exporter.yaml
  ```

- Modify path to `jmx_exporting` according your installation.
