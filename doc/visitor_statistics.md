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
final, but I don't want to forget it.

Here a query to create edges for first visits:

```cypher
MATCH (a:Post {uuid: "6f31d920-ef4a-459c-ab93-d4cb20096911"})
with a
MERGE (i:Initial_Visit)
with a, i
create (i)-[:navigated {time: 123412341234}]->(a)
WITH a, i
match (a)-[:relates_to]-(k:Keyword)
with a, i, k
CREATE (i)-[:navigated {time: 123412341234}]->(k)
```

And for Navigations from one post to another:

```cypher
MATCH (prev:Post {uuid: "6f31d920-ef4a-459c-ab93-d4cb20096911"}),(next:Post {uuid:"26e36435-84df-48e8-88b0-694fac37e4ee"})
with prev, next
create (prev)-[:navigated {time: 123412341234}]->(next)
WITH prev, next
match (prev)-[:relates_to]-(prevK:Keyword), (next)-[:relates_to]-(nextK:Keyword)
with prevK, nextK
CREATE (prevK)-[:navigated {time: 123412341234}]->(nextK)
```

This one is especially messy since it builds a cartesian product and
also may produce edges from nodes to themselves. And it would definitely create
wild amount of edges over time which is also not too great.

Maybe it would be a good idea to first store this navigation log in some kind of
time series database and only set unique edges with weights in the main graph
every few hours or something like that. The weights could be average navigations
per day over the last month or something like that. Updating the recommendations
frequently probably isn't that critical and that would reduce the load a great
deal probably. It could also be beneficial to have this data as time series, but
I have no plan yet for what exactly.

