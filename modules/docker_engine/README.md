# docker_engine

This module will monitor one or more Docker engines.

**Requirements:**
 * docker with enabled [`metric-address`](https://docs.docker.com/config/thirdparty/prometheus/)


It produces the following charts:

1. **Container Actions** in actions/s
 * changes
 * commits
 * create
 * delete
 * start

2. **Container States** in number of containers in state
 * running
 * paused
 * stopped
 
3. **Builder Builds Fails By Reason** in fails/s
 * build_canceled
 * build_target_not_reachable_error
 * command_not_supported_error
 * dockerfile_empty_error
 * dockerfile_syntax_error
 * error_processing_commands_error
 * missing_onbuild_arguments_error
 * unknown_instruction_error
 
4. **Health Checks** in events/s
 * fails

<br>

If docker in [swarm mode](https://docs.docker.com/engine/swarm/) and the instance is swarm manager additional charts will be added:

1. **Swarm Manager Leader** in bool
 * is_leader

2. **Swarm Manager Object Store** in count
 * nodes
 * services
 * tasks
 * networks
 * secrets
 * configs
 
3. **Swarm Manager Nodes Per State** in count
 * ready
 * down
 * unknown
 * disconnected
 
4. **Swarm Manager Tasks Per State** in count
 * running
 * failed
 * ready
 * rejected
 * starting
 * shutdown
 * new
 * orphaned
 * preparing
 * pending
 * complete
 * remove
 * accepted
 * assigned

### configuration

For all available options please see module [configuration file](https://github.com/netdata/go.d.plugin/blob/master/config/go.d/docker_engine.conf).
___

Needs only `url` to docker metric-address.

Here is an example for 2 docker instances:

```yaml
jobs:
  - name: local
    url : http://127.0.0.1:9323/metrics
      
  - name: remote
    url : http://100.64.0.1:9323/metrics
```

Without configuration, module attempts to connect to `http://127.0.0.1:9323/metrics`

---
