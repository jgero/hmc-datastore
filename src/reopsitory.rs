use gremlin_client::{GremlinClient, GremlinError, process::traversal::traversal};
use uuid::Uuid;

use crate::model::person::{NewPerson, Person};

pub trait Repository {
    fn new_person(&self, np: NewPerson) -> Result<Person, String>;
    fn get_persons(&self) -> Result<Vec<Person>, String>;
}

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

pub struct GremlinRepository {
    conn_str: &'static str
}

impl GremlinRepository {
    fn get_client(&self) -> Result<GremlinClient, GremlinError> {
        GremlinClient::connect(self.conn_str)
    }
}

impl Repository for GremlinRepository {
    fn new_person(&self, np: NewPerson) -> Result<Person, String> {
        let np = Person::create(np)?;
        let client = self.get_client().map_err(|_| "could not connect to gremlin".to_string())?;
        let g = traversal().with_remote(client);
        g.add_v("person")
            .property("name", np.name)
            .property("uuid", np.uuid.to_string())
            .property("created", np.created)
            .property("updated", np.updated)
            .property("updated_count", np.update_count)
            .next()
            .map_err(|_| "could not create person in database".to_string())?;
        todo!()
    }

    fn get_persons(&self) -> Result<Vec<Person>, String> {
        let client = self.get_client().map_err(|_| "could not connect to gremlin".to_string())?;
        let g = traversal().with_remote(client);
        let persons = g.v(()).has_label("person").to_list()
            .map_err(|_| "could not query persons".to_string())?
            .iter().map(|v| {
                Person {
                    name: v.property("name").unwrap().get::<String>().unwrap().to_string(),
                    uuid: Uuid::parse_str(v.property("uuid").unwrap().get::<String>().unwrap()).unwrap(),
                    created: v.property("created").unwrap().get::<String>().unwrap().parse().unwrap(),
                    updated: v.property("updated").unwrap().get::<String>().unwrap().parse().unwrap(),
                    update_count: v.property("update_count").unwrap().get::<String>().unwrap().parse().unwrap(),
                    keywords: Vec::new()
                }
            }).collect();
        Ok(persons)
    }
}

unsafe impl Sync for GremlinRepository {}
unsafe impl Send for GremlinRepository {}
