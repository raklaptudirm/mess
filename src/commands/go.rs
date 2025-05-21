use std::time;

use tetka::{
    games::{common::perft::perft, games::chess, interface::PositionType},
    uxi::{Bundle, Command, Flag, RunError, error},
};

use crate::core::Limits;

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

pub fn go() -> Command<Context> {
    Command::new(|bundle: Bundle<Context>| {
        lock! {
            bundle > ctx =>
            let position = ctx.position.clone(); // Get the position to search
            let mut searcher = ctx.searcher.clone(); // Get the previous search state
        }

        let mut nodes = 0;

        match parse_limits(&bundle, &position)? {
            // Search flags received, search the position.
            Config::Search(limits) => {
                let bestmove = searcher.search(&position, limits, &mut nodes);

                println!("bestmove {}", bestmove);

                lock! {
                    bundle > mut ctx =>
                    // Push the new search state to the context.
                    ctx.searcher = searcher;
                }

                Ok(())
            }

            // Perft flags received, run a perft on the position.
            Config::Perft(bulk, max_depth) => {
                for depth in 1..=max_depth {
                    let start = time::Instant::now();
                    let nodes = if bulk {
                        perft::<false, true, _>(position.clone(), depth)
                    } else {
                        perft::<false, false, _>(position.clone(), depth)
                    };
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

                Ok(())
            }
        }
    })
    // Flags for reporting the current time situation.
    .flag("binc", Flag::Single)
    .flag("winc", Flag::Single)
    .flag("btime", Flag::Single)
    .flag("wtime", Flag::Single)
    .flag("movestogo", Flag::Single)
    // Flags for imposing other search limits.
    .flag("depth", Flag::Single)
    .flag("nodes", Flag::Single)
    .flag("movetime", Flag::Single)
    // Flags for setting the search type.
    // .flag("ponder", Flag::Single)
    .flag("infinite", Flag::Single)
    // Flags for go perft command.
    .flag("perft", Flag::Single)
    .flag("bulk", Flag::Boolean)
    // This command should be run in a separate thread so that the Client
    // can still respond to and run other Commands while this one is running.
    .parallelize(true)
}

enum Config {
    Perft(bool, u8),
    Search(Limits),
}

fn parse_limits(bundle: &Bundle<Context>, position: &chess::Position) -> Result<Config, RunError> {
    ////////////////////////////////////////////
    // Check which of the limit flags are set //
    ////////////////////////////////////////////

    // Standard Time Control flags
    let binc = bundle.is_flag_set("binc");
    let winc = bundle.is_flag_set("winc");
    let btime = bundle.is_flag_set("btime");
    let wtime = bundle.is_flag_set("wtime");

    // Other limit flags
    let depth = bundle.is_flag_set("depth");
    let nodes = bundle.is_flag_set("nodes");
    let movetime = bundle.is_flag_set("movetime");

    // Infinite flag
    let infinite = bundle.is_flag_set("infinite");

    // Perft flags
    let perft = bundle.is_flag_set("perft");
    let bulk = bundle.is_flag_set("bulk");

    ///////////////////////////////////////////////////////
    // Ensure that the given flag configuration is valid //
    ///////////////////////////////////////////////////////

    let std_tc = btime || wtime || binc || winc;
    let oth_tc = depth || nodes || movetime;

    // Either all or none of the standard time control flags must be set.
    if std_tc && !(btime && wtime && binc && winc) {
        return error!("bad flag set: missing standard time control flags");
    }

    // No other time control flags may be set alongside 'infinite'.
    if infinite && (std_tc || oth_tc) {
        return error!("bad flag set: time control flags set alongside infinite");
    }

    // The bulk flag can't be set without the perft flag.
    if !perft && bulk {
        return error!("bad flag set: bulk flag set without perft flag");
    }

    // The perft flag can't be set alongside time control flags.
    if perft && (std_tc || oth_tc || infinite) {
        return error!("bad flag set: time control flags set alongside perft");
    }

    // A little utility macro to parse the given flag into the required type.
    macro_rules! get_flag {
        ($name:expr) => {
            match bundle.get_single_flag($name) {
                Some(value) => Some(value.parse()?),
                None => None,
            }
        };
    }

    ////////////////////////////////////////////
    // Parse the provided search/perft limits //
    ////////////////////////////////////////////

    if perft {
        Ok(Config::Perft(bulk, get_flag!("perft").unwrap()))
    } else {
        Ok(Config::Search(Limits {
            maxnodes: get_flag!("nodes"),
            maxdepth: get_flag!("depth"),
            movetime: if std_tc {
                let (time, incr) = match position.side_to_move() {
                    chess::Color::Black => ("btime", "binc"),
                    chess::Color::White => ("wtime", "winc"),
                };

                let time: u128 = get_flag!(time).unwrap();
                let incr: u128 = get_flag!(incr).unwrap();

                Some((time / 20 + incr / 2).max(1))
            } else {
                get_flag!("movetime")
            },
            movestogo: get_flag!("movestogo"),
        }))
    }
}
