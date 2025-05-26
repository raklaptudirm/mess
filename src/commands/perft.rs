use std::time;

use tetka::{
    games::{common::perft::perft, games::chess},
    uxi::{Bundle, Command, Flag, RunError, lock},
};

use super::Context;

pub fn perft_cmd() -> Command<Context> {
    Command::new(|bundle: Bundle<Context>| {
        lock! {
            bundle > ctx =>
            let position = ctx.position.clone();
        }

        let (depth, split, bulk) = parse_flags(&bundle)?;

        // If split perft is enabled, we just do the perft at the target depth.
        // Otherwise, we run perft for each depth starting from 1 to the target.
        let depths = match split {
            true => depth..=depth,
            false => 1..=depth,
        };

        for depth in depths {
            timed_perft(split, bulk, position.clone(), depth);
        }

        Ok(())
    })
    .flag("depth", Flag::Single)
    .flag("bulk", Flag::Boolean)
    .flag("split", Flag::Boolean)
    // This command should be run in a separate thread so that the Client
    // can still respond to and run other Commands while this one is running.
    .parallelize(true)
}

fn timed_perft(split: bool, bulk: bool, position: chess::Position, depth: u8) {
    let start = time::Instant::now();
    // Manual dynamic dispatch with the correct const generics.
    let nodes = match (split, bulk) {
        (true, true) => perft::<true, true, _>(position, depth),
        (true, false) => perft::<true, false, _>(position, depth),
        (false, true) => perft::<false, true, _>(position, depth),
        (false, false) => perft::<false, false, _>(position, depth),
    };
    let duration = start.elapsed();

    let time = duration.as_millis().max(1);
    let nps = nodes as u128 * 1000 / time;

    // Report the perft results.
    println!("info depth {depth} nodes {nodes} time {time} nps {nps}");
}

fn parse_flags(bundle: &Bundle<Context>) -> Result<(u8, bool, bool), RunError> {
    Ok((
        bundle.get_parsed_flag("depth")?.unwrap_or(6),
        bundle.is_flag_set("split"),
        bundle.is_flag_set("bulk"),
    ))
}
