use std::time;

use tetka::{
    games::{common::perft::perft, games::chess},
    uxi::{Bundle, Command, Flag, RunError},
};

use super::Context;

// TODO: Move these macros into UXI

macro_rules! lock {
    ($bundle:ident > $ctx:ident => $($stmt:stmt;)*) => {
        let $ctx = $bundle.lock();
        $(
            $stmt
        )*
        drop($ctx);
    };

    ($bundle:ident > mut $ctx:ident => $($stmt:stmt;)*) => {
        let mut $ctx = $bundle.lock();
        $(
            $stmt
        )*
        drop($ctx);
    };
}

pub fn perft_cmd() -> Command<Context> {
    Command::new(|bundle: Bundle<Context>| {
        lock! {
            bundle > ctx =>
            let position = ctx.position.clone(); // Get the position to search
        }

        let (depth, bulk, split) = parse_flags(&bundle)?;

        if split {
            if bulk {
                timed_perft::<true, true>(position.clone(), depth);
            } else {
                timed_perft::<true, false>(position.clone(), depth);
            }
        } else {
            let max_depth = depth;
            for depth in 1..=max_depth {
                if bulk {
                    timed_perft::<false, true>(position.clone(), depth);
                } else {
                    timed_perft::<false, false>(position.clone(), depth);
                }
            }
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

fn timed_perft<const SPLIT: bool, const BULK: bool>(position: chess::Position, depth: u8) {
    let start = time::Instant::now();
    let nodes = perft::<SPLIT, BULK, _>(position.clone(), depth);
    let duration = start.elapsed();

    let time = duration.as_millis().max(1);

    println!(
        "info depth {} nodes {} time {} nps {}",
        depth,
        nodes,
        time,
        1000 * nodes as u128 / time
    );
}

fn parse_flags(bundle: &Bundle<Context>) -> Result<(u8, bool, bool), RunError> {
    macro_rules! get_flag {
        ($name:expr) => {
            match bundle.get_single_flag($name) {
                Some(value) => Some(value.parse()?),
                None => None,
            }
        };
    }

    Ok((
        get_flag!("depth").unwrap_or(6),
        get_flag!("bulk").unwrap_or(false),
        get_flag!("split").unwrap_or(false),
    ))
}
