+++
title = "read"
next = "/usage/run"
prev = "/usage/list"
+++

Read the contents of a file on the given agent at the given path

```
mesos-cli read [OPTIONS] TASKID FILE
```

### Options

Option | Description
--- | ---
-f, --follow=false  | follow the content
-n, --nlines=0      | number of lines to read
-m, --master=""     | mesos master



### Example
```bash
mesos read 5b7a7d41-1f2a-4b6d-93fd-48354d7fa785-S0  /opt/mesos/0/slaves/5b7a7d41-1f2a-4b6d-93fd-48354d7fa785-S0/frameworks/5b7a7d41-1f2a-4b6d-93fd-48354d7fa785-0115/executors/807fe9fa-55df-40fa-ab4f-20e359d51d43/runs/4295fd85-5023-4b0d-a9c8-1d9bdce309ba/stdout

172.17.0.1 - - [08/Jan/2017:03:35:46 +0000] "GET / HTTP/1.1" 200 612 "http://localhost:8080/ui/" "Mozilla/5.0 (X11;...
...
```

