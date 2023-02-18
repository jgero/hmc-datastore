use gremlin_client::{GremlinClient, GremlinError, process::traversal::traversal, List, Vertex, VertexProperty};
use uuid::Uuid;

use crate::model::person::{NewPerson, Person};

pub struct InMemoryRepository {
}

impl InMemoryRepository {
    pub fn new() -> InMemoryRepository {
        InMemoryRepository{}
    }
}

impl Repository for InMemoryRepository {
    fn new_person(&self, np: NewPerson) -> Result<Person, String> {
        let p = Person::create(np)?;
        Ok(p) 
    }

    fn get_persons(&self) -> Result<Vec<Person>, String> {
        Ok(vec![Person::create(NewPerson { name: "john".to_string(), keywords: vec![] }).unwrap()])
    }
}

unsafe impl Sync for InMemoryRepository {}
unsafe impl Send for InMemoryRepository {}

