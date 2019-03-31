# dockerhub

This module will collect Docker Hub Repositories statistics.


It produces the following charts:

1. **Pulls Summary** in pulls
  * summary

2. **Pulls** in pulls
  * per repository
 
3. **Pulls Rate** in pulls/s
  * per repository

4. **Stars** in stars/s
  * per repository
  
5. **Current Status** in status
  * per repository
  
6. **Time Since Last Update** in seconds
  * per repository


### configuration

For all available options please see module [configuration file](https://github.com/netdata/go.d.plugin/blob/master/config/go.d/dockerhub.conf).
___

Needs only list of repositories.

Here is an example:

```yaml
jobs:
  - name: me
    repositories: ['me/repo1', 'me/repo2', 'me/repo3'] 
```
---
