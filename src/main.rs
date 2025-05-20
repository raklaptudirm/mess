use std::{str::FromStr, time};

use tetka::games::{common::perft::perft, games::chess, interface::PositionType};

fn main() {
    let position = chess::Position::from_str(chess::Position::STARTPOS).unwrap();

    let start = time::Instant::now();
    let nodes = perft::<true, false, _>(position, 6);
    let runtime = start.elapsed();

    println!(
        "nodes {} nps {}",
        nodes,
        (nodes as u128 * 1000) / runtime.as_millis()
    )
}
