use gremlin_client::{GremlinClient, GremlinError, Vertex, process::traversal::traversal};
use uuid::Uuid;

use crate::model::person::{NewPerson, Person};

use super::interface::{Repository, RepositoryError};

pub struct GremlinRepository {
    conn_str: &'static str
}

impl GremlinRepository {
    pub fn new(conn_str: &'static str) -> GremlinRepository {
        GremlinRepository { conn_str }
    }
    fn get_client(&self) -> Result<GremlinClient, GremlinError> {
        GremlinClient::connect(self.conn_str)
    }
}

impl Repository for GremlinRepository {
    fn new_person(&self, np: NewPerson) -> Result<Person, RepositoryError> {
        let np = Person::create(np).map_err(|_| RepositoryError::Datasource("weird stuff"))?;
        let client = self.get_client().map_err(|_| RepositoryError::Datasource("could not connect to gremlin"))?;
        let g = traversal().with_remote(client);
        g.add_v("person")
            .property("name", np.name.clone())
            .property("uuid", np.uuid.to_string())
            .property("created", np.created)
            .property("updated", np.updated)
            .property("update_count", np.update_count)
            .next()
            .map_err(|_| RepositoryError::Mutation("could not create person in database"))?;
        Ok(np)
    }

    fn get_persons(&self) -> Result<Vec<Person>, RepositoryError> {
        let client = self.get_client().map_err(|_| RepositoryError::Datasource("could not get gremlin client"))?;
        // TODO: don't just skip broken entries, handle them or at least log
        let persons = client.execute("g.V().hasLabel(l).toList()", &[("l",&"person")]).map_err(|_| RepositoryError::Query("could not execute persons query"))?
            .filter_map(Result::ok)
            .map(|f| f.take::<Vertex>())
            .filter_map(Result::ok)
            .collect::<Vec<Vertex>>()
            .iter()
            .map(|v| vertex_to_person(v))
            .flatten()
            .collect::<Vec<Person>>();
        Ok(persons)
    }
}

fn extract_vertex_property<T : gremlin_client::FromGValue>(key: &str, v: &Vertex) -> Option<T> {
    v.property(key).and_then(|v| Some( v.value() )).and_then(|v| match v.clone().take::<T>() {
        Ok(val) => Some(val),
        Err(_) => None
    })
}

fn vertex_to_person(v: &Vertex) -> Option<Person> {
    Some(
    Person { 
        name: extract_vertex_property::<String>("name", v)?,
        uuid: Uuid::parse_str(&extract_vertex_property::<String>("uuid", v)?).ok()?,
        keywords: Vec::new(),
        created: extract_vertex_property::<f64>("created", v)?,
        updated: extract_vertex_property::<f64>("updated", v)?,
        update_count: extract_vertex_property::<i32>("update_count", v)?,
    })
}

unsafe impl Sync for GremlinRepository {}
unsafe impl Send for GremlinRepository {}
