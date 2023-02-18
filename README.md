# Host my club datastore

## Database(s)

- run neo4j with `podman run -p 7474:7474 -p 7687:7687 -v $PWD/data:/data:Z neo4j`
- janusgraph with `podman run -p 8182:8182 janusgraph/janusgraph`

