#[derive(GraphQLObject)]
#[graphql(description="Keywords which link related content together")]
pub struct Keyword {
    pub value: String,
    pub usages: i32
}

#[derive(GraphQLInputObject)]
#[graphql(description="Ensure links from uuids to keywords. Exclusive removes all other links")]
pub struct SetKeywords {
    pub uuids: Vec<String>,
    pub keywords: Vec<String>,
    pub exclusive: bool
}
