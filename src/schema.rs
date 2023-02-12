use juniper::{FieldResult, EmptySubscription};

use crate::reopsitory::{Repository, GremlinRepository};

#[derive(GraphQLEnum)]
enum Episode {
    NewHope,
    Empire,
    Jedi,
}

#[derive(GraphQLObject)]
#[graphql(description="A humanoid creature in the Star Wars universe")]
struct Human {
    id: String,
    name: String,
    appears_in: Vec<Episode>,
    home_planet: String,
}

// There is also a custom derive for mapping GraphQL input objects.

#[derive(GraphQLInputObject)]
#[graphql(description="A humanoid creature in the Star Wars universe")]
struct NewHuman {
    name: String,
    appears_in: Vec<Episode>,
    home_planet: String,
}

// Now, we create our root Query and Mutation types with resolvers by using the
// graphql_object! macro.
// Objects can have contexts that allow accessing shared state like a database
// pool.

pub struct Context {
    // Use your real database pool here.
    // pool: DatabasePool,
    repo: Box<dyn Repository + Sync>
}

// To make our context usable by Juniper, we have to implement a marker trait.
impl juniper::Context for Context {}
pub struct Query;

#[graphql_object(context = Context)]
impl Query {
    fn apiVersion() -> &'static str {
        "1.0"
    }

    fn human(context: &Context, id: String) -> FieldResult<Human> {
        // // Get the context from the executor.
        // let context = executor.context();
        // // Get a db connection.
        // let connection = context.pool.get_connection()?;
        // // Execute a db query.
        // // Note the use of `?` to propagate errors.
        // let human = connection.find_human(&id)?;
        // Return the result.
        // Ok(human)
        todo!("implement this")
    }
}

pub struct Mutation;

#[graphql_object(context = Context)]
impl Mutation {
    fn createHuman(context: &Context, new_human: NewHuman) -> FieldResult<Human> {
        // let db = executor.context().pool.get_connection()?;
        // let human: Human = db.insert_human(&new_human)?;
        // Ok(human)
        todo!("implement this")
    }
}

// A root schema consists of a query and a mutation.
// Request queries can be executed against a RootNode.
type Schema = juniper::RootNode<'static, Query, Mutation, EmptySubscription<Context>>;

pub fn new_schema() -> Schema {
    Schema::new(Query, Mutation, EmptySubscription::new())
}
