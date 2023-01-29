# Host my club datastore

## Database(s)

- run neo4j with `podman run -p 7474:7474 -p 7687:7687 -v $PWD/data:/data:Z neo4j`

## TODO

- add unique constraints to neo4j uuid node attributes
- use unions for search results
- enable update mutations or change current one to upsert mutations

