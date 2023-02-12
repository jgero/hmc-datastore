use super::keyword::Keyword;

#[derive(GraphQLObject)]
pub struct Person {
    pub name: String,
    pub uuid: String,
    pub keywords: Vec<Keyword>,
    pub created: i32,
    pub updated: i32,
    pub update_count: i32
}

#[derive(GraphQLInputObject)]
pub struct NewPerson {
    pub name: String,
    pub keywords: Vec<String>,
}

#[derive(GraphQLInputObject)]
pub struct UpdatePerson {
    pub uuid: String,
    pub name: String,
}
