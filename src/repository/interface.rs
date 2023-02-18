use crate::model::person::{NewPerson, Person};

pub trait Repository {
    fn new_person(&self, np: NewPerson) -> Result<Person, RepositoryError>;
    fn get_persons(&self) -> Result<Vec<Person>, RepositoryError>;
}

pub enum RepositoryError {
    Datasource(&'static str),
    Query(&'static str),
    Mutation(&'static str)
}

impl RepositoryError {
    pub fn type_message(&self) -> String {
        match self {
            RepositoryError::Mutation(_) => "Mutation error".to_string(),
            RepositoryError::Query(_) => "Query error".to_string(),
            RepositoryError::Datasource(_) => "Datasource error".to_string(),
        }
    }
}

impl ToString for RepositoryError {
    fn to_string(&self) -> String {
        match self {
            RepositoryError::Mutation(details) 
                | RepositoryError::Query(details) 
                | RepositoryError::Datasource(details) => format!("{}: {}", self.type_message(), details),
        }
    }
}
