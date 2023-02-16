use juniper::{FieldResult, EmptySubscription, FieldError};

use crate::{reopsitory::Repository, model::person::{NewPerson, Person}};

pub struct Context {
    pub repo: Box<dyn Repository + Sync + Send>
}

// To make our context usable by Juniper, we have to implement a marker trait.
impl juniper::Context for Context {}
pub struct Query;

#[graphql_object(context = Context)]
impl Query {
    fn apiVersion() -> &'static str {
        "1.0"
    }

    fn persons(context: &Context) -> FieldResult<Vec<Person>> {
        context.repo.get_persons().or_else(|e| Err(FieldError::new(e.clone(), graphql_value!({ "internal_error": e }))))
    }
}

pub struct Mutation;

#[graphql_object(context = Context)]
impl Mutation {
    fn createPerson(context: &Context, new_person: NewPerson) -> FieldResult<Person> {
        context.repo.new_person(new_person).or_else(|e| Err(FieldError::new(e.clone(), graphql_value!({ "internal_error": e }))))
    }
}

// A root schema consists of a query and a mutation.
// Request queries can be executed against a RootNode.
type Schema = juniper::RootNode<'static, Query, Mutation, EmptySubscription<Context>>;

pub fn new_schema() -> Schema {
    Schema::new(Query, Mutation, EmptySubscription::new())
}
