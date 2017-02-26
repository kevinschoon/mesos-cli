+++
title = "tasks"
next = "/notice"
prev = "/usage/run"
+++

List and filter tasks on the target cluster.

    mesos-cli tasks [OPTIONS]

### Options

Option | Description
--- | ---
--master=""    |    Mesos master
--truncate=true|   Truncate long values
--task=""      |  Filter by task id
--fuzzy=true   | Fuzzy matching on string values
--state=[]     |filter by task state



### Example

    $ mesos tasks --state "TASK_FINISHED"

    ID        FRAMEWORK STATE         CPU MEM GPU DISK
    mesos-cli b620d6e2  TASK_FINISHED 0.1 128 0   32  
    mesos-cli b620d6e2  TASK_FINISHED 0.1 128 0   32  


