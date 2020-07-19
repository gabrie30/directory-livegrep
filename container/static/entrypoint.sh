#!/bin/bash

/root/directory-livegrep /data

/livegrep/bin/codesearch -index_only -dump_index /data/livegrep.idx /data/livegrep.json

/livegrep/bin/codesearch -load_index /data/livegrep.idx -reload_rpc -grpc "0.0.0.0:9999" &

/livegrep/bin/livegrep -connect "localhost:9999" -docroot /livegrep/web -listen "0.0.0.0:8910"