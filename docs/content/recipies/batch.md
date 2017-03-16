+++
title = "Batch Jobs"
chapter = false
prev = "/recipies"
next = "/notice"
+++

![image](http://i.giphy.com/3og0IO9FUngqYxdBnO.gif)

### Generate an optimized GIF on Mesos

In this example we generate a GIF from a video file and upload it to [giphy.com](https://giphy.com). We use [ffmpeg](http://ffmpeg.org/)
for optimization and conversion of the video file. The process is borrowed from [this](http://blog.pkh.me/p/21-high-quality-gif-with-ffmpeg.html) 
excellent tutorial where you can learn a little bit more about what is actually happening. The process can be broken down into four steps where each step is run within a seperate container and is dependent on the subsequent step. We will call each step in this process a *task* and call the 
collection of tasks a *job*. This very simple and linear process can be visuallized with a dag: 
![dag](http://localhost:1313/mesos-cli/img/linear.svg)

#### Storage considerations

One of the most challenging aspects of running containerized workloads is how to manage application state. Traditionally an application is run on a single server with some fixed storage device that retains all of it's data. When running in containerized environment applications often move across several machines and need to have access to consistent storage. There are many approaches for dealing with this in Mesos which we will not cover here. It is assumed going forward that all Mesos agents have access to a shared storage path at `/mnt`.


#### Tasks

Each task is described in the batch processing example [Mesosfile](https://raw.githubusercontent.com/vektorlab/mesos-cli/master/examples/batch_processing/Mesosfile) in the **mesos-cli** source directory. 

##### Download
We use the Mesos [fetcher](http://mesos.apache.org/documentation/latest/fetcher/) to download a video file from Wikimedia and then copy it onto shared storage.

``` yaml
- name: Download Source  # A friendly name given to each task
  command:
    environment:  # Array of environment variables exposed inside the container.
      variables:
      - name: OUTPUT_FILE
        value: /mnt/input.ogv
    uris:
      -
        # The URL to the source file we are converting
        value: https://upload.wikimedia.org/wikipedia/commons/2/2c/WorldSunshine.ogv
        # The name we give to the downloaded file
        output_file: input.ogv
    # Execute the value as a shell command e.g. /bin/sh -c $value
    shell: true
    # Copy the downloaded file onto our shared storage path
    value: cp -v $MESOS_SANDBOX/input.ogv $OUTPUT_FILE
    user: root
  container:
    type: MESOS
    volumes:
    # Expose a path on the host within the container
    - container_path: /mnt
      host_path: /mnt
      mode: RW
```

#### Generate a palette file

GIFs are limited to a palette of 256 colors which may not necessarily be optimized for the underlying images.
We use ffmpeg to generate a png palette file from the source video which is then used as a backing color palette 
of our GIF. The resulting palette file is saved to the shared `/mnt` path.

```yaml
- name: Generate Palette
  command:
    uris:
      -
        value: /mnt/input.ogv
    environment:
      variables:
        - name: INPUT_FILE
          value: /mnt/input.ogv
        - name: OUTPUT_FILE
          value: /mnt/palette.png
        - name: FILTERS
          value: fps=10,scale=1200:-1:flags=lanczos
    shell: true
    value: ffmpeg -v warning -i $INPUT_FILE -vf "$FILTERS,palettegen"  -y $OUTPUT_FILE
    user: root
  container:
    type: MESOS
    mesos:
      image:
        type: DOCKER
        docker:
          name: jrottenberg/ffmpeg:3.2
    volumes:
    - container_path: /mnt
      host_path: /mnt
      mode: RW
```

#### Generate the GIF from the video file with the reference palette

Next we take the video file, the generated palette, and create a GIF from our video file saving the output to the shared `/mnt` path.

``` yaml
- name: Transform to GIF
  command:
    uris:
      -
        value: /mnt/input.ogv
    environment:
      variables:
        - name: INPUT_FILE
          value: /mnt/input.ogv
        - name: PALETTE_FILE
          value: /mnt/palette.png
        - name: OUTPUT_FILE
          value: /mnt/output.gif
        - name: FILTERS
          value: fps=10,scale=1200:-1:flags=lanczos
    shell: true
    value: ffmpeg -v warning -i $INPUT_FILE -i $PALETTE_FILE -lavfi "$FILTERS [x]; [x][1:v] paletteuse" -y $OUTPUT_FILE
    user: root
  container:
    type: MESOS
    mesos:
      image:
        type: DOCKER
        docker:
          name: jrottenberg/ffmpeg:3.2
    volumes:
    - container_path: /mnt
      host_path: /mnt
      mode: RW
```

#### Upload the generated GIF to Giphy

In the final step we take the newly created GIF file and upload it to giphy.com. The API key used below is provided by their [beta API](https://github.com/Giphy/GiphyAPI). We also record the API response to a file so we can generate a link to the uploaded file.

``` yaml
# Upload the output to Giphy
- name: Upload to Giphy
  command:
    uris:
      -
        value: /mnt/output.gif
    environment:
      variables:
        - name: API_KEY
          value: dc6zaTOxFJmzC
        - name: OUTPUT_FILE
          value: /mnt/output.gif
    shell: true
    value:  curl -F "file=@$OUTPUT_FILE"  -F "api_key=$API_KEY" "http://upload.giphy.com/v1/gifs" > upload.json
    user: root
  container:
    type: MESOS
    volumes:
    - container_path: /mnt
      host_path: /mnt
      mode: RW

```

#### Run the job!

We can run the job from the mesos-cli examples [directory](https://github.com/vektorlab/mesos-cli/tree/master/examples/batch_processing).

The `Mesosfile` will be automatically chosen loaded, we use the `--sync` flag to ensure each task is run synchronously and without failure. 
We can also include the `--restart` which will retry each task automatically if it fails.

``` bash
$ cd mesos-cli/examples/batch_processing
$ ls -la
total 12
drwxr-xr-x 2 kevin kevin 4096 Mar 16 17:06 .
drwxr-xr-x 3 kevin kevin 4096 Mar 15 17:17 ..
-rw-r--r-- 1 kevin kevin 2466 Mar 16 16:59 Mesosfile
$ mesos run --sync --restart
2017-03-16T16:59:33.229+0700	INFO	CALL	{"type": "SUBSCRIBE", "framework": "&FrameworkInfo{User:root,Name:mesos-cli,ID:&FrameworkID{Value:,},FailoverTimeout:nil,Checkpoint:nil,Role:nil,Hostname:nil,Principal:nil,WebuiUrl:nil,Capabilities:[],Labels:nil,}"}
...
```

Now we can look at all of the tasks ran with **mesos-cli** or with the Mesos UI.
``` bash
$ mesos tasks --all
ID                                  	NAME            	FRAMEWORK	STATE        	CPU	MEM	GPU	DISK
16f66b16-0889-4f22-8c9f-d797a034e4ef	Download Source 	3ae4ade2 	TASK_FINISHED	0.1	64 	0  	64  
449f1384-16ca-4bc2-8378-e427d1a662ba	Generate Palette	3ae4ade2 	TASK_FINISHED	0.1	64 	0  	64  
bbf68b38-47e9-4cae-b7f2-5d1167fbcfa2	Transform To GIF	3ae4ade2 	TASK_FINISHED	0.1	64 	0  	64  
71228718-22fb-4453-8b0d-f2d836835b62	Upload to Giphy 	3ae4ade2 	TASK_FINISHED	0.1	64 	0  	64
```

Finally we can browse the sandbox directory of the final task and read it's output to determine the URL for Gliphy

``` bash
$ mesos list 3ae4ade2-a512-4819-8d32-ff7c9b4df0a0-S0 /mesos/agents/0/slaves/3ae4ade2-a512-4819-8d32-ff7c9b4df0a0-S0/frameworks/3ae4ade2-a512-4819-8d32-ff7c9b4df0a0-0000/executors/71228718-22fb-4453-8b0d-f2d836835b62/runs/ba88d0e7-27b1-43ba-b4d3-71a5ef707822
UID 	GID 	MODE 	MODIFIED	SIZE	PATH       
root	root	33188	TODO    	253 	output.gif 
root	root	33188	TODO    	247 	stderr     
root	root	33188	TODO    	246 	stdout     
root	root	33188	TODO    	251 	upload.json

$ mesos read 3ae4ade2-a512-4819-8d32-ff7c9b4df0a0-S0 /mesos/agents/0/slaves/3ae4ade2-a512-4819-8d32-ff7c9b4df0a0-S0/frameworks/3ae4ade2-a512-4819-8d32-ff7c9b4df0a0-0000/executors/71228718-22fb-4453-8b0d-f2d836835b62/runs/ba88d0e7-27b1-43ba-b4d3-71a5ef707822/upload.json |jq .
{
  "data": {
    "id": "3og0IO9FUngqYxdBnO" <-- Here is what we want
  },
  "meta": {
    "status": 200,
    "msg": "OK",
    "response_id": "58ca62bfa940de891bba28d1"
  }
}

```

Finally we can construct a URL: `https://giphy.com/gifs/3og0IO9FUngqYxdBnO` and... voilà! :

![output](https://media.giphy.com/media/3og0IO9FUngqYxdBnO/source.gif)


