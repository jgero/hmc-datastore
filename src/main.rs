#[macro_use] extern crate juniper;

use std::env;
use crate::schema::{Context, new_schema};
use crate::repository::gremlin::GremlinRepository;
use warp::{http::Response, Filter};

mod schema;
mod model;
mod repository;

#[tokio::main]
async fn main() {
    env::set_var("RUST_LOG", "warp_server");
    env_logger::init();

    let log = warp::log("warp_server");

    let homepage = warp::path::end().map(|| {
        Response::builder()
            .header("content-type", "text/html")
            .body(format!(
                "<html><h1>juniper_warp</h1><div>visit <a href=\"/graphiql\">/graphiql</a></html>"
            ))
    });

    log::info!("Listening on 127.0.0.1:8080");

    // let state = warp::any().map(move || Context { repo: Box::new(InMemoryRepository::new())  });
    let state = warp::any().map(move || Context { repo: Box::new(GremlinRepository::new("localhost"))  });
    let graphql_filter = juniper_warp::make_graphql_filter(new_schema(), state.boxed());

    warp::serve(
        warp::get()
            .and(warp::path("graphiql"))
            .and(juniper_warp::graphiql_filter("/graphql", None))
            .or(homepage)
            .or(warp::path("graphql").and(graphql_filter))
            .with(log),
    )
    .run(([127, 0, 0, 1], 8080))
    .await
}
