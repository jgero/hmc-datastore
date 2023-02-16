use std::time::{SystemTime, UNIX_EPOCH};
use uuid::Uuid;
use super::keyword::Keyword;

#[derive(Debug)]
pub enum Error {
    Uuid,
    Time
}

impl From<Error> for String {
    fn from(value: Error) -> Self {
        match value {
            Error::Uuid => "could not generatr uuid".to_string(),
            Error::Time => "could not determine current time".to_string()
        }
    }
}

#[derive(GraphQLObject)]
pub struct Person {
    pub name: String,
    pub uuid: Uuid,
    pub keywords: Vec<Keyword>,
    pub created: f64,
    pub updated: f64,
    pub update_count: i32
}

impl Person {
    pub fn create(np: NewPerson) -> Result<Person, Error> {
        let uuid = Uuid::new_v4();
        let created = match SystemTime::now().duration_since(UNIX_EPOCH) {
            Ok(created) => created.as_secs_f64(),
            Err(_) => return Err(Error::Time)
        };
        // let c = SystemTime::now().duration_since(UNIX_EPOCH)?.unwrap_or(0.0);
        Ok(Person { name: np.name, uuid, keywords: Vec::new(), created, updated: created, update_count: 0 })
    }
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
