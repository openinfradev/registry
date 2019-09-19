Taco-Registry Builder Project
=============================


Docker Build
------------

#### Git clone
<code> 
    $ git clone https://starlkj@tde.sktelecom.com/stash/scm/oreotools/taco-registry-builder.git
</code>

#### GCC
``` 
$ sudo apt -y update
$ sudo apt install -y build-essential
```

#### Docker build
> move taco-registry-builder directory
  $ make docker-build


Docker Deploy
-------------

#### Docker run
> $ docker run -d -p 4000:4000 \
    -v /var/run/docker.sock:/var/run/docker.sock \
    --restart=always --name builder \
    taco-registry/builder:v2 \
    -log.level=0 \
    -db.type=postgres \
    -db.host=exntu.kr \
    -db.port=25432 \
    -db.user=registry \
    -db.pass=registry1234\$\$ \
    -db.name=registry \
    -db.xarg= \
    -registry.name=taco-registry \
    -registry.insecure=true \
    -registry.endpoint=exntu.kr:25000 \
    -redis.endpoint=exntu.kr:26379 \
    -service.port=4000 \
    -service.tmp=/tmp
