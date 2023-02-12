use crate::model::{Foo, NewFoo};

pub trait Repository {
    fn new_person(&self, nf: NewFoo) -> Result<Foo, String>;
    fn get_foos(&self) -> Result<Vec<Foo>, String>;
}

pub struct GremlinRepository {}

impl Repository for GremlinRepository {
}

unsafe impl Sync for GremlinRepository {}
