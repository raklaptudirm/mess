use std::env;

use tetka::uxi::Client;

mod commands;
mod core;
mod options;

const IDENT: &str = concat!(env!("CARGO_PKG_NAME"), " v", env!("CARGO_PKG_VERSION"));
const AUTHOR: &str = env!("CARGO_PKG_AUTHORS");

fn main() {
    println!("{} by {}", IDENT, AUTHOR);

    let client = Client::new()
        .protocol("uci")
        .engine(IDENT)
        .author(AUTHOR)
        // Register engine options.
        .option("Hash", options::hash())
        .option("Threads", options::threads())
        // Register the custom commands.
        .command("d", commands::d())
        .command("go", commands::go())
        .command("bench", commands::bench())
        .command("perft", commands::perft_cmd())
        .command("protocol", commands::protocol())
        .command("position", commands::position())
        .command("ucinewgame", commands::ucinewgame());

    let args = env::args().skip(1).collect::<Vec<_>>().join(" ");
    if args.is_empty() {
        client.start();
    } else if let Err(err) = client.run_cmd_string(&args) {
        println!("{}", err);
    }
}
