use super::keyword::Keyword;

#[derive(GraphQLObject)]
pub struct Post {
    pub title: String,
    pub content: String,
    pub uuid: String,
    pub keywords: Vec<Keyword>,
    pub created: i32,
    pub updated: i32,
    pub update_count: i32
}

#[derive(GraphQLInputObject)]
pub struct NewPost {
    pub title: String,
    pub content: String,
}

#[derive(GraphQLInputObject)]
pub struct UpdatePost {
    pub uuid: String,
    pub name: String,
}
