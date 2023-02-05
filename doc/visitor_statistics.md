# Visitor statistics

To be able to recommend frequently visited content to users it needs to be
tracked how much a specific item was visited.

## What data is relevant?

In the end the reason for tracking what content is visited more than other is to
recommend more relevant content to the users. But the question is what
information is necessary to be able to deduct anything from it.

Simply counting visits on posts would be the easiest solution, but that sounds
like a system that would just recommend old content since that had the most time
to accumulate views. That's where the keywords could come into play. These are
supposed to group similar content together. Combined with the view-counting on
the posts the amount of views per topic could be extracted. Now it would be
possible to not only all-time popular posts, but also new posts on hot topics.

This can be taken a step further by also storing where users navigated after
seeing a post. But just storing any navigation the users make would produce a
lot of bad data, only navigations by which a user found what he was looking for
should be considered. Deciding which navigations are applicable for this is
hard.

## Storing the data

I was experimenting a bit with how inserting such navigation data could look
like. The following is just the result of some experimentation and nothing
final, but I don't want to forget it

```cypher
MERGE (i:Initial_Visit)
WITH i, ["test", "good stuff"] AS searchStrings
MATCH (k:Keyword)
WHERE k.value IN searchStrings
WITH i, k
CREATE (i)-[:navigated {time: 4312431}]->(k)
WITH i, ["a1cd824f-13af-459b-9108-66b80b098f9d"] AS searchStrings
MATCH (p:Post)
WHERE p.uuid IN searchStrings
WITH i, p
CREATE (i)-[:navigated {time: 43124312}]->(p)
```

Maybe it would be a good idea to first store this navigation log in some kind of
time series database and only update the main graph every few hours or something
like that. Updating the recommendations frequently probably isn't that critical
and that would reduce the load a great deal probably. It could also be
beneficial to have this data as time series, but I have no plan yet for what
exactly.

