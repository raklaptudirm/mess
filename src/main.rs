use std::env;

use tetka::uxi::Client;

mod commands;
mod core;
mod options;

fn main() {
    let client = Client::new()
        .protocol("uci")
        .engine("Mess v1.0.0")
        .author("Rak Laptudirm")
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

    let args = env::args()
        .skip(1)
        .reduce(|acc, e| format!("{} {}", acc, e))
        .unwrap_or("".to_string());
    let args = args.trim().to_string();
    if args.is_empty() {
        client.start();
    } else if let Err(err) = client.run_cmd_string(&args) {
        println!("{}", err);
    }
}
